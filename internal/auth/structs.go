package auth

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/OverlyDev/go-spotify/internal/logger"
)

var authFilename = ".auth"
var id = getUUID()

// Spotify response for client authentication
type AuthResponse struct {
	Token  string `json:"access_token"`
	Expire int    `json:"expires_in"`
}

// Holds received authentication tokens
type AuthStruct struct {
	ClientToken  string
	ClientExpire time.Time
	UserToken    string
	UserExpire   time.Time
}

// Write the stuct to a json file
func (a *AuthStruct) save() {
	jsonData, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		logger.ErrorLogger.Println("Failed to marshal auth data")
	}
	err = ioutil.WriteFile(authFilename, encrypt(id, jsonData), 0644)
	if err != nil {
		logger.ErrorLogger.Println("Failed to save auth file")
	}
}

// Load from a json file and populate the struct
func (a *AuthStruct) load() {
	data, err := ioutil.ReadFile(authFilename)
	if err != nil {
		logger.ErrorLogger.Println("Failed to read auth file")
	}
	e := json.Unmarshal(decrypt(id, data), a)
	if e != nil {
		logger.ErrorLogger.Println("Failed to load auth file")
	}
}

// Gets authentication tokens squared away
func (a *AuthStruct) Setup() {
	// No saved auth file
	if _, err := os.Stat(authFilename); err != nil {
		logger.DebugLogger.Println("No auth file")
		requestClientAuth()
		requestUserAuth()

		// Saved auth file
	} else {
		logger.DebugLogger.Println("Loading auth file")
		a.load()
	}

	// Check client token
	if a.ClientToken == "" || a.ClientExpire.Before(time.Now()) {
		logger.DebugLogger.Println("Getting new client token")
		requestClientAuth()
	}

	// Check user token
	if a.UserToken == "" || a.UserExpire.Before(time.Now()) {
		logger.DebugLogger.Println("Getting new user token")
		requestUserAuth()
	}
	logger.DebugLogger.Println("Saving auth file")
	a.save()
}
