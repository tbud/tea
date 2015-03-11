package tasks

import (
	"fmt"
	. "github.com/tbud/bud/context"
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
	ConfigFile  string // config file path
	CompileMode string // binary compile mode: debug, release. default is debug

	// use inside
	appConfig  Config // loaded config from app config file
	targetPath string // target abs path
	binPath    string // bin file abs path
}

func (c *commonCfg) Validate() (err error) {
	if len(c.BaseDir) == 0 {
		if c.BaseDir, err = os.Getwd(); err != nil {
			return err
		}
	}

	// validate config
	configFile := filepath.Join(c.BaseDir, c.ConfigFile)
	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("Config file not exist: %s, err: %v", configFile, err)
	}
	if c.appConfig, err = Load(configFile); err != nil {
		return err
	}

	// validate bin name
	if len(c.BinName) == 0 {
		c.BinName = c.appConfig.StringDefault("app.name", "sample")
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
		"configFile":  "conf/app.conf",
		"compileMode": "debug",
	})
}
