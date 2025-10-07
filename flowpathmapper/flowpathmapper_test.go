package flowpathmapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildClassNameFromPathSupportPSR2(t *testing.T) {
	basePath, className, err := pathToClassPath("/your/path/sites/dev/master-dev.neos-workplace.dev/Packages/Application/Ttree.FlowDebugProxyHelper/Classes/Ttree/FlowDebugProxyHelper/ProxyClassMapperComponent.php")
	assert.Equal(t, "/your/path/sites/dev/master-dev.neos-workplace.dev", basePath, "they should be equal")
	assert.Equal(t, "Ttree_FlowDebugProxyHelper_ProxyClassMapperComponent", className, "they should be equal")
	assert.Equal(t, nil, err, "there should not be an error")
}

func TestBuildClassNameFromPathSupportPSR4(t *testing.T) {
	basePath, className, err := pathToClassPath("/your/path/sites/dev/master-dev.neos-workplace.dev/Packages/Application/Ttree.FlowDebugProxyHelper/Classes/ProxyClassMapperComponent.php")
	assert.Equal(t, "/your/path/sites/dev/master-dev.neos-workplace.dev", basePath, "they should be equal")
	assert.Equal(t, "Ttree_FlowDebugProxyHelper_ProxyClassMapperComponent", className, "they should be equal")
	assert.Equal(t, nil, err, "there should not be an error")
}
