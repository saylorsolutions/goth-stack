package main

import (
	. "github.com/saylorsolutions/modmake"
)

func localBuild(base *Build) {
	local := NewBuild()
	local.AddNewStep("run", "Runs the application locally",
		Go().Run("./cmd/yourapp").
			CaptureStdin().
			Env("DBURL", "postgres://postgres:secretpassword@localhost:5432/postgres"),
	)

	base.ImportAndLink("local", local)
	base.Step("local:run").DependsOn(base.Generate())
}
