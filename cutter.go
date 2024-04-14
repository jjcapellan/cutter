package cutter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const VERSION uint16 = 1
const FILE_ID = "CUTTER"

type Header struct {
	Id      string // "CUTTER"
	Chunk   uint32
	Chunks  uint32
	Version uint16
}

func Cut(filePath string, folder string, chunks uint32) error {
	if chunks < 2 {
		return errors.New("number of chunks must be greater than 1")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	chunkSize, err := getChunkSize(file, chunks)
	if err != nil {
		return err
	}

	buf := make([]byte, chunkSize)
	fileName := filepath.Base(filePath)
	chunkBaseName := filepath.Join(folder, fileName)

	for i := 0; i < int(chunks); i++ {
		header := Header{FILE_ID, uint32(i), chunks, VERSION}
		chunkName := fmt.Sprintf("%s.p%d", chunkBaseName, i)

		err = writeChunk(file, buf, chunkName, header)
		if err != nil {
			return err
		}
	}

	return nil
}

func Join(chunkPath string, destFolder string) error {

	// Check file name of file 0, should be *.p0
	if !strings.HasSuffix(chunkPath, ".p0") {
		return errors.New("not file zero. Please, open the file with extension .p0")
	}

	if !fileExist(chunkPath) {
		return errors.New("no such file")
	}

	destName := filepath.Join(destFolder, strings.TrimSuffix(filepath.Base(chunkPath), ".p0"))
	chunkBaseName := strings.TrimSuffix(chunkPath, "0")
	chunk := 0
	isLast := false

	destFile, err := os.OpenFile(destName, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for !isLast {
		chunkPath = chunkBaseName + strconv.Itoa(chunk)
		isLast, err = copyChunk(destFile, chunkPath)
		if err != nil {
			return err
		}
		chunk++
	}

	return nil
}

func fileExist(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}

func copyChunk(destFile *os.File, chunkPath string) (bool, error) {
	headerSize := 16
	isLast := false
	bufferSize := 5 * 1024 * 1024 // 5Mb

	buf := make([]byte, headerSize)
	chunkFile, err := os.Open(chunkPath)
	if err != nil {
		return isLast, err
	}
	defer chunkFile.Close()
	chunkFile.Read(buf)

	header, err := readHeader(buf)
	if err != nil {
		return isLast, err
	}
	isLast = (header.Chunk >= header.Chunks-1)

	buf = make([]byte, bufferSize)

	for {
		n, err := chunkFile.Read(buf)
		if err != nil && err != io.EOF {
			return isLast, err
		}
		if n == 0 {
			break
		}
		_, err = destFile.Write(buf[:n])
		if err != nil {
			return isLast, err
		}
	}

	return isLast, nil
}

func readHeader(header []byte) (Header, error) {
	h := Header{}
	h.Id = string(header[:6])
	if h.Id != FILE_ID {
		return h, errors.New("unknown file type")
	}
	h.Chunk = binary.LittleEndian.Uint32(header[6:10])
	h.Chunks = binary.LittleEndian.Uint32(header[10:14])
	h.Version = binary.LittleEndian.Uint16(header[14:16])

	return h, nil
}

func getChunkSize(file *os.File, chunks uint32) (int64, error) {
	info, _ := file.Stat()
	fileSize := info.Size()
	size := int64(math.Ceil(float64(fileSize) / float64(chunks)))
	if size < 1 {
		return 0, errors.New("chunk size is less than 1 byte, try with less chunks")
	}
	return size, nil
}

func writeChunk(source *os.File, buf []byte, chunkName string, h Header) error {
	chunk, err := os.OpenFile(chunkName, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer chunk.Close()

	header := getHeaderBuf(h.Chunk, h.Chunks)
	_, err = chunk.Write(header)
	if err != nil {
		return err
	}

	n, err := source.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}

	chunk.Write(buf[:n])
	if err != nil {
		return err
	}
	return nil
}

func getHeaderBuf(chunk uint32, chunks uint32) []byte {
	size := 6 + 4 + 4 + 2
	h := make([]byte, size)
	copy(h[:6], FILE_ID)
	binary.LittleEndian.PutUint32(h[6:10], chunk)
	binary.LittleEndian.PutUint32(h[10:14], chunks)
	binary.LittleEndian.PutUint16(h[14:16], VERSION)
	return h
}
