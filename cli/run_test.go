package cli

import "testing"

func TestRun(t *testing.T) {
	cli := NewCLI()

	cli.TestCommand("run --end 100")
	cli.TestCommand("run --start 10 --end 100")
	cli.TestCommand("NewChainDump run --loop")

}
