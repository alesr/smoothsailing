package users

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
	_ "github.com/lib/pq"
)

const (
	insertUserQuery string = `INSERT INTO users 
(id, first_name, last_name, email, birth_date, password_hash, created_at) 
VALUES ($1, $2, $3,$4,$5, $6, $7);`

	getAllUsersQuery string = "SELECT id, first_name, last_name, email, birth_date, created_at FROM users;"

	getUserByIDQuery string = `SELECT id, first_name, last_name, email, birth_date, created_at 
	FROM users WHERE id=$1 LIMIT 1;`

	deleteUserQuery string = "DELETE FROM users WHERE id=$1;"

	getUserByEmailQuery string = `SELECT id, first_name, last_name, email, birth_date, password_hash, created_at 
	FROM users WHERE email=$1 LIMIT 1;`
)

func storeUser(ctx context.Context, usr User) error {
	if _, err := sqldb.Exec(
		ctx,
		insertUserQuery,
		usr.ID,
		usr.FirstName,
		usr.LastName,
		usr.Email,
		usr.BirthDate,
		usr.PasswordHash,
		usr.CreatedAt,
	); err != nil {
		return fmt.Errorf("could not exec insert user: %s", err)
	}
	return nil
}

func deleteUserByID(ctx context.Context, id string) error {
	if _, err := sqldb.Exec(ctx, deleteUserQuery, id); err != nil {
		return fmt.Errorf("could not exec delete user: %s", err)
	}
	return nil
}

func getUserByID(ctx context.Context, id string) (*User, error) {
	var usr User

	if err := sqldb.QueryRow(ctx, getUserByIDQuery, id).Scan(
		&usr.ID,
		&usr.FirstName,
		&usr.LastName,
		&usr.Email,
		&usr.BirthDate,
		&usr.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("could not query get user by id: %w", err)
	}
	return &usr, nil
}

func getAllUsers(ctx context.Context) ([]User, error) {
	rows, err := sqldb.Query(ctx, getAllUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("could not query get all users: %s", err)
	}

	var usrs []User

	for rows.Next() {
		var usr User
		if err := rows.Scan(
			&usr.ID,
			&usr.FirstName,
			&usr.LastName,
			&usr.Email,
			&usr.BirthDate,
			&usr.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("could not scan user from rows: %s", err)
		}
		usrs = append(usrs, usr)
	}
	return usrs, nil
}

func getByEmail(ctx context.Context, email string) (*User, error) {
	var usr User

	if err := sqldb.QueryRow(ctx, getUserByEmailQuery, email).Scan(
		&usr.ID,
		&usr.FirstName,
		&usr.LastName,
		&usr.Email,
		&usr.BirthDate,
		&usr.PasswordHash,
		&usr.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("could not query get user by email: %w", err)
	}
	return &usr, nil
}
