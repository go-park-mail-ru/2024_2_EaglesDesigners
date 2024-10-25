package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
)

const (
	invalidJSONError  = "Invalid format JSON"
	responseError     = "Failed to create response"
	userNotFoundError = "User not found"
)

type usecase interface {
	GetContacts(ctx context.Context, username string) (contacts []models.User, err error)
}

type token interface {
	// GetUserByJWT(ctx context.Context, cookies []*http.Cookie) (jwt.User, error)
	GetUserDataByJWT(cookies []*http.Cookie) (jwt.UserData, error)
}

type Delivery struct {
	usecase usecase
	token   token
	mu      sync.Mutex
}

func New(usecase usecase, token token) *Delivery {
	return &Delivery{
		usecase: usecase,
		token:   token,
	}
}

// LoginHandler godoc
// @Summary Get profile data
// @Description Get bio, avatar and birthdate of user.
// @Tags profile
// @Accept json
// @Produce json
// @Param credentials body models.GetProfileRequestDTO true "Credentials for get profile data"
// @Success 200 {object} models.GetProfileResponseDTO "Profile data found"
// @Failure 400 {object} responser.ErrorResponse "Invalid format JSON"
// @Failure 404 {object} responser.ErrorResponse "User not found"
// @Router /profile [get]
func (d *Delivery) GetContactsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := d.token.GetUserDataByJWT(r.Cookies())
	if err != nil {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	log.Println("hand before get con")

	contacts, err := d.usecase.GetContacts(ctx, user.Username)
	if err != nil {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	log.Println("hand after get con")

	response := models.GetContactsResponseDTO{
		Contacts: contacts,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		responser.SendErrorResponse(w, responseError, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// LoginHandler godoc
// @Summary Update profile data
// @Description Update bio, avatar, name or birthdate of user.
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param credentials body models.UpdateProfileRequestDTO true "Credentials for update profile data"
// @Success 200 {object} responser.SuccessResponse "Profile updated"
// @Failure 400 {object} responser.ErrorResponse "Invalid format JSON"
// @Failure 404 {object} responser.ErrorResponse "User not found"
// @Router /profile [put]
func (d *Delivery) AddContactHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// ctx := r.Context()

	// user, err := d.token.GetUserByJWT(ctx, r.Cookies())
	// if err != nil {
	// 	responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
	// 	return
	// }

	// var profile models.UpdateProfileRequestDTO

	// profile.Username = user.Username

	// if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
	// 	responser.SendErrorResponse(w, invalidJSONError, http.StatusBadRequest)
	// 	return
	// }

	// if err := d.usecase.UpdateProfile(ctx, profile); err != nil {
	// 	responser.SendErrorResponse(w, "Failed to update profile", http.StatusBadRequest)
	// 	return
	// }

	// responser.SendOKResponse(w, "Profile updated", http.StatusOK)
}
