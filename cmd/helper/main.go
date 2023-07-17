package main

import (
	"fmt"

	"github.com/OverlyDev/go-spotify/internal/auth"
	"github.com/OverlyDev/go-spotify/internal/settings"
)

func init() {
	settings.LoadHelperEnv()
}

func main() {
	fmt.Println("Hello, I am the helper!")

	auth.Auth.Setup()
}
