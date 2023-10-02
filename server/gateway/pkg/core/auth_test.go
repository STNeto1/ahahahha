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

func TestUpdateWithSuccess(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	err := core.CreateUser(db, "foo", "mail@mail.com", "102030")
	assert.NoError(t, err)

	user, err := core.AuthenticateUser(db, "mail@mail.com", "102030")
	assert.NotNil(t, user)
	assert.NoError(t, err)

	newName := "new name"
	newMail := "bar"
	newPassword := "102030"
	err = core.UpdateUser(db, user, &newName, &newMail, &newPassword)
	assert.NoError(t, err)
}

func TestUpdateWithNullFields(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	err := core.CreateUser(db, "foo", "mail@mail.com", "102030")
	assert.NoError(t, err)

	err = core.CreateUser(db, "foo", "mail2@mail.com", "102030")
	assert.NoError(t, err)

	user, err := core.AuthenticateUser(db, "mail@mail.com", "102030")
	assert.NotNil(t, user)
	assert.NoError(t, err)

	err = core.UpdateUser(db, user, nil, nil, nil)
	assert.NoError(t, err)
}

func TestUpdateWithAlreadyExistingMail(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	err := core.CreateUser(db, "foo", "mail@mail.com", "102030")
	assert.NoError(t, err)

	err = core.CreateUser(db, "foo", "mail2@mail.com", "102030")
	assert.NoError(t, err)

	user, err := core.AuthenticateUser(db, "mail@mail.com", "102030")
	assert.NotNil(t, user)
	assert.NoError(t, err)

	newName := "new name"
	newMail := "mail2@mail.com"
	newPassword := "102030"
	err = core.UpdateUser(db, user, &newName, &newMail, &newPassword)
	assert.Error(t, err, core.ErrUserAlreadyExists)
}

func TestDeleteUser(t *testing.T) {
	db := core.InitTempDB()
	defer db.Close()

	err := core.CreateUser(db, "foo", "mail@mail.com", "102030")
	assert.NoError(t, err)

	user, err := core.AuthenticateUser(db, "mail@mail.com", "102030")
	assert.NotNil(t, user)
	assert.NoError(t, err)

	err = core.DeleteUser(db, user)
	assert.NoError(t, err)

	user, err = core.AuthenticateUser(db, "mail@mail.com", "102030")
	assert.Nil(t, user)
	assert.Error(t, err, core.ErrUserDoesNotExists)

}
