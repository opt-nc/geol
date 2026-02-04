package main

import (
	"runtime"
	"runtime/debug"

	"github.com/opt-nc/geol/v2/cmd"
	"github.com/opt-nc/geol/v2/utilities"
)

// Variables injected at build time with ldflags
var (
	commit       = "none"
	date         = "unknown"
	builtBy      = "golang"
	goVersion    = runtime.Version()
	version      = "dev"
	platformOs   = runtime.GOOS
	platformArch = runtime.GOARCH
)

func main() {
	// Initialize version variables in utilities package
	utilities.Version = version
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "(devel)" {
				utilities.Version = info.Main.Version
			}
			for _, s := range info.Settings {
				switch s.Key {
				case "vcs.revision":
					commit = s.Value
				case "vcs.time":
					date = s.Value
				}
			}
		}
	}
	utilities.Commit = commit
	utilities.Date = date
	utilities.BuiltBy = builtBy
	utilities.GoVersion = goVersion
	utilities.PlatformOs = platformOs
	utilities.PlatformArch = platformArch
	cmd.Execute()
}
