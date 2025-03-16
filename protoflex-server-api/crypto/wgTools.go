package cryptoHelp

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// GenerateWGPrivateKey генерирует новый приватный ключ WireGuard.
// Функция выводит логи, проверяет ошибки и возвращает ключ.
func GenerateWGPrivateKey() (string, error) {
	log.Println("Генерируем WireGuard приватный ключ...")
	cmd := exec.Command("sudo", "wg", "genkey")
	// Используем CombinedOutput, чтобы получить вывод, включая ошибки
	output, err := cmd.CombinedOutput()
	log.Printf("Вывод команды wg genkey: %s", output)
	if err != nil {
		log.Printf("Ошибка при генерации приватного ключа: %v", err)
		return "", fmt.Errorf("failed to generate private key: %w", err)
	}
	privateKey := strings.TrimSpace(string(output))
	log.Printf("Сгенерированный приватный ключ: %s", privateKey)
	return privateKey, nil
}

// GenerateWGPublicKey генерирует публичный ключ из заданного приватного ключа.
// Функция выводит логи, проверяет выполнение команды и возвращает публичный ключ.
func GenerateWGPublicKey(privateKey string) (string, error) {
	log.Println("Генерируем WireGuard публичный ключ...")
	// Оборачиваем команду в bash, чтобы передать приватный ключ через echo
	cmd := exec.Command("sudo", "bash", "-c", fmt.Sprintf("echo %s | wg pubkey", privateKey))
	output, err := cmd.CombinedOutput()
	log.Printf("Вывод команды wg pubkey: %s", output)
	if err != nil {
		log.Printf("Ошибка при генерации публичного ключа: %v", err)
		return "", fmt.Errorf("failed to generate public key: %w", err)
	}
	publicKey := strings.TrimSpace(string(output))
	log.Printf("Сгенерированный публичный ключ: %s", publicKey)
	return publicKey, nil
}

// GenerateWGClientConfig генерирует конфигурацию WireGuard клиента в виде строки.
// Здесь никаких команд не выполняется, просто форматирование строки.
func GenerateWGClientConfig(peerPrivateKey, peerPublicKey, allowedIPs, serverEndpoint string) string {
	log.Println("Генерируем конфигурацию WireGuard клиента...")
	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s

[Peer]
PublicKey = %s
Endpoint = %s
AllowedIPs = %s
PersistentKeepalive = 25
`, peerPrivateKey, allowedIPs, peerPublicKey, serverEndpoint, allowedIPs)
	log.Printf("Сгенерированная конфигурация:\n%s", config)
	return config
}
