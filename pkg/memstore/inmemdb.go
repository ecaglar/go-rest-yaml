package memstore

import (
	"../logger"
	"../model"
	"strings"
)

//dedicated logger for storage operations
var db_logger *logger.AsyncLogger

//underlying structure to record key-value object
//value can be any type
type memDB struct {
	keyValDB map[string]interface{}
}

//CreateInMemDB creates the underlying storage and logger
func CreateInMemDB() *memDB {
	if db_logger != nil {
		defer db_logger.Log(logger.INFO, "In-Memory memstore has been created")
	}
	return &memDB{make(map[string]interface{})}
}

//Insert inserts given key val pair into the storage
func (db *memDB) Insert(key string, val interface{}) {
	db.keyValDB[key] = val
	if db_logger != nil {
		db_logger.Log(logger.INFO, "Value has been inserted to in-memory memstore with key: ", key)
	}
}

func (db *memDB) SetLogger(logger *logger.AsyncLogger) {
	db_logger = logger
}

//ReadWithParams queries the storage for objects match the given url query strings
func (db *memDB) ReadWithParams(params map[string][]string) []interface{} {

	var res []interface{}

	//If there is no search criteria then return all records
	//TODO: Paging should be done here for performance issues
	if len(params) == 0 {
		for _, model := range db.keyValDB {
			res = append(res, model)
		}
		return res
	}

	//If there is version in query then use it as key.
	//if there is version then there is no need to check other parameters as well
	if i, ok := params["version"]; ok {
		return append(res, db.Read(i[0]))
	}

	//check all records which match given query string
	for _, model := range db.keyValDB {
		if checkModelWithParams(model, params) {
			res = append(res, model)
		}
	}
	return res
}

//Read reads a record with the given key
func (db *memDB) Read(key string) interface{} {
	if val, ok := db.keyValDB[key]; ok {
		return val
	}
	return nil
}

//checkModelWithParams compares given object with the query string in case there is a match
//for query params title and description, it is checked if it is contained
//for other fields full string match is expected
//if title is "App v1.0.0" then a query with title=app will match
//if description is "This is a description for app" then description=for%20app will match
func checkModelWithParams(data interface{}, urlQuerystr map[string][]string) bool {

	metadata, ok := data.(model.Metadata)

	if ok {

		for param, value := range urlQuerystr {
			if (param != "maintainers.name" && param != "maintainers.email") && len(value) > 1 {
				return false
			}
			switch queryParam := strings.TrimSpace(param); queryParam {

			case "title":
				if !strings.Contains(metadata.Title, value[0]) {
					return false
				}
			case "description":
				if !strings.Contains(metadata.Description, value[0]) {
					return false
				}
			case "company":
				if metadata.Company != value[0] {
					return false
				}
			case "website":
				if metadata.Website != value[0] {
					return false
				}
			case "source":
				if metadata.Source != value[0] {
					return false
				}
			case "license":
				if metadata.License != value[0] {
					return false
				}
			case "maintainers.name":
				for _, name := range value {
					found := false
					for _, metadaMaintainer := range metadata.Maintainers {
						if metadaMaintainer.Name == name {
							found = true
						}
					}
					if found == false {
						return false
					}
				}
			case "maintainers.email":
				for _, email := range value {
					found := false
					for _, metadaMaintainer := range metadata.Maintainers {
						if metadaMaintainer.Email == email {
							found = true
						}
					}
					if found == false {
						return false
					}
				}
			}
		}
	}
	return true
}
