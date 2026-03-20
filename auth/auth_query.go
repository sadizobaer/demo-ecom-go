package auth

import (
	"context"
	"ecommerce/utilities"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TokenTableCreationQuery = `CREATE TABLE IF NOT EXISTS tokens (
	token_id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	token TEXT NOT NULL,
	refresh_token TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	);`

const UserTableCreationQuery = `CREATE TABLE IF NOT EXISTS users (	
	user_id SERIAL PRIMARY KEY,
	uuid UUID DEFAULT gen_random_uuid() UNIQUE,
	username VARCHAR(50) NOT NULL UNIQUE,
	email VARCHAR(100) NOT NULL UNIQUE,
	password TEXT NOT NULL,
	image_url TEXT,
	is_admin BOOLEAN DEFAULT FALSE,
	is_active BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const CreateUserQuery = `INSERT INTO users (username, email, password) 
VALUES ($1, $2, $3);`

const GetUserByEmailQuery = `SELECT user_id, username, email, password FROM users WHERE email = $1;`

const StoreTokenQuery = `INSERT INTO tokens (user_id, token, refresh_token) VALUES ($1, $2, $3);`

func RunUserTableCreationQuery(conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(), UserTableCreationQuery)
	if err != nil {
		return err
	}
	return nil
}

func RunTokenTableCreationQuery(conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(), TokenTableCreationQuery)
	if err != nil {
		return err
	}
	return nil
}

func CreateUserInDB(conn *pgxpool.Pool, user User) error {
	hashedPassword, err := utilities.HashPassword(user.Password)
	if err != nil {
		return err

	}

	_, err = conn.Exec(context.Background(), CreateUserQuery, user.Username, user.Email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(conn *pgxpool.Pool, email string) (*User, error) {
	row := conn.QueryRow(context.Background(), GetUserByEmailQuery, email)

	var user User
	err := row.Scan(&user.UserId, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func StoreTokenInDB(conn *pgxpool.Pool, userId int, token string, refreshToken string) error {
	_, err := conn.Exec(context.Background(), StoreTokenQuery, userId, token, refreshToken)
	if err != nil {
		return err
	}
	return nil
}
