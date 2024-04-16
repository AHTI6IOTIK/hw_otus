package arguments

import (
	"errors"
	"fmt"
)

var (
	invalidParamsTemplate = "%s cannot be options, received: %s"
)

type Args struct {
	envDir    string
	command   string
	arguments []string
	verbose   bool
	help      bool
}

func Parse(args []string) (*Args, error) {
	result := new(Args)

	if isHelp(args) {
		result.help = true

		return result, nil
	}

	if err := isValid(args); err != nil {
		return nil, InvalidArgumentsError{msg: err.Error()}
	}

	result.envDir = args[0]
	result.command = args[1]
	result.arguments = args[2:]

	return result, nil
}

func isValid(args []string) error {
	if len(args) < 2 {
		return errors.New("there must be more than 2 arguments")
	}

	if args[0][0] == '-' {
		return errors.New(
			fmt.Sprintf(
				invalidParamsTemplate,
				ArgsDescription.EnvDir.GetName(true),
				args[0],
			),
		)
	}

	if args[1][0] == '-' {
		return errors.New(
			fmt.Sprintf(
				invalidParamsTemplate,
				ArgsDescription.Command.GetName(true),
				args[1],
			),
		)
	}

	return nil
}

func isHelp(args []string) bool {
	if len(args) == 0 {
		return true
	}

	for _, arg := range args {
		if arg == "-h" {
			return true
		}
	}

	return false
}

func (a *Args) EnvDir() string {
	return a.envDir
}

func (a *Args) Command() string {
	return a.command
}

func (a *Args) FullCommand() []string {
	res := make([]string, 0, len(a.arguments)+1)

	res = append(res, a.command)
	res = append(res, a.arguments...)

	return res
}

func (a *Args) Arguments() []string {
	return a.arguments
}

func (a *Args) IsHelp() bool {
	return a.help
}

func (a *Args) IsVerbose() bool {
	return a.help
}
