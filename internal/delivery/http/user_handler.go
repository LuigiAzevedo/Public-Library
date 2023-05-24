package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	repoErr "github.com/LuigiAzevedo/public-library-v2/internal/database/repository"
	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
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
		http.Error(w, invalidUserID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u, err := h.UserUsecase.GetUser(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			if err == repoErr.ErrUserNotFound {
				http.Error(w, userNotFound, http.StatusNotFound)
			} else {
				http.Error(w, getUser, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, getUser, http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u entity.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := h.UserUsecase.CreateUser(ctx, &u)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			if err == repoErr.ErrAlreadyExists {
				http.Error(w, alreadyExists, http.StatusBadRequest)
			} else {
				http.Error(w, createUser, http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, createUser, http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u entity.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	u.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, invalidUserID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.UserUsecase.UpdateUser(ctx, &u)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			switch err {
			case repoErr.ErrUserNotFound:
				http.Error(w, userNotFound, http.StatusNotFound)
			case repoErr.ErrAlreadyExists:
				http.Error(w, alreadyExists, http.StatusBadRequest)
			default:
				http.Error(w, updateUser, http.StatusInternalServerError)
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
		http.Error(w, invalidUserID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.UserUsecase.DeleteUser(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())

		select {
		case <-ctx.Done():
			http.Error(w, timeout, http.StatusGatewayTimeout)
		default:
			if err == repoErr.ErrUserNotFound {
				http.Error(w, userNotFound, http.StatusNotFound)
			} else {
				http.Error(w, deleteUser, http.StatusInternalServerError)
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
