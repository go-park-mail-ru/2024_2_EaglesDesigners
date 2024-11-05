package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/csrf"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
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

	log.Println("Login delivery: пришел запрос на аутентификацию")

	var creds models.AuthReqDTO
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Println("Login delivery: не удалось распарсить json")
		responser.SendError(w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if d.usecase.Authenticate(ctx, creds.Username, creds.Password) {
		err := d.setTokens(w, r, creds.Username)
		if err != nil {
			log.Println("Login delivery: не удалось аутентифицировать пользователя")
			responser.SendError(w, "Invalid format JSON", http.StatusUnauthorized)
			return
		}

		log.Println("Login delivery: пользователь успешно аутентифицирован")

		responser.SendOK(w, "Authentication successful", http.StatusCreated)

	} else {
		log.Println("Login delivery: неверный логин или пароль")
		responser.SendError(w, "Incorrect login or password", http.StatusUnauthorized)
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

	log.Println("Register delivery: Пришел запрос на регистрацию")

	var creds models.RegisterReqDTO
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		responser.SendError(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if len(creds.Username) < 6 || len(creds.Password) < 8 || creds.Name == "" {
		responser.SendError(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if err := d.usecase.Registration(ctx, creds.Username, creds.Name, creds.Password); err != nil {
		responser.SendError(w, "A user with that username already exists", http.StatusConflict)
		return
	}

	log.Println("Register delivery: получение данных пользователя")

	d.setTokens(w, r, creds.Username)
	userData, err := d.usecase.GetUserDataByUsername(ctx, creds.Username)
	if err != nil {
		responser.SendError(w, "User failed to create", http.StatusBadRequest)
		return
	}

	log.Println("Register delivery: новый пользователь получен")

	userDataDTO := convertUserDataToDTO(userData)

	response := models.RegisterRespDTO{
		Message: "Registration successful",
		User:    userDataDTO,
	}

	log.Println("Register delivery: Пользователь успешно зарегистрирован")

	responser.SendStruct(w, response, http.StatusCreated)
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

	log.Println("Auth delivery: пришел запрос на авторизацию")

	user, ok := ctx.Value(models.UserKey).(jwt.User)
	if !ok {
		responser.SendError(w, "User not found", http.StatusNotFound)
		return
	}

	userData, err := d.usecase.GetUserDataByUsername(ctx, user.Username)
	if err != nil {
		log.Println("Auth delivery: не получилось получить данные пользователя")
		responser.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userDataDTO := convertUserDataToDTO(userData)

	log.Println("Auth delivery: пользователь успешно авторизован")

	response := models.AuthRespDTO{
		User: userDataDTO,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		responser.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

var errTokenExpired = errors.New("токен истек")

func (d *Delivery) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		user, err := d.token.IsAuthorized(ctx, r.Cookies())
		if err == errTokenExpired {
			log.Println("Auth middlware: токен истек, создается новый токен")
			user, err = d.token.GetUserByJWT(r.Context(), r.Cookies())

			if err != nil {
				responser.SendError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			d.setTokens(w, r, user.Username)
		}
		if err != nil {
			responser.SendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, models.UserKey, user)
		ctx = context.WithValue(ctx, models.MuxParamsKey, mux.Vars(r))

		r = r.WithContext(ctx)

		next(w, r)
	}
}

func (d *Delivery) Csrf(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("not_csrf")

		user := r.Context().Value(models.UserKey).(jwt.User)

		err := csrf.CheckCSRF(token, user.ID, user.Username)
		if err != nil {
			if err == errTokenExpired {
				responser.SendError(w, "csrf expired", http.StatusForbidden)
				return
			}
			responser.SendError(w, "Invalid csrf", http.StatusForbidden)
			return
		}

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
	tokenExists := false
	for _, cookie := range r.Cookies() {
		if cookie.Name == "access_token" {
			tokenExists = true
			break
		}
	}

	if !tokenExists {
		responser.SendError(w, "No access token found", http.StatusUnauthorized)
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

func (d *Delivery) setTokens(w http.ResponseWriter, r *http.Request, username string) error {
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

	csrf, err := csrf.CreateCSRF(token)
	if err != nil {
		log.Printf("Auth setTokens: не удалось создать csrf токен: %v", err)
		return err
	}

	w.Header().Set("X-CSRF-Token", csrf)

	return nil
}

func convertUserDataToDTO(userData models.UserData) models.UserDataRespDTO {
	return models.UserDataRespDTO{
		ID:        userData.ID,
		Username:  userData.Username,
		Name:      userData.Name,
		AvatarURL: userData.AvatarURL,
	}
}

func convertFromJWTUserData(userDataJWT jwt.UserData) models.UserDataRespDTO {
	return models.UserDataRespDTO{
		ID:       userDataJWT.ID,
		Username: userDataJWT.Username,
		Name:     userDataJWT.Name,
	}
}
