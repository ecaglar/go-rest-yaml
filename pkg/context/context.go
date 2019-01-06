/*
Package context defines an application context that can be passed as parameter as needed.
*/
package context

import (
	"../logger"
	"../memstore"
)

//AppContext defines pointers to storage and logger which
//all the packages use. Instead of passing all common attributes
//separately across calls, better to define a context and pass
//it around
type AppContext struct {
	Storage memstore.Storage
	Logger  *logger.AsyncLogger
}
