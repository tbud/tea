package tasks

import (
	. "github.com/tbud/bud/context"
	"github.com/tbud/tea"
	. "github.com/tbud/x/config"
	"os"
	"path/filepath"
	"runtime"
)

const (
	TEA_TASK_GROUP_NAME = "tea"
	TEA_TASK_GROUP      = Group(TEA_TASK_GROUP_NAME)
)

type commonCfg struct {
	BaseDir     string // run base path, default is os.Getwd
	BinName     string // compiled bin name
	TargetDir   string // default task work dir, default is target
	CompileMode string // binary compile mode: debug, release. default is debug

	// use inside
	targetPath string // target abs path
	binPath    string // bin file abs path
}

func (c *commonCfg) Validate() (err error) {
	if len(c.BaseDir) == 0 {
		if c.BaseDir, err = os.Getwd(); err != nil {
			return err
		}
	}

	// validate bin name
	if len(c.BinName) == 0 {
		c.BinName = tea.App.Name
	}

	// validate compileMode
	if c.CompileMode != "release" {
		c.CompileMode = "debug"
	}

	// init target path
	c.targetPath = filepath.Join(c.BaseDir, c.TargetDir)

	// init bin path
	c.binPath = filepath.Join(c.targetPath, c.CompileMode, c.BinName)
	if runtime.GOOS == "windows" {
		c.binPath += ".exe"
	}
	return nil
}

func init() {
	TaskConfig(TEA_TASK_GROUP_NAME, Config{
		"targetDir":   "target",
		"compileMode": "debug",
	})
}
