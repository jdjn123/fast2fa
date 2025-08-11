package sshutil

import (
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

// 通过 SSH 执行命令的辅助函数
func SSHExec(ip, user, password, cmd string) error {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	conn, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = io.Discard
	session.Stderr = io.Discard

	return session.Run(cmd)
}
