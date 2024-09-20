package interfaces

type ILoginService interface {
	Authenticate(username, password string) bool
}
