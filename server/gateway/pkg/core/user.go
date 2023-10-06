package core

import (
	"models"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

func SearchUsers(db *sqlx.DB, term *string) (*[]models.User, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().From("users").Select("*")

	if term != nil {
		sb.Where(sb.Like("name", "%"+*term+"%"))
	}

	_sql, args := sb.Build()

	rows, err := db.Queryx(_sql, args...)
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func GetUser(db *sqlx.DB, id string) (*models.User, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().From("users").Select("*")

	_sql, args := sb.Where(sb.Equal("id", id)).Build()

	rows, err := db.Queryx(_sql, args...)
	if err != nil {
		return nil, err
	}

	var user models.User
	for rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}
