package link_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yellowkode-academy/linkvault/internal/link"
)

func TestNewLink(t *testing.T) {
	l := link.NewLink("https://go.dev", "Go oficial", "go,docs")
	assert.Equal(t, "https://go.dev", l.URL)
	assert.Equal(t, "Go oficial", l.Title)
	assert.Equal(t, "go,docs", l.Tags)
	assert.False(t, l.CreatedAt.IsZero())
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		title   string
		wantErr bool
		errMsg  string
	}{
		{"valido", "https://go.dev", "Go", false, ""},
		{"url vazia", "", "Go", true, "url"},
		{"url invalida", "nao-e-url", "Go", true, "url"},
		{"title vazio", "https://go.dev", "", true, "title"},
		{"title longo", "https://go.dev", string(make([]byte, 201)), true, "title"},
		{"http valido", "http://localhost:8080", "Local", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link.NewLink(tt.url, tt.title, "")
			err := l.Validate()
			if tt.wantErr {
				require.Error(t, err)
				var ve *link.ValidationError
				require.ErrorAs(t, err, &ve)
				assert.Equal(t, tt.errMsg, ve.Field)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
