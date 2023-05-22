package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	"github.com/LuigiAzevedo/public-library-v2/internal/errs"
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
		http.Error(w, errs.ErrInvalidUserID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	l, err := h.LoanUsecase.SearchUserLoans(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if errors.Unwrap(err).Error() == errs.ErrNoLoansFound {
				http.Error(w, errs.ErrNoLoansFound, http.StatusNotFound)
			} else {
				http.Error(w, errs.ErrSearchUserLoans, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(l); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrSearchUserLoans, http.StatusInternalServerError)
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
		http.Error(w, errs.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.LoanUsecase.BorrowBook(ctx, req.UserID, req.BookID)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			error := strings.Split(err.Error(), ":")
			switch error[0] {
			case errs.ErrGetBook:
				http.Error(w, errs.ErrBookNotFound, http.StatusNotFound)
			case errs.ErrGetUser:
				http.Error(w, errs.ErrUserNotFound, http.StatusNotFound)
			case errs.ErrReturnBookFirst:
				http.Error(w, errs.ErrReturnBookFirst, http.StatusBadRequest)
			case errs.ErrBookUnavailable:
				http.Error(w, errs.ErrBookUnavailable, http.StatusNotFound)
			default:
				http.Error(w, errs.ErrBorrowBook, http.StatusInternalServerError)
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
		http.Error(w, errs.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.LoanUsecase.ReturnBook(ctx, req.UserID, req.BookID)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if err.Error() == errs.ErrLoanAlreadyReturned {
				http.Error(w, errs.ErrLoanAlreadyReturned, http.StatusNotFound)
			} else {
				http.Error(w, errs.ErrReturnBook, http.StatusInternalServerError)
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
