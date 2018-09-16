package defaults

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConsoleEnvironmentDefaults(t *testing.T) {
	defaults := ConsoleEnvironmentDefaults()
	logrus.Debugf("Loaded console defaults as %v", defaults)
	assert.Equal(t, "ROOT.war", defaults["SSO_OPENIDCONNECT_DEPLOYMENTS"])
	assert.Equal(t, "mavenUser", defaults["KIE_MAVEN_USER"], )
}

func TestServerEnvironmentDefaults(t *testing.T) {
	defaults := ServerEnvironmentDefaults()
	logrus.Debugf("Loaded server defaults as %v", defaults)
	assert.Equal(t, "ROOT.war", defaults["SSO_OPENIDCONNECT_DEPLOYMENTS"])
	assert.Equal(t, "mavenUser", defaults["RHPAMCENTR_MAVEN_REPO_USERNAME"])
}
