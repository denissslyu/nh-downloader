package utils

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// StructToMap convert a struct which all fields are string to a map[string]string
func StructToMap(obj interface{}) map[string]string {
	objValue := reflect.ValueOf(obj)

	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}
	if objValue.Kind() != reflect.Struct {
		return nil
	}
	objType := objValue.Type()

	result := make(map[string]string)
	for i := 0; i < objValue.NumField(); i++ {
		field := objType.Field(i)
		value := objValue.Field(i)
		tag := field.Tag.Get("json")

		if tag == "-" {
			continue // 忽略带有"-"标签的字段
		}

		key := field.Name
		if tag != "" {
			key = tag // 使用标签定义的名字作为键
		} else {
			key = strings.ToLower(key) // 将字段名转换为小写作为键
		}
		if valueStr, ok := value.Interface().(string); ok && valueStr != "" {
			result[key] = valueStr
		}
	}

	return result
}

func GetPrettyTitle(title string) string {
	title = removeEnclosedText(title, '[', ']')
	title = removeEnclosedText(title, '(', ')')

	// there are some repeated single brackets
	title = strings.ReplaceAll(title, "[", "")
	title = strings.ReplaceAll(title, "]", "")
	title = strings.ReplaceAll(title, "(", "")
	title = strings.ReplaceAll(title, ")", "")
	return strings.TrimSpace(title)
}

func removeEnclosedText(str string, startTag byte, endTag byte) string {
	stack := make([]int, 0)

	for i := 0; i < len(str); i++ {
		if str[i] == startTag {
			stack = append(stack, i)
		} else if str[i] == endTag && len(stack) > 0 {
			startIndex := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			str = str[:startIndex] + str[i+1:]
			i = startIndex - 1
		}
	}

	return str
}

func ZipFiles(sourceDir, destinationFile string) error {
	zipFile, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		zipFile, err := zipWriter.Create(info.Name())
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFile, file)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func SaveJsonToFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}
