package main

import "flag"

const (
	DefaultPemPath = "server.pem"
	DefaultKeyPath = "server.key"
	DefaultProto = HTTPSProtocol
	DefaultPort = "8888"
)

type Configuration struct {
	PemPath string
	KeyPath string
	Proto   string
	Port    string
}

var conf Configuration

func init() {
	flag.StringVar(&conf.PemPath, "pem", DefaultPemPath, "path to pem file")
	flag.StringVar(&conf.KeyPath, "key", DefaultKeyPath, "path to key file")
	flag.StringVar(&conf.Proto, "proto", DefaultProto, "Proxy protocol (http or https)")
	flag.StringVar(&conf.Port, "port", DefaultPort, "Proxy port")
	flag.Parse()
}