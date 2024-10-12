package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/utils"
)

type AuthController struct {
	authService  auth.AuthService
	tokenService auth.TokenService
	mu           sync.Mutex
}

func NewAuthController(authService auth.AuthService, tokenService auth.TokenService) *AuthController {
	return &AuthController{
		authService:  authService,
		tokenService: tokenService,
	}
}

// LoginHandler godoc
// @Summary User login
// @Description Authenticate a user with username and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body utils.AuthCredentials true "Credentials for login, including username and password"
// @Success 201 {object} utils.SuccessResponse "Authentication successful"
// @Failure 400 {object} utils.ErrorResponse "Invalid format JSON"
// @Failure 401 {object} utils.ErrorResponse "Incorrect login or password"
// @Router /login [post]
func (c *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds utils.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.SendErrorResponse(w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if c.authService.Authenticate(creds.Username, creds.Password) {
		err := c.setToken(w, creds.Username)
		if err != nil {
			utils.SendErrorResponse(w, "Invalid format JSON", http.StatusUnauthorized)
			return
		}

		utils.SendOKResponse(w, "Authentication successful", http.StatusCreated)

	} else {
		utils.SendErrorResponse(w, "Incorrect login or password", http.StatusUnauthorized)
	}
}

// @Summary Register a new user
// @Description Creates a new user with the provided credentials.
// @Tags auth
// @Accept json
// @Produce json
// @Param creds body utils.RegisterCredentials true "Registration information"
// @Success 201 {object} utils.RegisterResponse "Registration successful"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 409 {object} utils.ErrorResponse "A user with that username already exists"
// @Failure 400 {object} utils.ErrorResponse "User failed to create"
// @Router /signup [post]
func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	log.Println("Пришел запрос на регистрацию")

	var creds utils.RegisterCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if len(creds.Username) < 6 || len(creds.Password) < 8 || creds.Name == "" {
		utils.SendErrorResponse(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if err := c.authService.Registration(creds.Username, creds.Name, creds.Password); err != nil {
		utils.SendErrorResponse(w, "A user with that username already exists", http.StatusConflict)
	} else {
		c.setToken(w, creds.Username)
		userData, err := c.authService.GetUserDataByUsername(creds.Username)
		if err != nil {
			utils.SendErrorResponse(w, "User failed to create", http.StatusBadRequest)
			return
		}

		response := utils.RegisterResponse{
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
// @Success 200 {object} utils.AuthResponse "User data retrieved successfully"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized: token is invalid"
// @Router /auth [get]
func (c *AuthController) AuthHandler(w http.ResponseWriter, r *http.Request) {
	err := c.tokenService.IsAuthorized(r.Cookies())
	if err != nil {
		utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data, err := c.tokenService.GetUserDataByJWT(r.Cookies())
	log.Println("/auth cookie: ", data)
	if err != nil {
		utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := utils.AuthResponse{
		User: data,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (c *AuthController) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := c.isAuthorized(r)
		if err == errors.New("token expired") {
			log.Println("token expired, create new token")
			user, err := c.tokenService.GetUserByJWT(r.Cookies())
			if err != nil {
				utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			c.setToken(w, user.Username)
		}
		if err != nil {
			utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
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
// @Success 200 {object}  utils.SuccessResponse "Logout successful"
// @Failure 401 {object} utils.ErrorResponse "No access token found"
// @Router /logout [post]
func (c *AuthController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	tokenExists := false
	for _, cookie := range r.Cookies() {
		if cookie.Name == "access_token" {
			tokenExists = true
			break
		}
	}

	if !tokenExists {
		utils.SendErrorResponse(w, "No access token found", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "t",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	utils.SendOKResponse(w, "Logout successful", http.StatusOK)
}

func (c *AuthController) isAuthorized(r *http.Request) error {
	err := c.tokenService.IsAuthorized(r.Cookies())
	if err != nil {
		return err
	}

	return nil
}

func (c *AuthController) setToken(w http.ResponseWriter, username string) error {
	token, err := c.tokenService.CreateJWT(username)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	return nil
}
