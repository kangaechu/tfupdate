package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/minamijoyo/tfupdate/tfupdate"
)

func main() {
	filename := "./main.tf"

	updaterType := "terraform"
	name := ""
	version := "0.12.7"

	// updaterType := "provider"
	// name := "aws"
	// version := "2.23.0"

	option := tfupdate.NewOption(updaterType, name, version)
	err := updateFile(filename, option)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func updateFile(filename string, o tfupdate.Option) error {
	r, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %+v", err)
	}

	err = update(r, os.Stdout, filename, o)
	if err != nil {
		return err
	}

	return nil
}

func update(r io.Reader, w io.Writer, filename string, o tfupdate.Option) error {
	f, err := parseHCL(r, filename)
	if err != nil {
		return err
	}

	err = updateHCL(f, o)
	if err != nil {
		return err
	}

	err = writeHCL(f, w)
	if err != nil {
		return err
	}

	return nil
}

func parseHCL(r io.Reader, filename string) (*hclwrite.File, error) {
	src, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("failed to read input: err = %+v", err)
	}

	f, diags := hclwrite.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse file: %s", diags)
	}

	return f, nil
}

func writeHCL(f *hclwrite.File, w io.Writer) error {
	tokens := f.BuildTokens(nil)
	buf := hclwrite.Format(tokens.Bytes())

	fmt.Fprintln(w, string(buf))

	return nil
}

func updateHCL(f *hclwrite.File, o tfupdate.Option) error {
	u, err := tfupdate.NewUpdater(o)
	if err != nil {
		return err
	}

	return u.Update(f)
}
