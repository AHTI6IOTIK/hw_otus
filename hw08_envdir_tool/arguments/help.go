package arguments

import (
	"fmt"
	"log"
)

type argDesc struct {
	Name string
	Desc string
}

func (a *argDesc) GetName(isWrap bool) string {
	if isWrap {
		return fmt.Sprintf("[%s]", a.Name)
	}

	return a.Name
}

var ArgsDescription = struct {
	EnvDir    argDesc
	Command   argDesc
	Arguments argDesc
	Verbose   argDesc
	Help      argDesc
}{
	EnvDir: argDesc{
		Name: "env dir path",
		Desc: "the path to the directory with variables - required (/path/to/env/dir, ./env/dir, dir)",
	},
	Command: argDesc{
		Name: "command",
		Desc: "executable program - required",
	},
	Arguments: argDesc{
		Name: "command arguments",
		Desc: "arguments for the program being executed",
	},
	Verbose: argDesc{
		Name: "-vv",
		Desc: "verbose mode",
	},
	Help: argDesc{
		Name: "-h | help",
		Desc: "information about the program",
	},
}

func PrintHelp(name, version string) {
	log.Printf(
		"Help for command version %s\n%s %s %s %s \n\n",
		version,
		name,
		ArgsDescription.EnvDir.GetName(true),
		ArgsDescription.Command.GetName(true),
		ArgsDescription.Arguments.GetName(true),
	)

	log.Printf("\t[%s] %s\n", ArgsDescription.EnvDir.Name, ArgsDescription.EnvDir.Desc)
	log.Printf("\t[%s] %s\n", ArgsDescription.Command.Name, ArgsDescription.Command.Desc)
	log.Printf("\t[%s] %s\n", ArgsDescription.Arguments.Name, ArgsDescription.Arguments.Desc)
	log.Printf("\t[%s] %s\n", ArgsDescription.Verbose.Name, ArgsDescription.Verbose.Desc)
	log.Printf("\t%s %s\n\n", ArgsDescription.Help.GetName(true), ArgsDescription.Help.Desc)
}
