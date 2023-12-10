//go:build mage

package main

import (
	"fmt"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/mage/discovery"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/mage/gen"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/test"
	"github.com/guillaumedebavelaere/tezos-delegation/tools/mage/lint"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

const mongoContainerName = "delegation-mongodb"

var genFiles = []*gen.File{
	{
		Name:      "datastorer",
		Type:      gen.Mock,
		Dest:      "./pkg/tezos/datastore",
		Interface: []string{"Datastorer"},
		Pkg:       "github.com/guillaumedebavelaere/tezos-delegation/pkg/tezos/datastore",
	},
}

// Help prints the help message.
func Help() error {
	pterm.DefaultTable.WithHasHeader().WithRowSeparator("-").WithHeaderRowSeparator("-").WithData(pterm.TableData{
		{"Command", "Description", "Usage"},
		{"mage -l", "Print every available command", "mage -l"},
		{"help", "Show this help", "mage help"},
		{"all", "Runs clean, build, test:unit, lint", "mage all"},
		{"clean", "Clean every micro services and crons", "mage clean"},
		{"build", "Build every micro services and crons", "mage build"},
		{"lint", "Run all linters", "mage lint"},
		{"test:unit", "Run unit tests", "mage test:unit"},
		{"mongodb:start", "starts a MongoDB Docker container", "mage mongosb:start"},
		{"mongodb:stop", "stops a MongoDB Docker container", "mage mongodb:stop"},
		{"mongodb:status", "checks the status of the MongoDB Docker container", "mage mongodb:status"},
	}).Render()

	return nil
}

// Do all.
func All() error {
	if err := Clean(); err != nil {
		return err
	}

	if err := Gen(); err != nil {
		return err
	}

	if err := Build(); err != nil {
		return err
	}

	if err := Lint(); err != nil {
		return err
	}

	if err := Test.Unit(Test{}); err != nil {
		return err
	}

	return nil
}

// Build builds all services.
func Build() error {
	return executeToServices("build")
}

// Gen generate files.
func Gen() error {
	for _, file := range genFiles {
		pterm.Info.Printfln("Generating %s %s", file.Type, file.Name)

		if err := gen.Gen(file); err != nil {
			return err
		}

		pterm.Info.Printfln("Successfully generated %s %s", file.Type, file.Name)
	}

	return executeToServices("gen")
}

// Clean cleans all services.
func Clean() error {
	return executeToServices("clean")
}

// Lint runs all linters.
func Lint() error {
	pterm.Info.Println("Running golint")

	if err := lint.Go(".ci/lint.txt"); err != nil {
		return err
	}

	pterm.Success.Println("Successfully finished golint")
	pterm.Info.Println("Running golangci-lint")

	if err := lint.GolangCI(".ci/ci-lint.xml"); err != nil {
		return err
	}

	pterm.Success.Println("Successfully finished golangci-lint")

	return nil
}

// Test defines the main target.
type Test mg.Namespace

// Unit unit test all services.
func (t Test) Unit() error {
	pterm.Info.Println("Running unit tests for github.com/guillaumedebavelaere/tezos-delegation")

	if err := test.Unit("github.com/guillaumedebavelaere/tezos-delegation"); err != nil {
		return err
	}

	pterm.Success.Println("Successfully unit test service github.com/guillaumedebavelaere/tezos-delegation")

	return nil
}

type MongoDB mg.Namespace

// Start starts a MongoDB Docker container.
func (m MongoDB) Start() error {
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

// Stop stops and removes the MongoDB Docker container.
func (m MongoDB) Stop() error {
	// Define the Docker command to stop and remove the MongoDB container
	cmd := exec.Command("docker", "stop", mongoContainerName)
	if err := cmd.Run(); err != nil {
		pterm.Error.Printfln("failed to stop MongoDB container: %v", err)

		return err
	}

	cmd = exec.Command("docker", "rm", mongoContainerName)
	if err := cmd.Run(); err != nil {
		pterm.Error.Printfln("failed to remove MongoDB container: %v", err)

		return err
	}

	pterm.Info.Printfln("MongoDB container stopped and removed successfully.")
	return nil
}

// Status checks the status of the MongoDB Docker container.
func (m MongoDB) Status() error {
	// Define the Docker command to check the status of the MongoDB container
	cmd := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", mongoContainerName)

	// Run the Docker command and print the output
	output, err := cmd.Output()
	if err != nil {
		pterm.Error.Printfln("failed to check MongoDB container status: %v", err)
		return err
	}

	pterm.Info.Printfln("MongoDB container status: %s\n", output)
	return nil
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
