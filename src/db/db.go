package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func DBStart() *sql.DB {
	my_db, err := sql.Open("mysql", "root:1337@/anonstudchat")
	if err != nil {
		log.Panic(err)
	} else {
		err = my_db.Ping()
		if err != nil {
			log.Panic(err)
		}
	}
	return my_db
}
func FirstStart(user_id int, my_db *sql.DB) {
	stmtIns, err := my_db.Prepare("INSERT INTO users_info VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtIns.Exec(user_id,0)
	if err != nil {
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		panic(err.Error())
	}
}
