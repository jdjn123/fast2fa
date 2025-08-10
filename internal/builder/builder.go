package builder

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
)

func BuildTargetBinary(secret string) error {
    // 读取 target.go 模板
    targetCode, err := ioutil.ReadFile(filepath.Join("target", "target.go"))
    if err != nil {
        return err
    }

    // 替换 SECRET 占位符
    code := fmt.Sprintf(string(targetCode), secret)

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
