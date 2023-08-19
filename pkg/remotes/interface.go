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

package remotes

import (
	"context"
	"io"
)

// A remote should be considered an object store that is able to store all
// versions of the database that are uploaded to it, and be able to reference
// a specific version, including the latest version on that remote.
type Remote interface {
	PersistVersion(ctx context.Context, data io.WriteCloser) error
	GetVersion(ctx context.Context, version uint) (io.ReadCloser, error)
	GetLastVersion(ctx context.Context) (uint, error)
}
