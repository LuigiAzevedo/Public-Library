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
	ucErr "github.com/LuigiAzevedo/public-library-v2/internal/domain/usecase"
	"github.com/LuigiAzevedo/public-library-v2/internal/mock"
)

func TestSearchUserLoan(t *testing.T) {
	testCases := map[string]struct {
		ID            any
		buildStubs    func(uc *mock.MockLoanUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID: 1,
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					SearchUserLoans(gomock.Any(), gomock.Eq(1)).
					Times(1).
					Return([]*entity.Loan{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		"Invalid URL Param": {
			ID: "ID",
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					SearchUserLoans(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Not Found": {
			ID: 1,
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					SearchUserLoans(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, repoErr.ErrLoanNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			ID: 1,
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					SearchUserLoans(gomock.Any(), gomock.Any()).
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

			uc := mock.NewMockLoanUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/v1/loans/", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewLoanHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestBorrowBook(t *testing.T) {
	type testLoanRequest struct {
		UserID any `json:"user_id"`
		BookID any `json:"book_id"`
	}

	testCases := map[string]struct {
		body          testLoanRequest
		buildStubs    func(uc *mock.MockLoanUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			body: testLoanRequest{
				UserID: 1,
				BookID: 1,
			},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					BorrowBook(gomock.Any(), gomock.Eq(1), gomock.Eq(1)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		"Invalid Body": {
			body: testLoanRequest{
				UserID: "1",
				BookID: "1",
			},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					BorrowBook(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Book Not Found": {
			body: testLoanRequest{},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					BorrowBook(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(repoErr.ErrBookNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"User Not Found": {
			body: testLoanRequest{},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					BorrowBook(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(repoErr.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Book Already Borrowed": {
			body: testLoanRequest{},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					BorrowBook(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(ucErr.ErrReturnBookFirst)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Book Unavailable": {
			body: testLoanRequest{},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					BorrowBook(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(ucErr.ErrBookUnavailable)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			body: testLoanRequest{},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					BorrowBook(gomock.Any(), gomock.Any(), gomock.Any()).
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

			uc := mock.NewMockLoanUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/v1/loans/borrow", bytes.NewReader(data))
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewLoanHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestReturnBook(t *testing.T) {
	type testLoanRequest struct {
		UserID any `json:"user_id"`
		BookID any `json:"book_id"`
	}

	testCases := map[string]struct {
		body          testLoanRequest
		buildStubs    func(uc *mock.MockLoanUsecase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			body: testLoanRequest{
				UserID: 1,
				BookID: 1,
			},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					ReturnBook(gomock.Any(), gomock.Eq(1), gomock.Eq(1)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		"Invalid Body": {
			body: testLoanRequest{
				UserID: "1",
				BookID: "1",
			},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					ReturnBook(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Loan Already Returned": {
			body: testLoanRequest{},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					ReturnBook(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(ucErr.ErrLoanAlreadyReturned)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			body: testLoanRequest{},
			buildStubs: func(uc *mock.MockLoanUsecase) {
				uc.EXPECT().
					ReturnBook(gomock.Any(), gomock.Any(), gomock.Any()).
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

			uc := mock.NewMockLoanUsecase(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/v1/loans/return", bytes.NewReader(data))
			assert.NoError(t, err)

			router := chi.NewRouter()
			NewLoanHandler(router, uc)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
