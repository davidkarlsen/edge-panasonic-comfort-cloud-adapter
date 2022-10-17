package cmd

import (
	"github.com/futurehomeno/edge-panasonic-comfort-cloud-adapter/internal/config"
)

// ResetContainer resets service container for the testing purposes.
func ResetContainer() {
	services = &serviceContainer{}
}

// SetConfigService allows to inject config service into service container for the testing purposes.
func SetConfigService(cfgSrv *config.Service) {
	services.configService = cfgSrv
}
