package main

import (
	"github.com/nobleach/cobaltcms/internal/api"
	"github.com/nobleach/cobaltcms/internal/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// var k = koanf.New(".")

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// if err := k.Load(file.Provider("config/config.toml"), toml.Parser()); err != nil {
	// 	log.Fatal().Err(err)
	// } else {
	// 	log.Info().Msg("Loaded config")
	// }
	//
	store, err := storage.NewPostgresStore()

	if err != nil {
		log.Error().Err(err).Msg("Could not set up a connection to the data store")
	}

	// port := ":" + strconv.Itoa(k.Int("server.port"))
	port := ":8080"

	apiServer := api.NewApiServer(port, store)
	log.Info().Msgf("CobalCMS server is listening on port %s", port)
	apiServer.Run()
}
