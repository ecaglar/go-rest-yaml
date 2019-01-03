//package memstore defines high level interface for in-memory storage operations
package memstore

type Database interface {

	//Insert adds a key-value object into the in-memory storage
	Insert(key string, val interface{})

	//Read gets related object stored with the given key
	Read(key string) interface{}

	//ReadWithParams performs search using given parameters
	ReadWithParams(params map[string][]string) []interface{}
}
