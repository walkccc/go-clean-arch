// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: book.sql

package repository

import (
	"context"
)

const createBook = `-- name: CreateBook :one
INSERT INTO books (owner, name, language)
VALUES ($1, $2, $3)
RETURNING id, owner, name, language, created_at
`

type CreateBookParams struct {
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	Language string `json:"language"`
}

func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) (Book, error) {
	row := q.db.QueryRowContext(ctx, createBook, arg.Owner, arg.Name, arg.Language)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.Language,
		&i.CreatedAt,
	)
	return i, err
}

const deleteBook = `-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1
`

func (q *Queries) DeleteBook(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteBook, id)
	return err
}

const getBook = `-- name: GetBook :one
SELECT id, owner, name, language, created_at
FROM books
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetBook(ctx context.Context, id int64) (Book, error) {
	row := q.db.QueryRowContext(ctx, getBook, id)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.Language,
		&i.CreatedAt,
	)
	return i, err
}

const listBooks = `-- name: ListBooks :many
SELECT id, owner, name, language, created_at
FROM books
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListBooksParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListBooks(ctx context.Context, arg ListBooksParams) ([]Book, error) {
	rows, err := q.db.QueryContext(ctx, listBooks, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Book
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Name,
			&i.Language,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBook = `-- name: UpdateBook :one
UPDATE books
SET name = $2
WHERE id = $1
RETURNING id, owner, name, language, created_at
`

type UpdateBookParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateBook(ctx context.Context, arg UpdateBookParams) (Book, error) {
	row := q.db.QueryRowContext(ctx, updateBook, arg.ID, arg.Name)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.Language,
		&i.CreatedAt,
	)
	return i, err
}
