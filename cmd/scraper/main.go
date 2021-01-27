package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/fesyunoff/availability/pkg/api/scraper"
	"github.com/fesyunoff/availability/pkg/api/transport"
	"github.com/fesyunoff/availability/pkg/storage"
	"github.com/fesyunoff/availability/pkg/types"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
)

var (
	fileName = flag.String("file", "sites.txt", "name of file with servises list")
)

func main() {

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	// logger = kitlog.With(logger, "service", serviceName)
	logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.Caller(5))

	flag.Parse()
	c := &types.Config{
		Host:            "0.0.0.0",
		Port:            8991,
		FileName:        *fileName,
		Timeout:         10, //second
		Time:            60, //second
		HostDB:          "172.18.0.2",
		PortDB:          5432,
		UserDB:          "user",
		PasswordDB:      "pass",
		NameDB:          "postgres",
		SchemaName:      "availability",
		RespTableName:   "responces",
		ReqTableName:    "requests",
		UsrTableName:    "users",
		StatHoursBefore: 24,
		StatLimit:       50,
	}

	sites := getSitesFromFile(c.FileName)
	connStatement := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.HostDB, c.PortDB, c.UserDB, c.PasswordDB, c.NameDB)
	connDB, err := sql.Open("postgres", connStatement)
	if err != nil {
		panic(err)
	}
	defer connDB.Close()
	strg := storage.NewPostgreScrapeStorage(connDB, c)
	storage.PreparePostgresDB(connDB, c, sites)
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

	ch := make(chan int, 1)
	for _, site := range sites {

		go func(site string) {
			for {
				testSite(site, c, connDB)
				time.Sleep(time.Duration(c.Time) * time.Second)
			}
			ch <- 1
		}(site)
	}
	<-ch

	<-sigint

}

func testSite(site string, c *types.Config, db *sql.DB) bool {
	row := types.Row{}
	if site == "" {
		return false
	}
	row.Service = site
	client := http.Client{
		Timeout: time.Duration(c.Timeout) * time.Second,
	}

	if c.Proxy != "" {
		proxyUrl, _ := url.Parse(c.Proxy)
		transport := http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		client = http.Client{Transport: &transport}
	}

	site = "http://" + site
	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	start := time.Now()
	row.Date = start.Unix()
	resp, err := client.Do(req)
	if err != nil {
		row.Responce = false
		err = storage.WriteResponceToStorage(db, row)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
		return false
	}
	time := time.Since(start).Milliseconds()
	_, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		row.Responce = true
		row.Duration = time
		row.StatusCode = resp.StatusCode
		_ = storage.WriteResponceToStorage(db, row)
	}

	defer resp.Body.Close()
	return true
}

func getSitesFromFile(name string) (sites []string) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(0)
	}
	sites = strings.Split(string(data), "\n")
	return
}

func exitOnError(l kitlog.Logger, err error, msg string) {
	if err != nil {
		l.Log("err", err, "msg", msg)
		os.Exit(1)
	}
}
