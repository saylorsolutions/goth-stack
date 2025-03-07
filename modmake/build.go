package main

import (
	. "github.com/saylorsolutions/modmake"
	"github.com/saylorsolutions/modmake/pkg/git"
	"runtime"
)

const (
	appMainPath = "./cmd/yourapp"
	version     = "0.1.0"
)

var (
	appBuildPath = Path("./build/yourapp")
	appName      = "yourapp"
)

func init() {
	if runtime.GOOS == "windows" {
		appName += ".exe"
	}
}

func main() {
	Go().PinLatestV1(23)
	b := NewBuild()
	b.Tools().DependsOnRunner("install-templ", "", Go().Install("github.com/a-h/templ/cmd/templ@v0.3.833"))
	b.Generate().DependsOnRunner("gen-templ", "",
		Script(
			Exec("templ", "generate", "./cmd/yourapp/internal/templates"),
			Go().VetAll(),
		),
	)
	b.Build().DependsOnRunner("build-app", "", Script(
		RemoveDir(appBuildPath),
		MkdirAll(appBuildPath, 0755),
		Go().Build(appMainPath).
			Env("CGO_ENABLED", "0").
			StripDebugSymbols().
			SetVariable("main", "version", version).
			SetVariable("main", "gitHash", git.CommitHash()).
			OutputFilename(appBuildPath.Join(appName)),
	))

	localBuild(b)
	b.Execute()
}
