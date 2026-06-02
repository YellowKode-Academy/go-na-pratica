package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yellowkode-academy/linkvault/internal/link"
	"github.com/yellowkode-academy/linkvault/internal/storage"
)

func TestInMemoryRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("save e find by id", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		l, err := repo.Save(ctx, link.NewLink("https://go.dev", "Go", "go"))
		require.NoError(t, err)
		assert.Greater(t, l.ID, int64(0))

		found, err := repo.FindByID(ctx, l.ID)
		require.NoError(t, err)
		assert.Equal(t, l.URL, found.URL)
	})

	t.Run("find not found", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		_, err := repo.FindByID(ctx, 9999)
		assert.ErrorIs(t, err, storage.ErrLinkNotFound)
	})

	t.Run("url duplicada", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		repo.Save(ctx, link.NewLink("https://go.dev", "Go", ""))
		_, err := repo.Save(ctx, link.NewLink("https://go.dev", "Go2", ""))
		assert.ErrorIs(t, err, storage.ErrURLDuplicada)
	})

	t.Run("list", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		repo.Save(ctx, link.NewLink("https://go.dev", "Go", ""))
		repo.Save(ctx, link.NewLink("https://github.com", "GitHub", ""))
		links, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Len(t, links, 2)
	})

	t.Run("search", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		repo.Save(ctx, link.NewLink("https://go.dev", "Go oficial", "go,docs"))
		repo.Save(ctx, link.NewLink("https://github.com", "GitHub", "git"))
		links, err := repo.Search(ctx, "go")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(links), 1)
	})

	t.Run("delete", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		l, _ := repo.Save(ctx, link.NewLink("https://go.dev", "Go", ""))
		err := repo.Delete(ctx, l.ID)
		require.NoError(t, err)
		_, err = repo.FindByID(ctx, l.ID)
		assert.ErrorIs(t, err, storage.ErrLinkNotFound)
	})

	t.Run("delete not found", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		err := repo.Delete(ctx, 9999)
		assert.ErrorIs(t, err, storage.ErrLinkNotFound)
	})
}
