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
		buffer := make([]byte, 1)
		_, err = databaseFile.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		var table uint16
		if buffer[0] == 10 || buffer[0] == 13 {
			buffer = make([]byte, 7)
			_, err = databaseFile.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}
			table = (uint16(buffer[2]) << 8) | uint16(buffer[3])
		}

		fmt.Printf("database page size: %v\n", pageSize)
		fmt.Printf("number of tables: %v\n", table)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
