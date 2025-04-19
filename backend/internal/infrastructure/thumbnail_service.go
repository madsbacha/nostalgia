package infrastructure

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/buckket/go-blurhash"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"nostalgia/internal/core/port"
	"os"
	"os/exec"
	"path/filepath"
)

type ThumbnailService struct {
	ffmpegPath         string
	logger             *logrus.Entry
	temporaryDirectory string
}

func NewThumbnailService(ffmpegPath, temporaryDir string, logger *logrus.Entry) ThumbnailService {
	if _, err := os.Stat(temporaryDir); os.IsNotExist(err) {
		logger.Fatal("temporary directory does not exist")
	}

	return ThumbnailService{
		ffmpegPath:         ffmpegPath,
		temporaryDirectory: temporaryDir,
		logger:             logger,
	}
}

func (s ThumbnailService) getFFmpegArgs(inputFilePath string, outputFilePath string) []string {
	return []string{
		"-i",
		inputFilePath,
		"-ss",
		"00:00:01.000",
		"-vframes:v",
		"1",
		"-pred",
		"mixed",
		"-vf",
		"scale=320:-1",
		outputFilePath,
	}
}

func (s ThumbnailService) CreateTemporaryOutputFile(extension string) string {
	if extension[0] != '.' {
		s.logger.Fatalf("Invalid extension, missing dot: %s", extension)
	}
	id, err := uuid.NewRandom()
	if err != nil {
		s.logger.Fatalf("Failed to generate UUID: %s", err)
	}
	filename := fmt.Sprintf("thumbnail.%s%s", id, extension)
	return filepath.Join(s.temporaryDirectory, filename)
}

func (s ThumbnailService) GenerateFromVideo(ctx context.Context, filePath string) (*port.GeneratedThumbnail, error) {
	extension := ".png"
	mimeType := "image/png"
	outputFilePath := s.CreateTemporaryOutputFile(extension)

	args := s.getFFmpegArgs(filePath, outputFilePath)
	cmd := exec.Command(s.ffmpegPath, args...)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Printf("stdout of ffmpeg: %s\n", stdout)
		return nil, err
	}

	file, err := os.Open(outputFilePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			s.logger.Println(err.Error())
		}
	}(file)

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, file); err != nil {
		return nil, err
	}
	if err := os.Remove(outputFilePath); err != nil {
		return nil, err
	}

	bh, err := s.GenerateBlurHash(bytes.NewReader(buffer.Bytes()), extension)
	if err != nil {
		return nil, err
	}

	return &port.GeneratedThumbnail{
		File:      bytes.NewReader(buffer.Bytes()),
		MimeType:  mimeType,
		Extension: extension,
		Blurhash:  bh,
	}, nil
}

func (s ThumbnailService) GenerateBlurHash(file io.Reader, extension string) (string, error) {
	if extension[0] != '.' {
		s.logger.Fatalf("Invalid extension, missing dot: %s", extension)
	}

	var decoder func(io.Reader) (image.Image, error)
	if extension == ".png" {
		decoder = png.Decode
	} else if extension == ".jpg" || extension == ".jpeg" {
		decoder = jpeg.Decode
	} else {
		return "", errors.New("unsupported image format")
	}
	img, err := decoder(file)
	if err != nil {
		return "", err
	}
	return blurhash.Encode(4, 3, img)
}
