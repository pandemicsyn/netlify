package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/pandemicsyn/scratch/netlify/services/enrichment"
	"github.com/spf13/viper"
)

func configureLogging(v *viper.Viper) {
	level, err := log.ParseLevel(v.GetString("log_level"))
	if err != nil {
		log.Fatalln(err)
	}
	log.SetLevel(level)

	if v.GetString("log_format") == "text" {
		log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})
	} else if v.GetString("log_format") == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.Errorln("Error: log_type invalid, defaulting to text")
		log.SetFormatter(&log.TextFormatter{})
	}
	switch v.GetString("log_target") {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		log.Errorln("Error: log_target invalid, defaulting to Stdout")
		log.SetOutput(os.Stdout)
	}
}

func main() {

	v := viper.New()
	v.SetDefault("log_level", "info")
	v.SetDefault("log_format", "text")
	v.SetDefault("log_target", "stdout")
	v.SetEnvPrefix("churnapi")

	v.SetConfigName("churnapi")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/churnapi/")
	v.AddConfigPath("$HOME/.churnapi")
	v.ReadInConfig()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		project = "netlify-242319"
	}
	w, err := enrichment.New(project)
	if err != nil {
		log.Fatal(err)
	}
	if err := w.Receive(); err != nil {
		log.Fatal(err)
	}

	/*r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fixme"))
	})
	http.ListenAndServe(":3000", r)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Debugln(<-ch)
	//TODO: handle graceful shutdown
	log.Warnln("finished...exiting") */

}
