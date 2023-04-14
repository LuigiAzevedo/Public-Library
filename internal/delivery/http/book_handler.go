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

	b, err := h.BookUsecase.GetBook(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

	var b []*entity.Book
	if req.Title == "" {
		b, err = h.BookUsecase.ListBooks()
	} else {
		b, err = h.BookUsecase.SearchBooks(req.Title)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

	id, err := h.BookUsecase.CreateBook(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	err = h.BookUsecase.UpdateBook(b)
	if err != nil {
		if errors.Cause(err).Error() == "book not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *bookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	err = h.BookUsecase.DeleteBook(id)
	if err != nil {
		if errors.Cause(err).Error() == "book not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
