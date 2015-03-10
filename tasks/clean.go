package tasks

import (
	. "github.com/tbud/bud/context"
	"os"
)

type cleanTask struct {
	Patterns []string
}

func (c *cleanTask) Execute() error {
	for _, path := range c.Patterns {
		os.RemoveAll(path)
	}
	return nil
}

func (c *cleanTask) Validate() error {
	return nil
}

func init() {
	clean := &cleanTask{
		Patterns: []string{TEA_TARGET_PATH, "tmp", "routes"},
	}

	Task("clean", clean, TEA_TASK_GROUP, Usage("Use to clean tea application target path."))
}
