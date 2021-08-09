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
func (o *Options) GetNewestLocalDB(dir string) (name string, file *os.File, e error) {

	if o.DatabaseName == "" {

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
	} else {

		file, err := os.Open(fp.Join(dir, o.DatabaseName))
		if err != nil {
			return o.DatabaseName, nil, err
		}

		return o.DatabaseName, file, nil

	}

}

// Looks at the default remote in the options file and then returns the file info/modtime for
// the newest
func (o *Options) GetNewestRemoteDB(client *s3.S3) (name string, file *s3.ObjectVersion, e error) {

	for _, remote := range o.Remotes {
		if remote.IsDefault {

			in := &s3.ListObjectVersionsInput{
				Bucket: &remote.Bucket,
				Prefix: &o.DatabaseName,
			}

			out, err := client.ListObjectVersions(in)
			if err != nil {
				return "", nil, err
			}

			for _, object := range out.Versions {
				if *object.IsLatest == true {
					return *object.Key, object, nil
				}
			}

			return "", nil, errors.New("No version of your database found in bucket " + remote.Bucket)

		}
	}

	return "", nil, errors.New("No default remote defined, unable to search for newest file.")

}

// Main function that either triggers a pull of the newer version of the database from
// the default remote, or will upload the current version of the database locally to the remote.
func (o *OptionMeta) PushPull(client *s3.S3) error {

	go func(in string) (localname string, localfile *os.File, e error) {
		return o.Options.GetNewestLocalDB(in)
	}(os.Getenv("PWD"))

	go func(sess *s3.S3) (remotename string, remotefile *s3.ObjectVersion, e error) {
		return o.Options.GetNewestRemoteDB(sess)
	}(client)

	return nil
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
