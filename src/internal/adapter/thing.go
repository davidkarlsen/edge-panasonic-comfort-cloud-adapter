package adapter

import (
	"github.com/futurehomeno/cliffhanger/adapter"
	"github.com/futurehomeno/cliffhanger/adapter/service/numericsensor"
	"github.com/futurehomeno/cliffhanger/adapter/service/thermostat"
	"github.com/futurehomeno/cliffhanger/adapter/thing"
	"github.com/futurehomeno/edge-panasonic-comfort-cloud-adapter/ccontrol"
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"

	"github.com/futurehomeno/edge-panasonic-comfort-cloud-adapter/internal/config"
)

// NewThingFactory creates new instance of a thing factory.
func NewThingFactory(cfgSrv *config.Service) adapter.ThingFactory {
	return &thingFactory{
		cfgSrv: cfgSrv,
	}
}

// thingFactory is a private implementation of a thing factory service.
type thingFactory struct {
	cfgSrv *config.Service
}

// Create creates an instance of a thing using provided state.
func (f *thingFactory) Create(mqtt *fimpgo.MqttTransport, adapter adapter.ExtendedAdapter, thingState adapter.ThingState) (adapter.Thing, error) {
	thingConfig := &thing.ThermostatConfig{
		InclusionReport: &fimptype.ThingInclusionReport{},
		ThermostatConfig: &thermostat.Config{
			Specification: &fimptype.Service{
				Name: numericsensor.SensorTemp,
			},
			Controller: ccontrol.NewCloudControl(),
		},
		SensorTempConfig: &numericsensor.Config{
			Specification: &fimptype.Service{
				Name: numericsensor.SensorTemp,
			},
			Reporter: ccontrol.NewCloudControl(),
		},
	}
	return thing.NewThermostat(mqtt, thingConfig), nil
}
