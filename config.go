package main

type ConfigInfo struct {
	SourceFolder string `toml:"source"`
	DestinationFolder string `toml:"destination"`
}
