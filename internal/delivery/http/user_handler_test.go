package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	repoErr "github.com/LuigiAzevedo/public-library-v2/internal/database/repository"
	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	"github.com/LuigiAzevedo/public-library-v2/internal/mock"
)

func TestGetUser(t *testing.T) {
	testCases := map[string]struct {
		ID            any
		buildStubs    func(uc *mock.MockUserUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID: 1,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(1)).
					Times(1).
					Return(&entity.User{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		"Invalid URL Param": {
			ID: "ID",
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Not Found": {
			ID: 0,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, repoErr.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			ID: 1,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockUserUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/v1/users/", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewUserHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateUser(t *testing.T) {
	user := &entity.User{
		ID:       1,
		Username: "user123",
		Password: "password123",
		Email:    "user123@example.com",
	}

	testCases := map[string]struct {
		user          any
		buildStubs    func(uc *mock.MockUserUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			user: user,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(user)).
					Times(1).
					Return(user.ID, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		"Invalid Body": {
			user: "invalid body",
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Duplicated": {
			user: user,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(0, repoErr.ErrAlreadyExists)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Unexpected Error": {
			user: user,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(0, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockUserUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.user)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/v1/users/", bytes.NewReader(data))
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewUserHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	user := &entity.User{
		ID:       1,
		Username: "user123",
		Password: "password123",
		Email:    "user123@example.com",
	}

	testCases := map[string]struct {
		ID            any
		user          any
		buildStubs    func(uc *mock.MockUserUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID:   1,
			user: user,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(user)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		"Invalid Body": {
			ID:   1,
			user: "invalid body",
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Invalid URL Param": {
			ID:   "ID",
			user: user,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Not Found": {
			ID:   1,
			user: user,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repoErr.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Duplicated": {
			ID:   1,
			user: user,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repoErr.ErrAlreadyExists)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Unexpected Error": {
			ID:   1,
			user: user,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockUserUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.user)
			assert.NoError(t, err)

			url := fmt.Sprint("/v1/users/", tc.ID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewUserHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := map[string]struct {
		ID            any
		buildStubs    func(uc *mock.MockUserUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID: 1,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(1)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		"Invalid URL Param": {
			ID: "invalid",
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Not Found": {
			ID: 1,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repoErr.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			ID: 1,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockUserUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/v1/users/", tc.ID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewUserHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
