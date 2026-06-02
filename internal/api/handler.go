package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/yellowkode-academy/linkvault/internal/link"
	"github.com/yellowkode-academy/linkvault/internal/storage"
)

// LinkHandler implementa os handlers HTTP para links.
type LinkHandler struct {
	repo storage.LinkRepository
}

// NewLinkHandler cria um handler com o repositorio fornecido.
func NewLinkHandler(repo storage.LinkRepository) *LinkHandler {
	return &LinkHandler{repo: repo}
}

// RegisterRoutes registra os endpoints no mux.
func (h *LinkHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /links", h.list)
	mux.HandleFunc("POST /links", h.create)
	mux.HandleFunc("GET /links/{id}", h.getByID)
	mux.HandleFunc("DELETE /links/{id}", h.delete)
}

func (h *LinkHandler) health(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *LinkHandler) list(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	var (
		links []link.Link
		err   error
	)
	if query != "" {
		links, err = h.repo.Search(r.Context(), query)
	} else {
		links, err = h.repo.List(r.Context())
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if links == nil {
		links = []link.Link{}
	}
	respond(w, http.StatusOK, links)
}

func (h *LinkHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL   string `json:"url"`
		Title string `json:"title"`
		Tags  string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "json invalido")
		return
	}
	l := link.NewLink(req.URL, req.Title, req.Tags)
	if err := l.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.repo.Save(r.Context(), l)
	if err != nil {
		if errors.Is(err, storage.ErrURLDuplicada) {
			respondError(w, http.StatusConflict, "URL ja cadastrada")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, saved)
}

func (h *LinkHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "id invalido")
		return
	}
	l, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrLinkNotFound) {
			respondError(w, http.StatusNotFound, "link nao encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, l)
}

func (h *LinkHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "id invalido")
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, storage.ErrLinkNotFound) {
			respondError(w, http.StatusNotFound, "link nao encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respond(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	respond(w, code, map[string]string{"error": msg})
}
