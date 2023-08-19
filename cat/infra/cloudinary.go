package infra

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog/log"
)

func newCloudinary(c *config.Config) *cloudinary.Cloudinary {

	cl, err := cloudinary.NewFromParams(c.Cloudinary.CloudName, c.Cloudinary.ApiKey, c.Cloudinary.Secret)
	if err != nil {
		log.Fatal().Err(err).Msg("failed init cloudinary")
	}

	return cl
}
