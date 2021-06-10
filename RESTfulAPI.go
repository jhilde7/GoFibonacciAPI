package main

import (
	"log"
	"net/http"
	"fmt"
	"strconv"
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	host	 = "172.17.0.2"
	port	 = 5432
	user	 = "postgres"
	password = "password"
	dbname	 = "postgres"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=s% password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil{
		panic(err)
	}

	err = db.Ping()
	if err != nil{
		panic(err)
	}

	return db
}

func GETOrdinalHandler(w http.ResponseWriter, r *http.Request){
	val, err := strconv.Atoi(r.URL.Path[1:])
    if err != nil {
        fmt.Fprintf(w, "Supplied value %s is not a number", r.URL.Path[1:])
    } else {	
		ClearRecords()	
        fmt.Fprintf(w, "Fib(%s) == %d!", r.URL.Path[1:], FibSequencer(val))
    }
}

func GETUpToHandler(w http.ResponseWriter, r *http.Request){
	val, err := strconv.Atoi(r.URL.Path[8:])
	if err != nil {
        fmt.Fprintf(w, "Supplied value %s is not a number", r.URL.Path[1:])
    } else {

		db := OpenConnection()

		sqlStatement := "SELECT COUNT(RECORD) FROM Results WHERE RECORD < $1"
		_, err := db.Exec(sqlStatement, val)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		defer db.Close()
	}
}

func DELETEHandler(w http.ResponseWriter, r *http.Request){
	ClearRecords()
}

func FibSequencer(n int) int{
	if n <= 1 {
        return n
    }
	var fibonacci []int
	fibonacci = append(fibonacci, 0)
	fibonacci = append(fibonacci, 1)

	for i := 2; i <= n; i++ {
		fibonacci = append(fibonacci, fibonacci[i - 1] + fibonacci[i - 2])
		FibRecorder(fibonacci[i - 1] + fibonacci[i - 2])
	}

	return fibonacci[n]
}

func FibRecorder(val int) {
	db := OpenConnection()

	sqlStatement := "INSERT INTO Results (RECORD) VALUE($1)"
	_, err := db.Exec(sqlStatement, val)
	if err != nil {
		panic(err)
	}

	defer db.Close()
}

func ClearRecords(){
	db := OpenConnection()


	sqlStatement := "DELETE Results" 
	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}

	defer db.Close()
}

func main(){
	http.HandleFunc("/", GETOrdinalHandler)
	http.HandleFunc("/FibUpTo", GETUpToHandler)
	http.HandleFunc("/Clear", DELETEHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


