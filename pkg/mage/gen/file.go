package gen

const (
	// Mock define mock generation type.
	Mock = "MOCK"
)

// Type define generation type.
type Type string

// File represents a file to generate.
type File struct {
	Name      string
	Path      string
	Type      Type
	Dest      string
	Interface []string
	Pkg       string
}
