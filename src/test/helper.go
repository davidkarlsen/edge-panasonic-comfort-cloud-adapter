package test

import (
	"os"
	"path"
	"testing"

	"github.com/futurehomeno/cliffhanger/test/suite"

	"github.com/futurehomeno/edge-panasonic-comfort-cloud-adapter/cmd"
	"github.com/futurehomeno/edge-panasonic-comfort-cloud-adapter/internal/config"
)

func ServiceSetup(configSet string) suite.ServiceSetup {
	return func(t *testing.T) (service suite.Service, mocks []suite.Mock) {
		TearDown(configSet)(t)

		cfg := configSetup(t, configSet)

		app, err := cmd.Build(cfg)
		if err != nil {
			t.Fatalf("failed to build app: %s", err)
		}

		return app, nil
	}
}

func TearDown(configSet string) suite.Callback {
	return func(t *testing.T) {
		cmd.ResetContainer()

		err := os.RemoveAll(path.Join("../testdata/testing/", configSet, "/data/"))
		if err != nil {
			t.Fatalf("failed to clean up after previous tests: %s", err)
		}
	}
}

func configSetup(t *testing.T, configSet string) *config.Config {
	cfgSrv := config.NewConfigService(
		path.Join("../testdata/testing/", configSet),
	)

	cmd.SetConfigService(cfgSrv)

	err := cfgSrv.Load()
	if err != nil {
		t.Fatalf("failed to load configuration: %s", err)
	}

	return cfgSrv.Model().(*config.Config)
}
