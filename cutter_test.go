package cutter

import (
	"os"
	"reflect"
	"testing"
)

// (CUTTER(6 bytes) + PART(4 bytes) + PARTS(4 bytes) + VERSION(2 bytes)) <-- littleendian
var header0 = []byte{0x43, 0x55, 0x54, 0x54, 0x45, 0x52, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00}

const fileContent = "abcdefghijklmnopqrstuwxyz"

func checkHeader0(fileName string, t *testing.T) {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", fileName, err)
	}
	defer file.Close()

	buf := make([]byte, len(header0))
	_, err = file.Read(buf)
	if err != nil {
		t.Fatalf("Error reading file %s: %s", fileName, err)
	}

	if !reflect.DeepEqual(header0, buf) {
		t.Fatalf("Wrong header format in file %s: %#x\nriginal:%#x", fileName, buf, header0)
	}
}

func checkContent(file1Name string, file2Name string, t *testing.T) {
	file1, err := os.Open(file1Name)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", file1Name, err)
	}
	defer file1.Close()

	file2, err := os.Open(file2Name)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", file2Name, err)
	}
	defer file2.Close()

	info1, _ := file1.Stat()
	info2, _ := file2.Stat()
	file1Size := info1.Size()
	file2Size := info2.Size()

	headerSize := int64(len(header0))
	content1Size := file1Size - headerSize
	content2Size := file2Size - headerSize

	buf1 := make([]byte, content1Size)
	buf2 := make([]byte, content2Size)

	file1.ReadAt(buf1, headerSize)
	file2.ReadAt(buf2, headerSize)

	result := string(buf1) + string(buf2)

	if fileContent != result {
		t.Fatalf("Content is corrupted: %s", result)
	}

}

func TestCut(t *testing.T) {
	fileName := "./test_data/data.txt"
	destFolder := "test_data"
	var chunks uint32 = 2
	err := Cut(fileName, destFolder, chunks)
	if err != nil {
		t.Fatalf("Error cutting file %s: %s", fileName, err)
	}

	// File0 header
	checkHeader0(fileName+".p0", t)

	// Content
	checkContent(fileName+".p0", fileName+".p1", t)
}

func TestJoin(t *testing.T) {
	chunkName := "./test_data/data.txt.p0"
	destFolder := "test_data/join"
	err := Join(chunkName, destFolder)
	if err != nil {
		t.Fatalf("Error joining file %s: %s", chunkName, err)
	}

	buf1, _ := os.ReadFile("./test_data/data.txt")
	buf2, _ := os.ReadFile("./test_data/join/data.txt")

	if !reflect.DeepEqual(buf1, buf2) {
		t.Fatalf("Bad format of joined file")
	}
}
