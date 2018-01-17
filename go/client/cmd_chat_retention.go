// Copyright 2018 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package client

import (
	"fmt"
	"os"

	"github.com/keybase/cli"
	"github.com/keybase/client/go/chat"
	"github.com/keybase/client/go/libcmdline"
	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/chat1"
	"github.com/keybase/client/go/protocol/keybase1"
	isatty "github.com/mattn/go-isatty"
)

type CmdChatSetRetention struct {
	libkb.Contextified
	resolvingRequest chatConversationResolvingRequest
	hasTTY           bool
	convID           *chat1.ConversationID
	teamID           *keybase1.TeamID
	policy           *chat1.RetentionPolicy
}

func NewCmdChatSetRetentionRunner(g *libkb.GlobalContext) *CmdChatSetRetention {
	return &CmdChatSetRetention{
		Contextified: libkb.NewContextified(g),
	}
}

func newCmdChatSetRetention(cl *libcmdline.CommandLine, g *libkb.GlobalContext) cli.Command {
	return cli.Command{
		Name:         "retention-policy",
		Usage:        "Manage the chat retention policy for a conversation or team",
		ArgumentHelp: "[conversation]",
		Action: func(c *cli.Context) {
			cl.ChooseCommand(NewCmdChatSetRetentionRunner(g), "retention-policy", c)
		},
		Flags: append(getConversationResolverFlags(), []cli.Flag{
			cli.BoolFlag{
				Name:  "keep",
				Usage: `Keep messages indefinitely`,
			},
			cli.StringFlag{
				Name:  "age",
				Usage: `Delete messages after e.g. 2h, 3d, 1w `,
			},
			cli.BoolFlag{
				Name:  "inherit",
				Usage: `Use the team's policy for a channel`,
			},
		}...),
	}
}

func (c *CmdChatSetRetention) Run() (err error) {
	if c.resolvingRequest.TlfName != "" {
		err = annotateResolvingRequest(c.G(), &c.resolvingRequest)
		if err != nil {
			return err
		}
	}
	// TLFVisibility_ANY doesn't make any sense for send, so switch that to PRIVATE:
	if c.resolvingRequest.Visibility == keybase1.TLFVisibility_ANY {
		c.resolvingRequest.Visibility = keybase1.TLFVisibility_PRIVATE
	}

	if c.G().Standalone {
		switch c.resolvingRequest.MembersType {
		case chat1.ConversationMembersType_TEAM, chat1.ConversationMembersType_IMPTEAM:
			c.G().StartStandaloneChat()
		default:
			err = CantRunInStandaloneError{}
			return err
		}
	}

	resolver, err := newChatConversationResolver(c.G())
	if err != nil {
		return err
	}
	conversation, _, err := resolver.Resolve(ctx, c.resolvingRequest, chatConversationResolvingBehavior{
		CreateIfNotExists: false,
		MustNotExist:      false,
		Interactive:       c.hasTTY,
		IdentifyBehavior:  keybase1.TLFIdentifyBehavior_CHAT_CLI,
	})
	if err != nil {
		return err
	}
	switch conversation.Info.MembersType {
	case chat1.ConversationMembersType_TEAM:
		return c.showTeam(conversation)
	default:
		return c.showConv(conversation)
	}
}

func (c *CmdChatSetRetention) ParseArgv(ctx *cli.Context) (err error) {
	c.hasTTY = isatty.IsTerminal(os.Stdin.Fd())

	var tlfName string
	// Get the TLF name from the first position arg
	if len(ctx.Args()) >= 1 {
		tlfName = ctx.Args().Get(0)
	}
	if c.resolvingRequest, err = parseConversationResolvingRequest(ctx, tlfName); err != nil {
		return err
	}

	// Send a normal message.
	upto := ctx.Int("upto")
	if upto == 0 {
		cli.ShowCommandHelp(ctx, "delete-history-dev")
		return fmt.Errorf("upto must be > 0")
	}
	c.upto = chat1.MessageID(upto)

	return nil
}

func (c *CmdChatSetRetention) GetUsage() libkb.Usage {
	return libkb.Usage{
		Config: true,
		API:    true,
	}
}

func (c *CmdChatSetRetention) showConv(conversation *chat1.ConversationLocal) error {
	conversation.ConvRetention
	conversation.TeamRetention
	panic("@@@ TODO")
}

func (c *CmdChatSetRetention) showTeam(conversation *chat.ConversationLocal) error {
	conversation.ConvRetention
	conversation.TeamRetention
	panic("@@@ TODO")
}
