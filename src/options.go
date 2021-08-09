package src

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	fp "path/filepath"
	"regexp"
	"syscall"

	"gopkg.in/yaml.v3"

	"golang.org/x/term"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fire833/keepassxcync/src/config/debug"
)

type OptionMeta struct {
	// The exact file descriptor of the options file.
	File *os.File
	// Raw file data of the options file
	FileData []byte
	// Absolute path to options file
	FilePath string
	// Unmarshalled options information
	Options *Options
	// The exact file descriptor of the database to sync.
	// Db *os.File

	// List of s3 clients for each of the remotes that are defined in
	// the options file. Depending on mode there could be only 1 client in the object.
	Clients []Client
}

// Primary options struct, options.json/.yml/.yaml file unmarshalls into this
type Options struct {
	// Specifies the supposed name of the databse that you want to sync
	// with this binary. Should be a name
	DatabaseName string `json:"db_name" yaml:"db_name"`
	// Array of remotes that can be uploaded to/downloaded form.
	Remotes []Remote `json:"remotes" yaml:"remotes"`
}

// Describes a remote object, or a specific instance
type Remote struct {
	// Name of the remote endpoint.
	Name string `json:"name" yaml:"name"`
	// Actual URI of the endpoint.
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	// Specify the region of the endpoint to pass in to the sdk.
	Region string `json:"region" yaml:"region"`
	// Specify the bucket name to sync with
	Bucket string `json:"bucket" yaml:"bucket"`
	// API key id of the remote.
	Id string `json:"api_id" yaml:"api_id"`
	// API key of the remote.
	Key string `json:"api_key" yaml:"api_key"`
	// Only one of these remotes can be specified as default, but if this
	// bool is set to true, the binary will default to try and sync with
	// this remote as the first attempt. Only one remote in each file can be set
	// to default, otherwise the "default" remote will be inconsistently chosen.
	IsDefault bool `json:"default"`
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

// Import options from options.json file in the local directory
// of where the binary is.
func NewOptions() (o *OptionMeta) {

	o = &OptionMeta{}
	opts := &o.Options

	if CONFIG == "" {
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
					err3 := json.Unmarshal(file, opts)
					if err3 != nil && debug.DEBUG == true {
						fmt.Printf("Failed to unmarshall json file %v , error: %v", info.Name(), err3)
						continue
					} else if err3 != nil {
						continue
					}
				}
			case "options.yaml", "options.yml":
				{
					err4 := yaml.Unmarshal(file, opts)
					if err4 != nil && debug.DEBUG == true {
						fmt.Printf("Failed to unmarshall yaml file %v , error: %v", info.Name(), err4)
						continue
					} else if err4 != nil {
						continue
					}
				}
			default:
				{
					continue
				}
			}

			fd, err1 := os.OpenFile(fp.Join(os.Getenv("PWD"), r.Name()), os.O_CREATE|os.O_SYNC|os.O_RDWR, 0755)
			if err1 != nil && debug.DEBUG == true {
				fmt.Printf("Unable to open file %v, error: %v", r.Name(), err1)
			} else if err != nil {
				os.Exit(1)
			}
			o.FileData = file // Set the read data into the struct so that it
			// can be modified when saving the file later on.
			o.File = fd // Set the file descriptor into the struct for modification
			// as well, and to sync the status later on.
			o.FilePath = fp.Join(os.Getenv("PWD"), info.Name()) // Set the absolute path to the current
			// options file for synching later on.
			return o
		}

	} else {

		data, err := os.ReadFile(CONFIG)
		if err != nil {
			fmt.Println("Unable to read file " + CONFIG + ", error: " + err.Error())
			os.Exit(1)
		}

		switch fp.Ext(CONFIG) {
		case ".json":
			{
				err := json.Unmarshal(data, opts)
				if err != nil {
					fmt.Println("Error unmarshalling options file: " + err.Error())
					os.Exit(1)
				}
			}
		case ".yaml", ".yml":
			{
				err := yaml.Unmarshal(data, opts)
				if err != nil {
					fmt.Println("Error unmarshalling options file: " + err.Error())
					os.Exit(1)
				}
			}
		default:
			{
				fmt.Println("Unable to parse the options file you provided.")
				os.Exit(1)
			}
		}

		fd, err1 := os.OpenFile(CONFIG, os.O_CREATE|os.O_SYNC|os.O_RDWR, 0755)
		if err1 != nil && debug.DEBUG == true {
			fmt.Printf("Unable to open file %v, error: %v", CONFIG, err1)
		} else if err != nil {
			os.Exit(1)
		}
		o.FileData = data // Set the read data into the struct so that it
		// can be modified when saving the file later on.
		o.File = fd // Set the file descriptor into the struct for modification
		// as well, and to sync the status later on.
		o.FilePath = CONFIG // Set the absolute path to the current
		// options file for synching later on.
		return o
	}

	log.Fatal("No files in working directory that are options files.")
	return nil
}

