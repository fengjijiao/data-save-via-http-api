package coreLib

import "os"

func ExistsFile(filepath string) bool {
	// 文件不存在则返回error
	_, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetBaseDirPath() string {
	str, _ := os.Getwd()
	return str
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}