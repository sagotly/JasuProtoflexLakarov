package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetServerEndpoint() (string, error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return "", fmt.Errorf("failed to get external IP: %w", err)
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response from ifconfig.me: %w", err)
	}

	return strings.TrimSpace(string(ip)), nil
}
