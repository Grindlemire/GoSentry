package c

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Conf configuration for Oracle
type Conf struct {
	Flagged   []string `yaml:"flaggedPhrases"`
	Year      int      `yaml:"yearsBack"`
	Month     int      `yaml:"monthsBack"`
	Day       int      `yaml:"daysBack"`
	Files     []string `yaml:"filesToWatch"`
	ScanEvery string   `yaml:"scanEvery"`
	OutputDir string   `yaml:"outputDir"`
	Regexs    []string `yaml:"timeRegexPatterns"`
}

// GetConf gets configuration from the yaml file
func GetConf(path string) (c Conf, err error) {
	data, err := ioutil.ReadFile(path)
	if nil != err {
		return c, err
	}
	err = yaml.Unmarshal([]byte(data), &c)
	return c, err
}
