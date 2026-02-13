package models

type Book struct {
	ID         int
	Title      string
	Author     string
	Status     string // Available | Reserved | Borrowed
	ReservedBy int    // 0 if none
}
