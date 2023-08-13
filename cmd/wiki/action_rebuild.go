package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

// Rebuild
func Rebuild(c *cli.Context) error {
	if err := os.RemoveAll("./cache"); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := Refresh(c); err != nil {
		return err
	}
	if err := Update(c); err != nil {
		return err
	}
	return nil
}