// Save updates to disk.
func (o *OptionMeta) SaveOptions() error {
	return o.File.Sync()
}

func (o *OptionMeta) SyncOptionsToDisk() error {

	info, err := o.File.Stat()
	if err != nil {
		return err
	}
	var offset int64 = 0
	// var null byte

	switch fp.Ext(o.File.Name()) {
	case ".json":
		{
			data, err := json.Marshal(o.Options)
			if err != nil {
				return err
			}

			// If new file size is bigger than old file or the same, just write from index 0.
			if len(data) >= int(info.Size()) {

				o.File.WriteAt(data, offset)
				o.File.Sync()

				return nil

			} else {
				o.File.Close()

				err := os.Remove(o.FilePath)
				if err != nil {
					fmt.Println("Error removing options file for rewrite, exiting...")
					os.Exit(1)
				}

				file, err1 := os.Create(o.FilePath)
				if err1 != nil {
					fmt.Printf("Error with creating new options file for persisting options, here is your data: %v", o.Options)
				}

				_, err2 := file.Write(data)
				if err2 != nil {
					fmt.Printf("Error with writing to options file, here is your data: %v", o.Options)
				}

				file.Sync()
				file.Close()
				return nil
			}
		}
	case ".yml", ".yaml":
		{
			data, err := yaml.Marshal(o.Options)
			if err != nil {
				return err
			}
			// If new file size is bigger than old file or the same, just write from index 0.
			if len(data) >= int(info.Size()) {

				o.File.WriteAt(data, offset)
				o.File.Sync()

				return nil

			} else {
				o.File.Close()

				err := os.Remove(o.FilePath)
				if err != nil {
					fmt.Println("Error removing options file for rewrite, exiting...")
					os.Exit(1)
				}

				file, err1 := os.Create(o.FilePath)
				if err1 != nil {
					fmt.Printf("Error with creating new options file for persisting options, here is your data: %v", o.Options)
				}

				_, err2 := file.Write(data)
				if err2 != nil {
					fmt.Printf("Error with writing to options file, here is your data: %v", o.Options)
				}

				file.Sync()
				file.Close()
				return nil
			}
		}
	}

	return nil
}

// Adds new remote to current options object.
func (o *OptionMeta) AddRemote() {

	scan := bufio.NewScanner(os.Stdin)

	name := DefName(scan)

	fmt.Print("Enter remote URI: ")
	scan.Scan()
	endpoint := scan.Text()
	fmt.Print("Enter remote region: ")
	scan.Scan()
	region := scan.Text()
	fmt.Print("Enter remote bucket name: ")
	scan.Scan()
	bucket := scan.Text()
	fmt.Print("Enter remote Key ID: ")
	scan.Scan()
	id := scan.Text()
	fmt.Print("Enter remote key: ")
	key, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Error with reading password, please try again.")
		o.File.Close()
		os.Exit(1)
	}
	fmt.Print("Do you want this remote (" + name + ") to be set as the default remote? [y/N]")
	scan.Scan()

	yes, err := CheckUserYesNo(scan.Text(), false)
	if err != nil {
		fmt.Println(err)
		o.File.Close()
		os.Exit(1)
	}

	// Add the compiled remote to the list.
	o.Options.Remotes = append(o.Options.Remotes, Remote{
		Name:      name,
		Endpoint:  endpoint,
		Region:    region,
		Bucket:    bucket,
		Key:       string(key),
		Id:        id,
		IsDefault: yes,
	})

	fmt.Print("Test connectivity with this remote? [Y/n]")
	scan.Scan()

	yes, err2 := CheckUserYesNo(scan.Text(), true)
	if err2 != nil {
		fmt.Println(err2)
		o.File.Close()
		os.Exit(1)
	}

	switch yes {
	case true:
		{
			// Ping the remote to make sure that it works and you have proper access with those API keys.
			// Do this before persisting the remote to disk.
			client, _ := o.NewS3Client(name, false)

			out, err2 := client.Client.ListBuckets(&s3.ListBucketsInput{})
			switch {
			case err2 != nil:
				{
					fmt.Printf("Error with accessing new remote, here is the error: %v\n", err2)
					os.Exit(1)
				}
			case len(out.Buckets) == 0:
				{
					fmt.Println("No buckets are associated with these credentials.")
					os.Exit(1)
				}
			default:
				{
					for _, s3bucket := range out.Buckets {
						bname := *s3bucket.Name
						if bname == bucket {
							fmt.Printf("Bucket %v is accesible by this key, saving remote %v now.\n", bucket, name)
						} else {
							continue
						}
					}
				}
			}
		}
	case false:
		{
			err := o.SyncOptionsToDisk()
			if err != nil {
				fmt.Printf("Error with persisting changes to disk, aborting... Error: %v\n", err)
				o.File.Close()
				os.Exit(1)
			}
		}
	}
}

