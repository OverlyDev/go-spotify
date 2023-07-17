package settings

import (
	_ "embed"
	"os"

	"github.com/OverlyDev/go-spotify/internal/logger"
	"github.com/caarlos0/env/v9"
)

//go:embed .clientid
var clientId string

//go:embed .clientsecret
var clientSecret string

var UserAgent = "go-spotify/0.0.1"

var HelperSettings HelperEnv

// Settings for the Helper
type HelperEnv struct {
	Debug     bool   `env:"DEBUG" envDefault:"false"`
	ApiID     string `env:"API_ID,notEmpty,unset"`
	ApiSecret string `env:"API_SECRET,notEmpty,unset"`
}

// "load environment variables" for the helper
// will be more fleshed out later on
func LoadHelperEnv() {
	env.Parse(&HelperSettings)
	// if err := env.Parse(&HelperSettings); err != nil {
	// 	logger.ErrorLogger.Printf("%+v\n", err)
	// 	os.Exit(1)
	// }
	if HelperSettings.Debug {
		logger.EnableDebugLogging()
	}
	// logger.DebugLogger.Printf("Loaded helper settings: %#v\n", HelperSettings)

	// For now, pull in client id/secret from embed vars defined above
	HelperSettings.ApiID = clientId
	HelperSettings.ApiSecret = clientSecret

	// Make sure cliendId is provided
	if HelperSettings.ApiID == "" {
		logger.ErrorLogger.Println("Missing clientId")
		os.Exit(1)
	}

	// Make sure clientSecret is provided
	if HelperSettings.ApiSecret == "" {
		logger.ErrorLogger.Println("Missing clientSecret")
		os.Exit(1)
	}
}
