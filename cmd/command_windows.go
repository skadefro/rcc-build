package cmd

import (
	"os"
	"path/filepath"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/shell"
)

func osSpecificHolotreeSharing(enable bool) {
	if !enable {
		return
	}
	pathlib.ForceShared()
	parent := filepath.Dir(common.HoloLocation())
	_, err := pathlib.ForceSharedDir(parent)
	pretty.Guard(err == nil, 1, "Could not enable shared location at %q, reason: %v", parent, err)
	task := shell.New(nil, ".", "icacls", "C:/ProgramData/robocorp", "/grant", "BUILTIN\\Users:(OI)(CI)M", "/T", "/Q")
	_, err = task.Execute(false)
	pretty.Guard(err == nil, 2, "Could not set 'icacls' settings, reason: %v", err)
	err = os.WriteFile(common.SharedMarkerLocation(), []byte(common.Version), 0644)
	pretty.Guard(err == nil, 3, "Could not write %q, reason: %v", common.SharedMarkerLocation(), err)
}
