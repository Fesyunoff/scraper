package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fesyunoff/availability/pkg/storage"
	"github.com/fesyunoff/availability/pkg/types"
	"github.com/kelseyhightower/envconfig"
)

var (
	debug    = flag.Bool("debug", false, "debug")
	env      = flag.String("env", "local", "environment")
	fileName = flag.String("file", "sites.txt", "name of file with servises list")
)

func main() {

	flag.Parse()
	var c types.Config
	err := envconfig.Process(*env, &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.Time = 60    //second
	c.Timeout = 10 //second
	c.Debug = *debug
	fmt.Printf("%+v\n", c)
	sites := getSitesFromFile(*fileName)
	connStatement := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.HostDB, c.PortDB, c.UserDB, c.PasswordDB, c.NameDB)
	connDB, err := sql.Open("postgres", connStatement)
	if err != nil {
		panic(err)
	}
	defer connDB.Close()
	storage.PreparePostgresDB(connDB, &c, sites)

	ch := make(chan types.Row, 1)
	for _, site := range sites {

		go func(site string) {
			for {
				resp, row := testSite(site, &c)
				if resp {
					ch <- row
				}
				time.Sleep(time.Duration(c.Time) * time.Second)
			}
		}(site)
	}
	for {
		row := <-ch
		err := storage.WriteResponceToStorage(connDB, &c, row)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
	}

}
func testSite(site string, c *types.Config) (bool, types.Row) {
	row := types.Row{}
	if site == "" {
		return false, row
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

		return true, row
	}
	time := time.Since(start).Milliseconds()
	_, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		row.Responce = true
		row.Duration = time
		row.StatusCode = resp.StatusCode
	}

	defer resp.Body.Close()
	return true, row
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
