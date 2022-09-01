package syntax

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

// BANNER
//
// <!{S}yntax framework version="0.1.0">
//
// =======================================
const BANNER = "\n\n" +
	" \033[00;94m<\033[00m" +
	"\033[00;95m!{S}\033[00m" +
	"\033[00;94myntax Framework\033[00m" +
	"\033[00;37m version=\"0.1.0\"\033[00m" +
	"\033[00;94m>\033[00m" +
	"\n\n=======================================\n"

type Config struct {
	Dev        bool             `yaml:"dev"`
	LiveReload ConfigLiveReload `yaml:"live-reload"`
}

type ConfigLiveReload struct {
	Disabled  bool     `yaml:"disabled"`   // Allows you to disable LiveReload entirely
	Interval  int      `yaml:"interval"`   // Millis to wait on client to refresh when receive update. Defaults to `100`.
	Debounce  int      `yaml:"debounce"`   // Millis to wait before sending live reload events to the browser. Defaults to `0`.
	Pattern   []string `yaml:"pattern"`    // A list of patterns to trigger the live reloading. This option is required to enable any live reloading
	Endpoint  string   `yaml:"endpoint"`   // Endpoint of the live reload SSE event. Defaults to `dev.livereload`.
	ReloadCss bool     `yaml:"reload-css"` // If true, CSS changes will trigger a full page reload. Defaults to false.
}

var configLiveReloadPattern = []string{
	`.*\.(html|htm|js|css|png|jpeg|jpg|gif)$`,
}

// loadConfig get site configuration from config.yaml file
func loadConfig() *Config {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configFilePath := path.Join(pwd, "config.yaml")
	data, errLoadConfig := os.ReadFile(configFilePath)
	if errLoadConfig != nil {
		println(errLoadConfig)
	}

	config := &Config{}

	errUnmarshalYaml := yaml.Unmarshal(data, config)
	if errUnmarshalYaml != nil {
		failToStart(
			"Error processing configuration file",
			fmt.Sprintf("Check the file %s for possible errors.\n\n%s", configFilePath, errUnmarshalYaml.Error()),
		)
	}
	fmt.Printf("%+v\n", config)

	println(pwd)

	// @TODO: 1 Merge with environment
	// @TODO: 2 Merge with command line
	// https://docs.spring.io/spring-boot/docs/2.1.13.RELEASE/reference/html/boot-features-external-config.html
	// https://docs.spring.io/spring-boot/docs/current/reference/html/application-properties.html#appendix.application-properties.core

	return config
}
