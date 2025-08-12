package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

// 打印帮助信息
func init() {
	flag.Usage = printHelp
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}
	if os.Args[1] == "--help" {
		printHelp()
		os.Exit(0)
	}
}
func printHelp() {
	fmt.Println(`用法: manager [--hosts 文件名]
参数:
  --hosts   指定主机列表文件，格式为 ip,user,password（默认 hosts.csv）
  --help    显示本帮助信息`)
}

// 从文件读取主机列表
func readHostsFromFile(filename string) ([]Host, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []Host
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return nil, fmt.Errorf("格式错误，期望 ip,user,password: %s", line)
		}
		host := Host{
			IP:       strings.TrimSpace(parts[0]),
			User:     strings.TrimSpace(parts[1]),
			Password: strings.TrimSpace(parts[2]),
		}
		hosts = append(hosts, host)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return hosts, nil
}

func main() {
	// 定义命令行参数，默认 hosts.csv
	hostsFile := flag.String("hosts", "hosts.csv", "目标主机列表文件，格式 ip,user,password")
	flag.Parse()

	hosts, err := readHostsFromFile(*hostsFile)
	if err != nil {
		log.Fatalf("读取主机列表失败: %v", err)
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

	// 只在本地编译一次
	fmt.Println("[*] 本地编译目标机程序")
	if err := builder.BuildTargetBinary(secret.Secret()); err != nil {
		log.Fatalf("本地编译目标机程序失败: %v", err)
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
