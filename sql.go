package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	select_hash  *sql.Stmt
	insert_acc   *sql.Stmt
	update_hash  *sql.Stmt
	insert_admin *sql.Stmt
	delete_admin *sql.Stmt
	delete_acc   *sql.Stmt

	db *sql.DB
)

func openSqlDB(filename string) (err error) {
	db, err = sql.Open("sqlite3", filename)
	if err != nil {
		return
	}

	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS accounts (
	id 		PRIMARY_KEY INTEGER,
	email 	STRING UNIQUE,
	hash 	BLOB
);

CREATE TABLE IF NOT EXISTS admin (
	id 		PRIMARY_KEY INTEGER,
	email	STRING UNIQUE
);
`)

	if err != nil {
		return
	}

	rows, err := db.Query(`SELECT email FROM admin;`)

	if err != nil {
		return
	}

	defer rows.Close()

	sessions = map[string]Session{}

	for rows.Next() {
		email := ""
		if err = rows.Scan(&email); err != nil {
			return
		}
		sessions[email] = Session{
			"",
			time.Unix(0, 0),
			false,
		}
	}

	return rows.Err()
}

func sqlPrepareStmts() (err error) {
	select_hash, err = db.Prepare(`
SELECT hash FROM accounts WHERE email=?;
`)

	if err != nil {
		return
	}

	insert_acc, err = db.Prepare(`
INSERT INTO accounts (email,hash) VALUES(?,?);
`)

	if err != nil {
		return
	}

	update_hash, err = db.Prepare(`
UPDATE accounts SET hash=? WHERE email=?;
`)

	if err != nil {
		return
	}

	insert_admin, err = db.Prepare(`
INSERT INTO admin (email) VALUES(?);
`)
	if err != nil {
		return
	}

	delete_admin, err = db.Prepare(`
DELETE FROM admin WHERE email=?;
`)

	if err != nil {
		return
	}

	delete_acc, err = db.Prepare(`
DELETE FROM accounts WHERE email=?;
`)

	return
}

func SQLClose() {
	if err := db.Close(); err != nil {
		Log("Error closing DB:", err.Error())
	}
	if err := select_hash.Close(); err != nil {
		Log("Error closing prepared statement:", err.Error())
	}
	if err := insert_acc.Close(); err != nil {
		Log("Error closing prepared statement:", err.Error())
	}
	if err := update_hash.Close(); err != nil {
		Log("Error closing prepared statement:", err.Error())
	}
	if err := insert_admin.Close(); err != nil {
		Log("Error closing prepared statement:", err.Error())
	}
	if err := delete_admin.Close(); err != nil {
		Log("Error closing prepared statement:", err.Error())
	}
	if err := delete_acc.Close(); err != nil {
		Log("Error closing prepared statement:", err.Error())
	}
}
