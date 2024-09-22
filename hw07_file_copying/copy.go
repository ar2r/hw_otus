package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	copyBufferSize                  int64 = 32 * 1024
	ErrUnsupportedFile                    = errors.New("unsupported file")
	ErrOffsetExceedsFileSize              = errors.New("offset exceeds file size")
	ErrFromAndToPointsToTheSameFile       = errors.New("from and to points to the same file")
	ErrFromIsDirectory                    = errors.New("directory copying is not supported")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if isTheSameFile(fromPath, toPath) {
		return ErrFromAndToPointsToTheSameFile
	}

	// Открываем исходный файл
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	// Получаем информацию о файле
	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}

	// Проверяем, что это не директория
	if fromFileInfo.IsDir() {
		return ErrFromIsDirectory
	}

	if fromFileInfo.Size() == 0 && !fromFileInfo.Mode().IsRegular() {
		// Проверка для файлов, которые не поддерживают смещение (например /dev/urandom)
		// Этот способ самый компактный
		return ErrUnsupportedFile
	}

	// Проверяем, что offset не превышает размер файла
	if offset > fromFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// Открываем файл для записи
	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	// Устанавливаем смещение в исходном файле
	if _, err = fromFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	// Копируем данные с учетом лимита
	if limit == 0 || limit > fromFileInfo.Size()-offset {
		limit = fromFileInfo.Size() - offset
	}

	// Копирование с визуализацией процесса копирования
	if err := copyWithProgress(fromFile, toFile, limit); err != nil {
		return err
	}

	return nil
}

func copyWithProgress(fromFile *os.File, toFile *os.File, limit int64) error {
	var copied int64
	// Для срабатывания первого обновления прогресса добавляем секунду к текущему времени
	lastUpdate := time.Now().Add(time.Second)

	for copied < limit {
		nSize := copyBufferSize
		if limit-copied < copyBufferSize {
			nSize = limit - copied
		}

		n, err := io.CopyN(toFile, fromFile, nSize)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		copied += n
		if time.Since(lastUpdate) >= time.Second {
			printProgress(copied, limit)
			lastUpdate = time.Now()
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}
	printProgress(copied, limit) // Final update
	return nil
}

func printProgress(copied, total int64) {
	percent := float64(copied) / float64(total) * 100
	fmt.Printf("\rProgress: %.2f%% (%d/%d bytes)", percent, copied, total)
	if copied == total {
		// Пустая строка, т.к. прогресс бар закончился.
		fmt.Print("\n")
	}
}

func isTheSameFile(firstFilePath, secondFilePath string) bool {
	fromFileInfo, err := os.Stat(firstFilePath)
	if err != nil {
		return false
	}

	toFileInfo, err := os.Stat(secondFilePath)
	if err != nil {
		return false
	}

	return os.SameFile(fromFileInfo, toFileInfo)
}
