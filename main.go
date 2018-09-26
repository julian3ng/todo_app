package main;

import (
	"database/sql"
	"fmt"
	 _ "github.com/lib/pq"
)


const (
	host		= "127.0.0.1"
	port		= 5432
	user		= "todo_app"
	password	= "foo"
	dbname		= "todo"
)


type todo_db struct {
	handle *sql.DB
}

func (t *todo_db) check_error(err error) {
	if err != nil {
		t.handle.Close()
		panic(err)
	}
}


// Connects and tests the connection
// If the connection or the connection test fail, panic
func (t *todo_db) init() {
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
	fmt.Println("Successfully initialized")
}

func (t *todo_db) reset() {
	sqlStatement := `ALTER SEQUENCE todo_list_id_seq RESTART WITH 1`
	_, err := t.handle.Exec(sqlStatement)
	t.check_error(err)
	fmt.Println("Successfully reset")
}

func (t *todo_db) destroy() {
	t.handle.Close()
	fmt.Println("Successfully destroyed")
}

func (t *todo_db) insert(item string) {
	sqlStatement := `INSERT INTO todo_list (item) VALUES ($1)`
	_, err := t.handle.Exec(sqlStatement, item)
	t.check_error(err)
	fmt.Println("Successfully inserted")
}

func (t *todo_db) work(id int) {
	sqlStatement := `SELECT item FROM todo_list WHERE id = $1`
	var item string
	row := t.handle.QueryRow(sqlStatement, id)
	switch err := row.Scan(&item); err {
	case nil:
		fmt.Println("Got: ", item)
	case sql.ErrNoRows:
		fmt.Println("Got nothing.")
	default:
		panic(err)
	}
}

func (t *todo_db) delete(id int) {
	sqlStatement := `DELETE FROM todo_list WHERE id = $1`
	_, err := t.handle.Exec(sqlStatement, id)
	t.check_error(err)
	fmt.Println("Successfully deleted")
}

func main() {

	var tdb todo_db

	tdb.init()
	tdb.reset()
	tdb.insert("Eat")
	tdb.work(1)
	tdb.delete(1)
	tdb.destroy()
}
