package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/OverlyDev/go-spotify/internal/logger"
	"github.com/google/uuid"
)

// Generates a UUID based on username and hostname
func getUUID() string {
	var unique string

	user, err := user.Current()
	if err != nil || user.Name == "" {
		unique = "GenericUser"
	} else {
		unique = user.Name
	}

	host, err := os.Hostname()
	if err != nil {
		unique += "@generic-device-name"
	} else {
		unique += "@" + host
	}
	return strings.Replace(uuid.NewSHA1(uuid.NameSpaceDNS, []byte(unique)).String(), "-", "", -1)
}

// Encrypts text with key
func encrypt(key string, text []byte) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.ErrorLogger.Println("(encrypt) error creating cipher:", err)
	}

	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		logger.ErrorLogger.Println("(encrypt) error creating iv:", err)
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext
}

// Decrypts text with key
func decrypt(key string, text []byte) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.ErrorLogger.Println("(decrypt) error creating cipher:", err)
	}
	if len(text) < aes.BlockSize {
		logger.ErrorLogger.Println("(decrypt) ciphertext too short")
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		logger.ErrorLogger.Println("(decrypt) error decoding string:", err)
	}
	return data
}
