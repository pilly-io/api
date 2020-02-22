package tests

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pilly-io/api/internal/db"
)

// LoadJSON returns list or map contained in file at `path`
func LoadJSON(path string) interface{} {
	byteValue := LoadFile(path)

	var result interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

// LoadFile returns the bytes contained in a file
func LoadFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	return byteValue
}

// SetupDB connect database and returns it
func SetupDB() db.Database {
	database := db.NewBeegoDatabase(os.Getenv("PILLY_DB_URI"))
	database.Migrate()
	return database
}
