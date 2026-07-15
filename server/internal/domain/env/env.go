package env

import (
	"remoteesp/internal/domain/utils"

	"github.com/joho/godotenv"
)

type Env struct{}

var XToken string /// Security token of REST requests.

const (
	xTokenDefault = ""

	// ENVIRONTMENT FIELDS
	XTOKEN = "XTOKEN"
)

func Create() Env {
	return Env{}
}

func Apply() {
	env := Create()
	env.ReadCommandLine()
	env.LoadEnvirontments()
	//env.Output()
}

func (e *Env) LoadEnvirontments() {
	_ = godotenv.Load() // Safety load for manifest using (without error generation)
	XToken = utils.GetEnv("XTOKEN", XToken)
}
