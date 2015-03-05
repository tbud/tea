package tasks

import (
	. "github.com/tbud/bud/context"
	"github.com/tbud/tea"
)

type buildTask struct {
	DestPath string // build destination path
}

func (b *buildTask) Execute() error {
	println("hello world")
	return nil
}

func (b *buildTask) Validate() error {
	return nil
}

func init() {
	build := &buildTask{
		DestPath: TEA_TARGET_PATH,
	}

	Task("build", TEA_TASK_GROUP, build, Usage("Use to build tea framework application."))
}
