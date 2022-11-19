package utils

import "os"

func isExistErrorCheck(err error) bool {
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func IsExistDirOrFile(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return isExistErrorCheck(err)
	}
	return true
}

func IsDir(path string) bool {
	stats, err := os.Stat(path)
	return isExistErrorCheck(err) && stats.IsDir()
}

func IsFile(path string) bool {
	stats, err := os.Stat(path)
	return isExistErrorCheck(err) && !stats.IsDir()
}
