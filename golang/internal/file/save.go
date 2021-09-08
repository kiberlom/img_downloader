package file

import "os"

// сохранение файла изображения
func Save(b []byte) error {
	f, err := os.Create("img.jpeg")
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}
