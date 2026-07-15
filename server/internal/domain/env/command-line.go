/// Command-line service parameters.
/// Пакет предназyначен для получения конфигурации сервиса параметрами командной строки.

package env

import (
	"flag"
	"fmt"
)

func (e *Env) ReadCommandLine() {
	flag.StringVar(&XToken, "xToken", xTokenDefault, "xToken uses for security api requests")

	flag.Parse()
}

func (e *Env) Output() {
	fmt.Println("-xToken:", XToken)
}
