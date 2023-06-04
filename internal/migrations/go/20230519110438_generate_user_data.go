package _go

import (
	"SocialNetHL/internal/helper"
	"crypto/sha256"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pressly/goose"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	goose.AddMigration(upGenerateUserData, downGenerateUserData)
}

func upGenerateUserData(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	migrationsDir := helper.GetEnvValue("MIGR_DIR", "./internal/migrations")
	file, err := os.Open(migrationsDir + "/go/people.csv")
	if err != nil {
		panic(err)
	}
	chRecord := make(chan []string, 1_000_000)
	var wg sync.WaitGroup
	wg.Add(1)
	go readFile(chRecord, file, &wg)
	time.Sleep(1 * time.Second)
loop:
	for {
		select {
		case rec := <-chRecord:
			wg.Add(1)
			go insertRow(rec, &wg, tx)
		default:
			break loop
		}
	}
	wg.Wait()
	close(chRecord)
	return nil
}

func downGenerateUserData(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

func insertRow(record []string, wg *sync.WaitGroup, tx *sql.Tx) {
	//record := <-ch
	//fmt.Println(len(ch))
	defer wg.Done()
	query := `INSERT INTO social.users (id, first_name, second_name, age, birthdate, biography, city, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	id := uuid.Must(uuid.NewV4()).String()
	age, _ := strconv.Atoi(record[1])
	bDate := time.Now().Add(time.Duration(-age) * time.Hour * 8760)
	firstName := strings.Split(record[0], " ")[1]
	secondName := strings.Split(record[0], " ")[0]

	_, err := tx.Exec(query, id, firstName, secondName, age, bDate, "biography", record[2], fmt.Sprintf("%x", sha256.Sum256([]byte("123"))))
	if err != nil {
		fmt.Printf("unable to insert row: %v", err)
	}
}

func readFile(ch chan []string, file *os.File, wg *sync.WaitGroup) {
	defer file.Close()
	defer wg.Done()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3
	for {
		record, e := reader.Read()
		if e != nil {
			break
		}
		ch <- record
	}
}
