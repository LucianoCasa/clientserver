package fileutils

import (
	"os"
)

// Lê o conteúdo de um arquivo e retorna como string
func ReadFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Escreve conteúdo em um arquivo, substituindo se já existir
func WriteNewFile(newfile string, content string) error {
	file, err := os.Create(newfile)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
