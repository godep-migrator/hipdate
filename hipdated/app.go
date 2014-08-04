package main

import (
	"github.com/3onyc/hipdate/backends"
	"github.com/3onyc/hipdate/shared"
	"github.com/3onyc/hipdate/sources"
	"log"
	"sync"
)

type Application struct {
	Backend     backends.Backend
	Sources     []sources.Source
	wg          *sync.WaitGroup
	EventStream chan *shared.ChangeEvent
	sc          chan bool
}

func NewApplication(
	b backends.Backend,
	s []sources.Source,
	cce chan *shared.ChangeEvent,
	wg *sync.WaitGroup,
	sc chan bool,
) *Application {
	return &Application{
		Backend:     b,
		Sources:     s,
		EventStream: cce,
		wg:          wg,
		sc:          sc,
	}
}

func (a *Application) EventListener() {
	for {
		select {
		case ce := <-a.EventStream:
			log.Printf("Event received %v\n", ce)
			u := shared.Upstream("http://" + ce.IP + ":80")
			switch ce.Type {
			case "add":
				a.Backend.AddUpstream(ce.Host, u)
			case "remove":
				a.Backend.RemoveUpstream(ce.Host, u)
			case "stop":
				break
			}
		case <-a.sc:
			return
		}
	}
}

func (a *Application) startEventListener() {
	defer a.wg.Done()

	a.EventListener()
	log.Println("[app] stopped")
}

func (a *Application) Start() {
	log.Printf("Starting sources")
	for _, s := range a.Sources {
		go s.Start()
	}

	log.Printf("Starting main process")

	a.wg.Add(1)
	go a.startEventListener()

	a.wg.Wait()
	log.Printf("Stopping cleanly via EOF")
}
