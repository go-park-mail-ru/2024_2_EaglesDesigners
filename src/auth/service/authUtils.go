package service

import (
	"crypto/sha512"
	"encoding/hex"
)

func getSalt() []byte {
	return []byte{93, 108, 25, 43, 92, 102, 255, 179, 11, 87, 186, 198, 254, 160, 164, 56}
}

func hashPassword(password string) string {
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()
	passwordBytes = append(passwordBytes, getSalt()...)
	sha512Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}

func doPasswordsMatch(hashedPassword, currPassword string) bool {
	var currPasswordHash = hashPassword(currPassword)
	return hashedPassword == currPasswordHash
}
