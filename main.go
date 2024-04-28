package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 根据文件路径计算校验和
func calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New() // Change this to the hash algorithm you prefer
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// 检查文件是否在忽略列表中
func isIgnored(fileName string) bool {
	ignored := []string{"md5", "md4", "sha1", "sha256", "sha384", "sha512", "ripemd160", "panama", "tiger", "md2", "adler32", "crc32", "checksum"}
	for _, ignore := range ignored {
		if strings.EqualFold(fileName, ignore) {
			return true
		}
	}
	return false
}

func main() {
	// 获取当前文件夹路径
	dir, err := os.Getwd()
	fmt.Println("当前路径:", dir)
	if err != nil {
		fmt.Println("无法获取当前文件夹路径:", err)
		return
	}

	// 遍历当前文件夹下的所有文件
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("无法读取当前文件夹下的文件:", err)
		return
	}

	// 统计校验文件数量、无校验文件数量、校验正确数量、校验错误数量
	var checksumFiles, correctChecksum, incorrectChecksum int

	for _, file := range files {
		if file.IsDir() {
			continue // 跳过文件夹
		}
		fileName := file.Name()
		if fileName == "desktop.ini" {
			continue
		}
		fileFullName := filepath.Join(dir, fileName)
		if isIgnored(strings.TrimPrefix(filepath.Ext(fileFullName), ".")) && file.Size() < 1024 {
			continue // 跳过忽略列表中的文件
		}

		checksumFilePath := fileFullName + ".sha256"
		checksum, err := calculateChecksum(fileFullName)
		if err != nil {
			fmt.Printf("无法计算文件 %s 的校验和: %v\n", fileFullName, err)
			continue
		}

		// 检查是否存在校验文件
		if _, err := os.Stat(checksumFilePath); os.IsNotExist(err) {
			// 不存在
			// 生成校验文件
			fmt.Printf("生成 %s 的校验文件\n", fileName)
			if err := ioutil.WriteFile(checksumFilePath, []byte(checksum), 0644); err != nil {
				fmt.Printf("无法生成校验文件 %s: %v\n", checksumFilePath, err)
				continue
			}
			checksumFiles++
		} else {
			// 存在
			// 读取校验文件
			content, err := ioutil.ReadFile(checksumFilePath)
			if err != nil {
				fmt.Printf("无法读取校验文件 %s: %v\n", checksumFilePath, err)
				continue
			}
			if strings.TrimSpace(string(content)) == checksum {
				correctChecksum++
			} else {
				incorrectChecksum++
				fmt.Printf("%s 的校验文件错误\n", fileName)
			}
		}
	}

	// 输出统计结果
	fmt.Println("======================================")
	fmt.Printf("总共生成了 %d 个校验文件\n", checksumFiles)
	fmt.Printf("校验正确的文件数量为 %d\n", correctChecksum)
	fmt.Printf("校验错误的文件数量为 %d\n", incorrectChecksum)

}
