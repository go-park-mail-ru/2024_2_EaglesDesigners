package multiparthelper

import (
	"errors"
)

const (
	uploadPath = "/uploads/"
)

var ErrNotImage = errors.New("file is not image")
