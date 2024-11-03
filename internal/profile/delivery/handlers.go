package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	invalidJSONError  = "Invalid format JSON"
	responseError     = "Failed to create response"
	userNotFoundError = "User not found"
)

type usecase interface {
	UpdateProfile(ctx context.Context, profile models.UpdateProfileRequestDTO) error
	GetProfile(ctx context.Context, id uuid.UUID) (models.ProfileData, error)
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

// GetSelfProfileHandler godoc
// @Summary Get self profile data
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
func (d *Delivery) GetSelfProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("Profile delivery: пришел запрос на получение данных профиля")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	id := user.ID

	profileData, err := d.usecase.GetProfile(ctx, id)
	if err != nil {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	response := convertProfileDataToDTO(profileData)

	jsonResp, err := json.Marshal(response)
	if err != nil {
		responser.SendErrorResponse(w, responseError, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// GetProfileHandler godoc
// @Summary Get profile data
// @Description Get bio, avatar and birthdate of user.
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} models.GetProfileResponseDTO "Profile data found"
// @Failure 400 {object} responser.ErrorResponse "Invalid format JSON"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized"
// @Failure 404 {object} responser.ErrorResponse "User not found"
// @Router /profile/{userid} [get]
func (d *Delivery) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("Profile delivery: пришел запрос на получение данных профиля")

	vars := mux.Vars(r)
	userid := vars["userid"]

	id := uuid.MustParse(userid)

	profileData, err := d.usecase.GetProfile(ctx, id)
	if err != nil {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	response := convertProfileDataToDTO(profileData)

	jsonResp, err := json.Marshal(response)
	if err != nil {
		responser.SendErrorResponse(w, responseError, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// UpdateProfileHandler godoc
// @Summary Update profile data
// @Description Update bio, avatar, name or birthdate of user.
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param profile_data body models.UpdateProfileRequestDTO true "JSON representation of profile data"
// @Param avatar formData file false "User avatar image" example:"/2024_2_eaglesDesigners/uploads/avatar/f0364477-bfd4-496d-b639-d825b009d509.png"
// @Success 200 {object} responser.SuccessResponse "Profile updated"
// @Failure 400 {object} responser.ErrorResponse "Failed to update profile"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized"
// @Failure 404 {object} responser.ErrorResponse "User not found"
// @Router /profile [put]
func (d *Delivery) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()

	log.Println("Profile delivery: пришел запрос на обновление профиля")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	var profile models.UpdateProfileRequestDTO

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Profile delivery: не удалось распарсить запрос: ", err)
		responser.SendErrorResponse(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	profile.ID = user.ID

	jsonString := r.FormValue("profile_data")
	if jsonString != "" {
		if err := json.Unmarshal([]byte(jsonString), &profile); err != nil {
			responser.SendErrorResponse(w, invalidJSONError, http.StatusBadRequest)
			return
		}
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		responser.SendErrorResponse(w, "Failed to get avatar", http.StatusBadRequest)
		return
	}
	defer func() {
		if avatar != nil {
			avatar.Close()
		}
	}()

	if avatar != nil {
		profile.Avatar = &avatar
	}

	if err := d.usecase.UpdateProfile(ctx, profile); err != nil {
		responser.SendErrorResponse(w, "Failed to update profile", http.StatusBadRequest)
		return
	}

	responser.SendOKResponse(w, "Profile updated", http.StatusOK)
}

func convertProfileDataToDTO(profileData models.ProfileData) models.GetProfileResponseDTO {
	return models.GetProfileResponseDTO{
		Name:      profileData.Name,
		Bio:       profileData.Bio,
		AvatarURL: profileData.AvatarPath,
		Birthdate: profileData.Birthdate,
	}
}
