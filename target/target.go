package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
)

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

func main() {
    secret := "%s"

    runCmd("bash -c 'which google-authenticator || (apt-get update && apt-get install -y libpam-google-authenticator || yum install -y google-authenticator)'")

    homeDir, _ := os.UserHomeDir()
    os.WriteFile(homeDir+"/.google_authenticator", []byte(secret+"\n\" RATE_LIMIT 3 30\n\" DISALLOW_REUSE\n\" TOTP_AUTH\n"), 0600)

    runCmd("bash -c 'grep pam_google_authenticator.so /etc/pam.d/sshd || sed -i \"1iauth required pam_google_authenticator.so\" /etc/pam.d/sshd'")
    runCmd("bash -c 'grep ChallengeResponseAuthentication /etc/ssh/sshd_config && sed -i \"s/ChallengeResponseAuthentication no/ChallengeResponseAuthentication yes/\" /etc/ssh/sshd_config || echo \"ChallengeResponseAuthentication yes\" >> /etc/ssh/sshd_config'")
    runCmd("systemctl restart sshd")

    fmt.Println("2FA 配置完成")
}
