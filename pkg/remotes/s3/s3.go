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

package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Remote struct {
	cfg aws.Config

	s3client *s3.Client
}

func New() (*S3Remote, error) {
	c, e := config.LoadDefaultConfig(context.Background())
	if e != nil {
		return nil, e
	}

	client := s3.NewFromConfig(c)
	return &S3Remote{cfg: c, s3client: client}, nil
}

func (r *S3Remote) PersistVersion(ctx context.Context, data io.WriteCloser) error {
	_, e := r.s3client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{})
	if e != nil {
		return e
	}

	return nil
}

func (r *S3Remote) GetVersion(ctx context.Context, version uint) (io.ReadCloser, error) {
	return nil, nil
}

func (r *S3Remote) GetLastVersion(ctx context.Context) (uint, error) {
	return 0, nil
}
