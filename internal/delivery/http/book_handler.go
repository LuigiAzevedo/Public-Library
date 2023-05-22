package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	"github.com/LuigiAzevedo/public-library-v2/internal/errs"
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
		http.Error(w, errs.ErrInvalidBookID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	b, err := h.BookUsecase.GetBook(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if errors.Unwrap(err).Error() == errs.ErrBookNotFound {
				http.Error(w, errs.ErrBookNotFound, http.StatusNotFound)
			} else {
				http.Error(w, errs.ErrGetBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(b); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrGetBook, http.StatusInternalServerError)
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
		http.Error(w, errs.ErrWrongBodyTitle, http.StatusBadRequest)
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
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if errors.Unwrap(err).Error() == errs.ErrBookNotFound {
				http.Error(w, errs.ErrBookNotFound, http.StatusNotFound)
			} else {
				http.Error(w, errs.ErrSearchBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(b); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrSearchBook, http.StatusInternalServerError)
		return
	}
}

func (h *bookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var b entity.Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := h.BookUsecase.CreateBook(ctx, &b)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			http.Error(w, errs.ErrCreateBook, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrCreateBook, http.StatusInternalServerError)
		return
	}
}

func (h *bookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var b entity.Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	b.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrInvalidBookID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.BookUsecase.UpdateBook(ctx, &b)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if errors.Unwrap(err).Error() == errs.ErrBookNotFound {
				http.Error(w, errs.ErrBookNotFound, http.StatusNotFound)
			} else {
				http.Error(w, errs.ErrUpdateBook, http.StatusInternalServerError)
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
		http.Error(w, errs.ErrInvalidBookID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.BookUsecase.DeleteBook(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if errors.Unwrap(err).Error() == errs.ErrBookNotFound {
				http.Error(w, errs.ErrBookNotFound, http.StatusNotFound)
			} else {
				http.Error(w, errs.ErrDeleteBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
