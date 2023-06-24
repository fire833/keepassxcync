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

func TestYesOrNo(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    bool
		wantErr bool
	}{
		{
			name:    "1",
			s:       "Y",
			want:    true,
			wantErr: false,
		},
		{
			name:    "2",
			s:       "N",
			want:    false,
			wantErr: false,
		},
		{
			name:    "3",
			s:       "Sure",
			want:    false,
			wantErr: true,
		},
		{
			name:    "4",
			s:       "no",
			want:    false,
			wantErr: false,
		},
		{
			name:    "5",
			s:       "yes",
			want:    true,
			wantErr: false,
		},
		{
			name:    "6",
			s:       "YES",
			want:    true,
			wantErr: false,
		},
		{
			name:    "7",
			s:       "NOPE",
			want:    false,
			wantErr: false,
		},
		{
			name:    "8",
			s:       "8934tyjdh",
			want:    false,
			wantErr: true,
		},
		{
			name:    "9",
			s:       "734856hnasjkdg",
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YesOrNo(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("YesOrNo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("YesOrNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYesOrNoBufIO(t *testing.T) {
	tests := []struct {
		name    string
		s       *bufio.Scanner
		want    bool
		wantErr bool
	}{
		{
			name:    "1",
			s:       nil,
			want:    false,
			wantErr: true,
		},
		{
			name:    "2",
			s:       bufio.NewScanner(bytes.NewBuffer([]byte("Y\n"))),
			want:    true,
			wantErr: false,
		},
		{
			name:    "3",
			s:       bufio.NewScanner(bytes.NewBuffer([]byte("N\n"))),
			want:    false,
			wantErr: false,
		},
		{
			name:    "4",
			s:       bufio.NewScanner(bytes.NewBuffer([]byte("something\n"))),
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YesOrNoBufIO(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("YesOrNoBufIO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("YesOrNoBufIO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYesOrNoDefault(t *testing.T) {
	tests := []struct {
		name string
		s    string
		d    bool
		want bool
	}{
		{
			name: "1",
			s:    "Y",
			d:    false,
			want: true,
		},
		{
			name: "2",
			s:    "N",
			d:    true,
			want: false,
		},
		{
			name: "3",
			s:    "soem,eihfsd",
			d:    true,
			want: true,
		},
		{
			name: "4",
			s:    "gmnkdffiy",
			d:    false,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := YesOrNoDefault(tt.s, tt.d); got != tt.want {
				t.Errorf("YesOrNoDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYesOrNoBufIODefault(t *testing.T) {
	tests := []struct {
		name string
		s    *bufio.Scanner
		d    bool
		want bool
	}{
		{
			name: "1",
			s:    bufio.NewScanner(bytes.NewReader([]byte("Y\n"))),
			d:    false,
			want: true,
		},
		{
			name: "2",
			s:    bufio.NewScanner(bytes.NewReader([]byte("N\n"))),
			d:    true,
			want: false,
		},
		{
			name: "3",
			s:    bufio.NewScanner(bytes.NewReader([]byte("soomryhet\n"))),
			d:    true,
			want: true,
		},
		{
			name: "4",
			s:    bufio.NewScanner(bytes.NewReader([]byte("j54y89y\n"))),
			d:    false,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := YesOrNoBufIODefault(tt.s, tt.d); got != tt.want {
				t.Errorf("YesOrNoBufIODefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
