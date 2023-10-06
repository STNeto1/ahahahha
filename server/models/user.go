package models

type User struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	Role      string `db:"role"`
	CreatedAt string `db:"created_at"`
}
