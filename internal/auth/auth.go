package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/OverlyDev/go-spotify/internal/logger"
	"github.com/OverlyDev/go-spotify/internal/settings"
	"github.com/carlmjohnson/requests"
)

// Holds the auth information
var Auth AuthStruct

// Used for state generation
var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
var stateLength = 32
var generatedState = generateState(stateLength)

// URLs
var authUrl = "http://127.0.0.1:9090/auth"
var callbackUrl = "http://127.0.0.1:9090/callback"

// Spotify permission scopes
var scopes = []string{
	"playlist-read-private",       // Read access to user's private playlists
	"playlist-read-collaborative", // Include collaborative playlists when requesting a user's playlists
	"user-top-read",               // Read access to a user's top artists and tracks
	"user-read-recently-played",   // Read access to a user’s recently played tracks
	"user-library-read",           // Read access to a user's library
}

// Generates alphanumeric string for use in callback url state
func generateState(length int) string {
	ll := len(chars)
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%ll]
	}
	return string(b)
}

// Handles obtaining client authentication token
func requestClientAuth() {
	logger.DebugLogger.Println("Requesting client auth token from Spotify")
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", settings.HelperSettings.ApiID, settings.HelperSettings.ApiSecret)
	var res AuthResponse
	err := requests.URL("https://accounts.spotify.com/api/token").
		UserAgent(settings.UserAgent).
		ContentType("application/x-www-form-urlencoded").
		BodyBytes([]byte(body)).
		ToJSON(&res).
		Fetch(context.Background())

	if err != nil {
		logger.ErrorLogger.Printf("Failed to do client auth; err: %s", err)
		os.Exit(1)
	} else {
		logger.DebugLogger.Println("Successfully auth'd with Spotify")
		Auth.ClientToken = res.Token
		Auth.ClientExpire = time.Now().Add(time.Duration(res.Expire-60) * time.Second)
	}
}

// Handles obtaining user authentication token
func requestUserAuth() {
	logger.DebugLogger.Println("Requesting user auth token from Spotify")

	reqUrl := "https://accounts.spotify.com/authorize?response_type=token"
	reqUrl += "&client_id=" + settings.HelperSettings.ApiID
	reqUrl += "&scope=" + strings.Join(scopes, " ")
	reqUrl += "&redirect_uri=" + callbackUrl
	reqUrl += "&state=" + generatedState
	logger.DebugLogger.Println("Request URL:", reqUrl)

	logger.InfoLogger.Println("Please authenticate with Spotify in the browser window that opened")
	logger.InfoLogger.Println("If a window didn't open, manually navigate to the following link in your browser:")
	fmt.Printf("\n\t%s\n\n", authUrl)

	openUrlInDefaultBrowser(authUrl)

	processCallback(reqUrl)
	logger.DebugLogger.Println("Got token:", Auth.UserToken)
}

// Helper to open a url in the user's default browser
func openUrlInDefaultBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "openUrlInDefaultBrowser"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-openUrlInDefaultBrowser"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
