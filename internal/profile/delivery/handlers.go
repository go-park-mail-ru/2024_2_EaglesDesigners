package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/base64helper"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
	"github.com/google/uuid"
)

type usecase interface {
	UpdateProfile(ctx context.Context, profile models.UpdateProfileRequestDTO) error
	GetProfile(ctx context.Context, username string) (models.ProfileData, error)
}

type Delivery struct {
	usecase usecase
	mu      sync.Mutex
}

func New(usecase usecase) *Delivery {
	return &Delivery{
		usecase: usecase,
	}
}

func (d *Delivery) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var profile models.GetProfileRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		responser.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	username := profile.Username

	profileData, err := d.usecase.GetProfile(ctx, username)
	if err != nil {
		responser.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest) //-----------
		return
	}

	avatarUUID, err := uuid.Parse(*profileData.AvatarURL)
	if err != nil {
		responser.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest) //-----------
		return
	}

	avatarBase64, err := base64helper.ReadPhotoBase64(avatarUUID)

	response := models.GetProfileResponseDTO{
		Bio:          profileData.Bio,
		Birthdate:    profileData.Birthdate,
		AvatarBase64: &avatarBase64,
	}

	if err != nil {
		responser.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest) //-----------
		return
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		responser.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized) //-------------
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (d *Delivery) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()

	var profile models.UpdateProfileRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		responser.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if err := d.usecase.UpdateProfile(ctx, profile); err != nil {
		responser.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest) //-------------
		return
	}

}
