package tasks

import (
	. "github.com/tbud/bud/context"
	"os"
)

type cleanTask struct {
	commonCfg
	Patterns []string
}

func (c *cleanTask) Execute() (err error) {
	// remove target path
	if err = os.RemoveAll(c.targetPath); err != nil {
		return err
	}

	for _, path := range c.Patterns {
		if err = os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}

func (c *cleanTask) Validate() (err error) {
	if err = c.commonCfg.Validate(); err != nil {
		return err
	}

	return nil
}

func init() {
	clean := &cleanTask{}

	Task("clean", clean, TEA_TASK_GROUP, Usage("Use to clean tea application target path."))
}
