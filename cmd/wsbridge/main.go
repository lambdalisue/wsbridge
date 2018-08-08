package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/lambdalisue/wsbridge"
	"github.com/sirupsen/logrus"
)

var (
	// Version of wsbridge based on git tag
	Version string
)

func main() {
	os.Exit(run())
}

func run() int {
	const defaultHost = "127.0.0.1"
	const defaultPort = 8080
	const defaultPath = "/echo"
	const defaultScheme = "wss"

	var (
		version  = flag.Bool("v", false, "show version")
		host     = flag.String("host", defaultHost, "listening host name")
		port     = flag.Int("port", defaultPort, "listening host port")
		proxyURI = flag.String("proxy", "", "HTTP proxy used to connect the server")
	)
	flag.Usage = func() {
		fmt.Println("usage: wsbridge [options] {url}")
		fmt.Println("")
		fmt.Println("  Tiny websocket connection bridge server.")
		fmt.Println("")
		fmt.Println("  url: string")
		fmt.Println("\tthe server url (e.g. wss://echo.websocket.org/echo)")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		fmt.Println("version:", Version)
		return 0
	} else if flag.NArg() == 0 {
		flag.Usage()
		return 1
	}
	serverURI := flag.Arg(0)
	server, err := url.Parse(serverURI)
	if err != nil {
		fmt.Printf("Failed to parse '%s': %s\n", serverURI, err)
		return 1
	} else if err := validateURI(server); err != nil {
		fmt.Printf("Failed to parse '%s': %s\n", serverURI, err)
		return 1
	}

	config := wsbridge.NewConfig(
		*host,
		*port,
		*server,
	)
	if *proxyURI != "" {
		proxy, err := url.Parse(*proxyURI)
		if err != nil {
			fmt.Printf("Failed to parse '%s': %s\n", *proxyURI, err)
			return 1
		}
		config.Proxy = proxy
	}

	fmt.Println("*******************************************************************************************")
	fmt.Println("wsbridge - Tiny websocket connection bridge server")
	fmt.Println("")
	fmt.Printf("  listen: ws://%s:%d%s\n", config.Host, config.Port, config.Server.Path)
	fmt.Printf("  server: %s\n", config.Server.String())
	if config.Proxy != nil {
		fmt.Printf("  proxy:  %s\n", config.Proxy.String())
	} else {
		fmt.Printf("  proxy:  system\n")
	}
	fmt.Println("")
	fmt.Println("*******************************************************************************************")

	bridge := wsbridge.NewBridge(&config)
	if err := bridge.Start(); err != nil {
		logrus.Error(err)
		return 1
	}
	return 0
}

func validateURI(uri *url.URL) error {
	if uri.Scheme == "" {
		return fmt.Errorf(
			"url '%s' does not have scheme",
			uri.String(),
		)
	} else if !contains(uri.Scheme, []string{"ws", "wss", "http", "https"}) {
		return fmt.Errorf(
			"scheme '%s' is not supported. use ws, wss, http, or https",
			uri.Scheme,
		)
	} else if uri.Host == "" {
		return fmt.Errorf(
			"url '%s' does not have host",
			uri.String(),
		)
	} else if uri.Path == "" {
		return fmt.Errorf(
			"url '%s' does not have path",
			uri.String(),
		)
	}
	return nil
}

func contains(x string, a []string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
