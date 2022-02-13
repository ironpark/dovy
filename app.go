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
	return &App{
		dovey: nil,
	}
}

func (a *App) SendChatMessage(msg string) {
	a.dovey.SendChatMessage(msg)
}

func (a *App) OpenAuthorization() {
	a.dovey.OpenAuthorization()
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	dovey := dovy.NewDovey(ctx)
	go dovey.Serve()
	a.dovey = dovey
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
