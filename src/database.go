package src

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	fp "path/filepath"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Finds the newest version of the key database within the desired directory,
// defaults to the working directory of the binary.
func (o *OptionMeta) GetNewestLocalDB() (name string, file *os.File, e error) {

	switch {
	case o.Options.DatabaseName == "" && o.Options.DatabaseRegex == "":
		{

			var files map[int64]fs.DirEntry
			var times []int64

			entries, err := os.ReadDir(fp.Dir(o.FilePath))
			if err != nil {
				return "", nil, err
			}

			for _, entry := range entries {
				if fp.Ext(entry.Name()) == ".kdbx" && entry.Name() == o.Options.DatabaseName {
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
	case o.Options.DatabaseRegex != "" && o.Options.DatabaseName == "" || o.Options.DatabaseName != "" && o.Options.DatabaseRegex != "":
		{

			var files map[int64]fs.DirEntry
			var times []int64

			regex, err := regexp.Compile(o.Options.DatabaseRegex)
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
	case o.Options.DatabaseName != "" && o.Options.DatabaseRegex == "":
		{
			file, err := os.Open(fp.Join(fp.Dir(o.FilePath), o.Options.DatabaseName))
			if err != nil {
				return o.Options.DatabaseName, nil, err
			}

			return o.Options.DatabaseName, file, nil
		}
	}
	return "", nil, errors.New("Something weird happened...")
}

// Looks at the default remote in the options file and then returns the file info/modtime for
// the newest
func (o *OptionMeta) GetNewestRemoteDB() (name string, file *s3.ObjectVersion, client *Client, e error) {

	for _, remote := range o.Options.Remotes {
		if remote.IsDefault {

			conf := &aws.Config{
				Endpoint:    &remote.Endpoint,
				Credentials: credentials.NewStaticCredentials(remote.Id, remote.Key, ""),
				Region:      &remote.Region,
			}

			// Create the client here as it is needed and then will be garbage collected as needed too.
			cli := s3.New(session.New(conf))

			client := &Client{
				Client:        cli,
				RemoteOptions: remote,
			}

			in := &s3.ListObjectVersionsInput{
				Bucket: &remote.Bucket,
				Prefix: &o.Options.DatabaseName,
			}

			out, err := cli.ListObjectVersions(in)
			if err != nil {
				return "", nil, client, err
			}

			for _, object := range out.Versions {
				if *object.IsLatest == true {
					return *object.Key, object, client, nil
				}
			}

			return "", nil, client, errors.New("No version of your database found in bucket " + remote.Bucket)

		}
	}

	return "", nil, client, errors.New("No default remote defined, unable to search for newest file.")

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
	// Return client that should be used for future operations.
	Client *Client
}

// Main function that either triggers a pull of the newer version of the database from
// the default remote, or will upload the current version of the database locally to the remote.
// Designed to be used by default running of the binary or in "sync" mode.
func (o *OptionMeta) PushPull() error {

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

	go func() {
		remotename, remotefile, client, e := o.GetNewestRemoteDB()
		info := RemoteFileSumary{
			Name:   remotename,
			File:   remotefile,
			Error:  e,
			Client: client,
		}

		remoteinfo <- info
	}()

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
	// This is when the newest local database is newer than the remote one,
	// and thus it should be uploaded to the remote.
	case lmod.UnixNano() > rmod.UnixNano() && linf.Error == nil && rinf.Error == nil:
		{
			o.PushtoRemote(&rinf, &linf)
		}
	// This is when the newest local databse is older than the remote,
	// and thus the remote is newest and should be downloaded and added
	// to the directory with a timestamp
	case lmod.UnixNano() < rmod.UnixNano() && linf.Error == nil && rinf.Error == nil:
		{
			o.PulltoLocal(&rinf, &linf)
		}
	case lmod.UnixNano() == rmod.UnixNano() && linf.Error == nil && rinf.Error == nil:
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

func (o *OptionMeta) PulltoLocal(rinfo *RemoteFileSumary, linfo *LocalFileSummary) error {

	in := &s3.GetObjectInput{
		Key:       &rinfo.Name,
		VersionId: rinfo.File.VersionId,
		Bucket:    &rinfo.Client.RemoteOptions.Bucket,
	}

	out, err := rinfo.Client.Client.GetObject(in)
	defer out.Body.Close()
	if err != nil {
		return err
	}

	file, err1 := os.Create(linfo.File.Name())
	defer file.Close()
	if err1 != nil {
		return err1
	}

	_, err2 := io.Copy(file, out.Body)
	if err2 != nil {
		return err2
	}

	return nil

}

func (o *OptionMeta) PushtoRemote(rinfo *RemoteFileSumary, linfo *LocalFileSummary) error {

	in := &s3.PutObjectInput{
		Key:    &rinfo.Name,
		Bucket: &rinfo.Client.RemoteOptions.Bucket,
		Body:   linfo.File,
	}

	_, err := rinfo.Client.Client.PutObject(in)
	if err != nil {
		return err
	}

	return nil

}
