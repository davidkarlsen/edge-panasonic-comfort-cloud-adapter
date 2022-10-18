package app

import (
	"context"
	"fmt"
	"github.com/futurehomeno/cliffhanger/adapter"
	"github.com/futurehomeno/cliffhanger/app"
	"github.com/futurehomeno/cliffhanger/lifecycle"
	"github.com/futurehomeno/cliffhanger/manifest"
	"github.com/futurehomeno/edge-panasonic-comfort-cloud-adapter/ccontrol"
	log "github.com/sirupsen/logrus"

	"github.com/futurehomeno/edge-panasonic-comfort-cloud-adapter/internal/config"
)

// TODO: Application is your central service which communicates with the frontend application.
//  You may want to extend it with additional interfaces such as: InitializableApp, CheckableApp, LogginableApp or AuthorizableApp.

// New creates new instance of an application.
func New(
	cfgSrv *config.Service,
	appLifecycle *lifecycle.Lifecycle,
	manifestLoader manifest.Loader,
	adapter adapter.ExtendedAdapter,
) app.App {
	return &application{
		cfgSrv:         cfgSrv,
		appLifecycle:   appLifecycle,
		manifestLoader: manifestLoader,
		adapter:        adapter,
		cc:             ccontrol.NewCloudControl(),
	}
}

// application is a private implementation of the main application service.
type application struct {
	cfgSrv         *config.Service
	appLifecycle   *lifecycle.Lifecycle
	manifestLoader manifest.Loader
	adapter        adapter.ExtendedAdapter
	cc             *ccontrol.CloudControl
}

// GetManifest returns the manifest object based on current application state and configuration.
func (a application) GetManifest() (*manifest.Manifest, error) {
	appManifest, err := a.manifestLoader.Load()
	if err != nil {
		log.WithError(err).Error("application: failed to load the template")

		return nil, fmt.Errorf("failed to load the template")
	}

	// TODO: You may want to manipulate the manifest depending on current application state or available configuration.
	//  Good examples include modifying list of available devices or dynamic options based on API calls and application lifecycle.

	return appManifest, nil
}

// Configure performs update of the application state based on the provided configuration.
func (a application) Configure(model interface{}) error {
	cfg, ok := model.(*config.Config)
	if !ok {
		log.Errorf("application: invalid config received, should be of %T type, received %T instead", cfg, model)

		return fmt.Errorf("received an invalid configuration")
	}

	// TODO: You may want persist here specific configuration settings provided by the user or act upon them.
	//  Good examples include adding or removing devices from an adapter.

	return nil
}

// Uninstall performs all required cleaning up before uninstalling the application.
func (a application) Uninstall() error {
	err := a.adapter.DestroyAllThings()
	if err != nil {
		log.WithError(err).Errorf("application: failed to destroy all things")

		return fmt.Errorf("failed to destroy all things")
	}

	err = a.cfgSrv.Reset()
	if err != nil {
		log.WithError(err).Errorf("application: failed to reset configuration")

		return fmt.Errorf("failed to reset configuration")
	}

	a.appLifecycle.SetAppState(lifecycle.AppStateNotConfigured, nil)
	a.appLifecycle.SetConfigState(lifecycle.ConfigStateNotConfigured)
	a.appLifecycle.SetConnectionState(lifecycle.ConnStateDisconnected)
	a.appLifecycle.SetAuthState(lifecycle.AuthStateNotAuthenticated)

	return nil
}

func (a application) createThing() {
	err := a.adapter.CreateThing(ccontrol.ThermostatThingId, nil)
	if err != nil {
		log.Error(err)
	}
}

func (a application) Login(credentials *app.LoginCredentials) error {
	a.appLifecycle.SetAuthState(lifecycle.AuthStateInProgress)
	err := a.cc.Login(context.TODO(), credentials.Username, credentials.Password)
	if err != nil {
		a.appLifecycle.SetAuthState(lifecycle.AuthStateNotAuthenticated)
		return err
	}
	a.appLifecycle.SetAuthState(lifecycle.AuthStateAuthenticated)

	a.appLifecycle.SetConnectionState(lifecycle.ConnStateConnecting)
	err = a.cc.Init(context.TODO())
	if err != nil {
		a.appLifecycle.SetConnectionState(lifecycle.ConnStateDisconnected)
	} else {
		a.appLifecycle.SetConnectionState(lifecycle.ConnStateConnected)
	}

	a.createThing()

	return err
}

func (a application) Logout() error {
	log.Info("Logout, NOOP")
	return nil
}

func (a application) Check() error {
	online, err := a.cc.IsOnline()
	if online {
		a.appLifecycle.SetConnectionState(lifecycle.ConnStateConnected)
	} else {
		a.appLifecycle.SetConnectionState(lifecycle.ConnStateDisconnected)
	}

	return err
}
