package main

import (
	"strings"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

func watch(c config) error {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}

	const dir = "./"
	watcher.Add(dir)
	defer watcher.Close()

	interval := 1 * time.Second
	tick := time.Tick(interval)
	done := make(chan struct{}, 1)

	// a buffer for events
	var events []fsnotify.Event
	var loopErr error

	go func() {
	Loop:
		for {
			select {
			case event := <-watcher.Events:
				if !strings.HasSuffix(event.Name, ".go") {
					continue
				}
				if is(event, fsnotify.Create) || is(event, fsnotify.Write) {
					events = append(events, event)
				}
			case loopErr = <-watcher.Errors:
				done <- struct{}{}
				break Loop
			case <-tick:
				if len(events) == 0 {
					continue
				}
				// stop watching while gen'ing files
				watcher.Remove(dir)
				// gen the files
				run(c)
				// clear the buffer
				events = make([]fsnotify.Event, 0)
				// resume watching
				watcher.Add(dir)
			case <-done:
				break Loop
			}
		}
	}()

	<-done
	close(done)

	if loopErr != nil {
		return loopErr
	}

	return nil
}

func is(event fsnotify.Event, op fsnotify.Op) bool {
	return event.Op&op == op
}
