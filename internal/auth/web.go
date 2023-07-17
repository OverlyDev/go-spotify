package auth

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/OverlyDev/go-spotify/internal/logger"
)

var server = &http.Server{Addr: "0.0.0.0:9090"}
var response string

//go:embed authLayout.html
var authLayout string
var authHtml = template.Must(template.New("auth").Parse(authLayout))
var authPage = new(bytes.Buffer)

// Serves the auth page
func authHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write(authPage.Bytes())
}

//go:embed callbackLayout.html
var callbackLayout string
var callbackHtml = template.Must(template.New("callback").Parse(callbackLayout))
var callbackPage = new(bytes.Buffer)

// Serves the callback page
func callbackHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write(callbackPage.Bytes())
}

// Responds to the /token route
func tokenHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	go func() {
		time.Sleep(1 * time.Second)
		server.Shutdown(req.Context())
	}()

	d := json.NewDecoder(req.Body)

	t := struct {
		Url *string `json:"url"`
	}{}

	err := d.Decode(&t)
	if err != nil {
		logger.ErrorLogger.Println("Error decoding json:", err)
		os.Exit(1)
	}

	if t.Url == nil {
		logger.ErrorLogger.Println("No url in json")
		fmt.Printf("%#v\n", d)
		os.Exit(1)
	}

	logger.DebugLogger.Println("Got url:", *t.Url)
	Auth.UserToken, Auth.UserExpire = getTokenFromUrl(*t.Url)
}

// Regex things to get token + other data from response url
func getTokenFromUrl(url string) (string, time.Time) {
	var reState = regexp.MustCompile(fmt.Sprintf(`&state=(?P<state>[[:alnum:]]{%d})$`, stateLength))
	var reError = regexp.MustCompile(`error=(?P<error>\S*)&`)
	var reToken = regexp.MustCompile(`#access_token=(?P<token>[\w-]*)`)
	var reExpire = regexp.MustCompile(`&expires_in=(?P<expire>\d*)&`)
	var theToken string
	var expires = time.Now()

	// Check if state is included in the url
	if !reState.MatchString(url) {
		logger.ErrorLogger.Println("No state included in response! (something fishy is going on)")
		os.Exit(1)
	} else { // Validate the response's state
		state := reState.FindStringSubmatch(url)[1]
		if state != "" && state == generatedState {
			logger.DebugLogger.Println("State matches")
		} else {
			logger.ErrorLogger.Printf("State mismatch! (got: %s should be: %s)\n", state, generatedState)
			os.Exit(1)
		}
	}

	// Check for errors in response
	if reError.MatchString(url) {
		theError := reError.FindStringSubmatch(url)[1]
		logger.ErrorLogger.Println("Received error:", theError)
		os.Exit(1)
	}

	// Get the token from the response
	if reToken.MatchString(url) {
		theToken = reToken.FindStringSubmatch(url)[1]
	} else {
		logger.ErrorLogger.Println("Failed to get token! :(")
		os.Exit(1)
	}

	// Get the token expiration time from the response
	if reExpire.MatchString(url) {
		seconds, err := strconv.Atoi(reExpire.FindStringSubmatch(url)[1])
		if err != nil {
			logger.ErrorLogger.Println(err)
		} else {
			expires = expires.Add(time.Duration(seconds-60) * time.Second)
			logger.DebugLogger.Println("Token expires:", expires.String())
		}

	} else {
		logger.ErrorLogger.Println("Failed to get token expiration! :(")
		os.Exit(1)
	}

	return theToken, expires
}

// Data for the auth page
type AuthPageStruct struct {
	Url    string
	Scopes []string
}

// Wrapper around setting up various http handlers for callback things
func processCallback(redirectUrl string) {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	// Execute the template with the given redirectUrl, storing the final html in authPage bytes buffer
	err1 := authHtml.Execute(authPage, AuthPageStruct{Url: redirectUrl, Scopes: scopes})
	if err1 != nil {
		logger.ErrorLogger.Println("Error executing template:", err1)
	}

	err2 := callbackHtml.Execute(callbackPage, "")
	if err2 != nil {
		logger.ErrorLogger.Println("Error executing template:", err2)
	}

	go func() {
		defer wg.Done()
		http.HandleFunc("/auth", authHandler)
		http.HandleFunc("/callback", callbackHandler)
		http.HandleFunc("/token", tokenHandler)
		server.ListenAndServe()
	}()

	logger.DebugLogger.Println("Waiting for user to hit callback")
	wg.Wait()
}
