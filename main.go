// Copyright 2015 Dominique Feyer <dfeyer@ttree.ch>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	golog "log"

	"github.com/dfeyer/flow-debugproxy/config"

	"github.com/dfeyer/flow-debugproxy/logger"
	"github.com/dfeyer/flow-debugproxy/pathmapperfactory"
	"github.com/dfeyer/flow-debugproxy/pathmapping"
	"github.com/dfeyer/flow-debugproxy/xdebugproxy"

	// Register available path mapper
	_ "github.com/dfeyer/flow-debugproxy/dummypathmapper"
	_ "github.com/dfeyer/flow-debugproxy/flowpathmapper"

	"github.com/urfave/cli"

	"net"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "flow-debugproxy"
	app.Usage = "Flow Framework xDebug proxy"
	app.Version = "1.0.1"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "xdebug, l",
			Value: "127.0.0.1:9000",
			Usage: "Listen address IP and port number",
		},
		&cli.StringFlag{
			Name:  "ide, I",
			Value: "127.0.0.1:9010",
			Usage: "Bind address IP and port number",
		},
		&cli.StringFlag{
			Name:  "context, c",
			Value: "Development",
			Usage: "The context to run as",
		},
		&cli.StringFlag{
			Name:  "localroot, r",
			Value: "",
			Usage: "Local project root for remote debugging",
		},
		&cli.StringFlag{
			Name:  "framework",
			Value: "flow",
			Usage: "Framework support, currently on Flow framework (flow) or Dummy (dummy) is supported",
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "Verbose",
		},
		&cli.BoolFlag{
			Name:  "vv",
			Usage: "Very verbose",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Show debug output",
		},
	}

	app.Action = func(cli *cli.Context) error {
		c := &config.Config{
			Context:     cli.String("context"),
			Framework:   cli.String("framework"),
			LocalRoot:   strings.TrimRight(cli.String("localroot"), "/"),
			Verbose:     cli.Bool("verbose") || cli.Bool("vv"),
			VeryVerbose: cli.Bool("vv"),
			Debug:       cli.Bool("debug"),
		}

		log := &logger.Logger{
			Config: c,
		}

		laddr, raddr, listener, err := setupNetworkConnection(cli.String("xdebug"), cli.String("ide"))
		if err != nil {
			log.Warn(err.Error())
			os.Exit(1)
		}

		log.Info("Debugger from %v\nIDE      from %v\n", laddr, raddr)

		pathMapping := &pathmapping.PathMapping{}
		pathMapper, err := pathmapperfactory.Create(c, pathMapping, log)
		if err != nil {
			log.Warn(err.Error())
			os.Exit(1)
		}

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Warn("Failed to accept connection '%s'\n", err)
				continue
			}

			proxy := &xdebugproxy.Proxy{
				Lconn:      conn,
				Raddr:      raddr,
				PathMapper: pathMapper,
				Config:     c,
			}
			go proxy.Start()
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		golog.Fatal(err)
		return
	}
}

func setupNetworkConnection(xdebugAddr string, ideAddr string) (*net.TCPAddr, *net.TCPAddr, *net.TCPListener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", xdebugAddr)
	if err != nil {
		return nil, nil, nil, err
	}

	raddr, err := net.ResolveTCPAddr("tcp", ideAddr)
	if err != nil {
		return nil, nil, nil, err
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, nil, nil, err
	}

	return laddr, raddr, listener, nil
}
