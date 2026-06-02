package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yellowkode-academy/linkvault/internal/api"
	"github.com/yellowkode-academy/linkvault/internal/link"
	"github.com/yellowkode-academy/linkvault/internal/storage"
)

func setupServer(t *testing.T) (*httptest.Server, *storage.InMemoryRepository) {
	t.Helper()
	repo := storage.NewInMemoryRepository()
	mux := http.NewServeMux()
	api.NewLinkHandler(repo).RegisterRoutes(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, repo
}

func TestHealth(t *testing.T) {
	srv, _ := setupServer(t)
	resp, err := http.Get(srv.URL + "/health")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateAndListLinks(t *testing.T) {
	srv, _ := setupServer(t)

	body := `{"url":"https://go.dev","title":"Go oficial","tags":"go,docs"}`
	resp, err := http.Post(srv.URL+"/links", "application/json", bytes.NewBufferString(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created link.Link
	json.NewDecoder(resp.Body).Decode(&created)
	assert.Equal(t, "https://go.dev", created.URL)
	assert.Greater(t, created.ID, int64(0))

	resp, err = http.Get(srv.URL + "/links")
	require.NoError(t, err)
	var links []link.Link
	json.NewDecoder(resp.Body).Decode(&links)
	assert.Len(t, links, 1)
}

func TestCreateInvalidLink(t *testing.T) {
	srv, _ := setupServer(t)
	body := `{"url":"nao-e-url","title":"Test"}`
	resp, err := http.Post(srv.URL+"/links", "application/json", bytes.NewBufferString(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateDuplicateURL(t *testing.T) {
	srv, _ := setupServer(t)
	body := `{"url":"https://github.com","title":"GitHub"}`
	http.Post(srv.URL+"/links", "application/json", bytes.NewBufferString(body))
	resp, err := http.Post(srv.URL+"/links", "application/json", bytes.NewBufferString(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestGetByID(t *testing.T) {
	srv, repo := setupServer(t)
	saved, _ := repo.Save(context.Background(), link.NewLink("https://go.dev", "Go", ""))
	resp, err := http.Get(srv.URL + "/links/" + strconv.FormatInt(saved.ID, 10))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetByIDNotFound(t *testing.T) {
	srv, _ := setupServer(t)
	resp, err := http.Get(srv.URL + "/links/9999")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeleteLink(t *testing.T) {
	srv, repo := setupServer(t)
	saved, _ := repo.Save(context.Background(), link.NewLink("https://go.dev", "Go", ""))

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/links/"+strconv.FormatInt(saved.ID, 10), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	resp, err = http.Get(srv.URL + "/links/" + strconv.FormatInt(saved.ID, 10))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestSearchLinks(t *testing.T) {
	srv, repo := setupServer(t)
	repo.Save(context.Background(), link.NewLink("https://go.dev", "Go oficial", "go"))
	repo.Save(context.Background(), link.NewLink("https://github.com", "GitHub", "git"))

	resp, err := http.Get(srv.URL + "/links?q=go")
	require.NoError(t, err)
	var links []link.Link
	json.NewDecoder(resp.Body).Decode(&links)
	assert.GreaterOrEqual(t, len(links), 1)
}
