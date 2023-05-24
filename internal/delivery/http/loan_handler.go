package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	repoErr "github.com/LuigiAzevedo/public-library-v2/internal/database/repository"
	ucErr "github.com/LuigiAzevedo/public-library-v2/internal/domain/usecase"
	uc "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type loanHandler struct {
	LoanUsecase uc.LoanUsecase
}

// NewLoanHandler creates a new instance of loanHandler
func NewLoanHandler(r *chi.Mux, useCase uc.LoanUsecase) {
	handler := &loanHandler{
		LoanUsecase: useCase,
	}

	r.Route("/v1/loans", func(r chi.Router) {
		r.Get("/{id}", handler.SearchUserLoans)
		r.Post("/borrow", handler.BorrowBook)
		r.Post("/return", handler.ReturnBook)
	})
}

func (h *loanHandler) SearchUserLoans(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, invalidUserID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	l, err := h.LoanUsecase.SearchUserLoans(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			if err == repoErr.ErrLoanNotFound {
				http.Error(w, loanNotFound, http.StatusNotFound)
			} else {
				http.Error(w, searchUserLoans, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(l); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, searchUserLoans, http.StatusInternalServerError)
		return
	}
}

type LoanRequest struct {
	UserID int `json:"user_id"`
	BookID int `json:"book_id"`
}

func (h *loanHandler) BorrowBook(w http.ResponseWriter, r *http.Request) {
	var req LoanRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil && err != io.EOF {
		log.Error().Msg(err.Error())
		http.Error(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.LoanUsecase.BorrowBook(ctx, req.UserID, req.BookID)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			switch err {
			case repoErr.ErrBookNotFound:
				http.Error(w, bookNotFound, http.StatusNotFound)
			case repoErr.ErrUserNotFound:
				http.Error(w, userNotFound, http.StatusNotFound)
			case ucErr.ErrReturnBookFirst:
				http.Error(w, returnBookFirst, http.StatusBadRequest)
			case ucErr.ErrBookUnavailable:
				http.Error(w, bookUnavailable, http.StatusNotFound)
			default:
				http.Error(w, borrowBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *loanHandler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	var req LoanRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil && !errors.Is(err, io.EOF) {
		log.Error().Msg(err.Error())
		http.Error(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.LoanUsecase.ReturnBook(ctx, req.UserID, req.BookID)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			if err == ucErr.ErrLoanAlreadyReturned {
				http.Error(w, loanAlreadyReturned, http.StatusNotFound)
			} else {
				http.Error(w, returnBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
