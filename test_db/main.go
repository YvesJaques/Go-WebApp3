package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func getAllRowData(conn *sql.DB) error {
	rows, err := conn.Query("SELECT id, name, email FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	var id int
	var name string
	var email string
	for rows.Next() {
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("ID: %d, name: %s, email: %s\n", id, name, email)
	}

	if err != nil {
		log.Fatal("Error reading data", err)
	}
	return nil
}

func insertNewUser(conn *sql.DB, name string, email string, pw string, uType int) error {
	query := fmt.Sprintf(`INSERT INTO users (name, email, password, acct_created, last_login, user_type) VALUES
	('%s', '%s', '%s', current_timestamp, current_timestamp, %d)`, name, email, pw, uType)

	_, err := conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func getUserData(conn *sql.DB, id int) error {
	var name, email, pw, uType string

	query := fmt.Sprintf(`SELECT id, name, email, password, user_type FROM users WHERE id = %d`, id)

	row := conn.QueryRow(query)
	err := row.Scan(&id, &name, &email, &pw, &uType)

	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("ID: %d", id)
	fmt.Println("name: %s", name)
	fmt.Println("email: %s", email)
	fmt.Println("password: %s", pw)
	fmt.Println("user type: %d", uType)

	return nil
}

func main() {
	// Connect to postgres
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=blog_db user=postgres password=postgres")
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to connect to database: %v\n", err))
	}

	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to ping database: %v\n", err))
	}

	err = getAllRowData(conn)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to get all rows: %v\n", err))
	}

	fmt.Println("--------------")

	getUserData(conn, 1)
}
