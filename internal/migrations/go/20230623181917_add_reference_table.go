package _go

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddReferenceTable, downAddReferenceTable)
}

func upAddReferenceTable(tx *sql.Tx) error {
	query := `select create_reference_table('social.users');`
	_, err := tx.Exec(query)
	if err != nil {
		fmt.Printf("unable to insert row: %v", err)
	}
	return nil
}

func downAddReferenceTable(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
