package drive

import (
	"database/sql"

	"github.com/ruraomsk/TLServer/logger"

	_ "github.com/lib/pq"
)

type PgsRecord struct {
	Name      string `json:"name"`
	UID       uint64 `json:"uid"`
	Operation byte   `json:"op"`
	Value     []byte `json:"value"`
}
type DbData struct {
	name string
	conn *sql.Conn
}

var (
	mdb *sql.DB
)

const (
	Replace = iota
	Insert
	Delete
)

func ListerStart(in chan PgsRecord, stop chan interface{}) {
	for {
		select {
		case rec := <-in:
			logger.Info.Printf("in->%d %d", rec.UID, rec.Operation)
		case <-stop:

		}
	}
}
