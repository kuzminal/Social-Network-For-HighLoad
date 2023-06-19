package _go

import (
	"SocialNetHL/internal/helper"
	"bufio"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pressly/goose"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func init() {
	goose.AddMigration(upAddPostsAndFriends, downAddPostsAndFriends)
}

var (
	users = []string{
		"5dde64fb-c81f-47d0-8017-790255347edf",
		"9737b0a7-5ae1-43ce-8f41-541946cf2f06",
		"1eeb70a2-ca04-4ab9-baab-cc21d8a1fa52",
		"dc2e26ce-06dc-4e91-9400-4e51f06b3252",
		"5d4e2590-a63b-4f80-b21e-df3d58c655a8",
		"0a84fb9e-1302-4418-835c-a7dd1cc9a453",
		"9d3d40e6-638a-4f17-85c1-838dd7d180a9",
		"82a6e475-046b-4d25-9630-8c785817be8a",
		"6e680abd-061e-4931-94d2-205549d34755",
		"889e9dda-36bd-4422-9b71-e4ca2f99186f",
	}
)

func upAddPostsAndFriends(tx *sql.Tx) error {
	for _, user := range users {
		addFriends(tx, user)
	}

	migrationsDir := helper.GetEnvValue("MIGR_DIR", "./internal/migrations")
	file, err := os.Open(migrationsDir + "/go/posts.txt")
	if err != nil {
		panic(err)
	}
	chRecord := make(chan string, 10_000)
	var wg sync.WaitGroup
	wg.Add(1)
	go readTxtFile(chRecord, file, &wg)
	time.Sleep(1 * time.Second)
loop:
	for {
		select {
		case rec := <-chRecord:
			wg.Add(1)
			go insertPost(rec, &wg, tx)
		default:
			break loop
		}
	}
	wg.Wait()
	close(chRecord)
	return nil
}

func downAddPostsAndFriends(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

func readTxtFile(ch chan string, file *os.File, wg *sync.WaitGroup) {
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	defer wg.Done()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// do something with a line
		record := scanner.Text()
		ch <- record
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func insertPost(record string, wg *sync.WaitGroup, tx *sql.Tx) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(10-1) + 1

	query := `INSERT INTO social.posts (id, "text", author_user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id;`
	id := uuid.Must(uuid.NewV4()).String()

	_, err := tx.Exec(query, id, record, users[idx-1], time.Now())
	if err != nil {
		fmt.Printf("unable to insert row: %v", err)
	}
}

func addFriends(tx *sql.Tx, userId string) {
	query := `INSERT INTO social.friends (user_id, friend_id, created_at) VALUES ($1, $2, $3);`
	_, err := tx.Exec(query, "1", userId, time.Now())
	if err != nil {
		fmt.Printf("unable to insert row: %v", err)
	}
}
