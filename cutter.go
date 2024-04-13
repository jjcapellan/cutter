package cutter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
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
