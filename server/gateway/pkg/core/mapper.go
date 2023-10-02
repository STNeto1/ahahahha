package core

import "gateway/pkg/graph/model"

func MapUser(user *User) *model.User {
	return &model.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func MapUsers(users *[]User) []*model.User {
	var usersModel []*model.User

	for _, user := range *users {
		usersModel = append(usersModel, MapUser(&user))
	}

	return usersModel
}
