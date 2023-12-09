package discovery

import (
	"os"
	"strings"
)

// Services is a slice of strings.
type Services []string

// DiscoverServices returns services.
func DiscoverServices(path string) (Services, error) {
	services := Services{}

	files, err := os.ReadDir(path)
	if err != nil {
		return Services{}, err
	}

	for _, file := range files {
		if file.IsDir() && strings.Contains(file.Name(), "cron.") || strings.Contains(file.Name(), "service.") {
			services = append(services, file.Name())
		}
	}

	return services, nil
}
