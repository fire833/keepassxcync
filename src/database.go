package src

import (
	"errors"
	"io/fs"
	"os"
	fp "path/filepath"

	"github.com/aws/aws-sdk-go/service/s3"
)

// Finds the newest version of the key database within the desired directory,
// defaults to the working directory of the binary.
func (o *Options) GetNewestLocalDb(dir string) (name string, file *os.File, e error) {
	if dir == "" {
		dir = os.Getenv("PWD")
	}

	var files map[int64]fs.DirEntry
	var times []int64

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", nil, err
	}

	for _, entry := range entries {
		if fp.Ext(entry.Name()) == ".kdbx" && entry.Name() == o.DatabaseName {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			files[info.ModTime().UnixNano()] = entry
			times = append(times, info.ModTime().UnixNano())
		}
	}

	if len(files) == 0 {
		return "", nil, errors.New("No database files found in directory " + dir)
	}

	name = files[findBiggestTime(times)].Name()

	out, err := os.Open(fp.Join(dir, name))
	if err != nil {
		return files[findBiggestTime(times)].Name(), out, err
	}

	return name, out, nil

}

// Looks at the default remote in the options file and then returns the file info/modtime for
// the newest
func (o *Options) GetNewestRemoteDB(client *s3.S3) (name string, file *s3.Object, e error) {

	var defremote *Remote

	for _, remote := range o.Remotes {
		if remote.IsDefault {
			defremote = &remote
		}
	}

	in := &s3.ListObjectsV2Input{
		Bucket: &defremote.Bucket,
		Prefix: &o.DatabaseName,
	}

	out, err := client.ListObjectsV2(in)
	if err != nil {
		return "", nil, err
	}

	var files map[int64]*s3.Object
	var times []int64

	for _, object := range out.Contents {
		files[object.LastModified.UnixNano()] = object
		times = append(times, object.LastModified.UnixNano())
	}

	if len(files) == 0 {
		return "", nil, errors.New("No database files found in remote " + defremote.Name)
	}

	name = files[findBiggestTime(times)].String()

	return name, files[findBiggestTime(times)], nil

}

func findBiggestTime(array []int64) (largest int64) {
	largest = array[0]
	for _, num := range array {
		if num > largest {
			largest = num
		}
	}
	return largest
}
