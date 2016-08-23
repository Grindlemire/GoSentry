package watch

import (
	log "github.com/cihub/seelog"
	L "github.com/vrecan/life"
)

// Watch watches audit files and alerts if there are errors
type Watch struct {
	*L.Life
}

// NewWatch creates a new watch object
func NewWatch() (newWatch *Watch) {
	newWatch = &Watch{
		Life: L.NewLife(),
	}
	newWatch.SetRun(newWatch.run)
	return newWatch
}

// run starts watch
func (w Watch) run() {
	log.Info("Watcher running")
	for {
		select {
		case <-w.Life.Done:
			return
		}
	}
}

// Close satisfies the io.Closer interface for Life and Death
func (w Watch) Close() error {
	return nil
}
