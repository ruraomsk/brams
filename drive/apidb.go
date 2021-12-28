package drive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/brams/wal"
)

//Open открытие базы данных
//Открывает базу по имени и возвращает указатель на нее
//Если базы данных нет то возвращает ошибку
func Open(name string) (*Db, error) {
	dbs.RLock()
	defer dbs.RUnlock()
	if db, ok := dbs.dbs[name]; ok {
		return db, nil
	}
	return nil, fmt.Errorf("need create db %s", name)
}

//Drop удаление базы данных
func Drop(name string) {

	dbs.Lock()
	defer dbs.Unlock()
	db, ok := dbs.dbs[name]
	if !ok {
		return
	}
	delete(dbs.dbs, name)
	if db.fs {
		fname := path + name + ext
		_ = os.Remove(fname)
	}
}

//AddDb добавляет бд в пул бд
func AddDb(name string) error {
	dbs.Lock()
	defer dbs.Unlock()
	if _, ok := dbs.dbs[name]; ok {
		return fmt.Errorf("db %s is exist ", name)
	}
	fname := path + name + ext
	_, err := os.Stat(fname)
	if err != nil {
		return err
	}
	db := new(Db)
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &db)
	if err != nil {
		return err
	}
	db.fs = true
	dbs.dbs[name] = db
	return nil
}

//CreateDb cоздает бд и присваивает описание ключа
// где defkey массив имен переменных из value json
func CreateDb(name string, defkeys ...string) error {
	dbs.Lock()
	defer dbs.Unlock()
	if _, ok := dbs.dbs[name]; ok {
		return fmt.Errorf("db %s is exist ", name)
	}
	if len(defkeys) == 0 {
		return ErrWrongParameters
	}
	db := new(Db)
	db.Name = name
	db.Defkey = make([]string, 0)
	db.Defkey = append(db.Defkey, defkeys...)
	db.Values = make(map[string]*Value)
	db.UID = 0
	db.fs = true
	fname := path + name + ext
	_, err := os.Stat(fname)
	if err == nil {
		return fmt.Errorf("db file %s is exist the path %s", name, path)
	}
	buf, err := json.Marshal(&db)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fname, buf, os.FileMode(0644))
	if err != nil {
		return err
	}
	dbs.dbs[name] = db
	return nil
}
func CreateDbInMemory(name string, defkeys ...string) error {
	dbs.Lock()
	defer dbs.Unlock()
	if _, ok := dbs.dbs[name]; ok {
		return fmt.Errorf("db %s is exist ", name)
	}
	if len(defkeys) == 0 {
		return ErrWrongParameters
	}
	db := new(Db)
	db.Name = name
	db.Defkey = make([]string, 0)
	db.Defkey = append(db.Defkey, defkeys...)
	db.Values = make(map[string]*Value)
	db.fs = false
	db.UID = 0
	dbs.dbs[name] = db
	return nil
}
func GetListFilesDbs() []string {
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
func WorkerDB(step int, stop chan interface{}) {
	ticker := time.NewTicker(time.Duration(step * int(time.Second)))
	walChan = make(chan wal.WalRecord, 1000)
	go wal.ListerWalStart(walChan)
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
func saveDBs() {
	dbs.Lock()
	defer dbs.Unlock()
	for _, db := range dbs.dbs {
		err := db.SaveToFile()
		if err != nil {
			logger.Error.Printf("Сохранение БД %s %s", db.Name, err.Error())
		}
	}
}
func (db *Db) WriteRecord(value []byte) error {
	v := new(Value)
	v.Value = value
	fullkey, err := db.makeFullKeyOnValue(value)
	if err != nil {
		return err
	}
	db.RWMutex.RLock()
	old, is := db.Values[fullkey]
	var op byte
	if !is {
		//Insert
		op = wal.Insert
		db.UID++
		v.UID = db.UID
	} else {
		//Replace
		op = wal.Replace
		v.UID = old.UID
	}
	db.RWMutex.RUnlock()
	db.RWMutex.Lock()
	db.Values[fullkey] = v
	db.update = true
	db.RWMutex.Unlock()
	db.walOperation(op, fullkey, v)
	return nil
}
func (db *Db) DeleteRecord(keys ...interface{}) error {
	if len(keys) != len(db.Defkey) {
		return ErrWrongParameters
	}
	full, err := db.makeFullKey(keys)
	if err != nil {
		return err
	}
	db.RWMutex.RLock()
	value, is := db.Values[full]
	db.RWMutex.RUnlock()
	if !is {
		return ErrKeyNotFound
	}
	db.RWMutex.Lock()
	delete(db.Values, full)
	db.RWMutex.Unlock()
	db.walOperation(wal.Delete, full, value)
	return nil
}

func (db *Db) ReadRecord(keys ...interface{}) ([]byte, error) {
	if len(keys) != len(db.Defkey) {
		return make([]byte, 0), ErrWrongParameters
	}
	full, err := db.makeFullKey(keys)
	if err != nil {
		return make([]byte, 0), err
	}
	db.RWMutex.RLock()
	defer db.RWMutex.RUnlock()
	value, is := db.Values[full]
	if !is {
		return make([]byte, 0), ErrKeyNotFound
	}
	return value.Value, nil
}
func (db *Db) ReadListKeys(limit int, keys ...interface{}) ([]string, error) {
	db.RWMutex.RLock()
	defer db.RWMutex.RUnlock()
	if len(keys) > len(db.Defkey) {
		return make([]string, 0), ErrWrongParameters
	}
	return db.makeListKeys(limit, keys)

}
func (db *Db) ReadRecordFromList(key string) ([]byte, error) {
	db.RWMutex.RLock()
	defer db.RWMutex.RUnlock()
	value, is := db.Values[key]
	if !is {
		return make([]byte, 0), ErrKeyNotFound
	}
	return value.Value, nil
}
