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
	"bytes"
	"testing"
)

func TestUserInputNotNull(t *testing.T) {
	tests := []struct {
		name string
		s    *bufio.Scanner
		want string
	}{
		{
			name: "1",
			s:    bufio.NewScanner(bytes.NewReader([]byte("something\n"))),
			want: "something",
		},
		{
			name: "2",
			s:    bufio.NewScanner(bytes.NewReader([]byte("\n\n\n\nsomething else\n"))),
			want: "something else",
		},
		{
			name: "3",
			s:    bufio.NewScanner(bytes.NewReader([]byte("\n\n\n\n\n\n\n\nhi there\n"))),
			want: "hi there",
		},
		{
			name: "4",
			s:    bufio.NewScanner(bytes.NewReader([]byte("hi there\n"))),
			want: "hi there",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserInputNotNull(tt.s); got != tt.want {
				t.Errorf("UserInputNotNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInputNotNullDefault(t *testing.T) {
	tests := []struct {
		name string
		s    *bufio.Scanner
		def  string
		want string
	}{
		{
			name: "1",
			s:    bufio.NewScanner(bytes.NewReader([]byte("hi there\n"))),
			def:  "something",
			want: "hi there",
		},
		{
			name: "2",
			s:    bufio.NewScanner(bytes.NewReader([]byte("\n\nhi there\n"))),
			def:  "something",
			want: "something",
		},
		{
			name: "3",
			s:    bufio.NewScanner(bytes.NewReader([]byte("\n\nvalue\n"))),
			def:  "value2",
			want: "value2",
		},
		{
			name: "4",
			s:    bufio.NewScanner(bytes.NewReader([]byte("value\n"))),
			def:  "value2",
			want: "value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserInputNotNullDefault(tt.s, tt.def); got != tt.want {
				t.Errorf("UserInputNotNullDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInputOptional(t *testing.T) {
	tests := []struct {
		name string
		s    *bufio.Scanner
		want string
	}{
		{
			name: "1",
			s:    bufio.NewScanner(bytes.NewReader([]byte("\n"))),
			want: "",
		},
		{
			name: "2",
			s:    bufio.NewScanner(bytes.NewReader([]byte("something\n"))),
			want: "something",
		},
		{
			name: "3",
			s:    bufio.NewScanner(bytes.NewReader([]byte("\nsomething\n"))),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserInputOptional(tt.s); got != tt.want {
				t.Errorf("UserInputOptional() = %v, want %v", got, tt.want)
			}
		})
	}
}
