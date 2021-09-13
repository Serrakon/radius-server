package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/theaaf/radius-server/app"
)

var serveRADIUSCmd = &cobra.Command{
	Use:   "serve",
	Short: "runs the RADIUS server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			logrus.Info("signal received; shutting down")
			cancel()
		}()

		sharedSecret, _ := cmd.Flags().GetString("shared-secret")

		postgres, err1 := cmd.Flags().GetString("postgres")
		table, err2 := cmd.Flags().GetString("table")
		loginColumn, err3 := cmd.Flags().GetString("login")
		passwordColumn, err4 := cmd.Flags().GetString("password")
		if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
			return app.ServeRADIUSPostrges(ctx, sharedSecret, postgres, table, loginColumn, passwordColumn)
		}

		redis, err := cmd.Flags().GetString("redis")

		if err == nil {
			return app.ServeRADIUSRedis(ctx, sharedSecret, redis)
		}
		return nil
	},
}

func init() {
	serveRADIUSCmd.Flags().String("shared-secret", "", "the shared secret to use for mutual authentication (required)")
	serveRADIUSCmd.MarkFlagRequired("shared-secret")

	serveRADIUSCmd.Flags().String("redis", "", "Redis: server to use for storage")

	serveRADIUSCmd.Flags().String("postgres", "", "PostgreSQL: server to use for storage")
	serveRADIUSCmd.Flags().String("table", "", "PostgreSQL: users table name")
	serveRADIUSCmd.Flags().String("login", "", "PostgreSQL: login table column")
	serveRADIUSCmd.Flags().String("password", "", "PostgreSQL: password table column")

	rootCmd.AddCommand(serveRADIUSCmd)
}
