package main

import f "github.com/integrii/flaggy"

var OP string // Operation mode for the binary during this run
// var NewRemote string // String featuring the credentials
var NAME string // Name of remote to add/delete

var VERSION string // String to pass in the version to the binary at compiletime.

func init() {

	f.SetName("Keepassxcync")
	f.SetDescription("A portable binary to automatically sync your keepass/keepassxc databases to multiple remote clouds. ")
	f.SetVersion(VERSION)

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
}
