package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	maxFileSizeWithoutProgressBar int64 = 1024
	ErrUnsupportedFile                  = errors.New("unsupported file")
	ErrOffsetExceedsFileSize            = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открываем исходный файл
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	// Получаем информацию о файле
	fileInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}

	// Проверяем, что это не директория
	if fileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	if fileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}

	// Проверяем, что offset не превышает размер файла
	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// Открываем файл для записи
	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		err := toFile.Close()
		if err != nil {
			// IDE говорит, что нужно эту ошибку обработать.
			// Но в данном случае, если ошибка возникнет, то это будет критическая ошибка.
			// Тут нужно выкинуть панику или этого достаточно??
			fmt.Println("Error closing file:", err)
		}
	}()

	// Устанавливаем смещение в исходном файле
	if _, err = fromFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	// Копируем данные с учетом лимита
	if limit == 0 || limit > fileInfo.Size()-offset {
		limit = fileInfo.Size() - offset
	}

	// Копируем данные без прогресс бара:
	// - если размер копируемых данных меньше maxFileSizeWithoutProgressBar
	if limit < maxFileSizeWithoutProgressBar {
		_, err = io.CopyN(toFile, fromFile, limit)
		if err != nil && errors.Is(err, io.EOF) {
			return err
		}
		return nil
	}

	// Иначе: Копирование с визуализацией процесса копирования
	if err := copyWithProgres(fromFile, toFile, limit); err != nil {
		return err
	}

	return nil
}

func copyWithProgres(fromFile *os.File, toFile *os.File, limit int64) error {
	buf := make([]byte, 1024)
	var copied int64

	for copied < limit {
		bytesToRead := int64(len(buf))
		if limit-copied < bytesToRead {
			bytesToRead = limit - copied
		}

		n, err := fromFile.Read(buf[:bytesToRead])
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := toFile.Write(buf[:n]); err != nil {
			return err
		}

		copied += int64(n)
		printProgress(copied, limit)
	}
	return nil
}

func printProgress(copied, total int64) {
	percent := float64(copied) / float64(total) * 100
	fmt.Printf("\rProgress: %.2f%% (%d/%d bytes)", percent, copied, total)
	if copied == total {
		fmt.Println("\nCopy complete")
	}
	// time.Sleep(100 * time.Millisecond) // для плавного обновления прогресса
}
