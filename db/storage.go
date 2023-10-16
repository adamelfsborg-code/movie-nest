package db

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

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
