# sqlite-parser

A SQLite database file parser written in Go that implements basic database inspection commands.

## Overview

This project parses SQLite database files and extracts information about the database structure. It provides a command-line interface for inspecting SQLite databases without requiring the full SQLite library.

## Features

The parser currently supports two commands:

### `.dbinfo` Command
Displays basic database information:
- Database page size
- Number of tables in the database

### `.tables` Command
Lists all table names in the database, excluding system tables (like `sqlite_sequence`).

## Usage

```bash
go run app/main.go <database-file-path> <command>
```

### Examples

```bash
# Display database information
go run app/main.go sample.db .dbinfo

# List all tables
go run app/main.go sample.db .tables
```

## Implementation Details

### Database File Structure
The parser reads and interprets SQLite database files according to the [SQLite file format](https://www.sqlite.org/fileformat.html):

1. **File Header** (first 100 bytes):
   - Magic string ("SQLite format 3\000")
   - Page size (bytes 16-17)
   - File format version numbers
   - Other database metadata

2. **Page Structure**:
   - Page type (byte 100): interior/leaf, table/index
   - Number of cells (table entries)
   - Cell pointer array starting at byte 108

3. **Table Information**:
   - Stored as records in the database pages
   - Each table record includes: schema name, table name, root page, SQL create statement

### Parsing Logic

The parser implements:
- **Header Reading**: Extracts the 100-byte file header
- **Page Size Calculation**: Reads 2-byte big-endian page size from header
- **Page Loading**: Reads the entire first page based on calculated size
- **Table Counting**: Parses the page header to get table count
- **Table Name Extraction**: Parses table metadata records to extract table names
- **Serial Type Decoding**: Converts SQLite varint serial types to byte lengths

### Key Functions

```go
openDBFile(filePath string) *os.File
getFileHeader(file *os.File) []byte
getPageSize(header []byte) (uint16, error)
getPage(pageSize uint16, file *os.File) ([]byte, error)
getNumberOfTables(buffer []byte) uint16
getTableInfoAddr(tables uint16, page []byte) []uint16
parseTableData(tblInfo []byte) (string, string, string, error)
getSerialTypeFromVarInt(b byte) (byte, error)
```

## Data Structures

The parser works with:
- **Page headers** (bytes 100-107): page type and cell count
- **Cell pointer arrays** (starting at byte 108): 2-byte offsets to table records
- **Table records**: Variable-length records containing:
  - Record length (varint)
  - Row ID (varint)
  - Header length (varint)
  - Column serial types (varints)
  - Column data (variable length based on serial types)

## Serial Types

SQLite uses variable-length integer (varint) encoding for serial types. The parser handles:
- Type 0-4: Integers of various sizes (0, 1, 2, 3, 4 bytes)
- Type 5: 6-byte integer
- Type 6-7: 8-byte integer/float
- Type 8-9: Constants 0 and 1
- Type 12+: Text/blob with calculated length

## Current State

The parser currently:
- Reads SQLite database file headers
- Extracts page size information (big-endian 2-byte value)
- Parses the first database page (root page)
- Counts tables from page metadata
- Extracts and displays table names
- Filters out system tables (sqlite_sequence)

**Limitations**:
- Only parses the first page (root page) of the database
- Does not handle multi-page databases completely
- No support for queries beyond the two basic commands
- Does not parse table schemas or data records
- Limited error handling

**Dependencies**: None - uses only Go standard library (os, bufio, encoding/binary, bytes, strings, errors)

## Output Format

### `.dbinfo` Output
```
database page size: 4096
number of tables: 3
```

### `.tables` Output
```
table1 table2 table3
```
(Space-separated list of table names)

## Error Handling

The program uses `log.Fatal()` for critical errors:
- Unable to open database file
- Invalid file header
- Page read failures
- Table parsing errors

## Requirements

- Go 1.16 or later (no specific version in go.mod, module name not set)
- Valid SQLite database file

## Notes

- This is a learning/educational project for understanding SQLite file format
- Not suitable for production use
- Does not modify or write to database files (read-only)
- Assumes standard SQLite file format (version 3)
