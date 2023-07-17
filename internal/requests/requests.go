package requests

import (
	"context"

	"github.com/OverlyDev/go-spotify/internal/logger"
	"github.com/OverlyDev/go-spotify/internal/settings"
	"github.com/carlmjohnson/requests"
)

func MakeRequestStr(url, path string) string {
	var res string
	err := requests.URL(url).Path(path).UserAgent(settings.UserAgent).ToString(&res).Fetch(context.Background())
	if err != nil {
		logger.ErrorLogger.Printf("Failed to make request to: %s; err: %s", url+"/"+path, err)
	}
	return res
}
