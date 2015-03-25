package controllers

import (
	"github.com/tbud/bud/asset"
	. "github.com/tbud/tea/context"
	"path/filepath"
)

type Assets struct {
	*Context
}

func (a *Assets) At(path string, file string) {
	fp := filepath.Join(path, file)
	if _, err := asset.Open(fp); err != nil {
		Log.Error("Open file err: %v", err)
	}
}
