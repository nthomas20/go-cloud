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

	"github.com/urfave/cli/v2"
)

var (
	version        string
	buildDate      string
	appName        = "go-cloud"
	appDescription = "WebDav Server"

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
			pid, err = strconv.ParseInt(string(pidBytes), 10, 64)
		}
	}

	return state, pid
}

func registerCLI() ([]*cli.Command, []cli.Flag) {
	return []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Output version",
				Action: func(c *cli.Context) error {
					fmt.Println("Version:   ", version)
					fmt.Println("Build Date:", buildDate)

					return nil
				},
			},
			{
				Name:    "status",
				Aliases: []string{"t"},
				Usage:   "Retrieve status of the service daemon",
				Action: func(c *cli.Context) error {
					// Load configuration
					// loadEnvVars(c)

					if state, pid := alreadyRunning(); state == true {
						fmt.Println("The service is currently running with PID " + strconv.Itoa(int(pid)))

						// Load the http status report
						if response, err := http.Get("http://localhost:" + portInternal + "/status"); err == nil {
							defer response.Body.Close()
							html, _ := ioutil.ReadAll(response.Body)
							fmt.Println(string(html))
							return nil
						}

						return errors.New("Could not retrieve server status page. Check configuration")
					}

					return errors.New("The service is not currently running")
				},
			},
			{
				Name:    "stop",
				Aliases: []string{"k"},
				Usage:   "Terminate service daemon",
				Action: func(c *cli.Context) error {
					if state, pid := alreadyRunning(); state == true {
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

					return errors.New("The service is not currently running")
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

					// Load configuration
					// if err = loadEnvVars(c); err != nil {
					// 	return err
					// }

					// Impossible to run as daemon if app is not logged properly
					// Check relationship status
					if os.Args[len(os.Args)-1] == "CHILD_PROCESS" {
						// Listen in the child process
						// Launch the termination listener
						launchTerminationListener()

						err = launchApp(c)
					} else {
						// I am the parent
						// Check if pid file exists
						if state, _ := alreadyRunning(); state == true {
							err = errors.New("Service already running")
						} else {
							// Fork and get PID
							pid, err := syscall.ForkExec(os.Args[0], append(os.Args, []string{"CHILD_PROCESS"}...), &syscall.ProcAttr{Files: []uintptr{0, 1, 2}})

							if err != nil {
								return err
							}

							// The parent writes PID to file before dying
							if err := ioutil.WriteFile(pidFilename, []byte(strconv.Itoa(pid)), 0777); err != nil {
								return errors.New("Could not write pid file")
							}
						}
					}

					return err

				},
			},
		},

		[]cli.Flag{}
}

func launchTerminationListener() {
	signal.Notify(daemonChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case s := <-daemonChan:
			switch s {
			// TODO: Implement graceful termination
			case syscall.SIGINT, syscall.SIGTERM:
				log.Printf("Received %s signal, exiting", s.String())

				// Send kill signal to other services
				runnerChan <- true

				// Remove pid file
				os.Remove(pidFilename)

				// Wait a moment before death ☠️
				time.Sleep(2 * time.Second)

				// Exit
				os.Exit(1)
			}
		}
	}()
}

func launchApp(c *cli.Context) error {
	return nil
}

func main() {
	// Manage CLI switches
	// Setup command routes
	commands, flags := registerCLI()

	// Configure application
	mainApp := &cli.App{
		Name:  appName,
		Usage: appDescription,
		Action: func(c *cli.Context) error {
			return errors.New("Execute `./" + appName + " help` for options")
		},
		Commands: commands,
		Flags:    flags,
	}

	// Run the app
	err := mainApp.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
