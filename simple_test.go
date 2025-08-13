package main

import (
	"fmt"

	"github.com/jdjn123/fast2fa/internal/builder"
)

func main() {
	fmt.Println("开始测试...")

	zipSize := len(builder.TargetZip)
	fmt.Printf("Zip文件大小: %d 字节\n", zipSize)

	templateSize := len(builder.TargetGoTemplate)
	fmt.Printf("模板大小: %d 字符\n", templateSize)

	if zipSize > 0 {
		fmt.Println("✅ 嵌入功能正常！")
	} else {
		fmt.Println("❌ 嵌入功能异常！")
	}
}
