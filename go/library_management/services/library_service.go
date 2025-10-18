package services

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Ruth652/playground/go/library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)

	AddMember(member models.Member)
}

type Library struct {
	books   map[int]models.Book
	members map[int]*models.Member
}

func BookFromArgs(id int, title, author string) models.Book {
	return models.Book{ID: id, Title: title, Author: author, Status: "Available"}
}

func MemberFromArgs(id int, name string) models.Member {
	return models.Member{ID: id, Name: name, BorrowedBooks: []models.Book{}}
}

func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]*models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	book.Status = "Available"
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) error {
	b, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	if b.Status == "Borrowed" {
		return errors.New("cannot remove a borrowed book")
	}
	delete(l.books, bookID)
	return nil
}

func (l *Library) AddMember(member models.Member) {
	if member.BorrowedBooks == nil {
		member.BorrowedBooks = []models.Book{}
	}
	l.members[member.ID] = &member
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}
	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	l.books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	return nil
}
func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}
	found := -1
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			found = i
			break
		}
	}
	if found == -1 {
		return errors.New("this member did not borrow the specified book")
	}

	member.BorrowedBooks = append(member.BorrowedBooks[:found], member.BorrowedBooks[found+1:]...)
	book.Status = "Available"
	l.books[bookID] = book

	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	out := []models.Book{}
	for _, b := range l.books {
		if b.Status == "Available" {
			out = append(out, b)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}

func (l *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	member, ok := l.members[memberID]
	if !ok {
		return nil, errors.New("member not found")
	}
	return member.BorrowedBooks, nil
}

func (l *Library) PrintStatus() {
	fmt.Println("---- Library Books ----")
	keys := []int{}
	for id := range l.books {
		keys = append(keys, id)
	}
	sort.Ints(keys)
	for _, id := range keys {
		b := l.books[id]
		fmt.Printf("ID:%d | %s — %s | %s\n", b.ID, b.Title, b.Author, b.Status)
	}
	fmt.Println("---- Members ----")
	mkeys := []int{}
	for id := range l.members {
		mkeys = append(mkeys, id)
	}
	sort.Ints(mkeys)
	for _, id := range mkeys {
		m := l.members[id]
		fmt.Printf("Member %d: %s — Borrowed: %d books\n", m.ID, m.Name, len(m.BorrowedBooks))
	}
}
