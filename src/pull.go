package main

import (
	"fmt"
	"os"
	"os/exec"
)

func getDir(local, remote string) error {
	err := gitPull(remote)
	if err != nil {
		return err
	}
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
		err = getFile(pass, localPath, remotePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func getFile(pass, local, remote string) error {
	gpg := exec.Command("gpg", "--decrypt", "--passphrase", pass, "--batch")
	tar := exec.Command("tar", "-xzf", "-", local)
	remoteFile, err := os.OpenFile(remote, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}
	defer remoteFile.Close()
	gpg.Stdin = remoteFile
	gpg.Stderr = os.Stderr
	pipe, err := gpg.StdoutPipe()
	if err != nil {
		return err
	}
	tar.Stdin = pipe
	tar.Stdout = os.Stdout
	tar.Stderr = os.Stderr
	if err = gpg.Start(); err != nil {
		return err
	}
	if err = tar.Start(); err != nil {
		return err
	}
	if err = gpg.Wait(); err != nil {
		return err
	}
	if err = tar.Wait(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s -> %s\n", remote, local)
	return nil
}

func gitPull(remote string) error {
	push := exec.Command("sh", "-c", fmt.Sprintf("git -C %s pull -v", remote))
	push.Env = append(os.Environ(), "PATH=/data/data/com.termux/files/usr/bin:$PATH")
	push.Stdout = os.Stdout
	push.Stderr = os.Stderr
	if err := push.Run(); err != nil {
		return err
	}
	return nil
}
