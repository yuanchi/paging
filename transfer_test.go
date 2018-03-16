package paging

import (
	"testing"
	"fmt"
	"log"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func find(db *sql.DB, p *Paging) []User {
	nextTpl := `
		select *
		from user u
		where {{.Conds}}
		order by {{.Sorts}}
		limit {{.Limit}}
	`
	q := SortQueryByUniqueKey(p, nextTpl, "u")
	fmt.Println(q)	
	rows, err := db.Query(q)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User	
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			log.Fatalln(err)
		}
		users = append(users, user)	
	}
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}
	
	l := len(users)
	fmt.Printf("found count: %d\n", l)
	if l > 0 {
		fUser := users[0]
        	lUser := users[l-1]
        	fmt.Println("first:", fUser)
        	fmt.Println("last:", lUser)
	}
        return users
}

func findNext(db *sql.DB, p *Paging) {
	users := find(db, p)
	l := len(users)
	if l > 0 {
		findNext(db, &Paging{Limit: p.Limit, Next: p.Next, Id: p.Id, IdVal: users[l-1].Id, IdDesc: p.IdDesc})
	} 
}

func findPrev(db *sql.DB, p *Paging) {
	users := find(db, p)
	l := len(users)
	if l > 0 {
		findPrev(db, &Paging{Limit: p.Limit, Next: p.Next, Id: p.Id, IdVal: users[0].Id, IdDesc: p.IdDesc})
	}
}

func TestSortQueryByUniqueKey(t *testing.T) {
	
	db, err := sql.Open("mysql", "tester:gogogo@/test")
	if err != nil {
		log.Fatalln(err)	
	}
	defer db.Close()

	p := Paging{Limit: 3, Next: true, Id: "id", IdDesc: true}
	findNext(db, &p)
	
	//findPrev(db, &p)
}
/*
func TestSortQueryBy(t *testing.T) {
	tpl := `
		select *
		from message m
		where {{.Conds}}
		order by {{.Sorts}}
		limit {{.Limit}}
	`
	p1 := Paging {Limit: 10, Next: true, Id: "id", IdDesc: true}
	p1.Fields = append(p1.Fields, &FieldData{Unique: false, Desc: true, Name: "thumbs_up"})
	q1 := SortQueryBy(&p1, tpl, "m")
	fmt.Println("q1:", q1)

	p2 := Paging {Limit: 10, Next: true, Id: "id", IdVal: "10", IdDesc: true}
	p2.Fields = append(p2.Fields, &FieldData{Unique: true, Desc: true, Name: "id", Value: "10"})
	p2.Fields = append(p2.Fields, &FieldData{Unique: false, Desc: true, Name: "thumbs_up", Value: "98"})
	q2 := SortQueryBy(&p2, tpl, "m")
	fmt.Println("q2:", q2)

	p3 := p2
	p3.Fields = append(p3.Fields, &FieldData{Unique: false, Desc: true, Name: "Name", Value: "Mary"})
	q3 := SortQueryBy(&p3, tpl, "m")
	fmt.Println("q3:", q3)
}
*/
