package base64helper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestSavePng(t *testing.T) {
	bytes, err := ioutil.ReadFile("encoded-20241020143425.txt")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	id, err := SavePhotoBase64(string(bytes))

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if _, err := os.Stat("../../../images/" + id.String() + ".png"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		t.Fail()
	}

	os.Remove("../../../images/" + id.String() + ".png")
}
