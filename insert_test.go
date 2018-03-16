package paging 

import (
	"testing"

	"log"
	"fmt"	
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yuanchi/paging/util"
)

type User struct {
        Id, Name string
}

func InsertUsers(db *sql.DB, users []User) {
        insert := "INSERT INTO user (id, name) VALUES "
        for _, user := range users {
                id := util.CurrentTimef() + util.RandAlphabetic(5)
                id = util.PrepareString(id)
                name := util.PrepareString(user.Name)
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

func TestInsert(t *testing.T) {
	db, err := sql.Open("mysql", "tester:gogogo@/test")
	if err != nil {
		log.Fatal(err)	
	}
	defer db.Close()

	var users []User
	for i := 0; i < 3; i++ {
		users = append(users, User{Name: fmt.Sprintf("Name%d", i)})
	}
	InsertUsers(db, users)
}
