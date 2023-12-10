package gen

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"strings"
)

// Gen generate protobuf, mocks, and others.
func Gen(file *File) error {
	return genMock(file)
}

func genMock(file *File) error {
	return sh.RunV("mockgen",
		fmt.Sprintf("-destination=%s/mock/%s_mock.go", file.Dest, strings.ToLower(file.Name)),
		fmt.Sprintf("-package=mock_%s", file.Name),
		file.Pkg,
		strings.Join(file.Interface, ","),
	)
}
