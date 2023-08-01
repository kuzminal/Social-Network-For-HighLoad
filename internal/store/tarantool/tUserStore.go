package tarantool

import (
	"SocialNetHL/models"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"time"
)

type UserFromStore struct {
	UserId string
	User   models.UserInfo
}

type UserIdFromStore struct {
	UserId []string
}

func (t *TarantoolStore) SaveUser(ctx context.Context, user models.RegisterUser) (id string, err error) {
	id = uuid.Must(uuid.NewV4()).String()
	userInfo := models.UserInfo{}
	userInfo.Id = id
	userInfo.FirstName = user.FirstName
	userInfo.SecondName = user.SecondName
	bDate, _ := time.Parse("2006-01-02", user.Birthdate)
	userInfo.Age = calculateAge(bDate)
	userInfo.City = user.City
	userInfo.Biography = user.Biography
	userInfo.Birthdate = user.Birthdate
	userInfo.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
	var userId UserIdFromStore
	err = t.conn.CallTyped("create_user",
		[]interface{}{
			id,
			userInfo.FirstName,
			userInfo.SecondName,
			userInfo.Age,
			userInfo.Birthdate,
			userInfo.Biography,
			userInfo.City,
			userInfo.Password,
		},
		&userId)
	if err != nil {
		return "", err
	}
	if len(userId.UserId) != 1 {
		return "", errors.Errorf("Cannot create user with id: %s", userId)
	} else {
		return userId.UserId[0], nil
	}
}

func (t *TarantoolStore) LoadUser(ctx context.Context, id string) (usersInfo models.UserInfo, err error) {

	var user []models.UserInfo
	err = t.conn.CallTyped("get_user_by_id", []interface{}{id}, &user)
	if err != nil {
		return models.UserInfo{}, err
	}
	if len(user) == 1 {
		usersInfo = user[0]
	} else {
		errors.Errorf("Cannot find user with id: %s", id)
	}

	return usersInfo, nil
}

func calculateAge(bDate time.Time) int {
	curDate := time.Now()
	dur := curDate.Sub(bDate)
	return int(dur.Seconds() / 31207680)
}

func (t *TarantoolStore) SearchUser(ctx context.Context, request models.UserSearchRequest) (users []models.UserInfo, err error) {
	var user []models.UserInfo
	err = t.conn.CallTyped("search_user", []interface{}{request.FirstName, request.LastName}, &user)
	if err != nil {
		return []models.UserInfo{}, err
	}
	if len(user) == 1 && user[0].Id == "" {
		return users, nil
	}
	return user, nil
}

func (t *TarantoolStore) CheckIfExistsUser(ctx context.Context, userId string) (bool, error) {
	var res bool
	err := t.conn.CallTyped("check_user_exists", []interface{}{userId}, &res)
	if err != nil {
		return false, err
	}
	return res, nil
}
