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

func connect() (db, err) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s " +
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	defer db.Close()
	if err != nil {
		return db, err
	}
	
	if err = db.Ping(); err != nil {
		return db, err
	}

	fmt.Println("Successfully connected!")

	return db, nil
}

func main() {

	db, err := connect()

	if err != nil {
		panic(err)
	}
	

	sqlStatement := `
INSERT INTO users (age, email, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, 30, "jon@gcalhoun.io", "Jonathan", "Calhoun").Scan(&id)
	if err != nil {
		panic(err)
	}

	err = db.QueryRow(sqlStatement, 52, "bob@smith.io", "Bob", "Smith").Scan(&id)
	if err != nil {
		panic(err)
	}

	err = db.QueryRow(sqlStatement, 15, "jerryjr123@gmail.com", "Jerry", "Seinfeld").Scan(&id)
	if err != nil {
		panic(err)
	}	
	
	fmt.Println("New record ID is:", id)

	sqlStatement = `
UPDATE users
SET first_name = $2, last_name = $3
WHERE id = $1;`
	_, err = db.Exec(sqlStatement, 1, "NewFirst", "NewLast")
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully altered!")

	sqlStatement = `SELECT id, email FROM users WHERE id=$1`
	var email string

	row := db.QueryRow(sqlStatement, 3)
	
	switch err := row.Scan(&id, &email); err {
	case nil:
		fmt.Println(id, email)		
	case sql.ErrNoRows:
		fmt.Println("No rows returned!")
	default:
		panic(err)
	}


	fmt.Println("=====Multiple row query=====")
	rows, err := db.Query("SELECT id, first_name FROM users LIMIT $1", 3)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var firstName string
		err = rows.Scan(&id, &firstName)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, firstName)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("==========")

	sqlStatement = `
DELETE FROM users
WHERE id > $1`

	res, err := db.Exec(sqlStatement, 0)
	if err != nil {
		panic(err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Println(count)

	sqlStatement = `
ALTER SEQUENCE users_id_seq RESTART WITH 1`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}

	fmt.Println("Reset the id")
}
