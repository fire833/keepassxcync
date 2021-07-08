package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Primary options struct,
type Options struct {
	Remotes []Remote `json:"remote"` // Array of remotes that can be uploaded to/downloaded form.
}

type Remote struct {
	Name     string `json:"name"`     // Name of the remote endpoint
	Endpoint string `json:"endpoint"` // Actual URI of the endpoint
	Id       string `json:"api_id"`   // API key id of the remote.
	Key      string `json:"api_key"`  // API key of the remote
}

// Import options from options.json file in the local directory
// of where the binary is.
func NewOptions() (o *Options) {

	file, err := os.ReadDir(".") // read all files in the current working directory where the binary was executed.
	if err != nil {
		log.Fatalf("Failed to read working directory for options file: %v", err)
		os.Exit(1)
	}

	for _, r := range file {

		info, err1 := r.Info()
		if err1 != nil {
			log.Fatalf("Failed to get fileinfo from file: %v", err1)
			os.Exit(1)
		}
		if match, err2 := regexp.MatchString("^options.json$", info.Name()); match == true {
			file, err3 := os.ReadFile("./options.json")
			err4 := json.Unmarshal(file, o)

			if err2 != nil {
				fmt.Printf("Error matching regex: %v", err2)
			}
			if err3 != nil || err4 != nil {
				fmt.Printf("Error reading options file or unmarshalling into a struct: %v, %v", err3, err4)
			}
		}
	}
	log.Fatal("No files in working directory that are options.json files.")
	return nil
}

func NewEncryptedOptions() (o *Options) {
	return nil
}

func (o *Options) SaveOptions() {

}

// returns key value for specific provider.
func (o *Options) GetKey(i int) string {
	return o.Remotes[i].Key
}

// Returns id value for specific provider.
func (o *Options) GetId(i int) string {
	return o.Remotes[i].Id
}

// Returns endpoint value for specific provider.
func (o *Options) GetEndpoint(i int) string {
	return o.Remotes[i].Endpoint
}

// Adds new remote to current options object.
func (o *Options) AddRemote(name, endp, key, id string) {
	o.Remotes = append(o.Remotes, Remote{
		Name:     name,
		Endpoint: endp,
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

// Returns a new s3 Client according to the specified credentials integer.
func (o *Options) NewS3Client(i int) *s3.S3 {
	conf := &aws.Config{
		Endpoint:    &o.Remotes[i].Endpoint,
		Credentials: credentials.NewStaticCredentials(o.Remotes[i].Id, o.Remotes[i].Key, ""),
	}
	return s3.New(session.New(conf))

}
