package watch

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/cihub/seelog"
	CONF "github.com/grindlemire/GoSentry/c"
	L "github.com/vrecan/life"
)

// Watch watches audit files and alerts if there are errors
type Watch struct {
	*L.Life
	Files     []string
	Ticker    *time.Ticker
	OutputDir string
	Year      int
	Month     int
	Day       int
	Regexs    []*regexp.Regexp
	BadLines  []string
}

// NewWatch creates a new watch object
func NewWatch(c CONF.Conf) (newWatch *Watch, err error) {
	regexs := []*regexp.Regexp{}
	for _, rS := range c.Regexs {
		r, err := regexp.Compile(rS)
		if err != nil {
			return nil, log.Error("Regex Does not compile: ", err)
		}
		regexs = append(regexs, r)
	}

	duration, err := parseDuration(c.ScanEvery)
	if err != nil {
		return nil, log.Error("Error parsing scanEvery duration: ", err)
	}

	newWatch = &Watch{
		Life:      L.NewLife(),
		Files:     c.Files,
		Ticker:    time.NewTicker(duration),
		Year:      c.Year,
		Month:     c.Month,
		Day:       c.Day,
		OutputDir: c.OutputDir,
		Regexs:    regexs,
		BadLines:  c.Flagged,
	}

	newWatch.SetRun(newWatch.run)
	return newWatch, err

}

// run starts watch
func (w Watch) run() {
	log.Info("Watcher running")

	for {
		select {
		case <-w.Ticker.C:
			w.ReadFiles()
		case <-w.Life.Done:
			return
		}
	}
}

// ReadFiles reads all the specified files and parses them
func (w Watch) ReadFiles() {
	for _, fileStr := range w.Files {
		file, err := os.Open(fileStr)
		if err != nil {
			log.Error("Error opening file ", fileStr, ": ", err)
			continue
		}
		reader := bufio.NewReader(file)

		flagFound := 0
		flaggedLines := []string{}
		var line string

	readLoop:
		for {
			line, err = reader.ReadString('\n')
			if err != nil {
				break readLoop
			}

			inRange := w.isLineInRange(line)
			if inRange {
				if w.testLine(line) {
					flagFound++
					flaggedLines = append(flaggedLines, line)
				}
			}
		}

		totalName, err := writeFile(flaggedLines, w.OutputDir, file.Name())
		if err != nil {
			log.Error("Error Writing to output file: ", err)
			continue
		}

		if flagFound > 0 {
			log.Infof("Found %d flagged Phrases in %v. Wrote phrases to %v", flagFound, file.Name(), totalName)
		}
		file.Close() //Do this because defers in for loops is bad

	}
}

// writeFile writes flagged events to an output file
func writeFile(lines []string, outputDir, filePath string) (totalName string, err error) {
	if len(lines) == 0 {
		return
	}
	currTime := time.Now().Format("2006.01.02.15")
	fileName := filepath.Base(filePath)
	basePath := filepath.Dir(filePath)
	totalPath := outputDir + basePath
	totalName = totalPath + "/" + fileName + "." + currTime

	err = os.MkdirAll(totalPath, 0777)
	if err != nil {
		return totalName, log.Error("Error creating output dir: ", err)
	}

	file, err := os.Create(totalName)
	if err != nil {
		return totalName, log.Error("Error Opening OutputFile: ", err)

	}
	defer file.Close()

	for _, line := range lines {
		file.WriteString(line)
	}
	return totalName, nil
}

// testLine tests if any flagged strings are in the line
func (w Watch) testLine(line string) bool {
	for _, badLine := range w.BadLines {
		if strings.Contains(line, badLine) {
			return true
		}
	}
	return false

}

// isLineInRange parses the date on the string and checks if it is in the range to check
func (w Watch) isLineInRange(line string) (inRange bool) {
	found := false
	for _, r := range w.Regexs {

		dateStrSlice := r.FindAllString(line, -1)
		if len(dateStrSlice) == 0 {
			continue
		}
		found = true

		date, err := time.Parse("Jan 2 15:04:05", dateStrSlice[0])
		date = date.AddDate(time.Now().Year(), 0, 0) //Note this will break on new years every year
		if err != nil {
			log.Error("Error Parsing string to time in file: ", dateStrSlice[0])
		}

		end := time.Now()
		start := time.Now().AddDate(-w.Year, -w.Month, -w.Day)

		if inTimeRange(start, end, date) {
			return true
		}
	}
	if !found {
		log.Error("Unable to parse line to any regular expression")
	}
	return false
}

// inTimeRange tests whether the parsed time is within the range specified
func inTimeRange(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

// parseDuration parses the scanEvery option into a time
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
	case "h":
		return time.Duration(durationNum) * time.Hour, nil
	case "d":
		return time.Duration(durationNum) * 24 * time.Hour, nil
	case "M":
		return time.Duration(durationNum) * 24 * 30 * time.Hour, nil
	default:
		return 0, log.Error("Error Parsing Time unit for time: ", scanStr)
	}

}

// Close satisfies the io.Closer interface for Life and Death
func (w Watch) Close() error {
	return nil
}
