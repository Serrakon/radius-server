package app

import (
	"context"
)

func ServeRADIUSRedis(ctx context.Context, sharedSecret, redis string) error {
	radiusServer := &RADIUSServer{
		SharedSecret: []byte(sharedSecret),
		CredentialProvider: &RedisCredentialProvider{
			Redis: redis,
		},
	}
	if err := radiusServer.Start(); err != nil {
		return err
	}

	<-ctx.Done()
	radiusServer.Stop()
	return ctx.Err()
}

func ServeRADIUSPostrges(ctx context.Context, sharedSecret, postgres string, table string, loginColumn string, passwordColumn string) error {
	radiusServer := &RADIUSServer{
		SharedSecret: []byte(sharedSecret),
		CredentialProvider: &PostgresCredentialProvider{
			Postgres:       postgres,
			Table:          table,
			LoginColumn:    loginColumn,
			PasswordColumn: passwordColumn,
		},
	}
	if err := radiusServer.Start(); err != nil {
		return err
	}

	<-ctx.Done()
	radiusServer.Stop()
	return ctx.Err()
}
