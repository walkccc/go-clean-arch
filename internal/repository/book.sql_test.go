package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/walkccc/go-clean-arch/internal/util"
)

func createRandomBook(t *testing.T) Book {
	user := createRandomUser(t)

	arg := CreateBookParams{
		Owner:    user.Username,
		Name:     util.RandomBookName(),
		Language: util.RandomLanguage(),
	}

	book, err := testQueries.CreateBook(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, book)

	require.Equal(t, arg.Owner, book.Owner)
	require.Equal(t, arg.Name, book.Name)
	require.Equal(t, arg.Language, book.Language)

	require.NotZero(t, book.ID)
	require.NotZero(t, book.CreatedAt)
	return book
}

func TestCreateBook(t *testing.T) {
	createRandomBook(t)
}

func TestGetBook(t *testing.T) {
	book1 := createRandomBook(t)
	book2, err := testQueries.GetBook(context.Background(), book1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, book2)

	require.Equal(t, book1.ID, book1.ID)
	require.Equal(t, book1.Owner, book2.Owner)
	require.Equal(t, book1.Name, book2.Name)
	require.Equal(t, book1.Language, book2.Language)
	require.WithinDuration(t, book1.CreatedAt, book2.CreatedAt, time.Second)
}

func TestUpdateBook(t *testing.T) {
	book1 := createRandomBook(t)

	arg := UpdateBookParams{
		ID:   book1.ID,
		Name: util.RandomBookName(),
	}

	book2, err := testQueries.UpdateBook(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, book2)

	require.Equal(t, book1.ID, book1.ID)
	require.Equal(t, book1.Owner, book2.Owner)
	require.Equal(t, arg.Name, book2.Name)
	require.Equal(t, book1.Language, book2.Language)
	require.WithinDuration(t, book1.CreatedAt, book2.CreatedAt, time.Second)
}

func TestDeleteBook(t *testing.T) {
	book1 := createRandomBook(t)
	err := testQueries.DeleteBook(context.Background(), book1.ID)
	require.NoError(t, err)

	book2, err := testQueries.GetBook(context.Background(), book1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, book2)
}

func TestListBooks(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomBook(t)
	}

	arg := ListBooksParams{
		Limit:  5,
		Offset: 5,
	}

	books, err := testQueries.ListBooks(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, books, 5)

	for _, book := range books {
		require.NotEmpty(t, book)
	}
}
