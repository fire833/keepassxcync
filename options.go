package main

import "fmt"

type Options struct {
	Remotes []struct {
		Name     string `json:"name"`     // Name of the remote endpoint
		Endpoint string `json:"endpoint"` // Actual URI of the endpoint
		Id       string `json:"api_id"`   // API key id of the remote.
		Key      string `json:"api_key"`  // API key of the remote
	} `json:"remote"` // Array of remotes that can be
}

func NewOptions() {

}

func (o *Options) GetKey() interface{} {
	return o.Remotes
}

// Prints each of the remotes out to
func (o *Options) PrintRemotes(printkey bool) {
	for i, r := range o.Remotes {
		fmt.Println("Remote #%d: \n ", i)
		fmt.Println("Name: %s \n Remote Url: %s \n Remote key Id: %s \n", r.Name, r.Endpoint, r.Id)
		if printkey == true {
			fmt.Println("Remote API key: %s \n", r.Key)
		}
	}
}
