package delivery_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/delivery"
	mock_usecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/delivery/mocks"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetContactsHandler(t *testing.T) {
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
			name: "success",
			prepareCtx: func() context.Context {
				return context.WithValue(context.Background(), auth.UserKey, auth.User{ID: uuid.New()})
			},
			prepareMock: func() {
				contacts := []models.Contact{}
				contacts = append(contacts, models.Contact{
					ID:       uuid.New().String(),
					Username: "user22",
				})
				contacts = append(contacts, models.Contact{
					ID:       uuid.New().String(),
					Username: "user33",
				})

				mockUsecase.EXPECT().
					GetContacts(gomock.Any(), gomock.Any()).
					Return(contacts, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "user not found",
			prepareMock: func() {},
			prepareCtx: func() context.Context {
				return context.Background()
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMock()
			ctx := tt.prepareCtx()

			req := httptest.NewRequest(http.MethodGet, "/contacts/", bytes.NewBuffer(nil)).WithContext(ctx)
			recorder := httptest.NewRecorder()

			delivery.GetContactsHandler(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
		})
	}
}
