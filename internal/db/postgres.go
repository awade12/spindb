package db

type PostgresConfig struct {
	Name     string
	User     string
	Password string
	Port     int
	Version  string
	Public   bool
}
