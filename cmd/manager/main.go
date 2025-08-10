package main

import (
    "fmt"
    "log"

    "github.com/pquerna/otp"
    "github.com/pquerna/otp/totp"

    "github.com/jdjn123/fast2fa/internal/builder"
    "github.com/jdjn123/fast2fa/internal/sshutil"
)

type Host struct {
    IP       string
    User     string
    Password string
}

func main() {
    // 目标主机列表
    hosts := []Host{
        {"192.168.1.101", "root", "password1"},
        {"192.168.1.102", "root", "password2"},
    }

    // 生成统一的 TOTP Secret
    secret, _ := totp.Generate(totp.GenerateOpts{
        Issuer:      "MyServerCluster",
        AccountName: "admin",
        Algorithm:   otp.AlgorithmSHA1,
        Digits:      otp.DigitsSix,
        Period:      30,
    })
    fmt.Println("[+] 统一 TOTP Secret:", secret.Secret())

    // 编译目标机程序
    if err := builder.BuildTargetBinary(secret.Secret()); err != nil {
        log.Fatal("编译目标机程序失败:", err)
    }

    // 分发并执行
    for _, h := range hosts {
        fmt.Println("[*] 处理主机:", h.IP)
        if err := sshutil.ScpFile(h.IP, h.User, h.Password, "setup2fa", "/tmp/"); err != nil {
            log.Println("SCP 失败:", err)
            continue
        }
        if err := sshutil.SSHExec(h.IP, h.User, h.Password, "chmod +x /tmp/setup2fa && /tmp/setup2fa"); err != nil {
            log.Println("执行失败:", err)
            continue
        }
    }
}
