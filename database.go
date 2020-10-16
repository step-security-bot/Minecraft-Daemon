package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const serverTable = "servers"

const maxNameLength = 64
const maxDirectoryLength = 256
const maxJarFileLength = 64
const maxRunMemoryLength = 8
const maxStartMemoryLength = 8
const maxJavaArgsLength = 256
const maxNameLengthString = "64"
const maxDirectoryLengthString = "256"
const maxJarFileLengthString = "64"
const maxRunMemoryLengthString = "8"
const maxStartMemoryLengthString = "8"
const maxJavaArgsLengthString = "256"

const serverTableSchema = "id INT NOT NULL AUTO_INCREMENT, " +
	"name VARCHAR(" + maxNameLengthString + ") NOT NULL, " +
	"directory VARCHAR(" + maxDirectoryLengthString + ") NOT NULL, " +
	"jarfile VARCHAR(" + maxJarFileLengthString + ") NOT NULL, " +
	"runmemory VARCHAR(" + maxRunMemoryLengthString + ") NOT NULL, " +
	"startmemory VARCHAR(" + maxStartMemoryLengthString + ") NOT NULL, " +
	"javaargs VARCHAR(" + maxJavaArgsLengthString + "), " +
	"PRIMARY KEY (id)"

// id is omitted since it is generated by the database
const serverTableColumns = "name,directory,jarfile,runmemory,startmemory,javaargs"

type databaseServer struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Directory   string         `json:"directory"`
	JarFile     string         `json:"jar_file"`
	RunMemory   string         `json:"run_memory"`
	StartMemory string         `json:"start_memory"`
	JavaArgs    sql.NullString `json:"java_args"`
}

type responseServer struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	JarFile     string `json:"jar_file"`
	RunMemory   string `json:"run_memory"`
	StartMemory string `json:"start_memory"`
	JavaArgs    string `json:"java_args"`
}

var db *sqlx.DB

func databaseSetup() {
	db, err := sqlx.Open("mysql", config.Database.DatabaseUser+":"+config.Database.DatabasePassword+
		"@("+config.Database.DatabaseHost+":"+fmt.Sprint(config.Database.DatabasePort)+")"+"/")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE " + config.Database.DatabaseDatabase)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully created database..")
	}

	_, err = db.Exec("USE " + config.Database.DatabaseDatabase)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("CREATE TABLE " + serverTable + "(" + serverTableSchema + ")")
}

func connectDatabase() {
	tempdb, err := sqlx.Open("mysql", config.Database.DatabaseUser+":"+config.Database.DatabasePassword+
		"@("+config.Database.DatabaseHost+":"+fmt.Sprint(config.Database.DatabasePort)+")"+
		"/"+config.Database.DatabaseDatabase)
	if err != nil {
		panic(err)
	}

	db = tempdb

	// TODO: How to make this extend to life of daemon and still be closed properly
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func getSingleServerData(id int) (*responseServer, error) {
	result, err := db.Queryx("SELECT * FROM " + serverTable + " WHERE `id` = " + fmt.Sprint(id))
	if err != nil {
		return nil, err
	}

	defer result.Close()

	var temp databaseServer
	for result.Next() {
		if err := result.StructScan(&temp); err != nil {
			return nil, err
		}
	}

	server := responseServer{
		temp.ID,
		temp.Name,
		temp.Directory,
		temp.JarFile,
		temp.RunMemory,
		temp.StartMemory,
		temp.JavaArgs.String,
	}

	return &server, nil
}

func requestServerToString(server requestServer) string {
	if len(server.JavaArgs) > 0 {
		return "\"" + server.Name + "\",\"" + server.Directory + "\",\"" + server.JarFile + "\",\"" + server.RunMemory + "\",\"" + server.StartMemory + "\",\"" + server.JavaArgs + "\""
	}

	return "\"" + server.Name + "\",\"" + server.Directory + "\",\"" + server.JarFile + "\",\"" + server.RunMemory + "\",\"" + server.StartMemory + "\",null"
}

func addServerToDatabase(server requestServer) {
	_, err := db.Exec("INSERT INTO " + serverTable + "(" + serverTableColumns + ") " +
		"VALUES (" + requestServerToString(server) + ")")
	if err != nil {
		fmt.Println(err)
	}
}

func collectServerData() []responseServer {
	result, err := db.Queryx("SELECT * FROM " + serverTable)
	if err != nil {
		fmt.Println(err)
	}

	defer result.Close()

	var serverList []responseServer
	for result.Next() {
		var temp databaseServer
		if err := result.StructScan(&temp); err != nil {
			return nil
		}
		server := responseServer{
			temp.ID,
			temp.Name,
			temp.Directory,
			temp.JarFile,
			temp.RunMemory,
			temp.StartMemory,
			temp.JavaArgs.String,
		}
		serverList = append(serverList, server)
	}

	return serverList
}

func checkForDuplicateServer(name, directory string) (bool, error) {
	result, err := db.Queryx("SELECT * FROM " + serverTable + " WHERE `name` = \"" + name + "\"")
	if err != nil {
		return true, err
	}

	defer result.Close()

	var temp databaseServer
	for result.Next() {
		if err := result.StructScan(&temp); err != nil {
			return true, err
		}
	}

	if len(temp.Name) > 0 {
		return true, nil
	}

	result, err = db.Queryx("SELECT * FROM " + serverTable + " WHERE `directory` = \"" + directory + "\"")
	if err != nil {
		return true, err
	}

	defer result.Close()

	for result.Next() {
		if err := result.StructScan(&temp); err != nil {
			return true, err
		}
	}

	if len(temp.Name) > 0 {
		return true, nil
	}

	return false, nil
}
