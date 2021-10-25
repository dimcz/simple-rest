package util

func Setup(config *Config) {
	jwtSecret = []byte(config.SigningKey)
}
