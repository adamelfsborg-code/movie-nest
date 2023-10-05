package db

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

var Store *pg.DB

type Database struct {
	addr     string
	name     string
	user     string
	password string
}

func New(addr, name, user, password string) *Database {
	return &Database{
		addr:     addr,
		name:     name,
		user:     user,
		password: password,
	}
}

func (d *Database) CreateDatabase() {
	db := pg.Connect(&pg.Options{
		Addr:     d.addr,
		Database: d.name,
		User:     d.user,
		Password: d.password,
	})
	Store = db
}

type QueryLogger struct{}

func (*QueryLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	formattedQuery, err := q.FormattedQuery()
	if err != nil {
		return ctx, err
	}

	fmt.Println(string(formattedQuery))
	return ctx, nil
}

func (*QueryLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	return nil
}
