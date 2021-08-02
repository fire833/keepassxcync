module github.com/fire833/keepassxcync

go 1.16

require (
	github.com/aws/aws-sdk-go v1.39.2 // Used for s3 things
	github.com/integrii/flaggy v1.4.4 // Better flag package
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // Read password from terminal
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // Unmarshalling yaml stuffs
)
