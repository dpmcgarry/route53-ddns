package main

import "github.com/dpmcgarry/route53-ddns/cmd"

func main() {
	configureViper()
	configureLogging()
	cmd.Execute()
}
