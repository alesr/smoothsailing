package users

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
	_ "github.com/lib/pq"
)

const insertUserQuery string = `INSERT INTO users 
(id,first_name,last_name,email,birth_date,created_at) 
VALUES ($1, $2, $3,$4,$5, $6)`

func Insert(ctx context.Context, usr User) error {
	if _, err := sqldb.Exec(
		ctx,
		insertUserQuery,
		usr.id,
		usr.firstName,
		usr.lastName,
		usr.email,
		usr.birthDate,
		usr.createdAt,
	); err != nil {
		return fmt.Errorf("could not insert user: %s", err)
	}
	return nil
}
