package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/spf13/cobra"
)

// For usage, see https://huma.rocks/features/cli/#passing-options
type Options struct {
	Debug bool   `doc:"Enable debug logging"`
	Host  string `doc:"Hostname to listen on."`
	Port  int    `doc:"Port to listen on." short:"p" default:"8888"`
}

func main() {
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("Vote API", "1.0.0"))

	// TODO: register handlers for api

	// cli & env parsing for high level config and commands
	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", opts.Port),
			Handler: router,
		}

		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", opts.Port)
			server.ListenAndServe()
		})

		hooks.OnStop(func() {
			// Graceful shutdown :)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	})

	cli.Root().Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := api.OpenAPI().YAML()
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		},
	})

	// TODO: register more commands
	// i.e run db migrations, running tests(?), (de)registering admins(?)

	// When no commands are passed, this starts the server!
	cli.Run()
}
