package drive

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ruraomsk/TLServer/logger"
)

func saveDBs() {
	dbs.Lock()
	defer dbs.Unlock()
	for _, db := range dbs.dbs {
		err := db.saveToFile()
		if err != nil {
			logger.Error.Printf("Сохранение БД %s %s", db.Name, err.Error())
		}
	}
}

func workerDBFs() {
	ticker := time.NewTicker(time.Duration(sbrams.Step * int(time.Second)))
	for {
		select {
		case <-ticker.C:
			saveDBs()
		case <-stop:
			saveDBs()
			return
		}
	}
}
func workerPGS() {
	<-stop
	pgsStop <- 1
}

func getListFilesDbs() []string {
	list := make([]string, 0)
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return list
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		if strings.HasSuffix(dir.Name(), ext) {

			list = append(list, strings.TrimSuffix(dir.Name(), ext))
		}
	}
	return list
}
func StartBrams(dbstop chan interface{}) error {
	stop = dbstop
	if sbrams.FS == JSON {
		for _, db := range getListFilesDbs() {
			if err := addDbFromJson(db); err != nil {
				return err
			}
		}
		go workerDBFs()
		return nil
	}
	if sbrams.FS == PgSQL {
		for _, name := range GetListDBs() {
			if err := addDB(name); err != nil {
				return err
			}
		}

		pgsChan = make(chan PgsRecord, 1000)
		pgsStop = make(chan interface{})
		go ListerStart(pgsChan, pgsStop)
		go workerPGS()
	}
	return fmt.Errorf("need coorect parameter fs now is %s", sbrams.FS)

}
func addDB(name string) error {
	return nil
}

// go drive.WorkerDB(set.Step, )
// go netcom.ServerCommections(set.Port, time.Duration(set.Step*int(time.Second)), stop)
