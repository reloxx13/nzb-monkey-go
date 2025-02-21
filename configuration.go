package main

import (
	"sort"
	"strconv"

	"gopkg.in/ini.v1"
)

// configuration structure
type Configuration struct {
	General struct {
		Target            string `ini:"target"`
		Categorize        string `ini:"categorize"`
		Success_wait_time int    `ini:"success_wait_time"`
		Error_wait_time   int    `ini:"error_wait_time"`
		Debug             bool   `ini:"debug"`
	} `ini:"GENERAL"`
	Execute struct {
		Passtofile      bool   `ini:"passtofile"`
		Passtoclipboard bool   `ini:"passtoclipboard"`
		Nzbsavepath     string `ini:"nzbsavepath"`
		Dontexecute     bool   `ini:"dontexecute"`
		CleanUpEnable   bool   `ini:"clean_up_enable"`
		CleanUpMaxAge   int    `ini:"clean_up_max_age"`
	} `ini:"EXECUTE"`
	Sabnzbd struct {
		Host              string `ini:"host"`
		Port              int    `ini:"port"`
		Ssl               bool   `ini:"ssl"`
		SkipCheck         bool   `ini:"skip_check"`
		Nzbkey            string `ini:"nzbkey"`
		BasicauthUsername string `ini:"basicauth_username"`
		BasicauthPassword string `ini:"basicauth_password"`
		Basepath          string `ini:"basepath"`
		Category          string `ini:"category"`
		Addpaused         bool   `ini:"addpaused"`
	} `ini:"SABNZBD"`
	Nzbget struct {
		Host              string `ini:"host"`
		Port              int    `ini:"port"`
		Ssl               bool   `ini:"ssl"`
		SkipCheck         bool   `ini:"skip_check"`
		BasicauthUsername string `ini:"user"`
		BasicauthPassword string `ini:"pass"`
		Basepath          string `ini:"basepath"`
		Category          string `ini:"category"`
		Addpaused         bool   `ini:"addpaused"`
	} `ini:"NZBGET"`
	Synologyds struct {
		Host              string `ini:"host"`
		Port              int    `ini:"port"`
		Ssl               bool   `ini:"ssl"`
		SkipCheck         bool   `ini:"skip_check"`
		Username          string `ini:"user"`
		Password          string `ini:"pass"`
		BasicauthUsername string
		BasicauthPassword string
		Basepath          string `ini:"basepath"`
	} `ini:"SYNOLOGYDLS"`
	Nzbcheck struct {
		SkipFailed                bool    `ini:"skip_failed"`
		MaxMissingSegmentsPercent float64 `ini:"max_missing_segments_percent"`
		MaxMissingFiles           int     `ini:"max_missing_files"`
		BestNZB                   bool    `ini:"best_nzb"`
	} `ini:"NZBCheck"`
	Categories    map[string]string `ini:"-"` // will hold the categories regex patterns
	Searchengines []string          `ini:"-"` // will hold the search engines
	Directsearch  struct {
		Host           string `ini:"host"`
		Port           int    `ini:"port"`
		SSL            bool   `ini:"ssl"`
		Username       string `ini:"username"`
		Password       string `ini:"password"`
		Connections    int    `ini:"connections"`
		Hours          int    `ini:"hours"`
		ForwardHours   int    `ini:"forward_hours"`
		Step           int    `ini:"step"`
		Scans          int    `ini:"scans"`
		Skip           bool   `ini:"skip"`
		FirstGroupOnly bool   `ini:"first_group_only"`
	} `ini:"DIRECTSEARCH"`
}

// global configuration variable
var (
	conf Configuration
)

func loadConfig() {

	conf = Configuration{}

	iniOption := ini.LoadOptions{
		IgnoreInlineComment: true,
	}
	cfg, err := ini.LoadSources(iniOption, confPath)
	if err != nil {
		Log.Error("Unable to load configuration file '%s': %s", confPath, err.Error())
		exit(1)
	}

	err = cfg.MapTo(&conf)
	if err != nil {
		Log.Error("Unable to parse configuration file: %s", err.Error())
		exit(1)
	}

	// load categories
	conf.Categories = make(map[string]string)
	if cfg.HasSection("CATEGORIZER") {
		for _, key := range cfg.Section("CATEGORIZER").Keys() {
			conf.Categories[key.Name()] = key.Value()
		}
	}

	// load searchengines
	searchengines := make(map[string]int)
	if cfg.HasSection("SEARCHENGINES") {
		for _, key := range cfg.Section("SEARCHENGINES").Keys() {
			if _, ok := searchEngines[key.Name()]; ok {
				value, err := strconv.Atoi(key.Value())
				if err != nil {
					Log.Error("Unknown value for searchengine '%s' in configuration file: %s", key.Name(), key.Value())
					exit(1)
				}
				// only load the available searchengines to be used
				if _, ok := searchEngines[key.Name()]; ok && value != 0 {
					searchengines[key.Name()] = value
				}
			} else {
				Log.Error("Unknown searchengine '%s' in configuration file", key.Name())
				exit(1)
			}
		}
		// sort the searchengines
		engines := make([]string, 0, len(searchengines))
		for engine := range searchengines {
			engines = append(engines, engine)
		}
		sort.SliceStable(engines, func(i, j int) bool {
			return searchengines[engines[i]] < searchengines[engines[j]]
		})
		// add the searchengines to the config
		conf.Searchengines = engines
	}

	// check debug parameter
	if !args.Debug {
		args.Debug = conf.General.Debug
	}

}
