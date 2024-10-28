package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/base64helper"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
	"github.com/google/uuid"
)

const (
	invalidJSONError  = "Invalid format JSON"
	responseError     = "Failed to create response"
	userNotFoundError = "User not found"
)

type usecase interface {
	UpdateProfile(ctx context.Context, profile models.UpdateProfileRequestDTO) error
	GetProfile(ctx context.Context, username string) (models.ProfileData, error)
}

type token interface {
	GetUserByJWT(ctx context.Context, cookies []*http.Cookie) (jwt.User, error)
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
// @Security BearerAuth
// @Success 200 {object} models.GetProfileResponseDTO "Profile data found"
// @Failure 400 {object} responser.ErrorResponse "Invalid format JSON"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized"
// @Failure 404 {object} responser.ErrorResponse "User not found"
// @Router /profile [get]
func (d *Delivery) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	username := user.Username

	profileData, err := d.usecase.GetProfile(ctx, username)
	if err != nil {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	var avatarBase64 *string

	if profileData.AvatarURL != nil {
		avatarUUID, err := uuid.Parse(*profileData.AvatarURL)
		if err != nil {
			responser.SendErrorResponse(w, responseError, http.StatusBadRequest)
			return
		}

		avatarBase64 = new(string)
		*avatarBase64, err = base64helper.ReadPhotoBase64(avatarUUID)
		if err != nil {
			avatarBase64 = nil
		}
	}

	response := models.GetProfileResponseDTO{
		Name:         profileData.Name,
		Bio:          profileData.Bio,
		Birthdate:    profileData.Birthdate,
		AvatarBase64: avatarBase64,
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
// @Failure 401 {object} responser.ErrorResponse "Unauthorized"
// @Failure 404 {object} responser.ErrorResponse "User not found"
// @Router /profile [put]
func (d *Delivery) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()

	user, err := d.token.GetUserByJWT(ctx, r.Cookies())
	if err != nil {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	var profile models.UpdateProfileRequestDTO

	profile.Username = user.Username

	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		responser.SendErrorResponse(w, invalidJSONError, http.StatusBadRequest)
		return
	}

	if err := d.usecase.UpdateProfile(ctx, profile); err != nil {
		responser.SendErrorResponse(w, "Failed to update profile", http.StatusBadRequest)
		return
	}

	responser.SendOKResponse(w, "Profile updated", http.StatusOK)
}
