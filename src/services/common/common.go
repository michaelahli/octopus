package common

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type CommonSvc struct {
	config *viper.Viper
}

type Common interface {
	CommonHandler(w http.ResponseWriter, r *http.Request)
}

func New(config *viper.Viper) *CommonSvc {
	return &CommonSvc{config: config}
}

func (s *CommonSvc) CommonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from %s %s ", s.config.GetString("runtime.environment"), r.URL.Path)
}
