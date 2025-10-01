package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func showHelp() {
	fmt.Println("Usage: ")
	fmt.Println("pgitp pull [local] [remote]")
	fmt.Println("pgitp push [local] [remote]")
}

func promptPass() (string, error) {
	fmt.Printf("Password: ")
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkDeps() error {
	deps := []string{"git", "gpg", "tar"}
	missing := []string{}
	for _, bin := range deps {
		_, err := exec.LookPath(bin)
		if err != nil {
			missing = append(missing, bin)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing dependencies: %s", strings.Join(missing, ", "))
	}
	return nil
}
