package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"math/rand"
)

// reference: https://astaxie.gitbooks.io/build-web-application-with-golang/content/zh/09.5.html
// reference: https://astaxie.gitbooks.io/build-web-application-with-golang/content/zh/09.6.html

const (
	letterBytes   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	salt     = []byte("Way7Wid_MXHd0Q-V")
	salt1    = "!@#$%"
	salt2    = "^&*()"
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

func GetMD5Hash(plaintext string) string {
	md5Hash := md5.New()
	md5Hash.Write([]byte(salt1))
	md5Hash.Write([]byte(plaintext))
	md5Hash.Write([]byte(salt2))

	return hex.EncodeToString(md5Hash.Sum(nil))
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
