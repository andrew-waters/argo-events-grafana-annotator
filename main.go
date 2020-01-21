package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string `required:"true" envconfig:"MYSQL_HOST"`
	Port     int    `required:"true" envconfig:"MYSQL_PORT"`
	Username string `required:"true" envconfig:"MYSQL_USERNAME"`
	Password string `required:"true" envconfig:"MYSQL_PASSWORD"`
	Database string `required:"true" envconfig:"MYSQL_DATABASE"`
	Table    string `required:"true" envconfig:"MYSQL_TABLE"`
}

func timestamp() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05.999999")
}

func main() {

	var err error
	var c config

	err = envconfig.Process("notifier", &c)
	checkErr("Config Error", err)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", c.Username, c.Password, c.Host, c.Port, c.Database)

	db, err := sql.Open("mysql", dsn)
	checkErr("SQL Open", err)
	defer db.Close()

	t := timestamp()
	text := "Deployment happened"
	tags := "deployment"

	stmt, err := db.Prepare("INSERT annotations SET start=?, text=?, tags=?")
	checkErr("Prepare Statement", err)

	res, err := stmt.Exec(t, text, tags)
	checkErr("Exec Statement", err)

	id, err := res.LastInsertId()
	checkErr("Last Insert ID", err)

	log.Printf("Finished: %v", id)
}

func checkErr(prefix string, err error) {
	if err != nil {
		log.Fatal(fmt.Sprintf("%s:", prefix), err.Error())
		os.Exit(1)
	}
}
