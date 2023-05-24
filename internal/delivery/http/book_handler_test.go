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

func TestGetBook(t *testing.T) {
	testCases := map[string]struct {
		ID            any
		buildStubs    func(uc *mock.MockBookUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID: 1,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					GetBook(gomock.Any(), gomock.Eq(1)).
					Times(1).
					Return(&entity.Book{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		"Invalid URL Param": {
			ID: "ID",
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					GetBook(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Not Found": {
			ID: 1,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					GetBook(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, repoErr.ErrBookNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			ID: 1,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					GetBook(gomock.Any(), gomock.Any()).
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

			uc := mock.NewMockBookUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/v1/books/", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewBookHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestSearchBooks(t *testing.T) {
	type testSearchBookRequest struct {
		Title any `json:"title"`
	}

	testCases := map[string]struct {
		title         testSearchBookRequest
		buildStubs    func(uc *mock.MockBookUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK Search": {
			title: testSearchBookRequest{
				Title: "book title",
			},
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					SearchBooks(gomock.Any(), gomock.Eq("book title")).
					Times(1).
					Return([]*entity.Book{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		"OK List": {
			title: testSearchBookRequest{
				Title: "",
			},
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					ListBooks(gomock.Any()).
					Times(1).
					Return([]*entity.Book{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		"Invalid Body": {
			title: testSearchBookRequest{
				Title: 1,
			},
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					SearchBooks(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Not Found": {
			title: testSearchBookRequest{},
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					ListBooks(gomock.Any()).
					Times(1).
					Return(nil, repoErr.ErrBookNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			title: testSearchBookRequest{},
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					ListBooks(gomock.Any()).
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

			uc := mock.NewMockBookUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.title)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodGet, "/v1/books/", bytes.NewReader(data))
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewBookHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateBook(t *testing.T) {
	book := &entity.Book{
		ID:     1,
		Title:  "Book123",
		Author: "author123",
		Amount: 5,
	}

	testCases := map[string]struct {
		book          any
		buildStubs    func(uc *mock.MockBookUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			book: book,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					CreateBook(gomock.Any(), gomock.Eq(book)).
					Times(1).
					Return(book.ID, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		"Invalid Body": {
			book: "invalid body",
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					CreateBook(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Unexpected Error": {
			book: book,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					CreateBook(gomock.Any(), gomock.Any()).
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

			uc := mock.NewMockBookUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.book)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/v1/books/", bytes.NewReader(data))
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewBookHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateBook(t *testing.T) {
	book := &entity.Book{
		ID:     1,
		Title:  "Book123",
		Author: "author123",
		Amount: 5,
	}

	testCases := map[string]struct {
		ID            any
		book          any
		buildStubs    func(uc *mock.MockBookUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID:   1,
			book: book,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					UpdateBook(gomock.Any(), gomock.Eq(book)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		"Invalid Body": {
			ID:   1,
			book: "",
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					UpdateBook(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Invalid URL Param": {
			ID:   "ID",
			book: book,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					UpdateBook(gomock.Any(), gomock.Eq(book)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Not Found": {
			ID:   1,
			book: book,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					UpdateBook(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repoErr.ErrBookNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			ID:   1,
			book: book,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					UpdateBook(gomock.Any(), gomock.Any()).
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

			uc := mock.NewMockBookUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.book)
			assert.NoError(t, err)

			url := fmt.Sprint("/v1/books/", tc.ID)

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewBookHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteBook(t *testing.T) {
	testCases := map[string]struct {
		ID            any
		buildStubs    func(uc *mock.MockBookUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID: 1,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					DeleteBook(gomock.Any(), gomock.Eq(1)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		"Invalid URL Param": {
			ID: "ID",
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					DeleteBook(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Not Found": {
			ID: 1,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					DeleteBook(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repoErr.ErrBookNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			ID: 1,
			buildStubs: func(uc *mock.MockBookUsecase) {
				uc.EXPECT().
					DeleteBook(gomock.Any(), gomock.Any()).
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

			uc := mock.NewMockBookUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/v1/books/", tc.ID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewBookHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
