package pg

import (
	"SocialNetHL/models"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"os"
	"time"
)

func (pg *Postgres) SaveUser(ctx context.Context, user models.RegisterUser) (id string, err error) {
	query := `INSERT INTO social.users (id, first_name, second_name, age, birthdate, biography, city, password) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	id = uuid.Must(uuid.NewV4()).String()
	bDate, _ := time.Parse("2006-01-02", user.Birthdate)
	age := calculateAge(bDate)
	_, err = pg.db.Exec(ctx, query, id, user.FirstName, user.SecondName, age, bDate, user.Biography, user.City, fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))))
	if err != nil {
		return "", fmt.Errorf("unable to insert row: %w", err)
	}
	os.WriteFile("filename", []byte(id+"\n"), 0644)
	return id, nil
}

func (pg *Postgres) LoadUser(ctx context.Context, id string) (usersInfo models.UserInfo, err error) {
	query := `SELECT id, first_name, second_name, age, birthdate, biography, city, password FROM social.users WHERE id = $1`

	row := pg.db.QueryRow(ctx, query, id)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("unable to query users: %w", err)
	}
	var bDate time.Time
	user := models.UserInfo{}
	err = row.Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Age, &bDate, &user.Biography, &user.City, &user.Password)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("unable to scan row: %w", err)
	}
	user.Birthdate = bDate.Format("2006-01-02")
	return user, nil
}

func calculateAge(bDate time.Time) int {
	curDate := time.Now()
	dur := curDate.Sub(bDate)
	return int(dur.Seconds() / 31207680)
}

func (pg *Postgres) SearchUser(ctx context.Context, request models.UserSearchRequest) (users []models.UserInfo, err error) {
	query := `SELECT id, first_name, second_name, age, birthdate, biography, city, password FROM social.users WHERE first_name LIKE $1 AND second_name LIKE $2 ORDER BY id`
	//cont, cancel := context.WithTimeout(ctx, 2*time.Second)
	//defer cancel()
	rows, err := pg.db.Query(ctx, query, request.FirstName+"%", request.LastName+"%")
	defer rows.Close()
	if err != nil {
		return []models.UserInfo{}, fmt.Errorf("unable to query users: %w", err)
	}
	user := models.UserInfo{}
	var bDate time.Time
	for rows.Next() {
		//user := models.UserInfo{}
		err = rows.Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Age, &bDate, &user.Biography, &user.City, &user.Password)
		if err != nil {
			log.Printf("unable to scan row: %v", err)
		}
		user.Birthdate = bDate.Format("2006-01-02")
		users = append(users, user)
	}
	return users, nil
}

func (pg *Postgres) CheckIfExistsUser(ctx context.Context, userId string) (bool, error) {
	query := `SELECT id FROM social.users WHERE id=$1`
	row := pg.db.QueryRow(ctx, query, userId)

	var user string
	err := row.Scan(&user)
	if err != nil {
		return false, err
	}
	if len(user) == 0 {
		return false, nil
	}
	return true, nil
}
