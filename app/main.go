package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]
	page := getPageFromFile(databaseFilePath)
	switch command {
	case ".dbinfo":
		handleDBInfo(page)
	case ".tables":
		handleDotTablesCommand(page, handleDBInfo(page))
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}

func getPageFromFile(filePath string) []byte {
	var file = openDBFile(filePath)
	var header = getFileHeader(file)
	pageSizeInfoSlice, err := getPageSizeInfoSlice(header)
	if err != nil {
		log.Fatal(err)
	}
	pageSize, err := getPageSize(pageSizeInfoSlice)
	if err != nil {
		log.Fatal(err)
	}
	page, err := getPage(pageSize, file)
	if err != nil {
		log.Fatal(err)
	}
	return page
}
func handleDBInfo(page []byte) uint16 {
	numberOfTables := getNumberOfTables(page[100:108])
	fmt.Printf("database page size: %v\n", len(page))
	fmt.Printf("number of tables: %v\n", numberOfTables)
	return numberOfTables
}

func handleDotTablesCommand(page []byte, numberOfTables uint16) {
	tblAddrs := getTableInfoAddr(numberOfTables, page)
	tables := make([]string, 0)
	for _, el := range tblAddrs {
		_, tableName, _, err := parseTableData(page[el : el+uint16(page[el])+2])
		if err != nil {
			log.Fatal(err)
		}
		if tableName != "sqlite_sequence" {
			tables = append(tables, tableName)
		}
	}
	fmt.Println(strings.Join(tables, " "))
}

func openDBFile(filePath string) *os.File {
	databaseFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return databaseFile
}

func getFileHeader(databaseFile *os.File) []byte {
	header := make([]byte, 100)

	_, err := databaseFile.Read(header)
	if err != nil {
		log.Fatal(err)
	}
	return header
}
func getPageSizeInfoSlice(header []byte) ([]byte, error) {
	if len(header) == 100 {
		return header[16:18], nil
	} else {
		return []byte{}, errors.New("invalid header")
	}
}

func getPageSize(header []byte) (uint16, error) {
	var pageSize uint16
	err := binary.Read(bytes.NewReader(header), binary.BigEndian, &pageSize)
	if err != nil {
		return 0, err
	}
	return pageSize, nil
}

func getPage(pageSize uint16, file *os.File) ([]byte, error) {
	var page = make([]byte, pageSize)
	_, err := file.ReadAt(page, 0)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func getPageType(page []byte) byte {
	return page[100]
}
func getNumberOfTables(buffer []byte) uint16 {
	var table uint16
	if buffer[0] == 10 || buffer[0] == 13 {
		table = (uint16(buffer[3]) << 8) | uint16(buffer[4])
	}
	return table
}

func getTableInfoAddr(tables uint16, page []byte) []uint16 {
	if page[100] == 10 || page[100] == 13 {
		tableAddrsSlice := page[108 : 108+(2*tables)]
		var tblAddrs []uint16
		var tblAddr uint16
		for i := uint16(0); i < tables; i++ {
			_ = binary.Read(bytes.NewReader(tableAddrsSlice[2*i:(2*i)+2]), binary.BigEndian, &tblAddr)
			tblAddrs = append(tblAddrs, tblAddr)
		}
		//slices.Sort(tblAddrs)
		return tblAddrs
	}
	return nil
}

func parseTableData(tblInfo []byte) (string, string, string, error) {
	_ = tblInfo[0] //total Length
	_ = tblInfo[1] //Row Id
	tblHeaderLength := tblInfo[2]
	tblHeader := tblInfo[3 : 3+tblHeaderLength-1]
	columnDetails := make([]byte, 0)
	for i := 0; i < len(tblHeader); i++ {
		val, err := getSerialTypeFromVarInt(tblHeader[i])
		if err != nil {
			return "", "", "", err
		}
		columnDetails = append(columnDetails, val)
	}
	schemaName := string(tblInfo[int(3+tblHeaderLength-1+columnDetails[0]):int(3+tblHeaderLength-1+columnDetails[0]+columnDetails[1])])
	tableName := string(tblInfo[int(3+tblHeaderLength-1+columnDetails[0]+columnDetails[1]):int(3+tblHeaderLength-1+columnDetails[0]+columnDetails[1]+columnDetails[2])])
	createTableQuery := string(tblInfo[int(3+tblHeaderLength-1+columnDetails[0]+columnDetails[1]+columnDetails[2]):])
	return schemaName, tableName, createTableQuery, nil
	//fmt.Println(columnDetails)
}

func getSerialTypeFromVarInt(b byte) (byte, error) {
	switch b {
	case 0, 1, 2, 3, 4:
		return b, nil
	case 5:
		return 6, nil
	case 6, 7:
		return 8, nil
	case 8, 9:
		return 0, nil
	default:
		if b >= 12 {
			if b%2 == 0 {
				return (b - 12) / 2, nil
			} else {
				return (b - 13) / 2, nil
			}
		}
		return 0, errors.New("invalid byte")
	}
}
