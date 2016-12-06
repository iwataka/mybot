package main

import (
	"database/sql"
	"fmt"
)

var (
	db *sql.DB
)

func (c *DBConfig) insertImageAndResult(imageURL, result string) error {
	if c == nil {
		return nil
	}
	var err error
	if db == nil {
		driver := c.Driver
		dataSource := c.DataSource
		if driver == nil || dataSource == nil {
			return nil
		}
		db, err = sql.Open(*driver, *dataSource)
		if err != nil {
			return err
		}
	}
	table := c.VisionTable
	if table == nil {
		return nil
	}
	cmd := fmt.Sprintf("INSERT INTO %s VALUES ('%s', '%s')", *table, imageURL, result)
	_, err = db.Exec(cmd)
	if err != nil {
		return err
	}
	return nil
}
