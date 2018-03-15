package sql

import (
	"fmt"
	"log"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yuanchi/paging/sql/sqlutil"
	"github.com/yuanchi/paging/time/timeutil"
	"github.com/yuanchi/paging/rand/randutil"
)

type User struct {
	Id, Name string
}


func InsertUsers(db *sql.DB, users []User) {
	insert := "INSERT INTO user (id, name) VALUES "
	for _, user := range users {
		id := timeutil.CurrentTimef() + randutil.RandAlphabetic(5)
		id = sqlutil.PrepareString(id)
		name := sqlutil.PrepareString(user.Name)
		insert += fmt.Sprintf("(%s, %s),\n", id, name)
	}
	count := len(",\n")
	insert = insert[:len(insert)-count]
	fmt.Println("insert sql:\n", insert)

	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Rollback()
	result, err := tx.Exec(insert)
	if err != nil {
		log.Fatalln(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("rows affected count: %d", rows)
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
} 


