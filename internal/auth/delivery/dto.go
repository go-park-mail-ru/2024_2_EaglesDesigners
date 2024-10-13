package delivery

// @Schema
type AuthCredentials struct {
	Username string `json:"username" example:"user11"`
	Password string `json:"password"  example:"12345678"`
}

// @Schema
type RegisterCredentials struct {
	Username string `json:"username" example:"killer1994"`
	Name     string `json:"name" example:"Vincent Vega"`
	Password string `json:"password" example:"go_do_a_crime"`
}

// @Schema
type RegisterResponse struct {
	Message string   `json:"message" example:"Registration successful"`
	User    UserData `json:"user"`
}

// @Schema
type AuthResponse struct {
	User UserData `json:"user"`
}

// @Schema
type SignupResponse struct {
	Error  string `json:"error"`
	Status string `json:"status" example:"error"`
}

// @Schema
type User struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"mavrodi777"`
	Name     string `json:"name" example:"Мафиозник"`
	Password string `json:"password" example:"1234567890"`
	Version  int64  `json:"version" example:"1"`
}

type UserData struct {
	ID       int64  `json:"id" example:"2"`
	Username string `json:"username" example:"user12"`
	Name     string `json:"name" example:"Dr Peper"`
}
