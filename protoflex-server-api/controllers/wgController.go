package controllers

import (
	"fmt"
	"os/exec"
	cryptoHelp "protoflex-server-api/crypto"
)

type WgController struct {
	wgInterface   string
	allowedPeerIp string
}

func NewWGController(wgInterface string, allowedPeerIp string) *WgController {
	return &WgController{wgInterface: wgInterface, allowedPeerIp: allowedPeerIp}
}

func (wg *WgController) AddPeerToWGServer(serverEndpoint string) (string, error) {
	// Generate a private key the peer
	peerPrivateKey, err := cryptoHelp.GenerateWGPrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate private key for peer: %w", err)
	}

	// Generate the corresponding public key for the peer
	peerPublicKey, err := cryptoHelp.GenerateWGPublicKey(peerPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate public key for peer: %w", err)
	}

	// Add the peer to the WireGuard server using wg command
	cmd := exec.Command("wg", "set", wg.wgInterface, "peer", peerPublicKey, "endpoint", serverEndpoint+":51820", "allowed-ips", wg.allowedPeerIp)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to add peer to WireGuard server: %w, input: %s, output: %s", err, fmt.Sprint("wg", "set", wg.wgInterface, "peer", peerPublicKey, "endpoint", serverEndpoint+":51820", "allowed-ips", wg.allowedPeerIp), output)
	}

	// Generate the client configuration
	clientConfig := cryptoHelp.GenerateWGClientConfig(peerPrivateKey, peerPublicKey, wg.allowedPeerIp, serverEndpoint)
	return clientConfig, nil
}
