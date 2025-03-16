package init

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// InitWG initializes and configures a WireGuard server on the machine.
func Init(interfaceName, address, listenPort string) error {
	scriptPath := "./init.sh"

	// Execute the initialization script
	cmd := exec.Command("sudo", "bash", scriptPath, interfaceName, address, listenPort)

	// Redirect command output to both stdout and a buffer
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to attach to stdout: %w", err)
	}
	errOut, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to attach to stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start WireGuard initialization script: %w", err)
	}

	// Stream output to console
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, errOut)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("WireGuard initialization script failed: %w", err)
	}

	fmt.Println("WireGuard initialization completed successfully.")
	return nil
}
