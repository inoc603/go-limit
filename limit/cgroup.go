// +build linux

package limit

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/vbatts/go-cgroup"
)

type CgroupHelper struct {
	cgroup cgroup.Cgroup
}

func NewCgroupHelper(name string) *CgroupHelper {
	cgroup.Init()
	cg := cgroup.NewCgroup(name)
	cg.AddController("cpu")
	cg.AddController("memory")
	cg.Create()
	return &CgroupHelper{
		cgroup: cg,
	}

}

func (c *CgroupHelper) SetCPUPercentage(percentage int) error {
	if percentage > 100 || percentage < 0 {
		return errors.New("CPU Percentage out of range")
	}

	controllerCPU := c.cgroup.GetController("cpu")

	err := controllerCPU.SetValueInt64("cpu.cfs_period_us", 100*1e3)
	if err != nil {
		return err
	}

	err = controllerCPU.SetValueInt64("cpu.cfs_quota_us", int64(percentage)*1e3)
	if err != nil {
		return err
	}

	return c.cgroup.Modify()

}

var MemorySizeReg = regexp.MustCompile(`^(\d+)([kKmMgG]*)$`)

func (c *CgroupHelper) SetMemory(size string) error {
	if !MemorySizeReg.MatchString(size) {
		return errors.New("Invalid memory limit format")
	}

	controllerMemory := c.cgroup.GetController("memory")

	err := controllerMemory.SetValueString("memory.limit_in_bytes", size)
	if err != nil {
		return nil
	}

	err = controllerMemory.SetValueString("memory.memsw.limit_in_bytes", size)
	if err != nil {
		return nil
	}

	return c.cgroup.Modify()
}

func (c *CgroupHelper) AddTask(pid int) error {
	controllerMemory := c.cgroup.GetController("memory")
	controllerCPU := c.cgroup.GetController("cpu")
	task := strconv.Itoa(pid)

	err := controllerCPU.SetValueString("tasks", task)
	if err != nil {
		return err
	}

	err = controllerMemory.SetValueString("tasks", task)
	if err != nil {
		return err
	}

	return c.cgroup.Modify()
}

func (c *CgroupHelper) Delete() {
	c.cgroup.Delete()
}
