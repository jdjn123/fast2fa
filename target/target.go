package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// 执行命令的辅助函数
func runCmd(cmd string) {
	parts := strings.Split(cmd, " ")
	c := exec.Command(parts[0], parts[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Println("命令执行失败:", err)
		os.Exit(1)
	}
}

// 主函数，配置 2FA
func main() {
	secret := "%s"
	// 安装 Google Authenticator PAM 模块
	runCmd("bash -c 'which google-authenticator || (apt-get update && apt-get install -y git build-essential autoconf libtool automake libpam0g-dev || yum install -y git gcc make autoconf libtool pam-devel automake)'")
	runCmd("bash -c 'test -d /tmp/google-authenticator-libpam || git clone https://github.com/google/google-authenticator-libpam.git /tmp/google-authenticator-libpam'")
	runCmd("bash -c 'cd /tmp/google-authenticator-libpam && chmod +x bootstrap.sh && ./bootstrap.sh && ./configure && make && make install'")
	// 配置 PAM 和 SSHD
	homeDir, _ := os.UserHomeDir()
	os.WriteFile(homeDir+"/.google_authenticator", []byte(secret+"\n\" RATE_LIMIT 3 30\n\" DISALLOW_REUSE\n\" TOTP_AUTH\n"), 0600)

	runCmd("bash -c 'grep pam_google_authenticator.so /etc/pam.d/sshd || sed -i \"1iauth required pam_google_authenticator.so\" /etc/pam.d/sshd'")
	//判断是否是交互式shell，如果是则添加pam_permit.so
	runCmd("bash -c 'grep pam_google_authenticator.so /etc/pam.d/sshd || sed -i \"1iauth [success=1 default=ignore] pam_google_authenticator.so\" /etc/pam.d/sshd'")
	runCmd("bash -c 'grep pam_permit.so /etc/pam.d/sshd || sed -i \"2iauth required pam_permit.so\" /etc/pam.d/sshd'")
	// 修改sshd_config
	runCmd("bash -c 'grep ChallengeResponseAuthentication /etc/ssh/sshd_config && sed -i \"s/ChallengeResponseAuthentication no/ChallengeResponseAuthentication yes/\" /etc/ssh/sshd_config || echo \"ChallengeResponseAuthentication yes\" >> /etc/ssh/sshd_config'")
	runCmd("systemctl restart sshd")

	fmt.Println("2FA 配置完成")
}
