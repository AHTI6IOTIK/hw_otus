/*
Copyright © 2024 Pavel Sidorov <p.sidorov.dev@gmail.com>
*/
package main

import (
	"github.com/AHTI6IOTIK/hw_otus/hw08_envdir_tool/arguments"
	"log"
	"os"
)

var (
	version     = "1.0.0"
	executeName = ""
)

func main() {
	args, err := arguments.Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	if args.IsHelp() {
		arguments.PrintHelp(executeName, version)
		return
	}

	env, err := ReadDir(args.EnvDir())
	if err != nil {
		log.Fatalln(err)
	}

	exitCode := RunCmd(
		args.FullCommand(),
		env,
	)

	os.Exit(exitCode)
}

func init() {
	log.SetFlags(^log.Ltime & ^log.Ldate & ^log.Llongfile & ^log.Lshortfile & ^log.Lmicroseconds)
}
