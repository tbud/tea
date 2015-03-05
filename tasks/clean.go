package tasks

import (
	. "github.com/tbud/bud/context"
	"os"
)

type cleanTask struct {
	Includes []string
}

func (c *cleanTask) Execute() error {
	for _, path := range c.Includes {
		os.RemoveAll(path)
	}
	return nil
}

func (c *cleanTask) Validate() error {
	return nil
}

func init() {
	clean := &cleanTask{
		Includes: []string{TEA_TARGET_PATH, "tmp", "routes"},
	}

	Task("clean", clean, TEA_TASK_GROUP, Usage("Use to clean tea application target path."))
}
