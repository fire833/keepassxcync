module github.com/fire833/keepassxcync

go 1.20

require (
	github.com/aws/aws-sdk-go v1.44.289 // Used for s3 things
	github.com/integrii/flaggy v1.5.2 // Better flag package
	golang.org/x/term v0.9.0 // Read password from terminal
	gopkg.in/yaml.v3 v3.0.1 // Unmarshalling yaml stuffs
)

require (
	github.com/magefile/mage v1.15.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
)
