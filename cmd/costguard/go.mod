module github.com/tanay13/costguard/cmd/costguard

go 1.24.5

require github.com/tanay13/costguard/packages/mcp-server v0.0.0

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

replace github.com/tanay13/costguard/packages/mcp-server => ../../packages/mcp-server
