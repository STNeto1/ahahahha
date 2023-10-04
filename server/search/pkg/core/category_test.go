package core_test

import (
	searchpb "search/gen/protos"
	"search/pkg/core"
	"testing"

	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestFetchCategories(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	categories, err := core.FetchCategories(db)
	assert.NotNil(t, categories)
	assert.Len(t, *categories, 0)
	assert.Nil(t, err)
}

func TestGetCategoryWithSuccess(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	_sql, args := sqlbuilder.
		PostgreSQL.
		NewInsertBuilder().
		InsertInto("categories").
		Cols("id", "name", "slug", "parent_id").
		Values("foo", "foo", "foo", nil).
		Build()
	_, err := db.Exec(_sql, args...)
	assert.NoError(t, err)

	cat, err := core.GetCategory(db, "foo")
	assert.NotNil(t, cat)
	assert.NoError(t, err)
}

func TestGetCategoryWithInvalidId(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	cat, err := core.GetCategory(db, "bar")

	assert.Nil(t, cat)
	assert.Error(t, err)
}

func TestCreateCategoryWithSuccess(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	newCat, err := core.CreateCategory(db, &searchpb.CreateCategoryRequest{
		Name: "foo",
		Slug: "foo",
	})
	assert.NotNil(t, newCat)
	assert.Nil(t, err)
}

func TestCreateCategoryWithSlugAlreadyInUse(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	newCat, err := core.CreateCategory(db, &searchpb.CreateCategoryRequest{
		Name: "foo",
		Slug: "foo",
	})
	assert.NotNil(t, newCat)
	assert.Nil(t, err)

	invalidCat, err := core.CreateCategory(db, &searchpb.CreateCategoryRequest{
		Name: "foo",
		Slug: "foo",
	})
	assert.Nil(t, invalidCat)
	assert.Error(t, err)
}

func TestCreateCategoryWithInvalidParentId(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	invalidCat, err := core.CreateCategory(db, &searchpb.CreateCategoryRequest{
		Name:     "foo",
		Slug:     "foo",
		ParentId: getStrPointer("bar"),
	})
	assert.Nil(t, invalidCat)
	assert.Error(t, err)
}

func TestUpdateCategoryWithSuccess(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	newCat, err := core.CreateCategory(db, &searchpb.CreateCategoryRequest{
		Name: "foo",
		Slug: "foo",
	})
	assert.NotNil(t, newCat)
	assert.NoError(t, err)

	err = core.UpdateCategory(db, &searchpb.UpdateCategoryRequest{
		Id:   newCat.ID,
		Name: "bar",
		Slug: "bar",
	})
	assert.NoError(t, err)
}

func TestUpdateCategoryWithDifferentAndAlreadyInUseSlug(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	newCat, err := core.CreateCategory(db, &searchpb.CreateCategoryRequest{
		Name: "foo",
		Slug: "foo",
	})
	assert.NotNil(t, newCat)
	assert.Nil(t, err)

	anotherCat, err := core.CreateCategory(db, &searchpb.CreateCategoryRequest{
		Name: "bar",
		Slug: "bar",
	})
	assert.NotNil(t, anotherCat)
	assert.Nil(t, err)

	err = core.UpdateCategory(db, &searchpb.UpdateCategoryRequest{
		Id:   newCat.ID,
		Name: "bar",
		Slug: "bar",
	})
	assert.Error(t, err)
}

func TestDeleteCategoryWithSuccess(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	newCat, err := core.CreateCategory(db, &searchpb.CreateCategoryRequest{
		Name: "foo",
		Slug: "foo",
	})
	assert.NotNil(t, newCat)
	assert.Nil(t, err)

	err = core.DeleteCategory(db, &searchpb.DeleteCategoryRequest{
		Id: newCat.ID,
	})
	assert.NoError(t, err)

	cat, err := core.GetCategory(db, newCat.ID)
	assert.Nil(t, cat)
	assert.Error(t, err)
}

func TestDeleteCategoryWithInvalidId(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	err := core.DeleteCategory(db, &searchpb.DeleteCategoryRequest{
		Id: "bar",
	})
	assert.Error(t, err)
}

func getStrPointer(val string) *string {
	return &val
}
