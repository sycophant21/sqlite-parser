# sqlite-parser

> Parsing SQLite so you don't have to trust the library.

A from-scratch SQLite binary file format parser written in Go — no external SQLite library, just raw bytes.

## What it does

Reads a `.db` file directly and parses the SQLite page format by hand:
- Extracts the 100-byte file header to determine page size
- Walks the B-tree page structure to locate table cell addresses
- Decodes variable-length integers and record headers to extract schema data

Supported commands:

| Command    | Description                              |
|------------|------------------------------------------|
| `.dbinfo`  | Print page size and number of tables     |
| `.tables`  | Print names of all user-defined tables   |

## How it works

SQLite files are split into fixed-size pages. The first page begins with a 100-byte header that encodes the page size at bytes 16-17 (big-endian). After the header, interior/leaf table B-tree pages store a cell count and an array of 2-byte cell offsets. Each cell is a packed record containing variable-length integers (varints) for the header length and serial types, followed by the actual column data.

`sqlite-parser` decodes all of this manually using `encoding/binary` and a custom varint decoder.

## Tech stack

- **Go** (standard library only)

## Getting started

### Prerequisites

- Go 1.18+

### Run

```bash
go run app/main.go <path-to-database.db> .dbinfo
go run app/main.go <path-to-database.db> .tables
```

### Build

```bash
go build -o sqlite-parser app/main.go
./sqlite-parser sample.db .tables
```

## Example

```
$ go run app/main.go sample.db .dbinfo
database page size: 4096
number of tables: 3

$ go run app/main.go sample.db .tables
apples oranges employees
```

## Project structure

```
app/
└── main.go   # entry point + all parsing logic
```
