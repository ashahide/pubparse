package fileIO

import "os"

func MakeFile(path string) error {

	_, err := os.Stat(path)
	if os.IsExist(err) {
		// File Exists
		os.Remove(path)
		os.Create(path)
	} else if os.IsNotExist(err) {
		// File Doesn't Exist
		os.Create(path)
	} else if err != nil {
		return err
	}

	return nil
}
