package utils

import (
	"fmt"
	"os"
)

func WriteFile(file, data string, empty bool) {
	var flag int
	if empty {
		flag = os.O_CREATE | os.O_TRUNC | os.O_RDWR
	} else {
		flag = os.O_CREATE | os.O_APPEND | os.O_RDWR
	}

	f, err := os.OpenFile(file, flag, 0660)
	if err != nil {
		fmt.Printf("Cannot open file %s!\n", file)
		return
	}
	defer f.Close()
	f.WriteString(data)
}

func ReadFile(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Cannot open file %s!\n", file)
		return nil
	}
	data := make([]byte, 100)
	n, err := f.Read(data)
	if err != nil {
		fmt.Printf("read file error %s!\n", err)
		return nil
	}
	return data[:n]
}

func ReadFileToString(file string) string {
	return string(ReadFile(file))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
