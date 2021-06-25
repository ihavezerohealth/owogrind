// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"
	"time"

	"github.com/ihavezerohealth/owogrind/instance/scheduler"
)

const (
	prayCmdValue     = "owo pray"
	huntCmdValue     = "owo hunt"
	sellBaseCmdValue = "owo sell"
)

func sellCmdValue(amount, item string) string {
	return fmt.Sprintf("%v %v %v", sellBaseCmdValue, item, amount)
}

func shareCmdValue(amount, id string) string {
	return fmt.Sprintf("%v %v <@%v>", shareBaseCmdValue, amount, id)
}

// commands returns a command pointer slice with all commands that should be
// executed periodically. It contains all commands as configured.
func (in *Instance) newCmds() []*scheduler.Command {
	var cmds []*scheduler.Command
	if in.Features.Commands.Pray {
		cmds = append(cmds, &scheduler.Command{
			Value:       prayCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Pray) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Hunt {
		cmds = append(cmds, &scheduler.Command{
			Value:       huntCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Hunt) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.AutoSell

	for _, cmd := range in.Features.CustomCommands {
		// cmd.Value and cmd.Amount are not checked for correct values here
		// because they were checked when the application started using
		// cfg.Validate().
		cmds = append(cmds, &scheduler.Command{
			Value:    cmd.Value,
			Interval: time.Duration(cmd.Interval) * time.Second,
			Amount:   uint(cmd.Amount),
			CondFunc: func() bool {
				return cmd.PauseBelowBalance == 0 || in.balance >= cmd.PauseBelowBalance
			},
		})
	}
	return cmds
}

func (in *Instance) newAutoSellChain() *scheduler.Command {
	var cmds []*scheduler.Command
	for _, item := range in.Features.AutoSell.Items {
		cmds = append(cmds, &scheduler.Command{
			Value:    sellCmdValue("max", item),
			Interval: time.Duration(in.Compat.Cooldown.Sell) * time.Second,
		})
	}
	return in.newCmdChain(
		cmds,
		time.Duration(in.Features.AutoSell.Interval)*time.Second,
	)
}

func (in *Instance) newCmdChain(cmds []*scheduler.Command, chainInterval time.Duration) *scheduler.Command {
	for i := 0; i < len(cmds); i++ {
		if i != 0 {
			cmds[i-1].Next = cmds[i]
		}
		if i == len(cmds)-1 {
			cmds[i].Next = cmds[0]
			cmds[i].Interval = chainInterval
		}
	}
	return cmds[0]
}
