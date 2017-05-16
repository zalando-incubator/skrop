package main

import (
	"flag"
	"fmt"
	skropFilters "github.com/zalando-incubator/skrop/filters"
	"github.com/zalando/skipper"
	"github.com/zalando/skipper/filters"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/proxy"
)

const (
	addressFlag             = "address"
	verboseFlag             = "verbose"
	defaultAddress          = ":9090"
	routesFileFlag          = "routes-file"
	tlsCertFlag             = "tls-cert"
	tlsKeyFlag              = "tls-key"
	insecureFlag            = "insecure"
	experimentalUpgradeFlag = "experimental-upgrade"
)

const (
	usageHeader = `
  skrop – Skipper based media service using the vips library.

  https://github.com/zalando-incubator/skrop`

	addressUsage    = "network address that skoap should listen on"
	verboseUsage    = "enable verbose logging"
	routesFileUsage = `alternatively to the target address, it is possible to use a full
	eskip route configuration, and specify the auth() and authTeam()
	filters for the routes individually.
	See also: https://godoc.org/github.com/zalando/skipper/eskip`

	insecureUsage    = `when this flag set, skipper will skip TLS verification`
	certPathTLSUsage = "path of the certificate file"
	keyPathTLSUsage  = "path of the key"

	experimentalUpgradeUsage = "enable experimental feature to handle upgrade protocol requests"
)

var fs *flag.FlagSet

var (
	address             string
	certPathTLS         string
	insecure            bool
	keyPathTLS          string
	verbose             bool
	experimentalUpgrade bool
	routesFile          string
)

func usage() {
	fmt.Fprint(os.Stderr, usageHeader)
	fs.PrintDefaults()
}

func logUsage(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(-1)
}

func init() {
	fs = flag.NewFlagSet("flags", flag.ContinueOnError)
	fs.Usage = usage

	fs.StringVar(&address, addressFlag, defaultAddress, addressUsage)
	fs.BoolVar(&verbose, verboseFlag, false, verboseUsage)
	fs.StringVar(&routesFile, routesFileFlag, "", routesFileUsage)
	fs.BoolVar(&insecure, insecureFlag, false, insecureUsage)
	fs.StringVar(&certPathTLS, tlsCertFlag, "", certPathTLSUsage)
	fs.StringVar(&keyPathTLS, tlsKeyFlag, "", keyPathTLSUsage)
	fs.BoolVar(&experimentalUpgrade, experimentalUpgradeFlag, false, experimentalUpgradeUsage)

	err := fs.Parse(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(-1)
	}
}

func main() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if routesFile == "" {
		logUsage("A routes file needs to be specified.")
	}
	log.Debug(fmt.Sprintf("Using routes-file %s", routesFile))

	o := skipper.Options{
		Address:    address,
		RoutesFile: routesFile,
		CustomFilters: []filters.Spec{
			skropFilters.NewResize(),
			skropFilters.NewCrop(),
			skropFilters.NewCropByWidth(),
			skropFilters.NewCropByHeight(),
			skropFilters.NewResizeByWidth(),
			skropFilters.NewQuality(),
			skropFilters.NewLongerEdgeResize()},
		AccessLogDisabled:   true,
		ProxyOptions:        proxy.OptionsPreserveOriginal,
		CertPathTLS:         certPathTLS,
		KeyPathTLS:          keyPathTLS,
		ExperimentalUpgrade: experimentalUpgrade,
	}

	if insecure {
		o.ProxyOptions |= proxy.OptionsInsecure
	}

	err := skipper.Run(o)
	if err != nil {
		log.Fatal(err)
	}
}
