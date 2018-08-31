package config

//go:generate sh -c "CGO_ENABLED=0 go run .packr/packr.go ."

import (
	"github.com/gobuffalo/packr"
)

func ConfigPackr() packr.Box {
	box := packr.NewBox("./configs")
	box.AddString("common-env.json", "{\n  \"SSO_DISABLE_SSL_CERTIFICATE_VALIDATION\":\"FALSE\",\n  \"SSO_OPENIDCONNECT_DEPLOYMENTS\":\"ROOT.war\",\n  \"KIE_ADMIN_PWD\":\"RedHat\",\n  \"KIE_SERVER_CONTROLLER_PWD\":\"RedHat\",\n  \"KIE_SERVER_PWD\":\"RedHat\",\n  \"KIE_SERVER_USER\":\"executionUser\",\n  \"KIE_MBEANS\":\"enabled\",\n  \"KIE_SERVER_CONTROLLER_USER\":\"controllerUser\",\n  \"KIE_ADMIN_USER\":\"adminUser\"\n}")
	box.AddString("console-env.json", "{\n  \"PROBE_DISABLE_BOOT_ERRORS_CHECK\":\"true\",\n  \"PROBE_IMPL\":\"probe.eap.jolokia.EapProbe\",\n  \"KIE_MAVEN_USER\":\"mavenUser\"\n}")
	box.AddString("server-env.json", "{\n  \"KIE_SERVER_CONTROLLER_PROTOCOL\":\"ws\",\n  \"MAVEN_REPOS\":\"RHPAMCENTR,EXTERNAL\",\n  \"RHPAMCENTR_MAVEN_REPO_PASSWORD\":\"RedHat\",\n  \"RHPAMCENTR_MAVEN_REPO_USERNAME\":\"mavenUser\",\n  \"RHPAMCENTR_MAVEN_REPO_PATH\":\"/maven2/\"\n}")

	return box
}
