package main;

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"strings"
)


type todo_db struct {
	handle *sql.DB
}


// Connects and tests the connection
// If the connection or the connection test fail, panic
func (t *todo_db) init(host string, port int, user, password, dbname string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s " +
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	t.handle, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	
	if err = t.handle.Ping(); err != nil {
		t.handle.Close()
		panic(err)
	}
	fmt.Println("Initialized")
}

// Reset ids for debugging
func (t *todo_db) reset() {
	sqlStatement := `ALTER SEQUENCE todo_list_id_seq RESTART WITH 1`
	_, err := t.handle.Exec(sqlStatement)
	if  err != nil {
		panic(err)
	}
	fmt.Println("Reset")
}

// Close the handle!
func (t *todo_db) destroy() {
	t.handle.Close()
	fmt.Println("Destroyed")
}

// Insert item into the database
func (t *todo_db) insert(item string) {
	sqlStatement := `INSERT INTO todo_list (item) VALUES ($1) RETURNING item`
	var outitem string
	err := t.handle.QueryRow(sqlStatement, item).Scan(&outitem)
	switch err {
	case nil:
		fmt.Println("Inserted", outitem)
	case sql.ErrNoRows:
		fmt.Println("Nothing to insert.")
	default:
		panic(err)
	}
}

// Show all instances of item from the database
func (t *todo_db) display(item string) string{
	sqlStatement := `SELECT item FROM todo_list WHERE item = $1`
	rows, err := t.handle.Query(sqlStatement, item)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var b strings.Builder
	for rows.Next() {
		var it string
		err = rows.Scan(&it)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&b, "%s\n", it)
	}

	return b.String()
}

// Return a string representation of the entire list
func (t *todo_db) display_all() string {
	sqlStatement := `SELECT item FROM todo_list ORDER BY id`

	rows, err := t.handle.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	
	var b strings.Builder
	for rows.Next() {
		var item string
		err = rows.Scan(&item)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&b, "%s\n", item)
	}

	return b.String()
}

// Delete minimum id'd instance of item from the database 
func (t *todo_db) delete(item string) {
	sqlStatement := `DELETE FROM todo_list
WHERE item = $1 AND id = (SELECT min(id) FROM todo_list WHERE item = $1) RETURNING item`
	var outitem string
	err := t.handle.QueryRow(sqlStatement, item).Scan(&outitem)
	switch err {
	case nil:
		fmt.Println("Deleted", outitem)
	case sql.ErrNoRows:
		fmt.Println("Nothing to delete.")
	default:
		panic(err)
	}

	return
}

// Delete everything from the database
func (t *todo_db) clear() {
	_, err := t.handle.Exec(`DELETE FROM todo_list`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted everything.")
}
