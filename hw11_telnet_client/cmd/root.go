package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/AHTI6IOTIK/hw_otus/hw11_telnet_client/telnet"
	"github.com/spf13/cobra"
)

const (
	timeoutFlgName      = "timeout"
	timeoutFlgShortName = "t"
)

var (
	quit    chan os.Signal
	rootCmd = &cobra.Command{
		Use:   "go-telnet",
		Short: "Крайне примитивный TELNET клиент",
		Long:  `Крайне примитивный TELNET клиент (без поддержки команд, опций и протокола в целом)`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				_ = cmd.Help()
			}

			if err := cobra.MaximumNArgs(2)(cmd, args); err != nil {
				return err
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			timeout, err := cmd.Flags().GetString(timeoutFlgName)
			if err != nil {
				log.Fatalln(fmt.Errorf("getting the timeout flag: %w", err))
			}

			timeoutDuration, err := time.ParseDuration(timeout)
			if err != nil {
				log.Fatalln(fmt.Errorf("parsing the timeout flag: %w", err))
			}

			client := telnet.NewTelnetClient(
				net.JoinHostPort(args[0], args[1]),
				timeoutDuration,
				os.Stdin,
				os.Stdout,
			)

			err = client.Connect()
			if err != nil {
				log.Fatalln(err)
			}

			go func() {
				defer func() {
					client.Close()
					quit <- syscall.SIGINT
				}()
				err := client.Send()
				if err != nil {
					log.Println(err)
					return
				}
			}()

			go func() {
				defer func() {
					client.Close()
					quit <- syscall.SIGINT
				}()
				err := client.Receive()
				if err != nil {
					log.Println(err)
					return
				}
			}()
		},
	}
)

func Execute(ctx context.Context, quitCh chan os.Signal) {
	quit = quitCh
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.
		Flags().
		StringP(
			timeoutFlgName,
			timeoutFlgShortName,
			"10s",
			"Timeout for connecting to the server in seconds",
		)

	rootCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		fmt.Println(rootCmd.UsageString())
		os.Exit(0)
	})
}
