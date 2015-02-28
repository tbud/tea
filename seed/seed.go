package seed

import (
	"fmt"
	. "github.com/tbud/bud/seed"
	"go/build"
	"math/rand"
	"path/filepath"
	"time"
)

type teaSeed struct {
	AppName    string
	Secret     string
	appPath    string
	srcRoot    string
	importPath string

	step func(t *teaSeed, input string) (*Step, error)
}

func (t *teaSeed) Name() string {
	return "tea"
}

func (t *teaSeed) Description() string {
	return "Create a go framework use reactjs."
}

func (t *teaSeed) Start(args ...string) (step *Step, err error) {
	if err = t.checkGoPaths(); err != nil {
		return
	}

	if len(args) > 0 {
		if err = t.parseImportPath(args[0]); err != nil {
			return
		}
	}

	t.Secret = generateSecret()

	var msg string
	if len(t.appPath) > 0 {
		msg = fmt.Sprintf("The new application will be created in %s\nWhat is the application name? [%s]", t.appPath, t.AppName)
	} else {
		msg = t.emptyAppNameMsg()
	}

	t.step = stateStart
	step = &Step{Message: msg}
	return step, nil
}

func (t *teaSeed) NextStep(input string) (*Step, error) {
	if t.step != nil {
		return t.step(t, input)
	}
	return nil, nil
}

func (t *teaSeed) emptyAppNameMsg() string {
	return fmt.Sprintf("Please input application import path, it will created in %s", t.srcRoot)
}

func stateStart(t *teaSeed, input string) (step *Step, err error) {
	if len(input) > 0 {
		if err = t.parseImportPath(input); err != nil {
			return
		}
	} else {
		if len(t.AppName) == 0 {
			step = &Step{Message: t.emptyAppNameMsg()}
			return step, nil
		}
	}

	var srcPackage *build.Package
	srcPackage, err = build.Import("github.com/tbud/tea", "archetype", build.FindOnly)
	if err != nil {
		return
	}

	err = CreateArchetype(t.appPath, srcPackage.Dir, t)
	if err != nil {
		return
	}

	t.step = nil
	step = &Step{Message: fmt.Sprintf("OK, application %s is created.\n\nHave fun!", t.AppName)}
	return nil, nil
}

// lookup and set Go related variables
func (t *teaSeed) checkGoPaths() error {
	// lookup go path
	gopath := build.Default.GOPATH
	if gopath == "" {
		return fmt.Errorf("Abort: GOPATH environment variable is not set. " +
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.")
	}

	// set go src path
	t.srcRoot = filepath.Join(filepath.SplitList(gopath)[0], "src")
	return nil
}

// checking and setting application
func (t *teaSeed) parseImportPath(importPath string) (err error) {
	if filepath.IsAbs(importPath) {
		return fmt.Errorf("Abort: '%s' looks like a directory.  Please provide a Go import path instead.",
			importPath)
	}

	_, err = build.Import(importPath, "", build.FindOnly)
	if err == nil {
		return fmt.Errorf("Abort: Import path %s already exists.\n", importPath)
	}

	t.appPath = filepath.Join(t.srcRoot, filepath.FromSlash(importPath))
	t.AppName = filepath.Base(t.appPath)
	return nil
}

const alphaNumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func generateSecret() string {
	chars := make([]byte, 64)
	for i := 0; i < 64; i++ {
		chars[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
	}

	return string(chars)
}

func init() {
	rand.Seed(time.Now().UnixNano())

	Register(&teaSeed{})
}
