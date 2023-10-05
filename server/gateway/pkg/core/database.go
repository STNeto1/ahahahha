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
	usersTableSql, _ := sqlbuilder.NewCreateTableBuilder().
		CreateTable("users").
		IfNotExists().
		Define("id", "varchar(26)", "PRIMARY KEY").
		Define("name", "varchar(255)", "NOT NULL").
		Define("email", "varchar(255)", "NOT NULL UNIQUE").
		Define("password", "varchar(255)", "NOT NULL").
		Define("role", "varchar(255)", "NOT NULL", "DEFAULT 'user'").
		Define("created_at", "timestamp", "NOT NULL DEFAULT CURRENT_TIMESTAMP").
		Build()

	return []string{usersTableSql}
}
