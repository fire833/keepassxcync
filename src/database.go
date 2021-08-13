package src

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	fp "path/filepath"
	"regexp"

	"github.com/aws/aws-sdk-go/service/s3"
)

// Finds the newest version of the key database within the desired directory,
// defaults to the working directory of the binary.
func (o *Options) GetNewestLocalDB() (name string, file *os.File, e error) {

	switch {
	case o.DatabaseName == "" && o.DatabaseRegex == "":
		{

			var files map[int64]fs.DirEntry
			var times []int64

			entries, err := os.ReadDir(fp.Dir(o.FilePath))
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
				return "", nil, errors.New("No database files found in directory " + fp.Dir(o.FilePath))
			}

			name = files[findBiggestTime(times)].Name()

			out, err := os.Open(fp.Join(fp.Dir(o.FilePath), name))
			if err != nil {
				return files[findBiggestTime(times)].Name(), out, err
			}

			return name, out, nil
		}
	case o.DatabaseRegex != "" && o.DatabaseName == "" || o.DatabaseName != "" && o.DatabaseRegex != "":
		{

			var files map[int64]fs.DirEntry
			var times []int64

			regex, err := regexp.Compile(o.DatabaseRegex)
			if err != nil {
				fmt.Printf("Error with parsing your regex statements for finding")
			}

			entries, err := os.ReadDir(fp.Dir(o.FilePath))
			if err != nil {
				return "", nil, err
			}

			for _, entry := range entries {
				if regex.MatchString(entry.Name()) {
					info, err := entry.Info()
					if err != nil {
						continue
					}

					files[info.ModTime().UnixNano()] = entry
					times = append(times, info.ModTime().UnixNano())
				}
			}

			if len(files) == 0 {
				return "", nil, errors.New("No database files found in directory " + fp.Dir(o.FilePath))
			}

			name = files[findBiggestTime(times)].Name()

			out, err := os.Open(fp.Join(fp.Dir(o.FilePath), name))
			if err != nil {
				return files[findBiggestTime(times)].Name(), out, err
			}

			return name, out, nil

		}
	case o.DatabaseName != "" && o.DatabaseRegex == "":
		{
			file, err := os.Open(fp.Join(fp.Dir(o.FilePath), o.DatabaseName))
			if err != nil {
				return o.DatabaseName, nil, err
			}

			return o.DatabaseName, file, nil
		}
	}
	return "", nil, errors.New("Something weird happened...")
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

type LocalFileSummary struct {
	// Name of object
	Name string
	// Pointer to some file object that is either an s3 object version or an os.File
	File *os.File
	// Error if there was any with finding the operation.
	Error error
}

type RemoteFileSumary struct {
	// Name of object
	Name string
	// Pointer to some file object that is either an s3 object version or an os.File
	File *s3.ObjectVersion
	// Error if there was any with finding the operation.
	Error error
}

// Main function that either triggers a pull of the newer version of the database from
// the default remote, or will upload the current version of the database locally to the remote.
// Designed to be used by default running of the binary or in "sync" mode.
func (o *Options) PushPull() error {

	localinfo := make(chan LocalFileSummary)
	remoteinfo := make(chan RemoteFileSumary)

	go func(in string) {
		localname, localfile, e := o.GetNewestLocalDB()
		info := LocalFileSummary{
			Name:  localname,
			File:  localfile,
			Error: e,
		}

		localinfo <- info

	}(fp.Dir(o.FilePath))

	client, err5 := o.NewS3Client("", true)
	if err5 != nil {
		fmt.Printf("Error with creating default client: %v\n", err5)
	}

	o.DefClient = client

	go func(sess *s3.S3) {
		remotename, remotefile, e := o.GetNewestRemoteDB(sess)
		info := RemoteFileSumary{
			Name:  remotename,
			File:  remotefile,
			Error: e,
		}

		remoteinfo <- info
	}(o.DefClient.Client)

	// I hate all these variables, probably stupid to do all that crap to make the
	// local and remote search concurrent with each other, but whatever.
	linf := <-localinfo
	rinf := <-remoteinfo

	info, _ := linf.File.Stat()
	lmod := info.ModTime()
	rmod := *rinf.File.LastModified

	switch {
	case linf.Error != nil:
		{
			fmt.Printf("Error with getting local file information: %v", linf.Error)
		}
	case rinf.Error != nil:
		{
			fmt.Printf("Error with getting remote file information: %v", rinf.Error)
		}
	case lmod.UnixNano() > rmod.UnixNano():
		{

		}
	case lmod.UnixNano() < rmod.UnixNano():
		{

		}
	case lmod.UnixNano() == rmod.UnixNano():
		{
			fmt.Println("Latest version of local binary syncs with the remote binary.")
			os.Exit(0)
		}
	}

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
