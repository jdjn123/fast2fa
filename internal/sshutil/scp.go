package sshutil

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// 通过 SCP 传输文件的辅助函数
func ScpFile(ip, user, password, localFile, remotePath string) error {
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

	sess, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	src, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer src.Close()

	info, _ := src.Stat()
	go func() {
		w, _ := sess.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C0755 %d %s\n", info.Size(), filepathBase(localFile))
		io.Copy(w, src)
		fmt.Fprint(w, "\x00")
	}()
	return sess.Run("/usr/bin/scp -t " + remotePath)
}

// 获取文件名的辅助函数 (不使用path包以避免引入额外依赖)
func filepathBase(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
