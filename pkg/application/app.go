package application

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	Ctx context.Context
	Wg  sync.WaitGroup
}

func NewApp() (*App, <-chan struct{}) {
	appCtx, appShutdown := context.WithCancel(context.Background())
	app := &App{
		Ctx: appCtx,
		Wg:  sync.WaitGroup{},
	}
	appDone := make(chan struct{})

	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		<-done

		appShutdown()
		app.Wg.Wait()
		close(appDone)
	}()

	return app, appDone
}
