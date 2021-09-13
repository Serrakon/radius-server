package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresCredentialProvider struct {
	Postgres       string
	Table          string
	LoginColumn    string
	PasswordColumn string
}

func (p *PostgresCredentialProvider) CredentialsForIdentity(id string) (EAPCredentials, error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, p.Postgres)

	v, err := p.getUserPassword(ctx, db, id)
	if err != nil {
		return nil, err
	}
	return plaintextCredential(v), nil
}

func (p *PostgresCredentialProvider) getUserPassword(ctx context.Context, db *pgxpool.Pool, username string) (string, error) {
	query := fmt.Sprintf(`SELECT %v AS password FROM %v WHERE %v=$1`, p.PasswordColumn, p.Table, p.LoginColumn)
	rows, err := db.Query(ctx, query, username)
	if err != nil {
		return "", err
	}
	password := ""
	for rows.Next() {
		err := rows.Scan(&password)
		if err == nil {
			return password, nil
		} else {
			return "", err
		}
	}
	return "", errors.New("user not found")
}
