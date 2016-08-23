package watch

import (
	log "github.com/cihub/seelog"
	"github.com/fsnotify/fsnotify"
	L "github.com/vrecan/life"
)

// Watch watches audit files and alerts if there are errors
type Watch struct {
	*L.Life
	Watcher *fsnotify.Watcher
}

// NewWatch creates a new watch object
func NewWatch(auditFiles []string) (newWatch *Watch, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, file := range auditFiles {
		err = watcher.Add(file)
		if err != nil {
			log.Error("Error watching file ", file, ": ", err)
			return nil, err
		}
	}

	newWatch = &Watch{
		Life:    L.NewLife(),
		Watcher: watcher,
	}
	newWatch.SetRun(newWatch.run)
	return newWatch, err
}

// run starts watch
func (w Watch) run() {
	log.Info("Watcher running")

	for {
		select {
		case event := <-w.Watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Info("modified file:", event.Name)
			}
		case err := <-w.Watcher.Errors:
			log.Error("FSNotify Error:", err)
		case <-w.Life.Done:
			return
		}
	}
}

// Close satisfies the io.Closer interface for Life and Death
func (w Watch) Close() error {
	return nil
}
