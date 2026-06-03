package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/yellowkode-academy/linkvault/internal/link"
)

const schema = `
CREATE TABLE IF NOT EXISTS links (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	url        TEXT    NOT NULL UNIQUE,
	title      TEXT    NOT NULL DEFAULT '',
	tags       TEXT    NOT NULL DEFAULT '',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

// SQLiteRepository persiste links no SQLite.
type SQLiteRepository struct {
	db *sqlx.DB
}

// normalizeDSN converte um caminho simples para o formato URI exigido pelo modernc.org/sqlite v1.51+.
func normalizeDSN(dsn string) string {
	if dsn == ":memory:" || strings.HasPrefix(dsn, "file:") {
		return dsn
	}
	return "file:" + dsn
}

// NewSQLiteRepository abre o banco e aplica o schema.
func NewSQLiteRepository(dsn string) (*SQLiteRepository, error) {
	db, err := sqlx.Open("sqlite", normalizeDSN(dsn))
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("criar schema: %w", err)
	}
	return &SQLiteRepository{db: db}, nil
}

// Close fecha a conexao com o banco.
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

func (r *SQLiteRepository) Save(ctx context.Context, l link.Link) (link.Link, error) {
	query := `INSERT INTO links (url, title, tags, created_at) VALUES (:url, :title, :tags, :created_at)`
	res, err := r.db.NamedExecContext(ctx, query, l)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return link.Link{}, ErrURLDuplicada
		}
		return link.Link{}, fmt.Errorf("salvar link: %w", err)
	}
	id, _ := res.LastInsertId()
	l.ID = id
	return l, nil
}

func (r *SQLiteRepository) FindByID(ctx context.Context, id int64) (link.Link, error) {
	var l link.Link
	err := r.db.GetContext(ctx, &l, "SELECT * FROM links WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return link.Link{}, ErrLinkNotFound
	}
	if err != nil {
		return link.Link{}, fmt.Errorf("buscar link: %w", err)
	}
	return l, nil
}

func (r *SQLiteRepository) List(ctx context.Context) ([]link.Link, error) {
	var links []link.Link
	if err := r.db.SelectContext(ctx, &links, "SELECT * FROM links ORDER BY created_dat DESC"); err != nil {
		return nil, fmt.Errorf("listar links: %w", err)
	}
	return links, nil
}

func (r *SQLiteRepository) Search(ctx context.Context, query string) ([]link.Link, error) {
	var links []link.Link
	q := "%" + query + "%"
	sqlQuery := "SELECT * FROM links WHERE url LIKE ? OR title LIKE ? OR tags LIKE ? ORDER BY created_at DESC"
	if err := r.db.SelectContext(ctx, &links, sqlQuery, q, q, q); err != nil {
		return nil, fmt.Errorf("buscar links: %w", err)
	}
	return links, nil
}

func (r *SQLiteRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM links WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deletar link: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrLinkNotFound
	}
	return nil
}

func isUniqueConstraintErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") || strings.Contains(msg, "unique")
}
