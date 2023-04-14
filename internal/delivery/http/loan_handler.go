package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

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
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	u, err := h.LoanUsecase.SearchUserLoans(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

type LoanRequest struct {
	UserID int `json:"user_id"`
	BookID int `json:"book_id"`
}

func (h *loanHandler) BorrowBook(w http.ResponseWriter, r *http.Request) {
	var req LoanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if !errors.Is(err, io.EOF) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err := h.LoanUsecase.BorrowBook(req.UserID, req.BookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *loanHandler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	var req LoanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if !errors.Is(err, io.EOF) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err := h.LoanUsecase.ReturnBook(req.UserID, req.BookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
