package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerEnableRescueCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "enable-rescue [FLAGS] SERVER",
		Short:                 "Enable rescue for a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerEnableRescue),
	}
	cmd.Flags().String("type", "linux64", "Rescue type")
	cmd.Flags().IntSlice("ssh-key", nil, "ID of SSH key to inject (can be specified multiple times)")
	return cmd
}

func runServerEnableRescue(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	var (
		opts hcloud.ServerEnableRescueOpts
	)
	rescueType, _ := cmd.Flags().GetString("type")
	opts.Type = hcloud.ServerRescueType(rescueType)

	sshKeys, _ := cmd.Flags().GetIntSlice("ssh-key")
	for _, sshKey := range sshKeys {
		opts.SSHKeys = append(opts.SSHKeys, &hcloud.SSHKey{ID: sshKey})
	}

	result, _, err := cli.Client().Server.EnableRescue(cli.Context, server, opts)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(cli.Context, cli.Client(), result.Action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Password of server %s reset to: %s\n", idOrName, result.RootPassword)
	return nil
}