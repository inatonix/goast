package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type File struct {
	Path    string
	Content string
}

func addFileContent(messages *[]openai.ChatCompletionMessage, filePath string) error {
	//fmt.Println(filePath)
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	*messages = append(*messages, openai.ChatCompletionMessage{
		Role:    "system",
		Content: filePath + "\n```\n" + string(file) + "\n```\n",
	})

	return nil
}

func parseFileContents(outputFilePath string, output string) error {
	// 正規表現パターン
	pattern := "```.*\\n([\\s\\S]*?)\\n```"

	// 正規表現をコンパイル
	re := regexp.MustCompile(pattern)

	// マッチした部分を格納するスライス
	matches := re.FindAllStringSubmatch(output, -1)

	// マッチした内容を格納する配列
	var result string

	// マッチした部分を処理
	for _, match := range matches {
		if len(match) == 2 {
			contents := match[1]
			result = contents
		}
	}

	fmt.Println(outputFilePath)
	file, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open the file for output: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fprint(file, result)
	if err != nil {
		return fmt.Errorf("failed to write the file for output: %w", err)
	}

	return nil
}

func makeFileList() ([]string, error) {
	file, err := os.Open("gen/.genfile")
	if err != nil {
		return nil, fmt.Errorf("failed to open .read file: %w", err)
	}
	defer file.Close()

	var paths []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			paths = append(paths, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan .read file: %w", err)
	}

	return paths, nil
}
