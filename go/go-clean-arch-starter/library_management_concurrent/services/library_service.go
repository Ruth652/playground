package services

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Ruth652/playground/go/library_management/models"
)

type ReservationRequest struct {
	BookID   int
	MemberID int
	Result   chan error
}

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
	AddMember(member models.Member)
	ReserveBook(bookID int, memberID int) error
}

type Library struct {
	books           map[int]models.Book
	members         map[int]*models.Member
	mu              sync.Mutex
	ReservationChan chan ReservationRequest
}

func NewLibrary() *Library {
	return &Library{
		books:           make(map[int]models.Book),
		members:         make(map[int]*models.Member),
		ReservationChan: make(chan ReservationRequest, 10),
	}
}

func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()

	book.Status = "Available"
	book.ReservedBy = 0
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	if b.Status == "Borrowed" {
		return errors.New("cannot remove borrowed book")
	}
	delete(l.books, bookID)
	return nil
}

func (l *Library) AddMember(member models.Member) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.members[member.ID] = &member
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}

	if book.Status == "Reserved" && book.ReservedBy != memberID {
		return errors.New("book reserved by another member")
	}

	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	book.ReservedBy = 0
	l.books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)

	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	index := -1
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			index = i
			break
		}
	}
	if index == -1 {
		return errors.New("book not borrowed by this member")
	}

	member.BorrowedBooks = append(member.BorrowedBooks[:index], member.BorrowedBooks[index+1:]...)

	book.Status = "Available"
	book.ReservedBy = 0
	l.books[bookID] = book

	return nil
}

func (l *Library) ReserveBook(bookID int, memberID int) error {
	resultChan := make(chan error)

	req := ReservationRequest{
		BookID:   bookID,
		MemberID: memberID,
		Result:   resultChan,
	}

	l.ReservationChan <- req

	return <-resultChan
}

func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	var out []models.Book
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
	l.mu.Lock()
	defer l.mu.Unlock()

	member, ok := l.members[memberID]
	if !ok {
		return nil, errors.New("member not found")
	}
	return member.BorrowedBooks, nil
}

func (l *Library) PrintStatus() {
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Println("---- Books ----")
	for _, b := range l.books {
		fmt.Printf("ID:%d | %s | %s | %s\n",
			b.ID, b.Title, b.Author, b.Status)
	}
}
func (l *Library) ProcessReservation(req ReservationRequest) {

	l.mu.Lock()

	book, ok := l.books[req.BookID]
	if !ok {
		l.mu.Unlock()
		req.Result <- fmt.Errorf("book not found")
		return
	}

	if book.Status == "Borrowed" {
		l.mu.Unlock()
		req.Result <- fmt.Errorf("book already borrowed")
		return
	}

	if book.Status == "Reserved" {
		l.mu.Unlock()
		req.Result <- fmt.Errorf("book already reserved")
		return
	}

	book.Status = "Reserved"
	book.ReservedBy = req.MemberID
	l.books[req.BookID] = book

	l.mu.Unlock()

	req.Result <- nil

	// Auto-expire reservation
	go func(bookID int, memberID int) {
		time.Sleep(5 * time.Second)

		l.mu.Lock()
		defer l.mu.Unlock()

		b := l.books[bookID]

		if b.Status == "Reserved" && b.ReservedBy == memberID {
			b.Status = "Available"
			b.ReservedBy = 0
			l.books[bookID] = b
			fmt.Println("Reservation expired for book:", bookID)
		}
	}(req.BookID, req.MemberID)
}
