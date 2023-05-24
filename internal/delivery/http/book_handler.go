package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	repoErr "github.com/LuigiAzevedo/public-library-v2/internal/database/repository"
	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	uc "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type bookHandler struct {
	BookUsecase uc.BookUsecase
}

// NewBookHandler creates a new instance of bookHandler
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
		log.Error().Msg(err.Error())
		http.Error(w, invalidBookID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	b, err := h.BookUsecase.GetBook(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			if err == repoErr.ErrBookNotFound {
				http.Error(w, bookNotFound, http.StatusNotFound)
			} else {
				http.Error(w, getBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(b); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, getBook, http.StatusInternalServerError)
		return
	}
}

type SearchBookRequest struct {
	Title string `json:"title"`
}

func (h *bookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	var req SearchBookRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil && err != io.EOF {
		log.Error().Msg(err.Error())
		http.Error(w, wrongBodyTitle, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var b []*entity.Book
	if req.Title == "" {
		b, err = h.BookUsecase.ListBooks(ctx)
	} else {
		b, err = h.BookUsecase.SearchBooks(ctx, req.Title)
	}

	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			if err == repoErr.ErrBookNotFound {
				http.Error(w, bookNotFound, http.StatusNotFound)
			} else {
				http.Error(w, searchBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(b); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, searchBook, http.StatusInternalServerError)
		return
	}
}

func (h *bookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var b entity.Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := h.BookUsecase.CreateBook(ctx, &b)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			http.Error(w, createBook, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, createBook, http.StatusInternalServerError)
		return
	}
}

func (h *bookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var b entity.Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	b.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, invalidBookID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.BookUsecase.UpdateBook(ctx, &b)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			if err == repoErr.ErrBookNotFound {
				http.Error(w, bookNotFound, http.StatusNotFound)
			} else {
				http.Error(w, updateBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *bookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, invalidBookID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.BookUsecase.DeleteBook(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			if err == repoErr.ErrBookNotFound {
				http.Error(w, bookNotFound, http.StatusNotFound)
			} else {
				http.Error(w, deleteBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
