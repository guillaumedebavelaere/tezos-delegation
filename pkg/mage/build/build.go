package build

import "github.com/magefile/mage/sh"

// Build helpers function to execute a go build cmd.
func Build(src, dest string, env map[string]string, flags []string) error {
	args := []string{
		"build",
	}

	args = append(args, flags...)
	args = append(args, "-o")
	args = append(args, dest)
	args = append(args, src)

	return sh.RunWithV(env, "go",
		args...,
	)
}
