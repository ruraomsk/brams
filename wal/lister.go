package wal

import "github.com/ruraomsk/TLServer/logger"

type WalRecord struct {
	Name      string `json:"name"`
	UID       uint64 `json:"uid"`
	Operation byte   `json:"op"`
	Key       string `json:"key"`
	Value     []byte `json:"value"`
}

const (
	Replace = iota
	Insert
	Delete
)

func ListerWalStart(in chan WalRecord) {
	for {
		rec := <-in
		logger.Info.Printf("in->%d %d", rec.UID, rec.Operation)
	}
}
