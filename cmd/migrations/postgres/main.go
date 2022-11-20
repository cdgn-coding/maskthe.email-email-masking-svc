package main

import (
	"email-masking-svc/src/infrastructure/postgresql"
)

func main() {
	migrations := postgresql.NewMigrations()
	migrations.Apply()
}
