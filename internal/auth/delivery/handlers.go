package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"html"
	"log"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/validator"
	"github.com/gorilla/mux"
)

type usecase interface {
	Authenticate(ctx context.Context, username, password string) bool
	Registration(ctx context.Context, username, name, password string) error
	GetUserDataByUsername(ctx context.Context, username string) (models.UserData, error)
}

type token interface {
	CreateJWT(ctx context.Context, username string) (string, error)
	GetUserDataByJWT(cookies []*http.Cookie) (jwt.UserData, error)
	GetUserByJWT(ctx context.Context, cookies []*http.Cookie) (jwt.User, error)
	IsAuthorized(ctx context.Context, cookies []*http.Cookie) (jwt.User, error)
}

type Delivery struct {
	usecase usecase
	token   token
	mu      sync.Mutex
}

func NewDelivery(usecase usecase, token token) *Delivery {
	return &Delivery{
		usecase: usecase,
		token:   token,
	}
}

// LoginHandler godoc
// @Summary User login
// @Description Authenticate a user with username and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.AuthReqDTO true "Credentials for login, including username and password"
// @Success 201 {object} responser.SuccessResponse "Authentication successful"
// @Failure 400 {object} responser.ErrorResponse "Invalid format JSON"
// @Failure 401 {object} responser.ErrorResponse "Incorrect login or password"
// @Router /login [post]
func (d *Delivery) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на аутентификацию")

	var creds models.AuthReqDTO
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Errorf("не удалось распарсить json")
		responser.SendError(ctx, w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if err := validator.Check(creds); err != nil {
		log.Errorf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	if d.usecase.Authenticate(ctx, creds.Username, creds.Password) {
		err := d.setToken(w, r, creds.Username)
		if err != nil {
			log.Errorf("не удалось аутентифицировать пользователя")
			responser.SendError(ctx, w, "Invalid format JSON", http.StatusUnauthorized)
			return
		}

		log.Println("пользователь успешно аутентифицирован")

		responser.SendOK(w, "Authentication successful", http.StatusCreated)

	} else {
		log.Errorf("неверный логин или пароль")
		responser.SendError(ctx, w, "Incorrect login or password", http.StatusUnauthorized)
	}
}

// @Summary Register a new user
// @Description Creates a new user with the provided credentials.
// @Tags auth
// @Accept json
// @Produce json
// @Param creds body models.RegisterReqDTO true "Registration information"
// @Success 201 {object} models.RegisterRespDTO "Registration successful"
// @Failure 400 {object} responser.ErrorResponse "Invalid input data"
// @Failure 409 {object} responser.ErrorResponse "A user with that username already exists"
// @Failure 400 {object} responser.ErrorResponse "User failed to create"
// @Router /signup [post]
func (d *Delivery) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на регистрацию")

	var creds models.RegisterReqDTO
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		responser.SendError(ctx, w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if err := validator.Check(creds); err != nil {
		log.Errorf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	if err := d.usecase.Registration(ctx, creds.Username, creds.Name, creds.Password); err != nil {
		responser.SendError(ctx, w, "A user with that username already exists", http.StatusConflict)
		return
	}

	log.Println("получение данных пользователя")

	d.setToken(w, r, creds.Username)
	userData, err := d.usecase.GetUserDataByUsername(ctx, creds.Username)
	if err != nil {
		responser.SendError(ctx, w, "User failed to create", http.StatusBadRequest)
		return
	}

	log.Println("новый пользователь получен")

	userDataDTO := convertUserDataToDTO(userData)

	response := models.RegisterRespDTO{
		Message: "Registration successful",
		User:    userDataDTO,
	}

	if err := validator.Check(response); err != nil {
		log.Errorf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Println("пользователь успешно зарегистрирован")

	responser.SendStruct(ctx, w, response, http.StatusCreated)
}

// AuthHandler godoc
// @Summary Authenticate a user
// @Description Retrieve user data based on the JWT token present in the cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.UserDataRespDTO "User data retrieved successfully"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized: token is invalid"
// @Router /auth [get]
func (d *Delivery) AuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на авторизацию")

	user, ok := ctx.Value(models.UserKey).(jwt.User)
	if !ok {
		responser.SendError(ctx, w, "User not found", http.StatusNotFound)
		return
	}

	userData, err := d.usecase.GetUserDataByUsername(ctx, user.Username)
	if err != nil {
		log.Println("не получилось получить данные пользователя")
		responser.SendError(ctx, w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userDataDTO := convertUserDataToDTO(userData)

	response := models.AuthRespDTO{
		User: userDataDTO,
	}

	if err := validator.Check(response); err != nil {
		log.Errorf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Println("пользователь успешно авторизован")

	responser.SendStruct(ctx, w, response, http.StatusOK)
}

func (d *Delivery) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := d.token.IsAuthorized(ctx, r.Cookies())
		if err == errors.New("token expired") {
			log.Println("token expired, create new token")
			user, err = d.token.GetUserByJWT(r.Context(), r.Cookies())

			if err != nil {
				responser.SendError(ctx, w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			d.setToken(w, r, user.Username)
		}
		if err != nil {
			responser.SendError(ctx, w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, models.UserKey, user)
		ctx = context.WithValue(ctx, models.MuxParamsKey, mux.Vars(r))

		r = r.WithContext(ctx)

		next(w, r)
	}
}

// LogoutHandler godoc
// @Summary Log out a user
// @Description Invalidate the user's session by clearing the access token cookie.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object}  responser.SuccessResponse "Logout successful"
// @Failure 401 {object} responser.ErrorResponse "No access token found"
// @Router /logout [post]
func (d *Delivery) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на разлогин")

	tokenExists := false
	for _, cookie := range r.Cookies() {
		if cookie.Name == "access_token" {
			tokenExists = true
			break
		}
	}

	if !tokenExists {
		responser.SendError(ctx, w, "No access token found", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "t",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	log.Println("разлогин прошел успешно")

	responser.SendOK(w, "Logout successful", http.StatusOK)
}

func (c *Delivery) isAuthorized(r *http.Request) error {
	ctx := r.Context()
	_, err := c.token.IsAuthorized(ctx, r.Cookies())
	if err != nil {
		return err
	}

	return nil
}

func (d *Delivery) setToken(w http.ResponseWriter, r *http.Request, username string) error {
	ctx := r.Context()
	token, err := d.token.CreateJWT(ctx, username)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	return nil
}

func convertUserDataToDTO(userData models.UserData) models.UserDataRespDTO {
	var avatarURL *string
	if userData.AvatarURL != nil {
		avatarURL = new(string)
		*avatarURL = html.EscapeString(*userData.AvatarURL)
	}

	return models.UserDataRespDTO{
		ID:        userData.ID,
		Username:  html.EscapeString(userData.Username),
		Name:      html.EscapeString(userData.Name),
		AvatarURL: avatarURL,
	}
}

func convertFromJWTUserData(userDataJWT jwt.UserData) models.UserDataRespDTO {
	return models.UserDataRespDTO{
		ID:       userDataJWT.ID,
		Username: html.EscapeString(userDataJWT.Username),
		Name:     html.EscapeString(userDataJWT.Name),
	}
}
