// Copyright 2015 Dominique Feyer <dfeyer@ttree.ch>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pathmapperfactory

import (
	"github.com/dfeyer/flow-debugproxy/config"
	"github.com/dfeyer/flow-debugproxy/logger"
	"github.com/dfeyer/flow-debugproxy/pathmapping"
	"github.com/dfeyer/flow-debugproxy/xdebugproxy"

	"errors"
)

var pathMapperRegistry = map[string]xdebugproxy.XDebugProcessorPluginFactory{}

// Register a path mapper
func Register(f string, p xdebugproxy.XDebugProcessorPluginFactory) {
	pathMapperRegistry[f] = p
}

// Create return a pathmapper for the given framework
func Create(c *config.Config, p *pathmapping.PathMapping, l *logger.Logger) (xdebugproxy.XDebugProcessorPlugin, error) {
	if _, exist := pathMapperRegistry[c.Framework]; exist {
		pathmapperFactory := pathMapperRegistry[c.Framework]
		return pathmapperFactory.Create(c, l, p), nil
	}

	return nil, errors.New("unsupported framework")
}