func DefName(scan *bufio.Scanner) (name string) {

	if NAME == "" {
		fmt.Print("Enter name of remote: ")
		scan.Scan()
		name := scan.Text()
		return name
	} else {
		return NAME
	}

}

func (o *OptionMeta) UpdateRemote() {

	scan := bufio.NewScanner(os.Stdin)

	name := DefName(scan)

	fmt.Print("Enter remote URI (Press enter to keep the same): ")
	scan.Scan()
	endpoint := scan.Text()
	fmt.Print("Enter remote region (Press enter to keep the same): ")
	scan.Scan()
	region := scan.Text()
	fmt.Print("Enter remote bucket name (Press enter to keep the same): ")
	scan.Scan()
	bucket := scan.Text()
	fmt.Print("Enter remote Key ID (Press enter to keep the same): ")
	scan.Scan()
	id := scan.Text()
	fmt.Print("Enter remote key (Press enter to keep the same): ")
	key, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Error with reading password, please try again.")
		o.File.Close()
		os.Exit(1)
	}
	fmt.Print("Do you want this remote (" + name + ") to be set as the default remote? [y/N]")
	scan.Scan()
	yes, err := CheckUserYesNo(scan.Text(), false)
	if err != nil {
		fmt.Println(err)
		o.File.Close()
		os.Exit(1)
	}

	for _, remote := range o.Options.Remotes {
		if remote.Name == NAME {
			if endpoint != "" {
				remote.Endpoint = endpoint
			}
			if region != "" {
				remote.Region = region
			}
			if bucket != "" {
				remote.Bucket = bucket
			}
			if id != "" {
				remote.Id = id
			}
			if key != nil {
				remote.Key = string(key)
			}
			if yes {
				remote.IsDefault = true
			}
		}
	}

	err1 := o.SyncOptionsToDisk()
	if err1 != nil {
		fmt.Printf("Error with persisting changes to disk, aborting... Error: %v\n", err1)
		o.File.Close()
		os.Exit(1)
	}

}

func (o *OptionMeta) RemoveRemote() {

	if NAME == "" {
		fmt.Println("No remote defined to delete. Please define a remote to delete with -n flag.")
	}

	for i, remote := range o.Options.Remotes {
		if remote.Name == NAME {
			// Not necesarily efficient for doing a reslice at scale.
			o.Options.Remotes = append(o.Options.Remotes[:i], o.Options.Remotes[i+1:]...)

			err := o.SyncOptionsToDisk()
			if err != nil {
				fmt.Printf("Error with persisting changes to disk, aborting... Error: %v", err)
				os.Exit(1)
			}
			fmt.Printf("Removed remote %v.", NAME)

		} else {
			continue
		}
	}

}

// Prints each of the remotes out to stout.
func (o *OptionMeta) PrintRemotes(printkey bool) {
	for i, r := range o.Options.Remotes {
		fmt.Printf("Remote #%d: \n ", i)
		fmt.Printf("Name: %s \n Remote Url: %s \n Remote Key Id: %s \n Region: %s \n Bucket: %s \n", r.Name, r.Endpoint, r.Id, r.Region, r.Bucket)
		if printkey == true {
			fmt.Printf("Remote API key: %s \n", r.Key)
		}
	}
}

func (o *OptionMeta) S3ClientsAll() (c *Clients) {
	for _, remote := range o.Options.Remotes {
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
func (o *OptionMeta) NewS3Client(name string, def bool) (c *Client, e error) {

	for _, remote := range o.Options.Remotes {

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

func CheckUserYesNo(input string, defaultyes bool) (b bool, e error) {

	if defaultyes == false {
		if regexp.MustCompile(`[Nn].*`).MatchString(input) || input == "" {
			return false, nil
		} else if regexp.MustCompile(`[Yy].*`).MatchString(input) {
			return true, nil
		} else {
			return false, errors.New("Error with parsing your response, please try again.")
		}

	} else if defaultyes == true {
		if regexp.MustCompile(`[Yy].*`).MatchString(input) || input == "" {
			return true, nil
		} else if regexp.MustCompile(`[Nn].*`).MatchString(input) {
			return false, nil
		} else {
			return false, errors.New("Error with parsing your response, please try again.")
		}
	}

	return false, errors.New("Something weird happened.")
}
