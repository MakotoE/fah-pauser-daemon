package main

import (
	"flag"
	"github.com/MakotoE/go-fahapi"
	"github.com/go-yaml/yaml"
	"github.com/mitchellh/go-ps"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"time"
)

var verbose = flag.Bool("v", false, "verbose")

func main() {
	flag.Parse()

	config := readConfig()

	api, err := fahapi.NewAPI()
	if err != nil {
		log.Panicln(err)
	}

	defer api.Close()

	defer api.UnpauseAll() // Make sure FAH is unpaused in case of panic

	paused := false

	for {
		processes, err := ps.Processes()
		if err != nil {
			log.Panicln(err)
		}
		if *verbose {
			b := strings.Builder{}
			b.WriteString("current processes:\n")
			for _, process := range processes {
				b.WriteString(process.Executable() + "\n")
			}
			log.Printf(b.String())
		}

		if containsProcess(processes, config.PauseOn) {
			if !paused { // Found process; fah is unpaused
				if err := api.PauseAll(); err != nil {
					log.Panicln(err)
				}
				paused = true
				if *verbose {
					log.Println("paused")
				}
			}
		} else if paused { // No process found; fah is paused
			if err := api.UnpauseAll(); err != nil {
				log.Panicln(err)
			}
			paused = false
			if *verbose {
				log.Println("unpaused")
			}
		}

		time.Sleep(time.Minute * 5)
	}
}

type config struct {
	PauseOn []string `yaml:"PauseOn"`
}

func readConfig() *config {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(path.Join(u.HomeDir, ".config", "fah-pauser-daemon.yml"))
	if err != nil {
		log.Panicln(err)
	}

	result := &config{}
	if err := yaml.NewDecoder(file).Decode(result); err != nil {
		log.Panicln(err)
	}
	return result
}

// containsProcess returns true if processes contains an executable that matches any string in find.
func containsProcess(processes []ps.Process, find []string) bool {
	for _, process := range processes {
		for _, s := range find {
			if process.Executable() == s {
				return true
			}
		}
	}

	return false
}
