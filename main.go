/*
 * Filename: main.go
 * Author: Nathaniel Thomas
 */

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gitea.nthomas20.net/nathaniel/go-cloud/app/api"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/bootstrap"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/cmd"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/configuration"
	"gitea.nthomas20.net/nathaniel/go-cloud/app/jobs"
	"github.com/urfave/cli/v2"
)

var (
	version        string
	buildDate      string
	appName        = "go-cloud"
	appDescription = "GO Cloud Server"
	config         = configuration.NewConfiguration()

	portInternal = "8080"
	pidFilename  = ".go-cloud.pid"

	runnerChan = make(chan bool)
	daemonChan = make(chan os.Signal, 1)
)

func alreadyRunning() (bool, int64) {
	var (
		state = false
		pid   = int64(0)
	)

	if _, err := os.Stat(pidFilename); err == nil {
		state = true

		// Attempt to load the pid value
		if pidBytes, err := ioutil.ReadFile(pidFilename); err == nil {
			pid, _ = strconv.ParseInt(string(pidBytes), 10, 64)
		}
	}

	return state, pid
}

func registerCLI() []*cli.Command {
	return append(cmd.Commands(), []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Output version",
			Action: func(c *cli.Context) error {
				fmt.Println("Version:   ", version)
				fmt.Println("Build Date:", buildDate)

				if state, pid := alreadyRunning(); state {
					// Bootstrap Configuration
					bootstrap.SetupConfiguration()
					configuration.ReadConfiguration(config)
					fmt.Println("\nVersion information of currently running PID " + strconv.Itoa(int(pid)))

					// Load the http status report
					if response, err := http.Get("http://localhost:" + config.Port + "/version"); err == nil {
						defer response.Body.Close()
						html, _ := ioutil.ReadAll(response.Body)
						fmt.Println(string(html))
						return nil
					}

					return errors.New("could not retrieve server status page. check configuration")
				}

				return nil
			},
		},
		{
			Name:    "status",
			Aliases: []string{"t"},
			Usage:   "Retrieve status of the service daemon",
			Action: func(c *cli.Context) error {
				if state, pid := alreadyRunning(); state {
					// Bootstrap Configuration
					bootstrap.SetupConfiguration()
					configuration.ReadConfiguration(config)

					fmt.Println("The service is currently running with PID " + strconv.Itoa(int(pid)))

					// Load the http status report
					if response, err := http.Get("http://localhost:" + portInternal + "/status"); err == nil {
						defer response.Body.Close()
						html, _ := ioutil.ReadAll(response.Body)
						fmt.Println(string(html))
						return nil
					}

					return errors.New("could not retrieve server status page. check configuration")
				}

				return errors.New("the service is not currently running")
			},
		},
		{
			Name:    "stop",
			Aliases: []string{"k"},
			Usage:   "Terminate service daemon",
			Action: func(c *cli.Context) error {
				if state, pid := alreadyRunning(); state {
					// Send terminate signal to service
					// TODO: Implement graceful termination
					syscall.Kill(int(pid), syscall.SIGTERM)

					// Remove pid file
					os.Remove(pidFilename)

					fmt.Println("The service has been terminated")
					return nil
				}

				// Remove pid file
				os.Remove(pidFilename)

				return errors.New("the service is not currently running")
			},
		},
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "Start service daemon",
			Action: func(c *cli.Context) error {
				var (
					err error
				)

				// Check relationship status
				if os.Args[len(os.Args)-1] == "CHILD_PROCESS" {
					bootstrap.ConfigDirectory = os.Args[len(os.Args)-2]
					configuration.ReadConfiguration(config)
					// Listen in the child process
					// Launch the termination listener
					launchTerminationListener()

					err = launchApp(c)
				} else {
					// Bootstrap Configuration
					bootstrap.SetupConfiguration()
					configuration.ReadConfiguration(config)

					// I am the parent
					// Check if pid file exists
					if state, _ := alreadyRunning(); state {
						err = errors.New("service already running")
					} else {
						// Fork and get PID
						pid, err := syscall.ForkExec(os.Args[0], append(os.Args, []string{bootstrap.ConfigDirectory, "CHILD_PROCESS"}...), &syscall.ProcAttr{Files: []uintptr{0, 1, 2}})

						if err != nil {
							return err
						}

						// The parent writes PID to file before dying
						if err := ioutil.WriteFile(pidFilename, []byte(strconv.Itoa(pid)), 0777); err != nil {
							return errors.New("could not write pid file")
						}
					}
				}

				return err

			},
		},
	}...)
}

func launchTerminationListener() {
	signal.Notify(daemonChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Listen/wait for a signal
		s := <-daemonChan

		switch s {
		// TODO: Implement graceful termination
		case syscall.SIGINT, syscall.SIGTERM:
			log.Printf("Received %s signal, exiting", s.String())

			// Send kill signal to other services
			runnerChan <- true

			// Remove pid file
			os.Remove(pidFilename)

			// Wait a moment before death ??????
			time.Sleep(2 * time.Second)

			// Exit
			os.Exit(1)
		}
	}()
}

func launchApp(c *cli.Context) error {
	// Setup the app
	// var app api.API

	app := &api.Configuration{
		Configuration: config,
		Version:       version,
		BuildDate:     buildDate,
	}

	// Automatically reload configuration
	jobs.RefreshConfiguration(config, time.Duration(time.Second*60))

	// We follow this pattern to plan ahead to allow for control and internal status check listeners
	if success := app.Run(); success {
		// Give us a moment to bind to the port, or exit
		time.Sleep(2 * time.Second)

		log.Println("Service Started. Listening on port " + config.Port)

		// Run forever
		// TODO: How to fix this? (unnecessary assignment to the blank identifier (S1005))
		_ = <-runnerChan

		return nil
	}

	log.Println("could not launch service")
	return errors.New("could not launch service")
}

func main() {
	// Setup command routes and configure application
	mainApp := &cli.App{
		Name:  appName,
		Usage: appDescription,
		Action: func(c *cli.Context) error {
			return errors.New("Execute `./" + appName + " help` for options")
		},
		Commands: registerCLI(),
	}

	// Run the app
	err := mainApp.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
