package delivery_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/delivery"
	mock_usecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/delivery/mocks"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetProfileHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUsecase := mock_usecase.NewMockusecase(ctrl)
	delivery := delivery.New(mockUsecase)

	tests := []struct {
		name               string
		userID             string
		prepareMock        func()
		expectedStatusCode int
	}{
		{
			name:   "успешный запрос",
			userID: uuid.New().String(),
			prepareMock: func() {
				mockUsecase.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(models.ProfileData{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "невалидный UUID",
			userID: "невалидный-UUID",
			prepareMock: func() {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "пользователь не найден",
			userID: uuid.New().String(),
			prepareMock: func() {
				mockUsecase.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(models.ProfileData{}, errors.New("user not found"))
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMock()

			req := httptest.NewRequest(http.MethodGet, "/profile/{userid}", bytes.NewBuffer(nil))
			req = mux.SetURLVars(req, map[string]string{"userid": tt.userID})
			recorder := httptest.NewRecorder()

			delivery.GetProfileHandler(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
		})
	}
}

func TestGetSelfProfileHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUsecase := mock_usecase.NewMockusecase(ctrl)
	delivery := delivery.New(mockUsecase)

	tests := []struct {
		name               string
		prepareCtx         func() context.Context
		prepareMock        func()
		expectedStatusCode int
	}{
		{
			name: "успешный запрос",
			prepareCtx: func() context.Context {
				return context.WithValue(context.Background(), auth.UserKey, auth.User{ID: uuid.New()})
			},
			prepareMock: func() {
				mockUsecase.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(models.ProfileData{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "пользователь не найден в контексте",
			prepareCtx: func() context.Context {
				return context.Background()
			},
			prepareMock:        func() {},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "ошибка при получении профиля",
			prepareCtx: func() context.Context {
				return context.WithValue(context.Background(), auth.UserKey, auth.User{ID: uuid.New()})
			},
			prepareMock: func() {
				mockUsecase.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(models.ProfileData{}, errors.New("profile not found"))
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.prepareCtx()
			tt.prepareMock()

			req := httptest.NewRequest(http.MethodGet, "/profile", nil).WithContext(ctx)
			recorder := httptest.NewRecorder()

			delivery.GetSelfProfileHandler(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
		})
	}
}

func TestUpdateProfileHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUsecase := mock_usecase.NewMockusecase(ctrl)
	delivery := delivery.New(mockUsecase)

	tests := []struct {
		name               string
		prepareCtx         func() context.Context
		prepareRequest     func() *http.Request
		prepareMock        func()
		expectedStatusCode int
	}{
		{
			name: "успешное обновление профиля",
			prepareCtx: func() context.Context {
				return context.WithValue(context.Background(), auth.UserKey, auth.User{ID: uuid.New()})
			},
			prepareRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				profileData := `{"name": "kekabuba", "bio": "lule"}`
				_ = writer.WriteField("profile_data", profileData)

				part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
				part.Write([]byte("fake image data"))

				writer.Close()
				req := httptest.NewRequest(http.MethodPut, "/profile", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				return req
			},
			prepareMock: func() {
				mockUsecase.EXPECT().
					UpdateProfile(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "нет пользователя в контексте",
			prepareCtx: func() context.Context {
				return context.Background()
			},
			prepareRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodPut, "/profile", nil)
			},
			prepareMock:        func() {},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "некорректный JSON в поле profile_data",
			prepareCtx: func() context.Context {
				return context.WithValue(context.Background(), auth.UserKey, auth.User{ID: uuid.New()})
			},
			prepareRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				_ = writer.WriteField("profile_data", "{invalid-json}")
				writer.Close()

				req := httptest.NewRequest(http.MethodPut, "/profile", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				return req
			},
			prepareMock:        func() {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "валидация не прошла",
			prepareCtx: func() context.Context {
				return context.WithValue(context.Background(), auth.UserKey, auth.User{ID: uuid.New()})
			},
			prepareRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				_ = writer.WriteField("profile_data", `{"name": "&name", "bio": "lule"}`)
				writer.Close()

				req := httptest.NewRequest(http.MethodPut, "/profile", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				return req
			},
			prepareMock:        func() {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "ошибка в usecase.UpdateProfile",
			prepareCtx: func() context.Context {
				return context.WithValue(context.Background(), auth.UserKey, auth.User{ID: uuid.New()})
			},
			prepareRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				_ = writer.WriteField("profile_data", `{"name": "kekabuba", "bio": "lule"}`)
				writer.Close()

				req := httptest.NewRequest(http.MethodPut, "/profile", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				return req
			},
			prepareMock: func() {
				mockUsecase.EXPECT().
					UpdateProfile(gomock.Any(), gomock.Any()).
					Return(errors.New("something went wrong"))
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.prepareCtx()
			req := tt.prepareRequest().WithContext(ctx)
			tt.prepareMock()
			recorder := httptest.NewRecorder()

			delivery.UpdateProfileHandler(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
		})
	}
}
