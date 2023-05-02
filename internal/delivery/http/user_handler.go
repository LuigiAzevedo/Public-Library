package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	"github.com/LuigiAzevedo/public-library-v2/internal/errs"
	uc "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type userHandler struct {
	UserUsecase uc.UserUsecase
}

// NewUserHandler creates a new instance of userHandler
func NewUserHandler(r *chi.Mux, useCase uc.UserUsecase) {
	handler := &userHandler{
		UserUsecase: useCase,
	}

	r.Route("/v1/users", func(r chi.Router) {
		r.Get("/{id}", handler.GetUser)
		r.Post("/", handler.CreateUser)
		r.Put("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
	})
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrInvalidUserID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u, err := h.UserUsecase.GetUser(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, errs.ErrUserNotFound, http.StatusNotFound)
			} else {
				http.Error(w, errs.ErrGetUser, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrGetUser, http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u *entity.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := h.UserUsecase.CreateUser(ctx, u)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if strings.Contains(errors.Cause(err).Error(), "duplicate key value") {
				http.Error(w, errs.ErrAlreadyExists, http.StatusBadRequest)
			} else {
				http.Error(w, errs.ErrCreateUser, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrCreateUser, http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u *entity.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	u.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrInvalidUserID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.UserUsecase.UpdateUser(ctx, u)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if errors.Cause(err).Error() == errs.ErrUserNotFound {
				http.Error(w, errs.ErrUserNotFound, http.StatusNotFound)
				return
			} else if strings.Contains(errors.Cause(err).Error(), "duplicate key value") {
				http.Error(w, errs.ErrAlreadyExists, http.StatusBadRequest)
			} else {
				http.Error(w, errs.ErrUpdateUser, http.StatusInternalServerError)
				return
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, errs.ErrInvalidUserID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.UserUsecase.DeleteUser(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, errs.ErrTimeout, http.StatusGatewayTimeout)
		default:
			if errors.Cause(err).Error() == errs.ErrUserNotFound {
				http.Error(w, errs.ErrUserNotFound, http.StatusNotFound)
				return
			} else {
				http.Error(w, errs.ErrDeleteUser, http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
