package service

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"

	"github.com/guillaumedebavelaere/tezos-delegation/pkg/mage/build"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/mage/gen"
)

var (
	// Name defines service name.
	Name string
	// GenFiles defines file to generate inside service.
	GenFiles []*gen.File
)

// Clean cleans binaries and deliveries.
func Clean() error {
	pterm.Info.Printfln("Cleaning builds %s", Name)

	return sh.Rm("build")
}

// Build builds the service.
func Build() error {
	pterm.Info.Printfln("Building service %s", Name)

	if err := build.Build(
		fmt.Sprintf("cmd/%s/main.go", Name),
		fmt.Sprintf("build/%s", Name),
		map[string]string{
			"CGO_ENABLED": "0",
			"GOOS":        os.Getenv("GOOS"),
			"GOARCH":      os.Getenv("GOARCH"),
		},
		[]string{"-ldflags=-s"},
	); err != nil {
		return err
	}

	pterm.Success.Printfln("Successfully built service %s", Name)

	return nil
}

// Gen generate files.
func Gen() error {
	for _, file := range GenFiles {
		pterm.Info.Printfln("Generating service.%s %s %s", Name, file.Type, file.Name)

		if err := gen.Gen(file); err != nil {
			return err
		}

		pterm.Info.Printfln("Successfully generated service.%s %s %s", Name, file.Type, file.Name)
	}

	return nil
}

// Run target runs the service.
func Run() error {
	pterm.Info.Printfln("Running service %s", Name)

	err := sh.RunWithV(map[string]string{}, "go", "run", fmt.Sprintf("cmd/%s/main.go", Name))
	if err != nil {
		pterm.Error.Printfln("Failed to run service %s", Name)

		return err
	}

	pterm.Success.Printfln("Successfully ran service %s", Name)

	return nil
}
