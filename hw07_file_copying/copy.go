package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	copyBufferSize           int64 = 32 * 1024
	ErrUnsupportedFile             = errors.New("unsupported file")
	ErrOffsetExceedsFileSize       = errors.New("offset exceeds file size")
	ErrFromIsDirectory             = errors.New("directory copying is not supported")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromPath, _ = filepath.Abs(fromPath)
	toPath, _ = filepath.Abs(toPath)
	var originalToPath string

	// fromPath и toPath указывают на один и тот же файл.
	if isTheSameFile(fromPath, toPath) {
		originalToPath = toPath
		fileHash, err := calculateSHA256(fromPath)
		if err != nil {
			return ErrUnsupportedFile
		}
		// Подмена toPath на временный файл
		toPath = fmt.Sprintf("%s/%s_%d.tmp", filepath.Dir(toPath), fileHash, time.Now().UnixNano())
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

	if !fromFileInfo.Mode().IsRegular() {
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
	fmt.Printf("Copy %s to %s\n", fromPath, toPath)
	fmt.Printf("Offset %d Limit %d\n", offset, limit)
	if err := copyWithProgress(fromFile, toFile, limit); err != nil {
		os.Remove(toPath)
		return err
	}

	// Обработка копирования через промежуточный файл
	if originalToPath != "" {
		// Удаляем целевой файл
		if err := os.Remove(originalToPath); err != nil {
			// Удалить временный файл
			os.Remove(toPath)
			return err
		}

		if err := os.Rename(toPath, originalToPath); err != nil {
			// Удалить временный файл
			os.Remove(toPath)
			return err
		}
	}

	return nil
}

func copyWithProgress(fromFile *os.File, toFile *os.File, limit int64) error {
	var copied int64
	// Для срабатывания первого обновления прогресса добавляем секунду к текущему времени
	lastUpdate := time.Now().Add(time.Second)

	for copied < limit {
		nSize := min(copyBufferSize, limit-copied)

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

func calculateSHA256(filename string) (string, error) {
	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Создаем новый хеш SHA256
	hash := sha256.New()

	// Копируем содержимое файла в хеш
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Получаем SHA256 хеш в виде среза байт
	hashInBytes := hash.Sum(nil)

	// Преобразуем байты в строку в формате hex
	return fmt.Sprintf("%x", hashInBytes), nil
}
