package interfaces

import "github.com/go-park-mail-ru/2024_2_EaglesDesigner/login/model"

type ILoginRepository interface {
	FindByUsername(username string) model.User
}
