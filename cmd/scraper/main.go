package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/fesyunoff/availability/pkg/api/scraper"
	"github.com/fesyunoff/availability/pkg/api/transport"
	"github.com/fesyunoff/availability/pkg/storage"
	"github.com/fesyunoff/availability/pkg/types"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

var (
	debug = flag.Bool("debug", false, "debug")
	env   = flag.String("env", "local", "environment")
)

func main() {

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "service", "scraper")
	logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.Caller(5))

	flag.Parse()
	var c types.Config
	err := envconfig.Process(*env, &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.StatLimit = 50
	c.StatHoursBefore = 24
	c.Debug = *debug
	fmt.Printf("%+v\n", c)
	connStatement := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.HostDB, c.PortDB, c.UserDB, c.PasswordDB, c.NameDB)
	connDB, err := sql.Open("postgres", connStatement)
	if err != nil {
		panic(err)
	}
	defer connDB.Close()
	strg := storage.NewPostgreScrapeStorage(connDB, &c)
	svc := scraper.NewScraper(strg)
	svcHandlers, err := transport.MakeHandlerREST(svc)

	bindAddr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	r := mux.NewRouter().StrictSlash(true)

	exitOnError(logger, err, "failed create handlers")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	r.PathPrefix("/").Handler(svcHandlers)

	ln, err := net.Listen("tcp", bindAddr)
	exitOnError(logger, err, "failed create listener")
	defer ln.Close()

	_ = level.Info(logger).Log("msg", "server listen on "+ln.Addr().String())

	go func() {
		_ = http.Serve(ln, r)
	}()
	<-sigint

}

func exitOnError(l kitlog.Logger, err error, msg string) {
	if err != nil {
		l.Log("err", err, "msg", msg)
		os.Exit(1)
	}
}
