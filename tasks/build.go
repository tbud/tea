package tasks

import (
	"github.com/tbud/bud/builtin"
	. "github.com/tbud/bud/context"
	. "github.com/tbud/x/config"
	"os"
	"path/filepath"
	"text/template"
)

type buildTask struct {
	commonCfg

	// use inside
	mainSrcPath string
}

func (b *buildTask) Execute() (err error) {
	if err = b.writeMainSrc(); err != nil {
		return err
	}

	buildArgs := []string{"build"}
	if b.CompileMode == "release" {
		buildArgs = append(buildArgs, "-ldflags", "-w")
	}
	buildArgs = append(buildArgs, "-o", b.binPath, b.mainSrcPath)

	return builtin.Exec("go", buildArgs...)
}

func (b *buildTask) writeMainSrc() (err error) {
	if err = os.MkdirAll(filepath.Dir(b.mainSrcPath), 0744); err != nil {
		return err
	}

	// open main src file
	var mainFile *os.File
	if mainFile, err = os.Create(b.mainSrcPath); err != nil {
		return err
	}
	defer mainFile.Close()

	// execute template
	var temp *template.Template
	if temp, err = template.New("").Parse(Main); err != nil {
		return err
	}
	if err = temp.Execute(mainFile, nil); err != nil {
		return err
	}
	return nil
}

func (b *buildTask) Validate() (err error) {
	if err = b.commonCfg.Validate(); err != nil {
		return err
	}

	// init src path
	b.mainSrcPath = filepath.Join(b.targetPath, "src/main.go")
	return nil
}

func init() {
	build := &buildTask{}

	Task("build", TEA_TASK_GROUP, build, Usage("Use to build tea framework application."))

	Task("debug", TEA_TASK_GROUP, Usage("Use to build debug tea framework application."), func() error {
		return RunTask("tea.build", Config{"compileMode": "debug"})
	})

	Task("release", TEA_TASK_GROUP, Usage("Use to build release tea framework application."), func() error {
		return RunTask("tea.build", Config{"compileMode": "release"})
	})

	// Task("release", TEA_TASK_GROUP, build, Config{"compileMode": "release"}, Usage("Use to build release tea framework application."))
}

const Main = `// GENERATED CODE - DO NOT EDIT
package main

import (
	"flag"
	"github.com/tbud/tea"
)

var (
	port *int = flag.Int("port", 0, "By default, read from app.conf")
)

func main() {
	flag.Parse()

	tea.Run(*port)
}
`
