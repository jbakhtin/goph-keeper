package config

type DB struct {
	DSN    string `env:"DATABASE_DSN"`
	Driver string `env:"DATABASE_DRIVER"`
}

func (db *DB) SetDataBaseDSN(value string) {
	db.DSN = value
}

func (db DB) GetDataBaseDSN() string {
	return db.DSN
}

func (db *DB) SetDataBaseDriver(value string) {
	db.Driver = value
}

func (db DB) GetDataBaseDriver() string {
	return db.Driver
}
