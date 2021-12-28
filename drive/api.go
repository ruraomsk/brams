package drive

import "fmt"

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

//AddDb добавляет бд в пул бд
func AddDb(name string) error {

	return nil
}

//CreateDb cоздает бд и присваивает описание ключа
// где defkey массив имен переменных из value json
func CreateDb(name string, defkey []string) error {

	return nil
}
