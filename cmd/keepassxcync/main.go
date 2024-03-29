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
	"runtime"

	"github.com/fire833/keepassxcync/cmd/keepassxcync/app"
)

var (
	Version string = "unknown"         // String to pass in the version to the binary at compiletime.
	Commit  string = "unknown"         // Git commit version of this binary.
	Go      string = runtime.Version() // Go version at runtime.
	Os      string = runtime.GOOS      // operating system for this binary
	Arch    string = runtime.GOARCH    // architecture for this binary
)

func main() {
	cmd := app.NewKPXCCommand()
	cmd.Execute()
}
