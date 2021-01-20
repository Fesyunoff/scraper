package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/fesyunoff/availability/pkg/types"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteScraperStorage struct {
	Conn *sql.DB
}

var _ ScraperStorage = (*SQLiteScraperStorage)(nil)

func CreateSQLiteDB(name string, sites []string) (sqliteDB *sql.DB) {
	path := fmt.Sprintf("./%s", name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Creating %s...", name)
		file, err := os.Create(name)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Printf("%s created", name)
	} else {
		log.Printf("%s exist", name)
	}

	sqliteDB, _ = sql.Open("sqlite3", path)

	// do not forget close connection at func call place
	// defer sqliteDB.Close()

	createTables(sqliteDB, sites)

	return
}

func createTables(db *sql.DB, sites []string) {
	createTableReq := `CREATE TABLE IF NOT EXISTS availability (
		--"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		--"service" TEXT,
		"service" TEXT NOT NULL PRIMARY KEY,
		"date" INTEGER,
		"responce" INTEGER,
		"status" INTEGER,
		"duration" INTEGER
	  );`
	createNewTable(db, "availability", sites, createTableReq)

	// createUserTableReq := `CREATE TABLE IF NOT EXISTS users (
	// "user_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	// "name" TEXT,
	// "role" TEXT
	// );`
	// createNewTable(db, "users", createUserTableReq)
}

func createTableIfNotExist(db *sql.DB, name string, sites []string, createReq string) {
	existReq := fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = '%s';", name)
	row := db.QueryRow(existReq)
	var a int
	err := row.Scan(&a)
	if err != nil {
		log.Println(err.Error())
	}
	if a == 0 {
		log.Printf("Create table %s...", name)
		statement, err := db.Prepare(createReq)
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatalln(err.Error())
			log.Printf("ERROR: table %s not created", name)
		} else {
			log.Printf("Table '%s' created", name)
			err = prepareTable(db, sites)
			if err != nil {
				log.Fatal(err.Error())
			} else {
				log.Printf("Table '%s' prepeared", name)
			}
		}

	} else {
		log.Printf("table '%s' exist", name)
	}
}

func createNewTable(db *sql.DB, name string, sites []string, createReq string) {
	req := `DROP TABLE IF EXISTS availability;`
	statement, err := db.Prepare(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Create table %s...", name)
	statement, err = db.Prepare(createReq)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err.Error())
		log.Printf("ERROR: table %s not created", name)
	} else {
		log.Printf("Table '%s' created", name)
		err = prepareTable(db, sites)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			log.Printf("Table '%s' prepeared", name)
		}
	}
}

func prepareTable(db *sql.DB, sites []string) (err error) {
	for _, site := range sites {
		if site == "" {
			continue
		}
		req := `INSERT INTO availability(service) VALUES (?);`
		statement, err := db.Prepare(req)
		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(site)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	return
}

func WriteResponceToStorage(db *sql.DB, r types.Row) (err error) {
	req := `UPDATE availability SET date = ?, responce = ?, status = ?, duration = ? WHERE service = ?;`
	statement, err := db.Prepare(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(r.Date, r.Responce, r.StatusCode, r.Duration, r.Service)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return
}

func (s *SQLiteScraperStorage) DisplayServiceAvailability(db *sql.DB, site string) (r types.Row, err error) {
	req := fmt.Sprintf("SELECT * FROM availability WHERE service = '%s';", site)
	// fmt.Println(req)
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

func (s *SQLiteScraperStorage) DisplayServiceResponceTime(db *sql.DB, min bool) (r types.Row, err error) {
	var hyphen string
	if min {
		hyphen = "--"
	}
	req := fmt.Sprintf(`SELECT * FROM  availability
											WHERE duration != 0
											ORDER BY duration 
											%sDESC
											LIMIT 1;`, hyphen)
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

func NewSQLiteScrapeStorage(conn *sql.DB) *SQLiteScraperStorage {
	return &SQLiteScraperStorage{
		Conn: conn,
	}
}
