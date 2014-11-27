package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var fileList = []string{
	"pics/*.jpg",
	"pics/*.JPG",
	"pics/field/*.png",
	"pics/field/*.PNG",
	"pics/thumbnail/*.jpg",
	"pics/thumbnail/*.JPG",
	"script/c[0-9]*.lua",
	"script/c[0-9]*.LUA",
}

var dbName = flag.String("db", "cards.cdb", "name of the database file")
var listOnly = flag.Bool("list", false, "list the files without deleting")

func main() {
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbName)
	catch(err, "could not open cards.cdb")

	rows, err := db.Query("select id from datas")
	catch(err, "could not load cards.cdb")

	ids := make(map[string]bool)
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		catch(err, "column datas.id is bad formated")
		ids[id] = true
	}
	catch(rows.Err(), "database is corrupted")

	files := findFiles(fileList)
	deleteFiles(files, ids)

	fmt.Println("Done.")
}

func findFiles(patterns []string) []string {
	var list []string
	for _, pattern := range patterns {
		f, err := filepath.Glob(pattern)
		catch(err, "cannot list files")
		list = append(list, f...)
	}
	return list
}

func deleteFiles(files []string, ids map[string]bool) {
	for _, file := range files {
		cardID := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		cardID = strings.TrimPrefix(cardID, "c") // for scripts
		if !ids[cardID] {
			fmt.Println("unused", file)
			if !*listOnly {
				err := os.Remove(file)
				if err != nil {
					fmt.Println("cannot remove", file, "-", err)
				}
			}
		}
	}
}

func catch(err error, detail string) {
	if err != nil {
		fmt.Println("ERROR\t", detail)
		fmt.Println("\t", err)
		os.Exit(1)
	}
}
