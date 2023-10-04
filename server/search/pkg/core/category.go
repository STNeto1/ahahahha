package core

import (
	"context"
	"database/sql"
	"errors"
	"log"
	searchpb "search/gen/protos"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

func FetchCategories(db *sqlx.DB) (*[]Category, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().From("categories").Select("*")

	_sql, args := sb.Build()

	rows, err := db.Queryx(_sql, args...)
	if err != nil {
		return nil, err
	}

	data := []Category{}
	for rows.Next() {
		var row Category
		err := rows.Scan(&row.ID, &row.Name, &row.Slug, &row.ParentID, &row.CreatedAt)
		if err != nil {
			return nil, err
		}

		data = append(data, row)
	}

	return &data, nil
}

func GetCategory(db *sqlx.DB, id string) (*Category, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().From("categories").Select("*")
	_sql, args := sb.Where(sb.Equal("id", id)).Limit(1).Build()

	rows, err := db.Queryx(_sql, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCategoryDoesNotExists
		}

		return nil, err
	}

	if !rows.Next() {
		return nil, ErrCategoryDoesNotExists
	}

	var row Category
	if err := rows.Scan(&row.ID, &row.Name, &row.Slug, &row.ParentID, &row.CreatedAt); err != nil {
		return nil, err
	}

	return &row, nil
}

func CreateCategory(db *sqlx.DB, data *searchpb.CreateCategoryRequest) (*Category, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().From("categories")
	_sql, args := sb.Select("count(*)").Where(sb.Equal("slug", data.Slug)).Limit(1).Build()

	rows, err := tx.Query(_sql, args...)
	if err != nil {
		rollback(tx)

		return nil, err
	}

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			rollback(tx)

			return nil, err
		}
	}
	if count > 0 {
		rollback(tx)

		return nil, errors.New("slug already exists")
	}

	if data.ParentId != nil {
		_sql, args := sb.Select("count(*)").Where(sb.Equal("id", data.ParentId)).Limit(1).Build()

		_, err = tx.Query(_sql, args...)
		if err != nil {
			rollback(tx)

			return nil, err
		}

		for rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				rollback(tx)

				return nil, err
			}
		}

		if count == 0 {
			rollback(tx)

			return nil, errors.New("parent does not exist")
		}
	}

	newID := ulid.Make().String()

	_sql, args = sqlbuilder.
		PostgreSQL.
		NewInsertBuilder().
		InsertInto("categories").
		Cols("id", "name", "slug", "parent_id").
		Values(newID, data.Name, data.Slug, data.ParentId).
		Build()

	_, err = tx.Exec(_sql, args...)
	if err != nil {
		rollback(tx)

		return nil, err
	}

	return &Category{
		ID:       newID,
		Name:     data.Name,
		Slug:     data.Slug,
		ParentID: data.ParentId,
	}, tx.Commit()
}

func UpdateCategory(db *sqlx.DB, data *searchpb.UpdateCategoryRequest) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().From("categories").Select("*")

	_sql, args := sb.Where(sb.Equal("id", data.Id)).Limit(1).Build()
	rows, err := tx.Query(_sql, args...)
	if err != nil {
		rollback(tx)

		if err == sql.ErrNoRows {
			return ErrCategoryDoesNotExists
		}

		return err
	}

	var row Category
	for rows.Next() {
		err := rows.Scan(&row.ID, &row.Name, &row.Slug, &row.ParentID, &row.CreatedAt)
		if err != nil {
			rollback(tx)
			return err
		}
	}

	if err := rows.Close(); err != nil {
		rollback(tx)

		return err
	}

	if row.Slug != data.Slug {
		_sql, args = sb.Where(sb.Equal("slug", data.Slug)).Limit(1).Build()

		rows, err = tx.Query(_sql, args...)
		if err != nil && err != sql.ErrNoRows {
			rollback(tx)

			return err
		}

		if err := rows.Close(); err != nil {
			rollback(tx)

			return err
		}
	}

	if data.ParentId != nil && row.ParentID != data.ParentId {
		_sql, args := sb.Where(sb.Equal("id", data.ParentId)).Limit(1).Build()

		rows, err = tx.Query(_sql, args...)
		if err != nil {
			rollback(tx)

			if err == sql.ErrNoRows {
				return ErrCategoryDoesNotExists
			}

			return err
		}

		if err := rows.Close(); err != nil {
			rollback(tx)

			return err
		}
	}

	ub := sqlbuilder.
		PostgreSQL.
		NewUpdateBuilder().
		Update("categories")

	_sql, args = ub.
		Set(
			ub.Assign("name", data.Name),
			ub.Assign("slug", data.Slug),
			ub.Assign("parent_id", data.ParentId),
		).
		Where(ub.Equal("id", data.Id)).
		Build()

	_, err = tx.Exec(_sql, args...)
	if err != nil {
		rollback(tx)

		return err
	}

	return tx.Commit()
}

func DeleteCategory(db *sqlx.DB, data *searchpb.DeleteCategoryRequest) error {
	ctx := context.TODO()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().From("categories").Select("*")

	_sql, args := sb.Where(sb.Equal("id", data.Id)).Limit(1).Build()
	rows, err := tx.QueryContext(ctx, _sql, args...)
	if err != nil {
		rollback(tx)

		return err
	}

	if !rows.Next() {
		rollback(tx)

		return ErrCategoryDoesNotExists
	}
	if err := rows.Close(); err != nil {
		rollback(tx)

		return err
	}

	ub := sqlbuilder.
		PostgreSQL.
		NewDeleteBuilder().
		DeleteFrom("categories")

	_sql, args = ub.
		Where(ub.Equal("id", data.Id)).
		Build()

	_, err = tx.ExecContext(ctx, _sql, args...)
	if err != nil {
		rollback(tx)

		return err
	}

	return tx.Commit()
}

func rollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Println("failed to rollback", err)
	}
}
