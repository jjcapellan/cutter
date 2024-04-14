package cutter

import (
	"os"
	"reflect"
	"testing"
)

func TestCut(t *testing.T) {
	fileName := "./test_data/data.txt"
	destFolder := "test_data"
	var chunks uint32 = 2
	err := Cut(fileName, destFolder, chunks)
	if err != nil {
		t.Fatalf("Error cutting file %s: %s", fileName, err)
	}

	// Check p0
	buf1, _ := os.ReadFile("./test_data/data.txt.p0")
	buf2, _ := os.ReadFile("./test_data/model.p0")
	if !reflect.DeepEqual(buf1, buf2) {
		t.Fatalf("Bad format of chunk file (0)")
	}

	// Check p1
	buf1, _ = os.ReadFile("./test_data/data.txt.p1")
	buf2, _ = os.ReadFile("./test_data/model.p1")
	if !reflect.DeepEqual(buf1, buf2) {
		t.Fatalf("Bad format of chunk file (1)")
	}
}

func TestJoin(t *testing.T) {
	chunkName := "./test_data/model.p0"
	destFolder := "test_data/join"
	err := Join(chunkName, destFolder)
	if err != nil {
		t.Fatalf("Error joining file %s: %s", chunkName, err)
	}

	buf1, _ := os.ReadFile("./test_data/data.txt")
	buf2, _ := os.ReadFile("./test_data/join/model")

	if !reflect.DeepEqual(buf1, buf2) {
		t.Fatalf("Bad format of joined file")
	}
}
