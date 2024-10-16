package usecase

type User struct {
	ID       int64
	Username string
	Name     string
	Password string
	Version  int64
}

type UserData struct {
	ID       int64
	Username string
	Name     string
}
