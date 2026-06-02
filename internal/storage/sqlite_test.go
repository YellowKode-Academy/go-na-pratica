package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yellowkode-academy/linkvault/internal/link"
	"github.com/yellowkode-academy/linkvault/internal/storage"
)

func TestSQLiteRepository(t *testing.T) {
	ctx := context.Background()

	newRepo := func(t *testing.T) *storage.SQLiteRepository {
		t.Helper()
		repo, err := storage.NewSQLiteRepository(":memory:")
		require.NoError(t, err)
		t.Cleanup(func() { repo.Close() })
		return repo
	}

	t.Run("save e find by id", func(t *testing.T) {
		repo := newRepo(t)
		l, err := repo.Save(ctx, link.NewLink("https://go.dev", "Go", "go"))
		require.NoError(t, err)
		assert.Greater(t, l.ID, int64(0))

		found, err := repo.FindByID(ctx, l.ID)
		require.NoError(t, err)
		assert.Equal(t, "https://go.dev", found.URL)
	})

	t.Run("url duplicada", func(t *testing.T) {
		repo := newRepo(t)
		repo.Save(ctx, link.NewLink("https://go.dev", "Go", ""))
		_, err := repo.Save(ctx, link.NewLink("https://go.dev", "Go2", ""))
		assert.ErrorIs(t, err, storage.ErrURLDuplicada)
	})

	t.Run("not found", func(t *testing.T) {
		repo := newRepo(t)
		_, err := repo.FindByID(ctx, 9999)
		assert.ErrorIs(t, err, storage.ErrLinkNotFound)
	})

	t.Run("list e search", func(t *testing.T) {
		repo := newRepo(t)
		repo.Save(ctx, link.NewLink("https://go.dev", "Go oficial", "go"))
		repo.Save(ctx, link.NewLink("https://github.com", "GitHub", "git"))

		links, _ := repo.List(ctx)
		assert.Len(t, links, 2)

		found, _ := repo.Search(ctx, "go")
		assert.GreaterOrEqual(t, len(found), 1)
	})

	t.Run("delete", func(t *testing.T) {
		repo := newRepo(t)
		l, _ := repo.Save(ctx, link.NewLink("https://go.dev", "Go", ""))
		require.NoError(t, repo.Delete(ctx, l.ID))
		_, err := repo.FindByID(ctx, l.ID)
		assert.ErrorIs(t, err, storage.ErrLinkNotFound)
	})
}
