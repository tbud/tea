package tasks

import (
	. "github.com/tbud/bud/context"
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
	build := &buildTask{}

	Task("build", Group("tea"), build, Usage("Use to build tea framework application."))
}
