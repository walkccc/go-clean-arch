// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package repository

import (
	"context"
)

type Querier interface {
	CreateBook(ctx context.Context, arg CreateBookParams) (Book, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteBook(ctx context.Context, id int64) error
	GetBook(ctx context.Context, id int64) (Book, error)
	GetUser(ctx context.Context, username string) (User, error)
	ListBooks(ctx context.Context, arg ListBooksParams) ([]Book, error)
	UpdateBook(ctx context.Context, arg UpdateBookParams) (Book, error)
}

var _ Querier = (*Queries)(nil)
