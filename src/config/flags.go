package config

import (
	"runtime"

	f "github.com/integrii/flaggy"
)

/*

OP integers:

0: Default value, this runs on auto start-up, and it is when the binary tries to pull a newer version of the database down from cloud onto the local machine.
1: Mode where the binary finds the newest version of your local database, and uploads it to all available remotes.
2: Mode that creates a new remotes configuration, either encrypted or unencrypted.
3: Mode for listing remotes in the available configuration.
4: Mode that edits the configuration, either encrypted or unencrypted.
5: Mode for deleting a configuration of remotes, but could be done manually I suppose.
6. Mode for adding remotes to your configuration.

*/

var OP int = 0 // Operation mode for the binary during this run
// var NewRemote string // String featuring the credentials
var NAME string // Name of remote to add/delete

var Version string = "unknown"    // String to pass in the version to the binary at compiletime.
var Commit string = "unknown"     // Git commit version of this binary.
var Go string = runtime.Version() // Go version at runtime.
var Os string = runtime.GOOS      // operating system for this binary
var Arch string = runtime.GOARCH  // architecture for this binary

func init() {

	f.SetName("Keepassxcync")
	f.SetDescription("A portable binary to automatically sync your keepass/keepassxc databases to multiple remote clouds. ")
	f.SetVersion(Version + "\nGit Commit: " + Commit + "\nGo Version: " + Go + "\nOS: " + Os + "\nArchitecture: " + Arch)

	// remote subcommands.
	add := f.NewSubcommand("add")
	add.Description = "Use this subcommand if you want to add a remote to your current configuration."
	add.String(&NAME, "n", "name", "The name of the remote you want to add.")

	delete := f.NewSubcommand("delete")
	delete.Description = "Use this subcommand if you want to delete a remote from your current configuration."
	delete.String(&NAME, "n", "name", "The name of the remote you want to delete.")

	list := f.NewSubcommand("list")

	remote := f.NewSubcommand("remote")
	remote.Description = "Add, remove, and list the remotes that are configured to be written/read from when updating/uplloading your database."
	f.AttachSubcommand(remote, 1)

	remote.AttachSubcommand(add, 1)
	remote.AttachSubcommand(delete, 1)
	remote.AttachSubcommand(list, 1)

	f.Parse()

	switch {
	case remote.Used:
		{
			if list.Used {
				OP = 3
			} else if add.Used {
				OP = 6
			} else if delete.Used {
				OP = 5
			}
		}
	}

}
