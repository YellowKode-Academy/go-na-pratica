package link

import (
	"net/url"
	"time"
)

// Link representa um link salvo no LinkVault.
type Link struct {
	ID        int64     `db:"id"         json:"id"`
	URL       string    `db:"url"        json:"url"`
	Title     string    `db:"title"      json:"title"`
	Tags      string    `db:"tags"       json:"tags"`   // CSV separado por virgula
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// ValidationError representa erro de validacao com campo e motivo.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Reason
}

// NewLink cria um Link com CreatedAt preenchido.
func NewLink(rawURL, title, tags string) Link {
	return Link{
		URL:       rawURL,
		Title:     title,
		Tags:      tags,
		CreatedAt: time.Now(),
	}
}

// Validate valida os campos obrigatorios do Link.
func (l Link) Validate() error {
	if l.URL == "" {
		return &ValidationError{Field: "url", Reason: "obrigatoria"}
	}
	u, err := url.ParseRequestURI(l.URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return &ValidationError{Field: "url", Reason: "formato invalido"}
	}
	if l.Title == "" {
		return &ValidationError{Field: "title", Reason: "obrigatorio"}
	}
	if len(l.Title) > 200 {
		return &ValidationError{Field: "title", Reason: "maximo 200 caracteres"}
	}
	return nil
}
