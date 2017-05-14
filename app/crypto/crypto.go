package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"golang.org/x/crypto/scrypt"
)

// reference: https://astaxie.gitbooks.io/build-web-application-with-golang/content/zh/09.5.html
// reference: https://astaxie.gitbooks.io/build-web-application-with-golang/content/zh/09.6.html

var (
	commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	salt     = []byte("AES256Key-Squirrel32")
)

func EncryptPassword(password []byte) string {
	// Create cryptographic key
	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}

	// Create encrypt algorithm: AES
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Encrypt characters
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	cipherText := make([]byte, len(password))
	cfb.XORKeyStream(cipherText, password)

	return fmt.Sprintf("%x", cipherText)
}
