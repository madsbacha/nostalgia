package config

import "log"

type Jwt struct {
	HmacSecret []byte
}

func NewJwtConfig(hmacSecret string) Jwt {
	if len([]byte(hmacSecret))*8 < 256 {
		log.Panicln("HMAC secret must be at least 256 bits long")
	}

	return Jwt{
		HmacSecret: []byte(hmacSecret),
	}
}
