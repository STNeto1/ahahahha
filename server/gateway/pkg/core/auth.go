package core

import (
	"context"
	"database/sql"
	"log"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sqlx.DB, name, email, password string) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().
		From("users")
	_sql, args := sb.Select("count(*)").
		Where(sb.Equal("email", email)).
		Build()

	res, err := tx.Query(_sql, args...)
	if err != nil {
		rollback(tx)

		return err
	}

	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			rollback(tx)

			return err
		}
	}

	if count != 0 {
		rollback(tx)

		return ErrUserAlreadyExists
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		rollback(tx)

		return err
	}

	_sql, args = sqlbuilder.PostgreSQL.NewInsertBuilder().
		InsertInto("users").Cols("id", "name", "email", "password").
		Values(ulid.Make().String(), name, email, string(hashedPwd)).
		Build()

	_, err = tx.Exec(_sql, args...)
	if err != nil {
		rollback(tx)

		return err
	}

	return tx.Commit()
}

func AuthenticateUser(db *sqlx.DB, email, password string) (*User, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().From("users")
	_sql, args := sb.Select("*").
		Where(sb.Equal("email", email)).
		Limit(1).
		Build()

	res, err := db.Query(_sql, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserDoesNotExists
		}

		return nil, err
	}

	var user User
	for res.Next() {
		err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

func rollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Println("failed to rollback", err)
	}
}
