// Copyright (C) 2014 Adriano Soares
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
	defer db.Close()

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

	os.Chdir(filepath.Dir(*dbName))

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
			deleteFile(file)
		}
	}
}

func deleteFile(file string) {
	if *listOnly {
		return
	}

	err := os.Remove(file)
	if err != nil {
		fmt.Println("cannot remove", file, "-", err)
	}
}

func catch(err error, msg string) {
	if err != nil {
		fmt.Println("ERROR\t", msg)
		fmt.Println("\t", err)
		os.Exit(1)
	}
}
