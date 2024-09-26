package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}

		header := make([]byte, 100)

		_, err = databaseFile.Read(header)
		if err != nil {
			log.Fatal(err)
		}

		var pageSize uint16
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}
		buffer := make([]byte, pageSize)
		_, err = databaseFile.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		searchBuffer := []byte{67, 82, 69, 65, 84, 69}

		var table int = 0
		for i := 0; i < int(pageSize)-len(searchBuffer); i++ {
			slice := buffer[i : i+len(searchBuffer)]
			if bytes.Equal(slice, searchBuffer) {
				table++
			}
		}
		fmt.Printf("database page size: %v\n", pageSize)
		fmt.Printf("number of tables: %v\n", table)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
