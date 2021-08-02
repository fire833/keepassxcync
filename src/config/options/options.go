package options

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	fp "path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fire833/keepassxcync/src/config/debug"
)

type OptionMeta struct {
	// The exact file descriptor of the options file.
	File *os.File
	// Unmarshalled options information
	Options *Options
	// The exact file descriptor of the databse to sync.
	Db *os.File
}

// Primary options struct, options.json/.yml/.yaml file unmarshalls into this
type Options struct {
	// Array of remotes that can be uploaded to/downloaded form.
	Remotes []Remote `json:"remote"`
}

// Describes a remote object, or a specific instance
type Remote struct {
	// Name of the remote endpoint.
	Name string `json:"name"`
	// Actual URI of the endpoint.
	Endpoint string `json:"endpoint"`
	// Specify the region of the endpoint to pass in to the sdk.
	Region string
	// Specify the bucket name to sync with
	Bucket string
	// API key id of the remote.
	Id string `json:"api_id"`
	// API key of the remote.
	Key string `json:"api_key"`
	// Only one of these remotes can be specified as default, but if this
	// bool is set to true, the binary will default to try and sync with
	// this remote as the first attempt. Only one remote in each file can be set
	// to default, otherwise the "default" remote will be inconsistently chosen.
	IsDefault bool `json:"default"`
}

// Import options from options.json file in the local directory
// of where the binary is.
func NewOptions() (o *OptionMeta) {

	opt := o.Options

	file, err := os.ReadDir(os.Getenv("PWD")) // read all files in the current working directory where the binary was executed.
	if err != nil {
		fmt.Printf("Failed to read working directory for options file: %v\n", err)
		os.Exit(1)
	}

	for _, r := range file {

		info, err1 := r.Info()
		if err1 != nil && debug.DEBUG == true {
			fmt.Printf("Failed to get fileinfo from file: %v, error was: %v\n", info.Name(), err1)
			continue
		} else if err1 != nil {
			continue
		}

		file, err2 := os.ReadFile(fp.Join(os.Getenv("PWD"), info.Name()))
		if err2 != nil && debug.DEBUG == true {
			fmt.Printf("Failed to read file %v, error was: %v\n", info.Name(), err2)
		} else if err2 != nil {
			continue
		}

		switch info.Name() {
		case "options.json":
			{
				err3 := json.Unmarshal(file, opt)
				if err3 != nil && debug.DEBUG == true {
					fmt.Printf("Failed to unmarshall json file %v , error: %v", info.Name(), err3)
					continue
				} else if err3 != nil {
					continue
				}
				fd, err5 := os.OpenFile(fp.Join(os.Getenv("PWD"), r.Name()), os.O_CREATE|os.O_SYNC|os.O_RDWR, 0755)
				if err5 != nil && debug.DEBUG == true {
					fmt.Printf("Unable to open file %v, error: %v", r.Name(), err5)
				} else if err != nil {
					os.Exit(1)
				}
				o.File = fd
				return o
			}
		case "options.yaml", "options.yml":
			{
				err4 := yaml.Unmarshal(file, opt)
				if err4 != nil && debug.DEBUG == true {
					fmt.Printf("Failed to unmarshall yaml file %v , error: %v", info.Name(), err4)
					continue
				} else if err4 != nil {
					continue
				}
				fd, err6 := os.OpenFile(fp.Join(os.Getenv("PWD"), r.Name()), os.O_CREATE|os.O_SYNC|os.O_RDWR, 0755)
				if err6 != nil && debug.DEBUG == true {
					fmt.Printf("Unable to open file %v, error: %v", r.Name(), err6)
				} else if err != nil {
					os.Exit(1)
				}
				o.File = fd
				return o
			}
		}

	}
	log.Fatal("No files in working directory that are options files.")
	return nil
}

// List of all remotes and their s3 client for utilizing for synching.
type Clients struct {
	// Array of remotes.
	Remotes []Client
}

// A single remote config and its accompanying s3 client.
type Client struct {
	// Specific s3 client for this remote
	Client *s3.S3
	// Set of options for this remote
	RemoteOptions Remote
}

func NewEncryptedOptions() (o *OptionMeta) {
	return nil
}

// Save updates to disk.
func (o *OptionMeta) SaveOptions() error {
	return o.File.Sync()
}

// Returns key value for specific provider.
func (r *Remote) GetKey() string {
	return r.Key
}

// Returns id value for specific provider.
func (r *Remote) GetId() string {
	return r.Id
}

// Returns endpoint value for specific provider.
func (r *Remote) GetEndpoint() string {
	return r.Endpoint
}

// Returns name value for specific provider.
func (r *Remote) GetName() string {
	return r.Name
}

// Adds new remote to current options object.
func (o *Options) AddRemote(name, endp, region, bucket, key, id string) {
	o.Remotes = append(o.Remotes, Remote{
		Name:     name,
		Endpoint: endp,
		Region:   region,
		Bucket:   bucket,
		Key:      key,
		Id:       id,
	})
}

// Prints each of the remotes out to stout.
func (o *Options) PrintRemotes(printkey bool) {
	for i, r := range o.Remotes {
		fmt.Printf("Remote #%d: \n ", i)
		fmt.Printf("Name: %s \n Remote Url: %s \n Remote key Id: %s \n", r.Name, r.Endpoint, r.Id)
		if printkey == true {
			fmt.Printf("Remote API key: %s \n", r.Key)
		}
	}
}

func (o *Options) S3ClientsAll() (c *Clients) {
	for _, remote := range o.Remotes {
		conf := &aws.Config{
			Endpoint:    &remote.Endpoint,
			Credentials: credentials.NewStaticCredentials(remote.Id, remote.Key, ""),
			Region:      &remote.Region,
		}

		client := Client{}

		client.Client = s3.New(session.New(conf))
		client.RemoteOptions = remote

		c.Remotes = append(c.Remotes, client)
	}
	return c
}

// Returns a new s3 Client according to the specified remote name.
func (o *Options) NewS3Client(name string, def bool) (c *Client, e error) {

	for _, remote := range o.Remotes {

		conf := &aws.Config{
			Endpoint:    &remote.Endpoint,
			Credentials: credentials.NewStaticCredentials(remote.Id, remote.Key, ""),
			Region:      &remote.Region,
		}

		if def == false {
			if remote.Name == name {
				c.Client = s3.New(session.New(conf))
				c.RemoteOptions = remote
				return c, nil
			}
		} else if def == true {
			if remote.IsDefault == true {
				c.Client = s3.New(session.New(conf))
				c.RemoteOptions = remote
				return c, nil
			}
		}
	}
	return c, errors.New("Specified remote not found in config.")
}
