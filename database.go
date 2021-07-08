package main

import "io/fs"

type Database struct {
	Entry fs.DirEntry
}

func (d *Database) something() {

}
