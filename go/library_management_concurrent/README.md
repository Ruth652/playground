# Library Management System (Console)

A console-based Go library system demonstrating:

- Books & members management
- Borrowing, returning, and reserving books
- Reservation auto-expiry after 5 seconds
- Thread-safe operations using goroutines and channels

## Folder structure

    library_management/
    ├── main.go # Program entry point
    ├── controllers/ # Console interface logic
    ├── models/ # Book and Member structs
    ├── services/ # Library service implementing LibraryManager
    ├── concurrency/

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
