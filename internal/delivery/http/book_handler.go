package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	uc "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type bookHandler struct {
	BookUsecase uc.BookUsecase
}

// NewbookHandler creates a new instance of bookHandler
func NewBookHandler(r *chi.Mux, useCase uc.BookUsecase) {
	handler := &bookHandler{
		BookUsecase: useCase,
	}

	r.Route("/v1/books", func(r chi.Router) {
		r.Get("/{id}", handler.GetBook)
		r.Get("/", handler.SearchBooks)
		r.Post("/", handler.CreateBook)
		r.Put("/{id}", handler.UpdateBook)
		r.Delete("/{id}", handler.DeleteBook)
	})
}

func (h *bookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	b, err := h.BookUsecase.GetBook(ctx, id)
	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	json.NewEncoder(w).Encode(b)
}

type SearchBookRequest struct {
	Title string `json:"title"`
}

func (h *bookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	var req SearchBookRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	ctx := r.Context()
	var b []*entity.Book
	if req.Title == "" {
		b, err = h.BookUsecase.ListBooks(ctx)
	} else {
		b, err = h.BookUsecase.SearchBooks(ctx, req.Title)
	}

	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	json.NewEncoder(w).Encode(b)
}

func (h *bookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var b *entity.Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := h.BookUsecase.CreateBook(ctx, b)
	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *bookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var b *entity.Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.BookUsecase.UpdateBook(ctx, b)
	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusGatewayTimeout)
		default:
			if errors.Cause(err).Error() == "book not found" {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *bookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.BookUsecase.DeleteBook(ctx, id)
	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
