package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

func main() {
	migration := flag.String("f", "", "migration file name")
	dbType := flag.String("t", "postgres", "db type")
	dbUser := flag.String("u", "postgres", "db user")
	dbPass := flag.String("p", "", "db pass")
	dbHost := flag.String("h", "localhost", "db host")
	dbPort := flag.Int("h", 5432, "db name")
	dbName := flag.String("n", "", "db name")

	flag.Parse()

	fmt.Println(*dbName)
	connStr := fmt.Sprintf(`%s://%s:%s@%s:%d/%s?sslmode=disable`, *dbType, *dbUser, *dbPass, *dbHost, *dbPort, *dbName)
	db, err := sql.Open(*dbType, connStr)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	file, err := os.ReadFile(*migration)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(file))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Migration OK")
}
