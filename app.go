package main

import (
	"context"
	"dovey/dovy"
	"fmt"
)

// App struct
type App struct {
	ctx   context.Context
	dovey *dovy.Dovy
}

// NewApp creates a new App application struct
func NewApp() *App {
	do, _ := dovy.NewDovey()
	return &App{
		dovey: do,
	}
}

func (a *App) SendChatMessage(channel, msg string) {
	a.dovey.SendChatMessage(channel, msg)
}

func (a *App) OpenAuthorization() {
	a.dovey.OpenAuthorization()
}

func (a *App) Connect(channelName string) {
	a.dovey.Connect(channelName)
	return
}

func (a *App) IsAuthorized() bool {
	return a.dovey.IsAuthorized()
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.dovey.SetAppContext(ctx)
	a.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
