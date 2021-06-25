// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package config

import (
	"fmt"
	"regexp"
	"strings"
)

func (c Config) Validate() error {
	if len(c.Clusters) == 0 {
		return fmt.Errorf("clusters: no clusters, at least 1 is required")
	}
	for ck, cluster := range c.Clusters {
		if err := validateInstance(cluster.Master); err != nil {
			return fmt.Errorf("clusters[%v].master: %v", ck, err)
		}
		for i, instance := range cluster.Instances {
			if err := validateInstance(instance); err != nil {
				return fmt.Errorf("clusters[%v].instances[%v]: %v", ck, i, err)
			}
		}
	}
	if err := validateCompat(c.Compat); err != nil {
		return err
	}
	return nil
}

func validateInstance(instance Instance) error {
	if instance.Token == "" {
		return fmt.Errorf("no token")
	}
	if instance.ChannelID == "" {
		return fmt.Errorf("no channel id")
	}
	if !isValidID(instance.ChannelID) {
		return fmt.Errorf("invalid channel id")
	}
	if len(instance.Shifts) == 0 {
		return fmt.Errorf("no shifts")
	}
	if err := validateFeatures(instance.Features); err != nil {
		return err
	}
	if err := validateShifts(instance.Shifts); err != nil {
		return err
	}
	return nil
}

func validateFeatures(features Features) error {
	if features.AutoSell.Enable {
		if features.AutoSell.Interval < 0 {
			return fmt.Errorf("auto-sell interval must be greater than or equal to 0")
		}
		if len(features.AutoSell.Items) == 0 {
			return fmt.Errorf("auto-sell enabled but no items configured")
		}
	}
	if features.AutoShare.Enable {
		if features.AutoShare.Amount <= 0 {
			return fmt.Errorf("auto-share amount must be greater than 0")
		}
	}

	for i, cmd := range features.CustomCommands {
		if cmd.Value == "" {
			return fmt.Errorf("features.custom_commands[%v].value: no value", i)
		}
		if strings.Contains(cmd.Value, "owo hunt") {
			return fmt.Errorf("invalid custom command value: %v, this custom command is disallowed, use auto-gift instead", cmd.Value)
		}
		if strings.Contains(cmd.Value, "owoh") {
			return fmt.Errorf("invalid custom command value: %v, this custom command is disallowed, use auto-sell instead", cmd.Value)
		}
		if cmd.Amount < 0 {
			return fmt.Errorf("features.custom_commands[%v].amount: value must be greater than or equal to 0", i)
		}
	}
	return nil
}

func validateShifts(shifts []Shift) error {
	for _, shift := range shifts {
		if shift.State != ShiftStateActive && shift.State != ShiftStateDormant {
			return fmt.Errorf("invalid shift state: %v", shift.State)
		}
	}
	return nil
}

func validateCompat(compat Compat) error {
	if compat.Cooldown.Pray <= 0 {
		return fmt.Errorf("pray cooldown must be greater than 0")
	}
	if compat.Cooldown.Hunt <= 0 {
		return fmt.Errorf("hunt cooldown must be greater than 0")
	}
	if compat.Cooldown.Share <= 0 {
		return fmt.Errorf("share cooldown must be greater than 0")
	}
	if compat.AwaitResponseTimeout < 0 {
		return fmt.Errorf("await response timeout must be greater than 0")
	}
	return nil
}

func isValidID(id string) bool {
	return regexp.MustCompile(`^[0-9]+$`).Match([]byte(id))
}
