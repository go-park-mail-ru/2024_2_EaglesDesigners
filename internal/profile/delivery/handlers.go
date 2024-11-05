package delivery

import (
	"context"
	"encoding/json"
	"html"
	"net/http"
	"sync"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/validator"
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
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на получение данных о своем профиле")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	id := user.ID

	profileData, err := d.usecase.GetProfile(ctx, id)
	if err != nil {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	response := convertProfileDataToDTO(profileData)

	if err := validator.Check(response); err != nil {
		log.Errorf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Println("данные успешно отправлены")

	responser.SendStruct(ctx, w, response, http.StatusOK)
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
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на получение данных профиля")

	vars := mux.Vars(r)
	userid := vars["userid"]

	id := uuid.MustParse(userid)

	profileData, err := d.usecase.GetProfile(ctx, id)
	if err != nil {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	response := convertProfileDataToDTO(profileData)

	if err := validator.Check(response); err != nil {
		log.Errorf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Println("данные успешно отправлены")

	responser.SendStruct(ctx, w, response, http.StatusOK)
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
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на обновление профиля")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	var profile models.UpdateProfileRequestDTO

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Errorf("не удалось распарсить запрос: %v", err)
		responser.SendError(ctx, w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	profile.ID = user.ID

	jsonString := r.FormValue("profile_data")
	if jsonString != "" {
		if err := json.Unmarshal([]byte(jsonString), &profile); err != nil {
			responser.SendError(ctx, w, invalidJSONError, http.StatusBadRequest)
			return
		}
	}

	if err := validator.Check(profile); err != nil {
		log.Errorf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		responser.SendError(ctx, w, "Failed to get avatar", http.StatusBadRequest)
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
		responser.SendError(ctx, w, "Failed to update profile", http.StatusBadRequest)
		return
	}

	log.Println("профиль успешно обновлен")

	responser.SendOK(w, "Profile updated", http.StatusOK)
}

func convertProfileDataToDTO(profileData models.ProfileData) models.GetProfileResponseDTO {
	var safeName *string
	if profileData.Name != nil {
		safeName = new(string)
		*safeName = html.EscapeString(*profileData.Name)
	}

	var safeBio *string
	if profileData.Bio != nil {
		safeBio = new(string)
		*safeBio = html.EscapeString(*profileData.Bio)
	}

	var safeAvatarURL *string
	if profileData.AvatarPath != nil {
		safeAvatarURL = new(string)
		*safeAvatarURL = html.EscapeString(*profileData.AvatarPath)
	}

	return models.GetProfileResponseDTO{
		Name:      safeName,
		Bio:       safeBio,
		AvatarURL: safeAvatarURL,
		Birthdate: profileData.Birthdate,
	}
}
