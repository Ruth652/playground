# Library Management System (Console)

A simple console-based library management system in Go demonstrating:

- Structs, slices, and maps
- Interfaces and method receivers
- Basic error handling
- Console interaction for CRUD operations

## Folder structure

    library_management/
    ├── main.go # Program entry point
    ├── controllers/ # Console interface logic
    ├── models/ # Book and Member structs
    ├── services/ # Library service implementing LibraryManager

1. Navigate into the project folder:

```bash
cd go/library_management
```

2. Ensure the Go module is initialized (if first time):

```bash
go mod init github.com/Ruth652/playground/go/library_management
go mod tidy
```

3. Run the program

```bash
go run ./...
```
