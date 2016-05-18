package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/inoc603/go-limit/limit"
	"github.com/satori/go.uuid"
)

func main() {
	app := cli.NewApp()
	app.Name = "limit"
	app.Usage = "Limit cgroup usage"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "cpu",
			Usage: "CPU percentage limit",
			Value: 0,
		},
		cli.StringFlag{
			Name:  "memory, m",
			Usage: "Memory limit",
		},
	}
	app.Action = func(c *cli.Context) {
		if err := checkArgs(c); err != nil {
			fmt.Println(err.Error())
			return
		}

		uid := "limit-" + uuid.NewV4().String()

		cgroupHelper := limit.NewCgroupHelper(uid)
		defer cgroupHelper.Delete()

		if memory := c.String("memory"); memory != "" {
			err := cgroupHelper.SetMemory(memory)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

		if cpu := c.Int("cpu"); cpu != 0 {
			err := cgroupHelper.SetCPUPercentage(c.Int("cpu"))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

		// Add this process to the cgroup, so its child process can
		// be controlled
		err := cgroupHelper.AddTask(os.Getpid())
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Run command as child process
		cmd := exec.Command(c.Args()[0], c.Args()[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Start()

		// Wait for command to stop
		cmd.Wait()
	}

	app.Run(os.Args)
}

func checkArgs(c *cli.Context) error {

	if len(c.Args()) == 0 {
		return errors.New("No command")
	}

	if c.Int("cpu") > 100 || c.Int("cpu") < 0 {
		return errors.New("CPU percentage out of range")
	}

	memReg := regexp.MustCompile(`^(\d+)([kKmMgG]*)$`)
	if c.String("memory") != "" && !memReg.MatchString(c.String("memory")) {
		return errors.New("Invalid memory limit format")
	}

	return nil
}
