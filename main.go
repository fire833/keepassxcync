package main

import (
	"runtime"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fire833/keepassxcync/src"
	op "github.com/fire833/keepassxcync/src"
	"github.com/fire833/keepassxcync/src/config/debug"
	f "github.com/integrii/flaggy"
)

/*

OP integers:

0: Default value, this tries to sync the newest version of the database either from the default remote
	to your local machine or if the local version is newer, upload that to the default remote.
1: Try to upload your local database to all configured remotes.

// Operate on the current configuration
2: List remotes from current config.
3: Add remotes to curret config
4: Delete remotes from current config
5: Update remotes from current config

6: Try to sync your database with the database of a specified remote.

*/

var OP int = 0 // Operation mode for the binary during this run

var DBNAME string // name of database file to interact upon
var DEFAULT bool  // Var to pass around whether or not to set a remote as the default value.

var Version string = "unknown"    // String to pass in the version to the binary at compiletime.
var Commit string = "unknown"     // Git commit version of this binary.
var Go string = runtime.Version() // Go version at runtime.
var Os string = runtime.GOOS      // operating system for this binary
var Arch string = runtime.GOARCH  // architecture for this binary

var globalS3 *s3.S3

func main() {

	f.SetName("Keepassxcync")
	f.SetDescription("A portable binary to automatically sync your keepass/keepassxc databases to multiple remote clouds. ")
	f.SetVersion(Version + "\nGit Commit: " + Commit + "\nGo Version: " + Go + "\nOS: " + Os + "\nArchitecture: " + Arch)

	// remote subcommands.
	add := f.NewSubcommand("add")
	add.Description = "Use this subcommand if you want to add a remote to your current configuration."
	add.String(&src.NAME, "n", "name", "The name of the remote you want to add.")
	add.String(&src.CONFIG, "f", "file", "Specify the config file/path to config file you want to utilize.")

	delete := f.NewSubcommand("delete")
	delete.Description = "Use this subcommand if you want to delete a remote from your current configuration."
	delete.String(&src.NAME, "n", "name", "The name of the remote you want to delete.")
	delete.String(&src.CONFIG, "f", "file", "Specify the config file/path to config file you want to utilize.")

	list := f.NewSubcommand("list")
	list.Description = "Use this subcommand if you want to list available remotes."
	list.String(&src.CONFIG, "f", "file", "Specify the config file/path to config file you want to utilize.")

	update := f.NewSubcommand("update")
	update.Description = "Use this subcommand to update the values of a specific remote."
	update.String(&src.NAME, "n", "name", "The name of the remote you want to update.")
	update.Bool(&DEFAULT, "d", "set-default", "Specify that this remote should be set as the default remote.")
	update.String(&src.CONFIG, "f", "file", "Specify the config file/path to config file you want to utilize.")

	remote := f.NewSubcommand("remote")
	remote.Description = "Add, remove, list, and update the remotes that are configured to be written/read from when updating/uploading your database."
	f.AttachSubcommand(remote, 1)

	all := f.NewSubcommand("all")
	all.Description = "Use this to specify synching up to all of the remotes in your options file."

	sync := f.NewSubcommand("sync")
	sync.Description = "Sync up your database with specified remote. It either pulls a newer version from the specified remote or pushes your local version if it is the most up-to-date."
	sync.String(&src.NAME, "n", "name", "Specify the name of the remote you want to sync with, otherwise defaults to the remote with default bool set to true.")
	sync.String(&DBNAME, "f", "file", "Specify the specific database file to sync up.")
	sync.String(&src.CONFIG, "c", "config", "Specify the config file/path to config file you want to utilize.")
	f.AttachSubcommand(sync, 1)
	sync.AttachSubcommand(all, 1)

	remote.AttachSubcommand(add, 1)
	remote.AttachSubcommand(delete, 1)
	remote.AttachSubcommand(list, 1)
	remote.AttachSubcommand(update, 1)

	f.Parse()

	opts := op.NewOptions()

	switch {
	case remote.Used:
		{
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

				}
			case update.Used:
				{

				}
			}
		}
	case sync.Used:
		{

		}
	}

}
