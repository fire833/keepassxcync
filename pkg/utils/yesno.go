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

package utils

import (
	"bufio"
	"errors"
	"regexp"
)

var (
	yes *regexp.Regexp = regexp.MustCompile("^[Yy].*")
	no  *regexp.Regexp = regexp.MustCompile("^[Nn].*")
)

// Wrapper around YesOrNobufIO, but if there is an error with parsing, then this method returns
// the default value instead of failing.
func YesOrNoBufIODefault(s *bufio.Scanner, d bool) bool {
	if yn, e := YesOrNoBufIO(s); e == nil {
		return yn
	} else {
		return d
	}
}

// Wrapper around YesOrNo, but gets text from a bufio.Scanner interface.
func YesOrNoBufIO(s *bufio.Scanner) (bool, error) {
	if s == nil {
		return false, errors.New("bufio scanner is nil")
	}

	if s.Scan() {
		return YesOrNo(s.Text())
	} else {
		return false, errors.New("no bytes left to scan")
	}
}

// Wrapper around YesOrNo, but if there is an error with parsing, then this method returns
// the default value instead of failing.
func YesOrNoDefault(s string, d bool) bool {
	if yn, e := YesOrNo(s); e == nil {
		return yn
	} else {
		return d
	}
}

// Parses a string from the user as whether or not the user is specifying yes or no.
// Or returns an error if a boolean decision cannot be evaluated.
func YesOrNo(s string) (bool, error) {
	if yes.Match([]byte(s)) {
		return true, nil
	} else if no.Match([]byte(s)) {
		return false, nil
	} else {
		return false, errors.New("unable to parse whether the use wants yes or no")
	}
}
