package channels

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/cmd/cmdutils"
	"github.com/danielkrainas/shex/manager"
	"github.com/danielkrainas/shex/mods"
)

func init() {
	cmd.Register("channels", Info)
}

var (
	Info = &cmd.Info{
		Use:   "channels",
		Short: "channel operations",
		Long:  "Perform operations on channels.",
		SubCommands: []*cmd.Info{
			{
				Use:   "add <name> <endpoint>",
				Short: "add a remote channel",
				Long:  "Add a remote channel with the specified name and endpoint.",
				Run:   cmd.ExecutorFunc(addChannel),
				Flags: []*cmd.Flag{
					{
						Long:        "protocol",
						Short:       "p",
						Description: "set the protocol to use with the channel",
					},
				},
			},

			{
				Use:   "remove <name>",
				Short: "remove a channel",
				Long:  "Remove a channel with the specified name.",
				Run:   cmd.ExecutorFunc(removeChannel),
			},
			{
				Use:   "list",
				Short: "list available remote channels",
				Long:  "List available remote channels.",
				Run:   cmd.ExecutorFunc(listChannels),
			},
		},
	}
)

/* Add Channel Command */
func addChannel(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("argument missing: name")
	} else if len(args) < 2 {
		return errors.New("argument missing: endpoint")
	}

	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	alias := strings.ToLower(args[0])
	endpoint := args[1]
	if ch, ok := m.Channels()[alias]; ok {
		fmt.Printf("Overriding %s (=%s:%s)\n", alias, ch.Protocol, ch.Endpoint)
	}

	ch := &mods.Channel{
		Alias:    alias,
		Endpoint: endpoint,
	}

	protocol := ctx.Value("flags.protocol").(string)
	if len(protocol) < 1 {
		protocol = "http"
	}

	if protocol != "http" && protocol != "https" {
		fmt.Printf("unknown protocol: %s\n", protocol)
		return nil
	} else {
		ch.Protocol = protocol
	}

	if err := m.AddChannel(ch); err != nil {
		fmt.Printf("error adding channel: %v\n", err)
		return nil
	}

	fmt.Printf("channel added: %s\n", ch.Alias)
	return nil
}

/* Remove Channel Command */
func removeChannel(ctx context.Context, args []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("argument missing: name")
	}

	alias := args[0]
	var ch *mods.Channel
	if ch, err = m.RemoveChannel(alias); err != nil {
		fmt.Printf("error removing channel: %v\n", err)
		return nil
	}

	if ch == manager.DefaultChannel {
		if err := m.SaveConfig(); err != nil {
			fmt.Printf("error saving config: %v\n", err)
			return nil
		}
	}

	fmt.Printf("channel removed: %s => %s\n", ch.Alias, ch.Endpoint)
	return nil
}

/* List Channels Command */
func listChannels(ctx context.Context, _ []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	format := "%15s  %10s   %s\n"
	fmt.Printf(format, "alias", "protocol", "endpoint")
	fmt.Printf(format, "==========", "========", "==========")
	for _, ch := range m.Channels() {
		fmt.Printf(format, ch.Alias, ch.Protocol, ch.Endpoint)
	}

	return nil
}
