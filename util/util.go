package util

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"
)

func GenerateConfigWithPath(templatePath, outputPath string) error {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	// 读取模板文件(go 1.16版本下的读取方式)
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %v", err)
	}

	// 替换占位符
	updatedContent := strings.ReplaceAll(string(content), "{{project_path}}", currentDir)

	// 写入新文件
	err = ioutil.WriteFile(outputPath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("写入新配置文件失败: %v", err)
	}
	fmt.Printf("生成配置文件成功: %s\n", outputPath)
	return nil
}

func GenerateConfigWithPathSafe(templatePath, outputPath string) error {
	// 检查模板文件是否存在
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("模板文件不存在: %s", templatePath)
	}

	// 验证是否为YAML文件
	if !strings.HasSuffix(templatePath, ".yaml") && !strings.HasSuffix(templatePath, ".yml") {
		return fmt.Errorf("模板文件不是YAML格式: %s", templatePath)
	}

	// 检查输出路径的目录是否存在，如果不存在则创建
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 检查输出文件是否已存在
	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("输出文件已存在: %s", outputPath)
	}

	// 执行生成操作
	return GenerateConfigWithPath(templatePath, outputPath)
}
