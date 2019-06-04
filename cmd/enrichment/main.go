package main

import (
	"database/sql"
	"os"

	"github.com/pandemicsyn/netlify/services/enrichment"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var db *sql.DB

func main() {

	v := viper.New()
	v.SetDefault("db", "postgres://postgres@localhost:32768/postgres?sslmode=disable")
	v.SetEnvPrefix("enrichment")
	v.BindEnv("db")

	v.SetConfigName("enrichment")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/enrichment/")
	v.AddConfigPath("$HOME/.enrichment")
	v.ReadInConfig()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		project = "netlify-242319"
	}

	var err error
	db, err = sql.Open("postgres", v.GetString("db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println(PrepareDB(db, false))

	w, err := enrichment.New(project, db)
	if err != nil {
		log.Fatal(err)
	}
	if err := w.Receive(); err != nil {
		log.Fatal(err)
	}
}
