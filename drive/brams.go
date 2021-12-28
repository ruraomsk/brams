package drive

import (
	"context"
	"os"
	"sync"
)

var (
	dbs struct {
		sync.RWMutex
		dbs map[string]*Db
	}
	mutex = &sync.RWMutex{}
)

type Db struct {
	sync.RWMutex
	name       string
	fk         *os.File
	fv         *os.File
	fs         *os.File
	vals       map[string]*Value
	cancelSync context.CancelFunc
	storemode  int
	update     bool
}
type Value struct {
	val []byte
}

func init() {
	dbs.dbs = make(map[string]*Db)

}
