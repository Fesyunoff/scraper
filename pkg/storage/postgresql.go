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
	Conf *types.Config
}

var _ ScraperStorage = (*PostgreScraperStorage)(nil)

func PreparePostgresDB(db *sql.DB, c *types.Config, sites []string) {

	req := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS  %s;`, c.SchemaName)
	result, err := db.Exec(req)
	CheckError(result, err, "Schema created")

	req = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s_temp (
	"service" VARCHAR(20) NOT NULL 
	);`, c.SchemaName, c.RespTableName)
	result, err = db.Exec(req)
	CheckError(result, err, "Temporary responce table created")

	req = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
	"id" SERIAL PRIMARY KEY,
	"service" VARCHAR(20) NOT NULL UNIQUE,
	"date" INTEGER,
	"responce" BOOLEAN,
	"status" SMALLINT,
	"duration" SMALLINT
	);`, c.SchemaName, c.RespTableName)
	result, err = db.Exec(req)
	CheckError(result, err, "Responce table created")

	result, err = prepareTable(db, sites, c)
	CheckError(result, err, "Responce table prepeared")

	req = fmt.Sprintf(`DROP TABLE IF EXISTS %s.%s_temp;`,
		c.SchemaName, c.RespTableName)
	result, err = db.Exec(req)
	CheckError(result, err, "Temporary responce table dropped")

	req = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
	"id" SERIAL PRIMARY KEY,
	"date" INTEGER,
	"service" VARCHAR(20),
	"user" SMALLINT
	);`, c.SchemaName, c.ReqTableName)
	result, err = db.Exec(req)
	CheckError(result, err, "Request table created")

	err = createUsersTable(db, c)
	CheckError(result, err, "User table created")
}

func WriteResponceToStorage(db *sql.DB, r types.Row) (err error) {
	req := fmt.Sprintf(`UPDATE %s.%s SET 
	date = %d, responce = %t, status = %d, duration = %d WHERE service = '%s';`,
		"availability", "responces",
		r.Date, r.Responce, r.StatusCode, r.Duration, r.Service)
	_, err = db.Exec(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return
}

// return true if admin
func (s *PostgreScraperStorage) DisplayServiceAvailability(db *sql.DB, site string) (r types.Row, err error) {
	req := fmt.Sprintf("SELECT * FROM %s.%s WHERE service = '%s';", "availability", "responces", site)
	row, err := db.Query(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&r.Id, &r.Service, &r.Date, &r.Responce, &r.StatusCode, &r.Duration)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	err = writeRequest(db, r.Service)
	if err != nil {
		log.Fatalln(err.Error())
	}
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
	row, err := db.Query(req)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&r.Id, &r.Service, &r.Date, &r.Responce, &r.StatusCode, &r.Duration)
		if err != nil {
			log.Fatal(err)
		}
	}

	return
}

// Dislay last day and 50 servises with max count requests if default
func (s *PostgreScraperStorage) DisplayStatistics(db *sql.DB, c *types.Config, hours int64, limit int) (out []types.Stat, err error) {
	if hours == 0 {
		hours = c.StatHoursBefore
	}
	date := time.Now().Unix() - hours*3600
	fmt.Println(date)
	if limit == 0 {
		limit = c.StatLimit
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
	return
}

func (s *PostgreScraperStorage) ReturnUsersRole(db *sql.DB, id int) (admin bool, err error) {
	req := fmt.Sprintf(`SELECT admin FROM %s.%s 
			WHERE id=%d;`,
		"availability", "users", id)
	row := db.QueryRow(req)

	err = row.Scan(&admin)
	if err != nil {
		log.Println(err)
	}
	return
}

func writeRequest(db *sql.DB, service string) (err error) {
	date := time.Now()
	req := fmt.Sprintf(`INSERT INTO %s.%s(date, service) VALUES (%d, '%s');`,
		"availability", "requests", date.Unix(), service)
	_, err = db.Exec(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return
}

func CheckError(result sql.Result, err error, msg string) {
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		rows, _ := result.RowsAffected()
		log.Printf("%s. %d rows affected.", msg, rows)
	}
}
func prepareTable(db *sql.DB, sites []string, c *types.Config) (result sql.Result, err error) {
	for _, site := range sites {
		if site == "" {
			continue
		}

		req := fmt.Sprintf(`INSERT INTO %s.%s_temp (service) VALUES ('%s');`,
			c.SchemaName, c.RespTableName, site)
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
		c.SchemaName, c.RespTableName)
	result, err = db.Exec(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return
}
func createUsersTable(db *sql.DB, c *types.Config) (err error) {

	req := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
	"id" SERIAL PRIMARY KEY,
	"admin" BOOLEAN
	);`, c.SchemaName, c.UsrTableName)
	_, err = db.Exec(req)
	if err != nil {
		log.Fatalln(err.Error())
	}

	req = fmt.Sprintf(`SELECT id FROM %s.%s 
											LIMIT 1;`, "availability", "users")
	row := db.QueryRow(req)
	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	if id == 0 {
		req := fmt.Sprintf(`INSERT INTO %s.%s (admin) VALUES 
			(true), (false), (false), (false)
			;`, c.SchemaName, c.UsrTableName)
		fmt.Println(req)
		_, err = db.Exec(req)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	return
}

func NewPostgreScrapeStorage(conn *sql.DB, conf *types.Config) *PostgreScraperStorage {
	return &PostgreScraperStorage{
		Conn: conn,
		Conf: conf,
	}
}
