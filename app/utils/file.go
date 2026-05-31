package utils

import "os"

func RedirectToFile(content, file string, append bool) error {
	f, err := OpenFile(file, append)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func OpenFile(name string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	return os.Create(name)
}
