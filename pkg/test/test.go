package test

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"os"
)

// Unit helpers function to execute test command.
func Unit(pkg string) error {
	if err := goTestWithCoverage(pkg, true); err != nil {
		return err
	}

	return excludeUnwantedCodeFromCoverage()
}

func goTestWithCoverage(pkg string, isUnitOnly bool) error {
	//nolint:gofumpt
	if err := os.MkdirAll(".ci", 0700); err != nil {
		return err
	}

	out, errOut := sh.Output("go",
		"test",
		isShort(isUnitOnly),
		"-json",
		"-race",
		"-covermode=atomic",
		fmt.Sprintf("-coverpkg=%s/...", pkg),
		"-coverprofile",
		".ci/coverage.txt",
		"./...",
	)

	//nolint:gofumpt
	if err := os.WriteFile(".ci/tests.jsonl", []byte(out), 0600); err != nil {
		return err
	}

	if err := sh.RunV("gotestsum",
		"--junitfile",
		".ci/tests.xml",
		"--raw-command",
		"cat",
		".ci/tests.jsonl",
	); err != nil {
		return err
	}

	return errOut
}

func isShort(isUnitOnly bool) string {
	if isUnitOnly {
		return "-short"
	}

	// mage does not like empty command
	// use a generic useful command to substitute to `-short` for that purpose
	return "-failfast"
}

func excludeUnwantedCodeFromCoverage() error {
	pterm.Info.Printfln("Excluding mocks and generated code from coverage")

	excludes := []string{
		`/_mock\.go/d`,
		`/tdata/d`,
		`/\.pb\.go/d`,
		`/\.pb.gw\.go/d`,
		`/mage/d`,
	}

	for _, exclude := range excludes {
		if err := sh.Run("sed",
			"-i",
			exclude,
			".ci/coverage.txt",
		); err != nil {
			return err
		}
	}

	return sh.RunV("go",
		"tool",
		"cover",
		"-func",
		".ci/coverage.txt",
	)
}
