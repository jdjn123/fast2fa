package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// BuildTargetBinary 替换 SECRET 占位符，并编译为 setup2fa
func BuildTargetBinary(secret string) error {
	// 将嵌入的zip文件写入临时文件
	if err := ioutil.WriteFile("google-authenticator.zip", TargetZip, 0644); err != nil {
		return fmt.Errorf("写入zip文件失败: %v", err)
	}
	defer os.Remove("google-authenticator.zip")

	// 替换 SECRET 占位符
	code := fmt.Sprintf(TargetGoTemplate, secret)

	// 写入临时 main.go
	if err := ioutil.WriteFile("main.go", []byte(code), 0644); err != nil {
		return err
	}
	defer os.Remove("main.go")

	// 编译为 setup2fa
	cmd := exec.Command("go", "build", "-o", "setup2fa", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
