package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/fesyunoff/availability/pkg/types"
	_ "github.com/lib/pq"
)

type PostgreScraperStorage struct {
	Conn *sql.DB
}

var _ ScraperStorage = (*PostgreScraperStorage)(nil)

func PreparePostgresDB(db *sql.DB, name string, sites []string) {

	req := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS  %s;`, "availability")
	result, err := db.Exec(req)
	CheckError(result, err, "Schema created")

	req = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s_temp (
	"service" VARCHAR(20) NOT NULL 
	);`, "availability", "responces")
	result, err = db.Exec(req)
	CheckError(result, err, "Temporary table 'responces' created")

	req = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
	"id" SERIAL PRIMARY KEY,
	"service" VARCHAR(20) NOT NULL UNIQUE,
	"date" INTEGER,
	"responce" BOOLEAN,
	"status" SMALLINT,
	"duration" SMALLINT
	);`, "availability", "responces")
	result, err = db.Exec(req)
	CheckError(result, err, "Table 'responces' created")

	result, err = prepareTable(db, sites)
	CheckError(result, err, "Table 'responces' prepeared")

	req = fmt.Sprintf(`DROP TABLE IF EXISTS %s.%s_temp;`, "availability", "responces")
	result, err = db.Exec(req)
	CheckError(result, err, "Temporary table 'responces' dropped")

	req = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
	"id" SERIAL PRIMARY KEY,
	"date" INTEGER,
	"service" VARCHAR(20),
	"user" SMALLINT
	);`, "availability", "requests")
	result, err = db.Exec(req)
	CheckError(result, err, "Table requests created")

}
func CheckError(result sql.Result, err error, msg string) {
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		rows, _ := result.RowsAffected()
		log.Printf("%s. %d rows affected.", msg, rows)
	}
}
func prepareTable(db *sql.DB, sites []string) (result sql.Result, err error) {
	for _, site := range sites {
		if site == "" {
			continue
		}

		req := fmt.Sprintf(`INSERT INTO %s.%s_temp (service) VALUES ('%s');`,
			"availability", "responces", site)
		_, err = db.Exec(req)
		if err != nil {
			log.Fatalln(err.Error())
		}

	}
	req := fmt.Sprintf(`INSERT INTO %[1]s.%[2]s (service)
						SELECT DISTINCT service FROM %[1]s.%[2]s_temp
						WHERE  service NOT IN (
							SELECT service 
							FROM %[1]s.%[2]s
						);`,
		"availability", "responces")
	result, err = db.Exec(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return
}

func WriteResponceToStorage(db *sql.DB, r types.Row) (err error) {
	req := fmt.Sprintf(`UPDATE %s.%s SET 
	date = %d, responce = %t, status = %d, duration = %d WHERE service = '%s';`,
		"availability", "responces",
		r.Date, r.Responce, r.StatusCode, r.Duration, r.Service)
	// fmt.Println(req)
	_, err = db.Exec(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return
}

func writeRequest(db *sql.DB, service string) (err error) {
	date := time.Now()
	req := fmt.Sprintf(`INSERT INTO %s.%s(date, service) VALUES (%d, '%s');`,
		"availability", "requests", date.Unix(), service)
	fmt.Println(req)
	_, err = db.Exec(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return
}

func (s *PostgreScraperStorage) DisplayServiceAvailability(db *sql.DB, site string) (r types.Row, err error) {
	req := fmt.Sprintf("SELECT * FROM %s.%s WHERE service = '%s';", "availability", "responces", site)
	row, err := db.Query(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&r.Service, &r.Date, &r.Responce, &r.StatusCode, &r.Duration)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	err = writeRequest(db, r.Service)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(r)
	return
}

func (s *PostgreScraperStorage) DisplayServiceResponceTime(db *sql.DB, min bool) (r types.Row, err error) {
	var hyphen string
	if min {
		hyphen = "--"
	}
	req := fmt.Sprintf(`SELECT * FROM %s.%s 
											WHERE duration != 0
											ORDER BY duration 
											%sDESC
											LIMIT 1;`, "availability", "responces", hyphen)
	fmt.Println(req)
	row, err := db.Query(req)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&r.Service, &r.Date, &r.Responce, &r.StatusCode, &r.Duration)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(r)
	}

	return
}

// last day and 50 servises with max count requests
func (s *PostgreScraperStorage) DisplayStatistics(db *sql.DB, hours int64, limit int) (out []types.Stat, err error) {
	if hours == 0 {
		hours = 24
	}
	date := time.Now().Unix() - hours*3600
	fmt.Println(date)
	if limit == 0 {
		limit = 50
	}
	req := fmt.Sprintf(`SELECT service, COUNT(service) as cs FROM  %s.%s 
			WHERE date > %d
			GROUP BY service 
			ORDER BY cs DESC 
			LIMIT %d
			;`,
		"availability", "requests",
		date, limit)
	row, err := db.Query(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer row.Close()
	for row.Next() {
		s := types.Stat{}
		err = row.Scan(&s.Service, &s.Count)
		if err != nil {
			log.Fatalln(err.Error())
		}
		out = append(out, s)
	}
	fmt.Println(out)

	return
}

func NewPostgreScrapeStorage(conn *sql.DB) *PostgreScraperStorage {
	return &PostgreScraperStorage{
		Conn: conn,
	}
}
