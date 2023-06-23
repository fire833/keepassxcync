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

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go/service/s3"
	pkg "github.com/fire833/keepassxcync/pkg"
	"github.com/fire833/keepassxcync/pkg/config/debug"
	f "github.com/integrii/flaggy"
)

var OP int = 0 // Operation mode for the binary during this run

var (
	DBNAME  string // name of database file to interact upon
	DEFAULT bool   // Var to pass around whether or not to set a remote as the default value.
)

var (
	Version string = "unknown"         // String to pass in the version to the binary at compiletime.
	Commit  string = "unknown"         // Git commit version of this binary.
	Go      string = runtime.Version() // Go version at runtime.
	Os      string = runtime.GOOS      // operating system for this binary
	Arch    string = runtime.GOARCH    // architecture for this binary
)

var globalS3 *s3.S3

func main() {
	f.SetName("Keepassxcync")
	f.SetDescription("A portable binary to automatically sync your keepass/keepassx/keepassxc databases to multiple remote clouds. ")
	f.SetVersion(Version + "\nGit Commit: " + Commit + "\nGo Version: " + Go + "\nOS: " + Os + "\nArchitecture: " + Arch)

	// remote subcommands.
	add := f.NewSubcommand("add")
	add.Description = "Use this subcommand if you want to add a remote to your current configuration."
	add.String(&pkg.NAME, "n", "name", "The name of the remote you want to add.")
	add.String(&pkg.CONFIG, "f", "file", "Specify the config file/path to config file you want to utilize.")

	delete := f.NewSubcommand("delete")
	delete.Description = "Use this subcommand if you want to delete a remote from your current configuration."
	delete.String(&pkg.NAME, "n", "name", "The name of the remote you want to delete.")
	delete.String(&pkg.CONFIG, "f", "file", "Specify the config file/path to config file you want to utilize.")

	list := f.NewSubcommand("list")
	list.Description = "Use this subcommand if you want to list available remotes."
	list.String(&pkg.CONFIG, "f", "file", "Specify the config file/path to config file you want to utilize.")

	update := f.NewSubcommand("update")
	update.Description = "Use this subcommand to update the values of a specific remote."
	update.String(&pkg.NAME, "n", "name", "The name of the remote you want to update.")
	update.Bool(&DEFAULT, "d", "set-default", "Specify that this remote should be set as the default remote.")
	update.String(&pkg.CONFIG, "f", "file", "Specify the config file/path to config file you want to utilize.")

	remote := f.NewSubcommand("remote")
	remote.Description = "Add, remove, list, and update the remotes that are configured to be written/read from when updating/uploading your database."
	f.AttachSubcommand(remote, 1)

	all := f.NewSubcommand("all")
	all.Description = "Use this to specify synching up to all of the remotes in your options file."

	sync := f.NewSubcommand("sync")
	sync.Description = "Sync up your database with specified remote. It either pulls a newer version from the specified remote or pushes your local version if it is the most up-to-date."
	sync.String(&pkg.NAME, "n", "name", "Specify the name of the remote you want to sync with, otherwise defaults to the remote with default bool set to true.")
	sync.String(&DBNAME, "f", "file", "Specify the specific database file to sync up.")
	sync.String(&pkg.CONFIG, "c", "config", "Specify the config file/path to config file you want to utilize.")
	f.AttachSubcommand(sync, 1)
	sync.AttachSubcommand(all, 1)

	remote.AttachSubcommand(add, 1)
	remote.AttachSubcommand(delete, 1)
	remote.AttachSubcommand(list, 1)
	remote.AttachSubcommand(update, 1)

	f.Parse()

	opts := pkg.NewOptions()
	defer opts.File.Close()

	if remote.Used {
		switch {
		case list.Used:
			{
				if debug.DEBUG == true {
					opts.PrintRemotes(true)
				} else {
					opts.PrintRemotes(false)
				}
			}
		case add.Used:
			{
				opts.AddRemote()
			}
		case delete.Used:
			{
				opts.RemoveRemote()
			}
		case update.Used:
			{
				opts.UpdateRemote()
			}
		}
	} else if sync.Used {
		// sync stuff up
		err := opts.PushPull()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err1 := opts.PushPull()
		if err1 != nil {
			fmt.Println(err1)
		}
	}

	os.Exit(0)
}
