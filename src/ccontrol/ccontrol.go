package ccontrol

import (
	"context"
	"fmt"
	"github.com/futurehomeno/cliffhanger/adapter/service/numericsensor"
	"github.com/futurehomeno/cliffhanger/adapter/service/thermostat"
	cc "github.com/hacktobeer/go-panasonic/cloudcontrol"
	"github.com/hacktobeer/go-panasonic/types"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

const (
	ThermostatThingId = "panasonic-indoor-thermostat"
)

var instance *CloudControl

type CloudControl struct {
	cc *cc.Client
}

func NewCloudControl() *CloudControl {
	if instance == nil {
		client := cc.NewClient("")
		instance = &CloudControl{&client}
	}

	return instance
}

func (r *CloudControl) Login(ctx context.Context, username string, password string) error {
	_, err := r.cc.CreateSession(username, password)
	return err
}

func (r *CloudControl) ListDevices(ctx context.Context) []string {
	var devices []string
	groups, err := r.cc.GetGroups()
	if err == nil {
		devices = lo.Map(groups.Groups[0].Devices, func(device types.Device, index int) string {
			return device.DeviceGUID
		})
	}

	return devices
}

func (r *CloudControl) IsOnline() (bool, error) {
	device, err := r.cc.GetDeviceStatus()
	if err != nil {
		return false, err
	}

	return device.Parameters.Online, err
}

func (r *CloudControl) Init(ctx context.Context) error {
	groups, err := r.cc.GetGroups()
	if err != nil {
		return err
	}
	if groups.GroupCount == 1 {
		deviceGuid := groups.Groups[0].Devices[0].DeviceGUID
		r.cc.SetDevice(deviceGuid)
		log.Info("set active device", deviceGuid)
	}

	return nil
}

func (r *CloudControl) SetDevice(ctx context.Context, deviceGuid string) {
	r.cc.SetDevice(deviceGuid)
}

func (r *CloudControl) SetTemp(ctx context.Context, temp float64) error {
	_, err := r.cc.SetTemperature(temp)
	return err
}

/////controller

// SetThermostatMode sets a new thermostat mode.
func (r *CloudControl) SetThermostatMode(mode string) error {
	//
	return nil
}

// SetThermostatSetpoint sets a setpoint for a particular mode.
func (r *CloudControl) SetThermostatSetpoint(mode string, value float64, unit string) error {
	return nil
}

// ThermostatModeReport returns a current mode information.
func (r *CloudControl) ThermostatModeReport() (mode string, err error) {
	device, err := r.cc.GetDeviceStatus()
	if err != nil {
		return "", err
	}

	switch device.Parameters.Operate {
	case 0:
		mode = thermostat.ModeOff
	case 1:
		mode = "on"
	default:
		mode = "unknown"
	}

	return mode, err
}

// ThermostatSetpointReport returns a current setpoint for given mode.
func (r *CloudControl) ThermostatSetpointReport(mode string) (value float64, unit string, err error) {
	device, err := r.cc.GetDeviceStatus()
	if err != nil {
		return 0, thermostat.UnitC, err
	}

	return device.Parameters.InsideTemperature, thermostat.UnitC, err
}

// ThermostatStateReport returns a current state of the thermostat.
func (r *CloudControl) ThermostatStateReport() (string, error) {
	device, err := r.cc.GetDeviceStatus()
	if err != nil {
		return "", err
	}

	var state string
	switch device.Parameters.OperationMode {
	case 0:
		state = "auto"
	case 1:
		state = "dry"
	case 2:
		state = "cool"
	case 3:
		state = thermostat.StateHeat
	case 4:
		state = "dry"
	default:
		state = "unknown"
	}

	return state, nil
}

/// reporter

func (r *CloudControl) NumericSensorReport(unit string) (float64, error) {
	if unit == numericsensor.UnitC {
		device, err := r.cc.GetDeviceStatus()
		if err == nil {
			return device.Parameters.InsideTemperature, nil
		}

		return 0, err
	}

	return 0, fmt.Errorf("unkown unit: %v", unit)
}
