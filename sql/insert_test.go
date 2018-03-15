package sql

import (
	"log"
	"fmt"
	"testing"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

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
