package sdk

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type ClientWithIP struct {
	Client *ClientWithResponses
	IP     string
}

func CreateClient(logger *log.Logger, ip string) *ClientWithIP {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: jar,
	}
	serverAddr := viper.GetString("tui.server")
	client, err := NewClientWithResponses(serverAddr, WithHTTPClient(httpClient))
	if err != nil {
		logger.Fatal("Failed to create client", "err", err)
	}

	return &ClientWithIP{
		Client: client,
		IP:     ip,
	}
}
