package utils

import (
	shell "github.com/ipfs/go-ipfs-api"
	gohttp "net/http"
	"time"
)

func NewShellDbchain(url string) *shell.Shell {
	c := &gohttp.Client{
		Transport: &gohttp.Transport{
			Proxy:             gohttp.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
		Timeout: 2 * time.Second,
	}

	return shell.NewShellWithClient(url, c)
}
