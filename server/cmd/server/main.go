package main

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/nobleach/cobaltcms/internal/api"
	"github.com/nobleach/cobaltcms/internal/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var k = koanf.New(".")

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	if err := k.Load(file.Provider("config/config.toml"), toml.Parser()); err != nil {
		log.Fatal().Err(err)
	} else {
		log.Info().Msg("Loaded config")
	}

	store, err := storage.NewPostgresStore()

	if err != nil {
		log.Error().Err(err).Msg("Could not set up a connection to the data store")
	}

	port := k.String("server.port")

	apiServer := api.NewApiServer(k, store)
	log.Info().Msgf("CobalCMS server is listening on port %s", port)
	apiServer.Run()
}
