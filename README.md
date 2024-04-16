# CUTTER
Cutter is a lightweight Go package designed to split and merge files effortlessly.

## How to use
```go
package main

import (
    "fmt"
    "github.com/jjcapellan/cutter"
)

func main() {
    // Define the source file path and the folder where chunks will be saved
    filePath := "large_file.dat"
    folder := "chunks"

    // Split the file into 5 chunks
    err := cutter.Cut(filePath, folder, 5)
    if err != nil {
        fmt.Println("Error splitting the file:", err)
        return
    }
    fmt.Println("File splitted successfully into chunks.")

    // Join the chunks into a single file
    chunkPath := "chunks/large_file.dat.p0" // Using the first chunk as reference
    destFolder := "reconstructed_file" // You could use "" for current folder
    err = cutter.Join(chunkPath, destFolder)
    if err != nil {
        fmt.Println("Error joining file chunks:", err)
        return
    }
    fmt.Println("File chunks joined successfully.")
}
```