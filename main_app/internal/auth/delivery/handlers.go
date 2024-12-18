package delivery

import (
	"context"
	"errors"
	"html"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	"go.octolab.org/pointer"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/csrf"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/responser"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/validator"
	authv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/protos/gen/go/authv1"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go

type Delivery struct {
	client authv1.AuthClient
	mu     sync.Mutex
}

func New(client authv1.AuthClient) *Delivery {
	return &Delivery{
		client: client,
	}
}

func init() {
	prometheus.MustRegister(requestAuthDuration)
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
	log.Info("Метрики для авторизации зарегистрированы")
}

var requestAuthDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "request_auth_duration_seconds",
	},
	[]string{"method"},
)

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
// @Router /login [post].
func (d *Delivery) LoginHandler(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "LoginHandler")
	}()

	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на аутентификацию")

	var creds models.AuthReqDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &creds); err != nil {
		log.Errorf("не удалось распарсить json")
		responser.SendError(ctx, w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if err := validator.Check(creds); err != nil {
		log.Errorf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	grpcResp, err := d.client.Authenticate(
		ctx,
		&authv1.AuthRequest{
			Username: creds.Username,
			Password: creds.Password,
		})
	if err != nil {
		log.Errorf("не удалось аутентифицировать пользователя")
		responser.SendError(ctx, w, "Invalid format JSON", http.StatusUnauthorized)
	}

	if grpcResp.GetIsAuthenticated() {
		err := d.setTokens(w, r, creds.Username)
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
// @Router /signup [post].
func (d *Delivery) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "RegisterHandler")
	}()

	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на регистрацию")

	var creds models.RegisterReqDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &creds); err != nil {
		responser.SendError(ctx, w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if err := validator.Check(creds); err != nil {
		log.Errorf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	_, err := d.client.Registration(
		ctx,
		&authv1.RegistrationRequest{
			Username: creds.Username,
			Name:     creds.Username,
			Password: creds.Password,
		})
	if err != nil {
		responser.SendError(ctx, w, "A user with that username already exists", http.StatusConflict)
		return
	}

	log.Println("получение данных пользователя")

	d.setTokens(w, r, creds.Username)
	userData, err := d.client.GetUserDataByUsername(ctx, &authv1.GetUserDataByUsernameRequest{Username: creds.Username})
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

	// responser.SendStruct(ctx, w, response, http.StatusCreated)
	jsonResp, err := easyjson.Marshal(response)
	responser.SendJson(ctx, w, jsonResp, err, http.StatusCreated)
}

// AuthHandler godoc
// @Summary Authenticate a user
// @Description Retrieve user data based on the JWT token present in the cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.UserDataRespDTO "User data retrieved successfully"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized: token is invalid"
// @Router /auth [get].
func (d *Delivery) AuthHandler(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "AuthHandler")
	}()

	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на авторизацию")

	user, ok := ctx.Value(models.UserKey).(models.User)
	if !ok {
		responser.SendError(ctx, w, "User not found", http.StatusNotFound)
		return
	}

	userData, err := d.client.GetUserDataByUsername(ctx, &authv1.GetUserDataByUsernameRequest{Username: user.Username})
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

	// responser.SendStruct(ctx, w, response, http.StatusOK)
	jsonResp, err := easyjson.Marshal(response)
	responser.SendJson(ctx, w, jsonResp, err, http.StatusOK)
}

var errTokenExpired = errors.New("токен истек")

func (d *Delivery) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := d.parseCookies(r.Cookies())
		if err != nil {
			log.Println("не получилось получить токен")
			responser.SendError(ctx, w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := d.client.IsAuthorized(ctx, &authv1.Token{Token: token})
		if err == errTokenExpired {
			log.Println("токен истек, создается новый токен")
			d.setTokens(w, r, user.Username)
		}

		if err != nil && err != errTokenExpired {
			log.Println("токен невалиден")
			responser.SendError(ctx, w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, models.UserKey, convertFromGRPCUser(user))
		ctx = context.WithValue(ctx, models.MuxParamsKey, mux.Vars(r))

		r = r.WithContext(ctx)

		next(w, r)
	}
}

func (d *Delivery) Csrf(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// token := r.Header.Get("not_csrf")

		// user := r.Context().Value(models.UserKey).(jwt.User)

		// err := csrf.CheckCSRF(token, user.ID, user.Username)
		// if err != nil {
		// 	if err == errTokenExpired {
		// 		responser.SendError(context.Background(), w, "csrf expired", http.StatusForbidden)
		// 		return
		// 	}
		// 	responser.SendError(context.Background(), w, "Invalid csrf", http.StatusForbidden)
		// 	return
		// }

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
// @Router /logout [post].
func (d *Delivery) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "LogoutHandler")
	}()

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
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	log.Println("разлогин прошел успешно")

	responser.SendOK(w, "Logout successful", http.StatusOK)
}

func (d *Delivery) setTokens(w http.ResponseWriter, r *http.Request, username string) error {
	ctx := r.Context()
	grcpResp, err := d.client.CreateJWT(ctx, &authv1.CreateJWTRequest{Username: username})
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    grcpResp.GetToken(),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	csrf, err := csrf.CreateCSRF(grcpResp.GetToken())
	if err != nil {
		log.Printf("Auth setTokens: не удалось создать csrf токен: %v", err)
		return err
	}

	w.Header().Set("X-CSRF-Token", csrf)
	w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")

	return nil
}

func (d *Delivery) parseCookies(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie.Value, nil
		}
	}
	return "", errors.New("cookie does not exist")
}

func convertUserDataToDTO(userData *authv1.GetUserDataByUsernameResponse) models.UserDataRespDTO {
	return models.UserDataRespDTO{
		ID:        uuid.MustParse(userData.GetID()),
		Username:  html.EscapeString(userData.GetUsername()),
		Name:      html.EscapeString(userData.GetName()),
		AvatarURL: pointer.ToStringOrNil(userData.GetAvatarURL()),
	}
}

func convertFromGRPCUser(user *authv1.UserJWT) models.User {
	return models.User{
		ID:       uuid.MustParse(user.GetID()),
		Username: user.GetUsername(),
		Name:     user.GetName(),
		Password: user.GetPassword(),
		Version:  user.GetVersion(),
	}
}
