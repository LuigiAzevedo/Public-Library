package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &entity.User{
		ID:        1,
		Username:  "user135",
		Password:  "secret",
		Email:     "user135@email.com",
		UpdatedAt: time.Time{},
		CreatedAt: time.Now(),
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "amount", "updated_at", "created_at"}).
			AddRow(user.ID, user.Username, user.Password, user.Email, user.UpdatedAt, user.CreatedAt)

		mock.ExpectPrepare("SELECT \\* FROM users WHERE id = ").
			ExpectQuery().
			WithArgs(user.ID).
			WillReturnRows(rows)

		gotUser, err := repo.Get(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user, gotUser)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM users WHERE id = ").
			WillReturnError(sql.ErrConnDone)

		gotBook, err := repo.Get(context.Background(), user.ID)
		assert.Error(t, err)
		assert.Empty(t, gotBook)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Scan Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM users WHERE id = ").
			ExpectQuery().
			WithArgs(0)

		gotBook, err := repo.Get(context.Background(), 0)
		assert.Error(t, err)
		assert.Empty(t, gotBook)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM users WHERE id = ").
			ExpectQuery().
			WithArgs(user.ID).
			WillReturnError(sql.ErrNoRows)

		gotUser, err := repo.Get(context.Background(), user.ID)
		assert.Error(t, err)
		assert.Empty(t, gotUser)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &entity.User{
		Username: "user135",
		Password: "secret",
		Email:    "user135@email.com",
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(1)

		mock.ExpectPrepare("INSERT INTO users").
			ExpectQuery().
			WithArgs(user.Username, user.Password, user.Email).
			WillReturnRows(rows)

		id, err := repo.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO users").
			WillReturnError(sql.ErrConnDone)

		id, err := repo.Create(context.Background(), user)
		assert.Error(t, err)
		assert.Empty(t, id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO users").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		id, err := repo.Create(context.Background(), user)
		assert.Error(t, err)
		assert.Empty(t, id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &entity.User{
		ID:       1,
		Username: "user135",
		Password: "secret",
		Email:    "user135@email.com",
	}

	t.Run("OK", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE users").
			ExpectExec().
			WithArgs(user.Username, user.Password, user.Email, user.ID).
			WillReturnResult(sqlmock.NewResult(int64(user.ID), 1))

		err := repo.Update(context.Background(), user)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE users").
			WillReturnError(sql.ErrConnDone)

		err := repo.Update(context.Background(), user)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Exec Failed", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE users").
			ExpectExec().
			WithArgs(user.Username, user.Password, user.Email, user.ID).
			WillReturnError(sql.ErrConnDone)

		err := repo.Update(context.Background(), user)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE users").
			ExpectExec().
			WithArgs(user.Username, user.Password, user.Email, user.ID).
			WillReturnResult(sqlmock.NewResult(int64(user.ID), 0))

		err := repo.Update(context.Background(), user)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("OK", func(t *testing.T) {
		mock.ExpectPrepare("DELETE FROM users WHERE id =").
			ExpectExec().
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		err := repo.Delete(context.Background(), 1)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("DELETE FROM users WHERE id =").
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Exec Failed", func(t *testing.T) {
		mock.ExpectPrepare("DELETE FROM users WHERE id =").
			ExpectExec().
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("DELETE FROM users WHERE id =").
			ExpectExec().
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(int64(1), 0))

		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
