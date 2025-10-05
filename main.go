// Copyright 2015 Dominique Feyer <dfeyer@ttree.ch>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
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
			Value: "Development:9003",
			Usage: "Listen address IP and port number",
		},
		&cli.StringFlag{
			Name:  "ide, I",
			Value: "127.0.0.1:9010",
			Usage: "Bind address IP and port number",
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
			Context:     "",
			Framework:   cli.String("framework"),
			LocalRoot:   strings.TrimRight(cli.String("localroot"), "/"),
			Verbose:     cli.Bool("verbose") || cli.Bool("vv"),
			VeryVerbose: cli.Bool("vv"),
			Debug:       cli.Bool("debug"),
		}

		log := &logger.Logger{
			Config: c,
		}

		listener, raddr, err := setupNetworkConnection(strings.Split(cli.String("xdebug"), ","), cli.String("ide"))
		if err != nil {
			log.Warn(err.Error())
			os.Exit(1)
		}

		for _, listenerWithContext := range listener {
			log.Info("Debugger from %v for context %v", listenerWithContext.addr, listenerWithContext.context)
		}
		log.Info("IDE      from %v", raddr)
		if c.Verbose {
			log.Info("Framework     %v", c.Framework)
			log.Info("Local Root    %v", c.LocalRoot)
			log.Info("Verbose       %v", c.Verbose)
			log.Info("Very Verbose  %v", c.VeryVerbose)
			log.Info("Debug         %v", c.Debug)
		}

		connections := make(chan *xdebugproxy.Proxy)
		for _, listenerWithContext := range listener {
			listenerWithContext := listenerWithContext
			originalConfig := *c
			proxyConfig := originalConfig // copy config
			proxyConfig.Context = listenerWithContext.context

			pathMapping := &pathmapping.PathMapping{}
			pathMapper, err := pathmapperfactory.Create(&proxyConfig, pathMapping, log)
			if err != nil {
				log.Warn(err.Error())
				os.Exit(1)
			}

			go func() {
				for {
					conn, err := listenerWithContext.listener.AcceptTCP()
					if err != nil {
						log.Warn("Failed to accept connection '%s'\n", err)
						continue
					}

					connections <- &xdebugproxy.Proxy{
						Lconn:      conn,
						Raddr:      raddr,
						PathMapper: pathMapper,
						Config:     &proxyConfig,
					}
				}
			}()
		}

		for {
			proxy := <-connections
			go proxy.Start()
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		golog.Fatal(err)
		return
	}
}

type listenerWithContext struct {
	addr     *net.TCPAddr
	listener *net.TCPListener
	context  string
}

func splitContextAndPort(contextWithPort string) (string, string, error) {
	contextAndPort := strings.Split(contextWithPort, ":")
	if len(contextAndPort) != 2 {
		return "", "", fmt.Errorf("could not parse port and context information '%s'", contextWithPort)
	}
	return contextAndPort[0], contextAndPort[1], nil
}

func setupNetworkConnection(xdebugAddr []string, ideAddr string) ([]*listenerWithContext, *net.TCPAddr, error) {
	listener := make([]*listenerWithContext, 0, len(xdebugAddr))
	for _, portWithContext := range xdebugAddr {
		context, portStr, err := splitContextAndPort(portWithContext)
		if err != nil {
			return nil, nil, err
		}

		addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+portStr)
		if err != nil {
			return nil, nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, nil, err
		}

		listener = append(listener, &listenerWithContext{
			addr:     addr,
			listener: l,
			context:  context,
		})
	}

	raddr, err := net.ResolveTCPAddr("tcp", ideAddr)
	if err != nil {
		return nil, nil, err
	}

	return listener, raddr, nil
}
