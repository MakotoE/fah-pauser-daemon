package main

import (
	"encoding/json"
	"flag"
	"github.com/MakotoE/go-fahapi"
	"github.com/go-yaml/yaml"
	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	"log"
	"net"
	"os"
	"os/user"
	"path"
	"syscall"
	"time"
)

var verbose = false

func main() {
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()

	config := readConfig()

	var api *fahapi.API
	for {
		a, err := fahapi.Dial(fahapi.DefaultAddr)
		if err != nil {
			if e, ok := errors.Cause(err).(*net.OpError); ok {
				if syscallErr, ok := e.Err.(*os.SyscallError); ok && syscallErr.Err == syscall.ECONNREFUSED {
					log.Println("connection refused; trying again after a bit")
					time.Sleep(time.Second * 30)
					continue
				}
			}
			log.Panicln(err)
		}
		api = a
		break
	}
	defer api.Close()

	if err := api.UnpauseAll(); err != nil {
		log.Panicln(err)
	}

	paused := false

	for {
		processes, err := ps.Processes()
		if err != nil {
			log.Panicln(err)
		}
		if verbose {
			var processStr []string
			for _, process := range processes {
				processStr = append(processStr, process.Executable())
			}
			s, err := json.Marshal(processStr)
			if err != nil {
				log.Panicln(err)
			}
			log.Printf("current processes: %s\n", string(s))
		}

		if containsProcess(processes, config.PauseOn) {
			if !paused { // Found process; fah is unpaused
				if err := api.PauseAll(); err != nil {
					log.Panicln(err)
				}
				paused = true
				if verbose {
					log.Println("pausing fah")
				}
			}
		} else if paused { // No process found; fah is paused
			if err := api.UnpauseAll(); err != nil {
				log.Panicln(err)
			}
			paused = false
			if verbose {
				log.Println("unpausing fah")
			}
		}

		time.Sleep(time.Minute * 2)
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
	defer file.Close()

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
