package watch

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/cihub/seelog"
	L "github.com/vrecan/life"
)

var badLines []string

func init() {
	badLines = []string{
		"session opened for user root",
		"Invalid user",
		"POSSIBLE BREAK-IN ATTEMPT!",
		"Failed password for root",
		"Successful su for root",
		"useradd",
		"gpasswd",
	}
}

// Watch watches audit files and alerts if there are errors
type Watch struct {
	*L.Life
	Files     []string
	Ticker    *time.Ticker
	OutputDir string
	Year      int
	Month     int
	Day       int
}

// NewWatch creates a new watch object
func NewWatch(auditFiles []string, tickerDuration time.Duration, outputDir string, year, month, day int) (newWatch *Watch, err error) {
	newWatch = &Watch{
		Life:      L.NewLife(),
		Files:     auditFiles,
		Ticker:    time.NewTicker(tickerDuration),
		Year:      year,
		Month:     month,
		Day:       day,
		OutputDir: outputDir,
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
				if testLine(line) {
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
func testLine(line string) bool {
	for _, badLine := range badLines {
		if strings.Contains(line, badLine) {
			return true
		}
	}
	return false

}

// isLineInRange parses the date on the string and checks if it is in the range to check
func (w Watch) isLineInRange(line string) (inRange bool) {
	r, err := regexp.Compile("(^Jan|^Feb|^Mar|^May|^Jun|^Jul|^Aug|^Sep|^Oct|^Nov|^Dec)\\s+[0-9]{1,2}\\s+[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}")
	if err != nil {
		log.Error("Regex Does not compile: ", err)
		return
	}

	dateStrSlice := r.FindAllString(line, -1)
	if len(dateStrSlice) == 0 {
		log.Error("Date Failed to parse in regular expression: ", line)
		return
	}

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
	return false
}

func inTimeRange(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

// Close satisfies the io.Closer interface for Life and Death
func (w Watch) Close() error {
	return nil
}
