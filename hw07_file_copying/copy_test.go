package main

import (
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func TestCopy(t *testing.T) {
	// Временный файл для одновременного чтения и записи
	randomInputOutputName := getRandomFileName()
	file, _ := os.Create(randomInputOutputName)
	file.Close()

	tests := []struct {
		title            string
		offset           int64
		limit            int64
		inputFilePath    string
		outputFilePath   string
		expectedFilePath string
		expectedError    error
	}{
		{
			title:            "Весь файл без ограничений",
			offset:           0,
			limit:            0,
			inputFilePath:    "testdata/input.txt",
			expectedFilePath: "testdata/input.txt",
		},
		{
			title:            "Offset 100",
			offset:           100,
			limit:            0,
			inputFilePath:    "testdata/input.txt",
			expectedFilePath: "testdata/out_offset100_limit0.txt",
		},
		{
			title:            "Limit 1000",
			offset:           0,
			limit:            1000,
			inputFilePath:    "testdata/input.txt",
			expectedFilePath: "testdata/out_offset0_limit1000.txt",
		},
		{
			title:            "Offset 100 Limit 1000",
			offset:           100,
			limit:            1000,
			inputFilePath:    "testdata/input.txt",
			expectedFilePath: "testdata/out_offset100_limit1000.txt",
		},
		{
			title:            "Limit превышает размер файла",
			offset:           0,
			limit:            math.MaxInt64,
			inputFilePath:    "testdata/input.txt",
			expectedFilePath: "testdata/input.txt",
		},
		{
			title:         "Ошибка чтения файла неизвестной длины (не поддерживает смещение)",
			offset:        0,
			limit:         0,
			inputFilePath: "/dev/urandom",
			expectedError: ErrUnsupportedFile,
		},
		{
			title:          "Ошибка на чтение и запись указан один файл",
			offset:         0,
			limit:          0,
			inputFilePath:  randomInputOutputName,
			outputFilePath: randomInputOutputName,
			expectedError:  ErrFromAndToPointsToTheSameFile,
		},
		{
			title:         "Ошибка Offset превышает размер файла",
			offset:        math.MaxInt64,
			limit:         0,
			inputFilePath: "testdata/input.txt",
			expectedError: ErrOffsetExceedsFileSize,
		},
		{
			title:         "Ошибка чтения директории",
			offset:        0,
			limit:         0,
			inputFilePath: "/tmp",
			expectedError: ErrFromIsDirectory,
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			outputFilePath := ""
			if tc.outputFilePath != "" {
				outputFilePath = tc.outputFilePath
			} else {
				outputFilePath = getRandomFileName()
			}

			defer func() {
				if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
					// Файл не существует и удалять его не нужно
					return
				}

				if err := os.Remove(outputFilePath); err != nil {
					t.Errorf("Файл с результатом копирования найден, но удалить не удалось: %v", err)
				}
			}()

			copyErr := Copy(tc.inputFilePath, outputFilePath, tc.offset, tc.limit)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, copyErr)
			} else if tc.expectedFilePath == "" {
				if copyErr != nil {
					t.Errorf("Функция Copy вернула непредвиденную ошибку: %v", copyErr)
				}
				require.FileExists(t, outputFilePath)
				resultData, _ := os.ReadFile(outputFilePath)
				expectedData, _ := os.ReadFile(tc.expectedFilePath)
				assert.Equal(t, expectedData, resultData)
			}
		})
	}
}

func getRandomFileName() string {
	return os.TempDir() + "/copy_result_" + RandString(10)
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
