package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	select_hash *sql.Stmt
	insert_acc  *sql.Stmt
	update_hash *sql.Stmt

	db *sql.DB
)

func openSqlDB(filename string) (err error) {
	db, err = sql.Open("sqlite3", filename)
	if err != nil {
		return
	}

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS accounts (
	id 		PRIMARY_KEY INTEGER,
	email 	STRING UNIQUE,
	hash 	BLOB
);
`)

	return
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

	return
}
