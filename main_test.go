package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/newtonproject/newchain-dump/cli"
)

// test address
const taddr = "0xDB2C9C06E186D58EFe19f213b3d5FaF8B8c99481"

func getTempFile() (string, func()) {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}

	file := dir + string(os.PathSeparator) + "lumen-integration-test.json"

	return file, func() {
		logrus.Debugf("cleaning up temp file: %s", file)
		os.RemoveAll(dir)
	}
}

func run(cli *cli.CLI, command string) string {
	fmt.Printf("$ ./%s %s\n", cli.Name, command)
	got := cli.TestCommand(command)
	fmt.Printf("%s\n", got)
	return strings.TrimSpace(got)
}

func runArgs(cli *cli.CLI, args ...string) string {
	fmt.Printf("$ ./%s %s\n", cli.Name, strings.Join(args, " "))
	got := cli.Embeddable().Run(args...)
	fmt.Printf("%s\n", got)
	return strings.TrimSpace(got)
}

func expectOutput(t *testing.T, cli *cli.CLI, want string, command string) {
	got := run(cli, command)

	if got != want {
		t.Errorf("(%s) wrong output: want %v, got %v", command, want, got)
	}
}

func newCLI() (*cli.CLI, func()) {
	_, cleanupFunc := getTempFile()

	dpos := cli.NewCLI()
	dpos.TestCommand("version")
	run(dpos, fmt.Sprintf("version"))

	return dpos, cleanupFunc
}

// Create new funded test account
func TestAll(t *testing.T) {
	cli, _ := newCLI()
	run(cli, "init") // no password

	run(cli, "run --end 100")
}
