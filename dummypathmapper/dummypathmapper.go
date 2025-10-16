// Copyright 2015 Dominique Feyer <dfeyer@ttree.ch>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dummypathmapper

import (
	"github.com/dfeyer/flow-debugproxy/config"
	"github.com/dfeyer/flow-debugproxy/logger"
	"github.com/dfeyer/flow-debugproxy/pathmapperfactory"
	"github.com/dfeyer/flow-debugproxy/pathmapping"
	"github.com/dfeyer/flow-debugproxy/xdebugproxy"
)

const framework = "dummy"

func init() {
	p := &PathMapperFactory{}
	pathmapperfactory.Register(framework, p)
}

type PathMapperFactory struct{}

func (p *PathMapperFactory) Create(c *config.Config, l *logger.Logger, m *pathmapping.PathMapping) xdebugproxy.XDebugProcessorPlugin {
	return &PathMapper{
		config:      c,
		logger:      l,
		pathMapping: m,
	}
}

// PathMapper handle the mapping between real code and proxy
type PathMapper struct {
	config      *config.Config
	logger      *logger.Logger
	pathMapping *pathmapping.PathMapping
}

// ApplyMappingToTextProtocol change file path in xDebug text protocol
func (p *PathMapper) ApplyMappingToTextProtocol(message []byte) ([]byte, error) {
	return message, nil
}

// ApplyMappingToXML change file path in xDebug XML protocol
func (p *PathMapper) ApplyMappingToXML(message []byte) ([]byte, error) {
	return message, nil
}
