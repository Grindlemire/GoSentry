package start

import (
	"io"
	"os"
	"regexp"
	"strconv"
	SYS "syscall"
	"time"

	log "github.com/cihub/seelog"
	CONF "github.com/grindlemire/GoSentry/c"
	"github.com/grindlemire/GoSentry/watch"
	"github.com/vrecan/death"
)

// Run Entry Point for GoSentry
func Run() {

	c, err := CONF.GetConf("c.yml")
	if nil != err {
		log.Critical(err)
		os.Exit(2)
	}

	log.Info("Starting GoSentry Service")

	var goRoutines []io.Closer
	death := death.NewDeath(SYS.SIGINT, SYS.SIGTERM)

	duration, err := parseDuration(c.ScanEvery)
	if err != nil {
		log.Error("Error parsing duration: ", err)
	}

	watch, err := watch.NewWatch(c.Files, duration, c.OutputDir, c.Year, c.Month, c.Day)
	if err != nil {
		log.Error("Unable to create Watcher: ", err)
		return
	}
	watch.Start()
	goRoutines = append(goRoutines, watch)

	death.WaitForDeath(goRoutines...)

}

func parseDuration(scanStr string) (duration time.Duration, err error) {
	r, err := regexp.Compile("(?P<number>[0-9])+\\s*(?P<unit>[s|m|d|M])")
	if err != nil {
		log.Error("Regex Does not compile: ", err)
		return
	}

	matches := r.FindStringSubmatch(scanStr)
	if len(matches) != 3 {
		return 0, log.Error("Incorrect parsing of duration string ", scanStr)
	}

	durationNum, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, log.Error("Number not in duration: ", matches[1])
	}
	durationUnit := matches[2]

	switch durationUnit {
	case "s":
		return time.Duration(durationNum) * time.Second, nil
	case "m":
		return time.Duration(durationNum) * time.Minute, nil
	case "d":
		return time.Duration(durationNum) * 24 * time.Hour, nil
	case "M":
		return time.Duration(durationNum) * 24 * 30 * time.Hour, nil
	default:
		return 0, log.Error("Error Parsing Time unit for time: ", scanStr)
	}

}
