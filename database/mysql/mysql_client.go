package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

var (
	mu      sync.RWMutex
	clients map[string]*sql.DB = make(map[string]*sql.DB)
)

func GetClient(dsn, dbname string) (*sql.DB, error) {
	mu.RLock()
	defer mu.RUnlock()
	if c, ok := clients[dbname]; ok {
		if err := c.Ping(); err != nil {
			return nil, err
		}
		return c, nil
	}
	dsnR := fmt.Sprintf(dsn, dbname)
	db, err := sql.Open("mysql", dsnR)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	clients[dbname] = db
	return db, nil
}
