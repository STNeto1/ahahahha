package core

type Category struct {
	ID        string  `db:"id"`
	Name      string  `db:"name"`
	Slug      string  `db:"slug"`
	CreatedAt string  `db:"created_at"`
	ParentID  *string `db:"parent_id"`
}
