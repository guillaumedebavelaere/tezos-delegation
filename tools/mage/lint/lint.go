package lint

import (
	"os"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
)

// Go helper function for golint.
func Go(output string) error {
	//nolint:gofumpt
	if err := os.MkdirAll(".ci", 0700); err != nil {
		return err
	}

	goList, err := sh.Output("go", "list", "./...")
	if err != nil {
		return err
	}

	args := []string{"-set_exit_status"}
	sGoList := strings.Split(goList, "\n")
	newList := []string{}
	// Remove imock
	for _, gl := range sGoList {
		if !strings.Contains(gl, "imock") {
			newList = append(newList, gl)
		}
	}

	args = append(args, newList...)

	ret, err := sh.Output("golint", args...)

	//nolint:gofumpt
	if err := os.WriteFile(output, []byte(ret), 0600); err != nil {
		return err
	}

	if len(ret) > 0 {
		pterm.DefaultBasicText.Println(ret)
	}

	return err
}

// GolangCI helper function for golangci-lint.
func GolangCI(output string) error {
	//nolint:gofumpt
	if err := os.MkdirAll(".ci", 0700); err != nil {
		return err
	}

	out, errOut := sh.Output("golangci-lint",
		"run",
		"--out-format",
		"checkstyle",
	)

	//nolint:gofumpt
	if err := os.WriteFile(output, []byte(out), 0600); err != nil {
		return err
	}

	if len(out) > 0 {
		pterm.DefaultBasicText.Println(out)
	}

	return errOut
}
