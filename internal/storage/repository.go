package storage

import (
	"context"
	"errors"

	"github.com/yellowkode-academy/linkvault/internal/link"
)

// ErrLinkNotFound e retornado quando o link nao existe.
var ErrLinkNotFound = errors.New("link nao encontrado")

// ErrURLDuplicada e retornado quando a URL ja esta cadastrada.
var ErrURLDuplicada = errors.New("URL ja cadastrada")

// LinkRepository define operacoes de persistencia de links.
type LinkRepository interface {
	Save(ctx context.Context, l link.Link) (link.Link, error)
	FindByID(ctx context.Context, id int64) (link.Link, error)
	List(ctx context.Context) ([]link.Link, error)
	Search(ctx context.Context, query string) ([]link.Link, error)
	Delete(ctx context.Context, id int64) error
}
