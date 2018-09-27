package main;

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	host		= "127.0.0.1"
	port		= 5432
	user		= "todo_app"
	password	= "foo"
	dbname		= "todo"
)


func main() {
	var tdb todo_db
	tdb.init(host, port, user, password, dbname)
	defer tdb.destroy()	
//	tdb.reset()
	
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if s := strings.Split(r.URL.Path, "/")[1]; s != "" {
				fmt.Fprintf(w, "Insert " + s)
				tdb.insert(s)
			} else {
				fmt.Fprintf(w, "No item given")
			}
		case "DELETE":
			if s := strings.Split(r.URL.Path, "/")[1]; s != "" {
				fmt.Fprintf(w, "Delete " + s)
				tdb.delete(s)
			} else {
				tdb.clear()
			}
		case "GET":
			if s := strings.Split(r.URL.Path, "/")[1]; s != "" {
				fmt.Fprintf(w, tdb.display(s))
			} else {
				fmt.Fprintf(w, tdb.display_all())
			}
		case "PUT":
			fmt.Fprintf(w, "Put " + strings.Split(r.URL.Path, "/")[1])
		}
	})

	http.ListenAndServe(":9000", nil)
}
