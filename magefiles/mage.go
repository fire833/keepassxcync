//go:build mage

/*
*	Copyright (C) 2022  Kendall Tauser
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
	"time"

	"github.com/magefile/mage/sh"
)

func Build() error {
	version := os.Getenv("VERSION")
	commit, e := sh.Output("git", "rev-parse", "HEAD")
	if e != nil {
		return e
	}
	date := time.Now().String()

	return sh.Run("go", "build", "-o",
		"bin/keepassxcync", "-trimpath", "-ldflags",
		fmt.Sprintf(`-X 'main.Version=%s' -X 'main.Commit=%s' -X 'main.BuildTime=%s' -aslr`, version, commit, date),
		"cmd/keepassxcync/main.go",
	)
}
