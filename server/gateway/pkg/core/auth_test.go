package core_test

import (
	"gateway/pkg/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserWithSuccess(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	err := core.CreateUser(db, "foo", "mail@mail.com", "102030")
	assert.Nil(t, err)
}

func TestNotCreateUserWithEmailReadyInUse(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	err := core.CreateUser(db, "foo", "mail@mail.com", "102030")
	assert.Nil(t, err)

	err = core.CreateUser(db, "bar", "mail@mail.com", "102030")
	assert.Equal(t, core.ErrUserAlreadyExists, err)
}

func TestAuthenticateWithSuccess(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	err := core.CreateUser(db, "foo", "mail@mail.com", "102030")
	assert.Nil(t, err)

	user, err := core.AuthenticateUser(db, "mail@mail.com", "102030")
	assert.Nil(t, err)
	assert.Equal(t, "foo", user.Name)
}

func TestFailAuthenticationWithInvalidEmail(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	user, err := core.AuthenticateUser(db, "invalid@mail.com", "102030")
	assert.Nil(t, user)
	assert.Error(t, err, core.ErrUserDoesNotExists)
}

func TestFailAuthenticationWithInvalidPassword(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	err := core.CreateUser(db, "foo", "mail@mail.com", "102030")
	assert.Nil(t, err)

	user, err := core.AuthenticateUser(db, "mail@mail.com", "10203040")
	assert.Nil(t, user)
	assert.Error(t, err, core.ErrUserDoesNotExists)
}
