package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Ruth652/playground/go/library_management/models"
	"github.com/Ruth652/playground/go/library_management/services"
)

type Controller struct {
	lib *services.Library
	in  *bufio.Reader
}

func NewController(lib *services.Library) *Controller {
	return &Controller{
		lib: lib,
		in:  bufio.NewReader(os.Stdin),
	}
}

func (c *Controller) Run() {
	for {
		c.printMenu()
		choice := c.readLine("Enter choice: ")
		switch choice {
		case "1":
			c.addBook()
		case "2":
			c.removeBook()
		case "3":
			c.addMember()
		case "4":
			c.borrowBook()
		case "5":
			c.returnBook()
		case "6":
			c.listAvailableBooks()
		case "7":
			c.listBorrowedBooks()
		case "8":
			c.lib.PrintStatus()
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown option. Try again.")
		}
		fmt.Println()
	}
}

func (c *Controller) printMenu() {
	fmt.Println("=== Library Management ===")
	fmt.Println("1) Add Book")
	fmt.Println("2) Remove Book")
	fmt.Println("3) Add Member")
	fmt.Println("4) Borrow Book")
	fmt.Println("5) Return Book")
	fmt.Println("6) List Available Books")
	fmt.Println("7) List Borrowed Books by Member")
	fmt.Println("8) Print Library Status (debug)")
	fmt.Println("0) Exit")
}

func (c *Controller) readLine(prompt string) string {
	fmt.Print(prompt)
	text, _ := c.in.ReadString('\n')
	return strings.TrimSpace(text)
}

func (c *Controller) addBook() {
	idStr := c.readLine("Book ID (int): ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	title := c.readLine("Title: ")
	author := c.readLine("Author: ")
	book := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: "Available",
	}
	c.lib.AddBook(book)
	fmt.Println("Book added.")
}

func (c *Controller) removeBook() {
	idStr := c.readLine("Book ID to remove: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	err = c.lib.RemoveBook(id)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book removed.")
	}
}

func (c *Controller) addMember() {
	idStr := c.readLine("Member ID (int): ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	name := c.readLine("Name: ")
	member := models.Member{
		ID:            id,
		Name:          name,
		BorrowedBooks: []models.Book{},
	}
	c.lib.AddMember(member)
	fmt.Println("Member added.")
}

func (c *Controller) borrowBook() {
	bookStr := c.readLine("Book ID to borrow: ")
	bookID, err := strconv.Atoi(bookStr)
	if err != nil {
		fmt.Println("Invalid book ID")
		return
	}
	memberStr := c.readLine("Member ID: ")
	memberID, err := strconv.Atoi(memberStr)
	if err != nil {
		fmt.Println("Invalid member ID")
		return
	}
	if err := c.lib.BorrowBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book borrowed.")
	}
}

func (c *Controller) returnBook() {
	bookStr := c.readLine("Book ID to return: ")
	bookID, err := strconv.Atoi(bookStr)
	if err != nil {
		fmt.Println("Invalid book ID")
		return
	}
	memberStr := c.readLine("Member ID: ")
	memberID, err := strconv.Atoi(memberStr)
	if err != nil {
		fmt.Println("Invalid member ID")
		return
	}
	if err := c.lib.ReturnBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book returned.")
	}
}

func (c *Controller) listAvailableBooks() {
	books := c.lib.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}
	fmt.Println("Available books:")
	for _, b := range books {
		fmt.Printf("ID:%d | %s — %s\n", b.ID, b.Title, b.Author)
	}
}

func (c *Controller) listBorrowedBooks() {
	memberStr := c.readLine("Member ID: ")
	memberID, err := strconv.Atoi(memberStr)
	if err != nil {
		fmt.Println("Invalid member ID")
		return
	}
	books, err := c.lib.ListBorrowedBooks(memberID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(books) == 0 {
		fmt.Println("No borrowed books for this member.")
		return
	}
	fmt.Printf("Borrowed books for member %d:\n", memberID)
	for _, b := range books {
		fmt.Printf("ID:%d | %s — %s\n", b.ID, b.Title, b.Author)
	}
}
