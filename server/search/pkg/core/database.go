package core

import (
	"log"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=auction sslmode=disable")
	if err != nil {
		log.Fatalln("failed to connect", err)
	}

	for _, schema := range getSchemas() {
		_, err := db.Exec(schema)
		if err != nil {
			log.Fatalln("failed to create table", err)
		}
	}

	return db
}

func InitTempDB() *sqlx.DB {
	// This commented like is to create a temporary database in /tmp
	// db, err := sqlx.Connect("sqlite3", fmt.Sprintf("file:/tmp/%s.db", ulid.MustNew(ulid.Now(), nil).String()))
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.Fatalln("failed to connect", err)
	}

	for _, schema := range getSchemas() {
		_, err := db.Exec(schema)
		if err != nil {
			log.Fatalln("failed to create table", err)
		}
	}

	return db
}

func getSchemas() []string {
	categoriesTableSql, _ := sqlbuilder.NewCreateTableBuilder().
		CreateTable("categories").
		IfNotExists().
		Define("id", "varchar(26)", "PRIMARY KEY").
		Define("name", "varchar(255)", "NOT NULL").
		Define("slug", "varchar(255)", "NOT NULL UNIQUE").
		Define("parent_id", "varchar(26)", "REFERENCES categories(id)").
		Define("created_at", "timestamp", "NOT NULL DEFAULT CURRENT_TIMESTAMP").
		Build()

	return []string{categoriesTableSql}
}
