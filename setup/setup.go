package setup

type SetupBrams struct {
	DbPath    string    `toml:"dbpath"`
	LogPath   string    `toml:"logpath"`
	Step      int       `toml:"step"`
	FS        string    `toml:"fs"`
	TCPServer TCPServer `toml:"tcp"`
	Chan      bool      `toml:"chan"`
	PgDB      PgDB      `toml:"PGdB"`
}

type TCPServer struct {
	Start   bool `toml:"start"`
	Timeout int  `toml:"timeout"`
	Port    int  `toml:"port"`
}

//DataBase настройки базы данных postresql
type PgDB struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBname   string `toml:"dbname"`
}

var Set SetupBrams
