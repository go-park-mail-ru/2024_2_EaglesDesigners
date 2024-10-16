package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	usecaseDto "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/usecase"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
)

type usecase interface {
	Authenticate(username, password string) bool
	Registration(username, name, password string) error
	GetUserDataByUsername(username string) (usecaseDto.UserData, error)
}

type token interface {
	CreateJWT(username string) (string, error)
	GetUserDataByJWT(cookies []*http.Cookie) (jwt.UserData, error)
	GetUserByJWT(cookies []*http.Cookie) (jwt.User, error)
	IsAuthorized(cookies []*http.Cookie) error
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
// @Param credentials body AuthCredentials true "Credentials for login, including username and password"
// @Success 201 {object} responser.SuccessResponse "Authentication successful"
// @Failure 400 {object} responser.ErrorResponse "Invalid format JSON"
// @Failure 401 {object} responser.ErrorResponse "Incorrect login or password"
// @Router /login [post]
func (d *Delivery) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		responser.SendErrorResponse(w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if d.usecase.Authenticate(creds.Username, creds.Password) {
		err := d.setToken(w, creds.Username)
		if err != nil {
			responser.SendErrorResponse(w, "Invalid format JSON", http.StatusUnauthorized)
			return
		}

		responser.SendOKResponse(w, "Authentication successful", http.StatusCreated)

	} else {
		responser.SendErrorResponse(w, "Incorrect login or password", http.StatusUnauthorized)
	}
}

// @Summary Register a new user
// @Description Creates a new user with the provided credentials.
// @Tags auth
// @Accept json
// @Produce json
// @Param creds body RegisterCredentials true "Registration information"
// @Success 201 {object} RegisterResponse "Registration successful"
// @Failure 400 {object} responser.ErrorResponse "Invalid input data"
// @Failure 409 {object} responser.ErrorResponse "A user with that username already exists"
// @Failure 400 {object} responser.ErrorResponse "User failed to create"
// @Router /signup [post]
func (d *Delivery) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	log.Println("Пришел запрос на регистрацию")

	var creds RegisterCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		responser.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if len(creds.Username) < 6 || len(creds.Password) < 8 || creds.Name == "" {
		responser.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if err := d.usecase.Registration(creds.Username, creds.Name, creds.Password); err != nil {
		responser.SendErrorResponse(w, "A user with that username already exists", http.StatusConflict)
	} else {
		d.setToken(w, creds.Username)
		userDataUC, err := d.usecase.GetUserDataByUsername(creds.Username)
		if err != nil {
			responser.SendErrorResponse(w, "User failed to create", http.StatusBadRequest)
			return
		}

		userData := convertFromUsecaseUserData(userDataUC)

		response := RegisterResponse{
			Message: "Registration successful",
			User:    userData,
		}

		log.Println("Пользователь успешно зарегистрирован")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		jsonResp, _ := json.Marshal(response)
		w.Write(jsonResp)
	}
}

// AuthHandler godoc
// @Summary Authenticate a user
// @Description Retrieve user data based on the JWT token present in the cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} AuthResponse "User data retrieved successfully"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized: token is invalid"
// @Router /auth [get]
func (d *Delivery) AuthHandler(w http.ResponseWriter, r *http.Request) {
	err := d.token.IsAuthorized(r.Cookies())
	if err != nil {
		responser.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dataJWT, err := d.token.GetUserDataByJWT(r.Cookies())
	log.Println("/auth cookie: ", dataJWT)
	if err != nil {
		responser.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := convertFromJWTUserData(dataJWT)

	response := AuthResponse{
		User: data,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		responser.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (d *Delivery) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := d.isAuthorized(r)
		if err == errors.New("token expired") {
			log.Println("token expired, create new token")
			user, err := d.token.GetUserByJWT(r.Cookies())
			if err != nil {
				responser.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			d.setToken(w, user.Username)
		}
		if err != nil {
			responser.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
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
		responser.SendErrorResponse(w, "No access token found", http.StatusUnauthorized)
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

	responser.SendOKResponse(w, "Logout successful", http.StatusOK)
}

func (c *Delivery) isAuthorized(r *http.Request) error {
	err := c.token.IsAuthorized(r.Cookies())
	if err != nil {
		return err
	}

	return nil
}

func (d *Delivery) setToken(w http.ResponseWriter, username string) error {
	token, err := d.token.CreateJWT(username)
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

func convertFromUsecaseUserData(userDataUC usecaseDto.UserData) UserData {
	return UserData{
		ID:       userDataUC.ID,
		Username: userDataUC.Username,
		Name:     userDataUC.Name,
	}
}

func convertFromJWTUserData(userDataJWT jwt.UserData) UserData {
	return UserData{
		ID:       userDataJWT.ID,
		Username: userDataJWT.Username,
		Name:     userDataJWT.Name,
	}
}
