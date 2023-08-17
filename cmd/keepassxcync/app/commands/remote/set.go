/*
*	Copyright (C) 2023 Kendall Tauser
*
*	This program is free software; you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation; either version 2 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be useful,
*	but WITHOUT ANY WARRANTY; without even the implied warranty of
*	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*	GNU General Public License for more details.
*
*	You should have received a copy of the GNU General Public License along
*	with this program; if not, write to the Free Software Foundation, Inc.,
*	51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package remote

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewSETCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set",
		Aliases: []string{},
		Short:   "",
		Long:    "",
		Version: "0.0.1",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	set := pflag.NewFlagSet("set", pflag.ExitOnError)

	cmd.Flags().AddFlagSet(set)
	cmd.AddCommand()

	return cmd
}
