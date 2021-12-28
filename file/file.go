package file

import "os"

//检查路径是否存在

func ExistsDir(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建目录

func CreateDir(dir string) error {
	return os.MkdirAll(dir, 0777)

}
