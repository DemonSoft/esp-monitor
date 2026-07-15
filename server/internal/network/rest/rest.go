package rest

import (
	"remoteesp/internal/domain/utils"
	"remoteesp/internal/model"

	"github.com/joho/godotenv"
)

const defaultPort = "80"

type Mqtt interface {
	SendDeviceConfig(ssdp string, cfg model.ConfigPayload)
}

type Rest struct {
	cfg  Config
	db   Database
	mqtt Mqtt
}

type Config struct {
	Port string
}

type Response[T any] struct {
	Data  *T            `json:"data"`
	Error *RequestError `json:"error"`
}

type RequestError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Caution string `json:"caution"`
}

type Empty struct{}

func Create(c Config, db Database, mqtt Mqtt) *Rest {
	return &Rest{cfg: c, db: db, mqtt: mqtt}
}

func (r *Rest) Start() {
	router := r.Router()
	router.Run(":" + r.cfg.Port)
}

func CreateCfg() Config {
	_ = godotenv.Load() // Safety load for manifest using (without error generation)

	port := utils.GetEnv("REST_PORT", defaultPort)

	return Config{
		Port: port,
	}
}
