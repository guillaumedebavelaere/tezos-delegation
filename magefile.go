//go:build mage

package main

import (
	"fmt"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/mage/discovery"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

// Help prints the help message.
func Help() error {
	pterm.DefaultTable.WithHasHeader().WithRowSeparator("-").WithHeaderRowSeparator("-").WithData(pterm.TableData{
		{"Command", "Description", "Usage"},
		{"mage -l", "Print every available command", "mage -l"},
		{"help", "Show this help", "mage help"},
		{"build", "Build every micro services and crons", "mage build"},
		{"mongoDBStart", "starts a MongoDB Docker container", "mage mongoDBStart"},
		{"mongoDBStop", "stops a MongoDB Docker container", "mage mongoDBStop"},
		{"mongoDBStatus", "checks the status of the MongoDB Docker container", "mage mongoDBStatus"},
	}).Render()

	return nil
}

// Build builds all services.
func Build() error {
	return executeToServices("build")
}

func executeToServices(cmd string) error {
	services, err := discovery.DiscoverServices("./")
	if err != nil {
		return err
	}

	for _, service := range services {
		if err := sh.RunV(
			"mage", "-d",
			fmt.Sprintf("%s", service),
			cmd); err != nil {
			return err
		}
	}

	return nil
}

// Run 'mage mongoDBStart' to start the MongoDB Docker container.
func MongoDBStart() error {
	// Define the Docker Compose command to start the MongoDB container
	cmd := exec.Command(
		"docker-compose", "-f", "dev-tools/docker-compose.yml",
		"up", "-d", "mongodb",
		"--build", "--force-recreate", "--remove-orphans",
	)

	// Set the command's output to the current os.Stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the Docker Compose command
	if err := cmd.Run(); err != nil {
		pterm.Error.Printfln("failed to start MongoDB container: %v", err)

		return err
	}

	pterm.Info.Printfln("MongoDB container started successfully.")

	return nil
}

// Run 'mage mongoDBStop' to stop and remove the MongoDB Docker container.
func MongoDBStop() error {
	// Define the Docker command to stop and remove the MongoDB container
	cmd := exec.Command("docker", "stop", "delegation-mongodb")
	if err := cmd.Run(); err != nil {
		pterm.Error.Printfln("failed to stop MongoDB container: %v", err)

		return err
	}

	cmd = exec.Command("docker", "rm", "delegation-mongodb")
	if err := cmd.Run(); err != nil {
		pterm.Error.Printfln("failed to remove MongoDB container: %v", err)

		return err
	}

	pterm.Info.Printfln("MongoDB container stopped and removed successfully.")
	return nil
}

// Run 'mage mongoDBStatus' to check the status of the MongoDB Docker container.
func MongoDBStatus() error {
	// Define the Docker command to check the status of the MongoDB container
	cmd := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", "delegation-mongodb")

	// Run the Docker command and print the output
	output, err := cmd.Output()
	if err != nil {
		pterm.Error.Printfln("failed to check MongoDB container status: %v", err)
		return err
	}

	pterm.Info.Printfln("MongoDB container status: %s\n", output)
	return nil
}
