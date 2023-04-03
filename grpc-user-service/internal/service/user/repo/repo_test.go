package repo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func setup(t *testing.T) (*repo, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewRepo(sqlxDB, 5)

	return r, mock
}

func TestRepo_CreateUser(t *testing.T) {
	r, dbMock := setup(t)
	ctx := context.Background()

	u := model.User{
		ID:       0,
		Username: "Vasyan123",
		Email:    "Vasya@pro.rab",
		Name:     "Vasiliy",
		Surname:  "Pupkevich",
	}

	query := `INSERT INTO users (username,email) VALUES ($1,$2) RETURNING id`
	dbMock.ExpectExec(query).
		WithArgs(u.Username, u.Email).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err1 := r.CreateUser(ctx, &u)

	require.NoError(t, err1)

	err2 := dbMock.ExpectationsWereMet()
	require.NoError(t, err2)

}
