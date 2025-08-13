package sshutil

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jdjn123/fast2fa/internal/builder"
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
	fmt.Println("[+] 连接成功")
	sess, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	fmt.Println("[+] 创建会话成功")
	src, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer src.Close()
	fmt.Println("[+] 打开文件成功")
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

// ScpGoogleAuthenticator 将 google-authenticator.zip 复制到远程主机的 /tmp 目录
func ScpGoogleAuthenticator(ip, user, password string) error {
	// 使用嵌入的 zip 文件数据
	zipData := builder.TargetZip

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "google-authenticator-*.zip")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// 写入 zip 数据到临时文件
	if _, err := tmpFile.Write(zipData); err != nil {
		return fmt.Errorf("写入临时文件失败: %v", err)
	}

	// 确保数据写入磁盘
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("同步临时文件失败: %v", err)
	}

	// 关闭文件以便 SCP 可以读取
	tmpFile.Close()

	fmt.Printf("[+] 正在将 google-authenticator.zip 复制到 %s:/tmp/\n", ip)

	// 直接使用 SCP 协议传输，指定正确的远程文件名
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

	src, err := os.Open(tmpFile.Name())
	if err != nil {
		return err
	}
	defer src.Close()

	info, _ := src.Stat()
	go func() {
		w, _ := sess.StdinPipe()
		defer w.Close()
		// 指定远程文件名为 google-authenticator.zip
		fmt.Fprintf(w, "C0755 %d google-authenticator.zip\n", info.Size())
		io.Copy(w, src)
		fmt.Fprint(w, "\x00")
	}()
	return sess.Run("/usr/bin/scp -t /tmp/")
}
