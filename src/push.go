package main

import (
	"fmt"
	"os"
	"os/exec"
)

func putDir(local, remote string) error {
	pass, err := promptPass()
	if err != nil {
		return err
	}
	localFiles, err := os.ReadDir(local)
	if err != nil {
		return err
	}
	for _, file := range localFiles {
		localPath := fmt.Sprintf("%s/%s", local, file.Name())
		remotePath := fmt.Sprintf("%s/%s", remote, file.Name())
		err = putFile(pass, localPath, remotePath)
		if err != nil {
			return err
		}
	}
	err = gitPush(remote)
	return err
}

func putFile(pass, local, remote string) error {
	tar := exec.Command("tar", "-czf", "-", local)
	gpg := exec.Command("gpg", "--symmetric", "--cipher-algo", "AES256", "--armor", "--batch", "--yes", "--passphrase", pass)
	tar.Stderr = os.Stderr
	pipe, err := tar.StdoutPipe()
	if err != nil {
		return err
	}
	remoteFile, err := os.OpenFile(remote, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer remoteFile.Close()
	gpg.Stdin = pipe
	gpg.Stdout = remoteFile
	gpg.Stderr = os.Stderr
	if err = tar.Start(); err != nil {
		return err
	}
	if err = gpg.Start(); err != nil {
		return err
	}
	if err = tar.Wait(); err != nil {
		return err
	}
	if err = gpg.Wait(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s -> %s\n", local, remote)
	return nil
}

func gitPush(remote string) error {
	add := exec.Command("git", "-C", remote, "add", "-A")
	add.Stdout = os.Stdout
	add.Stderr = os.Stderr
	if err := add.Run(); err != nil {
		return err
	}
	commit := exec.Command("git", "-C", remote, "commit", "-m", "WIP")
	commit.Stdout = os.Stdout
	commit.Stderr = os.Stderr
	if err := commit.Run(); err != nil {
		return err
	}
	push := exec.Command("git", "-C", remote, "push", "-v")
	push.Stdout = os.Stdout
	push.Stderr = os.Stderr
	if err := push.Run(); err != nil {
		return err
	}
	return nil
}
