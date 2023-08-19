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

package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type KeepassxCyncConfig struct {
	filePath string      `json:"-" yaml:"-"`
	perms    fs.FileMode `json:"-" yaml:"-"`

	Remotes        []*KeepassxCyncRemote   `json:"remotes" yaml:"remotes"`
	ActiveRemote   string                  `json:"activeRemote" yaml:"activeRemote"`
	Databases      []*KeepassxCyncDatabase `json:"dbs" yaml:"dbs"`
	ActiveDatabase string                  `json:"activeDb" yaml:"activeDb"`
}

type KeepassxCyncRemote struct {
	Name string `json:"name" yaml:"name"`
}

type KeepassxCyncDatabase struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}

func Load(path string) (*KeepassxCyncConfig, error) {
	data, e := os.ReadFile(path)
	if e != nil {
		return nil, e
	}

	conf := &KeepassxCyncConfig{}

	switch filepath.Ext(path) {
	case ".json":
		if e := json.Unmarshal(data, conf); e != nil {
			return nil, e
		}
	case ".yaml", ".yml":
		if e := yaml.Unmarshal(data, conf); e != nil {
			return nil, e
		}
	default:
		return nil, errors.New("extension must be .json or .yaml")
	}

	conf.filePath = path
	return conf, nil
}

func (c *KeepassxCyncConfig) Flush() error {
	switch filepath.Ext(c.filePath) {
	case ".json":
		if bytes, e := json.MarshalIndent(c, "", "	"); e != nil {
			return e
		} else {
			return os.WriteFile(c.filePath, bytes, c.perms)
		}
	case ".yaml", ".yml":
		if bytes, e := yaml.Marshal(c); e != nil {
			return e
		} else {
			return os.WriteFile(c.filePath, bytes, c.perms)
		}
	default:
		return errors.New("extension must be .json or .yaml")
	}
}
