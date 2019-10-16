package command

import (
	"fmt"
	"strings"

	"github.com/minamijoyo/tfupdate/tfupdate"
	flag "github.com/spf13/pflag"
)

// TerraformCommand is a command which update version constraints for terraform.
type TerraformCommand struct {
	Meta
	version     string
	path        string
	recursive   bool
	ignorePaths []string
}

// Run runs the procedure of this command.
func (c *TerraformCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("terraform", flag.ContinueOnError)
	cmdFlags.StringVarP(&c.version, "version", "v", "latest", "A new version constraint")
	cmdFlags.BoolVarP(&c.recursive, "recursive", "r", false, "Check a directory recursively")
	cmdFlags.StringArrayVarP(&c.ignorePaths, "ignore-path", "i", []string{}, "A regular expression for path to ignore")

	if err := cmdFlags.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("failed to parse arguments: %s", err))
		return 1
	}

	if len(cmdFlags.Args()) != 1 {
		c.UI.Error(fmt.Sprintf("The command expects 1 argument, but got %#v", cmdFlags.Args()))
		c.UI.Error(c.Help())
		return 1
	}

	c.path = cmdFlags.Arg(0)

	option, err := tfupdate.NewOption("terraform", "", c.version, c.recursive, c.ignorePaths)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	err = tfupdate.UpdateFileOrDir(c.Fs, c.path, option)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

// Help returns long-form help text.
func (c *TerraformCommand) Help() string {
	helpText := `
Usage: tfupdate terraform [options] <PATH>

Arguments
  PATH               A path of file or directory to update

Options:
  -v  --version      A new version constraint
  -r  --recursive    Check a directory recursively (default: false)
  -i  --ignore-path  A regular expression for path to ignore
                     If you want to ignore multiple directories, set the flag multiple times.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns one-line help text.
func (c *TerraformCommand) Synopsis() string {
	return "Update version constraints for terraform"
}
