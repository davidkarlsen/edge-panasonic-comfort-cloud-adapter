package main_test

import (
	"testing"

	"github.com/futurehomeno/cliffhanger/test/suite"

	"github.com/futurehomeno/edge-panasonic-comfort-cloud-adapter/test"
)

func TestEnergyGuard(t *testing.T) {
	s := &suite.Suite{
		Cases: []*suite.Case{
			{
				Name:     "Configuration",
				Setup:    test.ServiceSetup("not_configured"),
				TearDown: []suite.Callback{test.TearDown("not_configured")},
				Nodes: []*suite.Node{
					{
						Name:    "Configure log level",
						Command: suite.StringMessage("pt:j1/mt:cmd/rt:ad/rn:panasonic_comfort_cloud/ad:1", "cmd.log.set_level", "panasonic_comfort_cloud", "warn"),
						Expectations: []*suite.Expectation{
							suite.ExpectString("pt:j1/mt:evt/rt:ad/rn:panasonic_comfort_cloud/ad:1", "evt.log.level_report", "panasonic_comfort_cloud", "warning"),
						},
					},
				},
			},
		},
	}

	s.Run(t)
}
