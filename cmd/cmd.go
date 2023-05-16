package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Cmder interface {
	FlagSet() *pflag.FlagSet
	Cmd() *cobra.Command
	SetRun(fc RunFunc)
}

type Config struct {
	CmdUse   string
	CmdShort string
	CmdArgs  cobra.PositionalArgs
}
type CmderWrapper struct {
	conf *Config
	cmd  *cobra.Command
}
type RunFunc func(args []string)

func New(conf *Config) Cmder {
	return &CmderWrapper{
		conf: conf,
		cmd: &cobra.Command{
			Use:   conf.CmdUse,
			Short: conf.CmdShort,
			Args:  conf.CmdArgs,
		},
	}
}

func NewWithRunFunc(conf *Config, fc RunFunc) Cmder {
	c := New(conf)
	c.SetRun(fc)
	return c
}

func (c *CmderWrapper) SetRun(fc RunFunc) {
	c.cmd.Run = func(cmd *cobra.Command, args []string) {
		fc(args)
	}
}

func (c *CmderWrapper) FlagSet() *pflag.FlagSet {
	return c.cmd.Flags()
}

func (c *CmderWrapper) Cmd() *cobra.Command {
	return c.cmd
}
