package plugins_test

import (
	"testing"

	"github.com/jenkins-x-plugins/jx-kube-test/pkg/plugins"
	"github.com/stretchr/testify/assert"
)

func TestConftestPlugin(t *testing.T) {
	t.Parallel()

	plugin := plugins.CreateConftestPlugin(plugins.ConftestVersion)

	assert.Equal(t, plugins.ConftestPluginName, plugin.Name, "plugin.Name")
	assert.Equal(t, plugins.ConftestPluginName, plugin.Spec.Name, "plugin.Spec.Name")

	foundLinux := false
	foundMac := false
	foundWindows := false
	foundArm := false
	for _, b := range plugin.Spec.Binaries {
		switch b.Goarch {
		case "arm64":
			switch b.Goos {
			case "Linux":
				foundArm = true
				assert.Equal(t, "https://github.com/open-policy-agent/conftest/releases/download/v"+plugins.ConftestVersion+"/conftest_"+plugins.ConftestVersion+"_Linux_arm64.tar.gz", b.URL, "URL for linux arm binary")
				t.Logf("found linux binary URL %s", b.URL)
			}

		case "amd64":
			switch b.Goos {
			case "Darwin":
				foundMac = true
				assert.Equal(t, "https://github.com/open-policy-agent/conftest/releases/download/v"+plugins.ConftestVersion+"/conftest_"+plugins.ConftestVersion+"_Darwin_x86_64.tar.gz", b.URL, "URL for mac binary")
				t.Logf("found mac binary URL %s", b.URL)
			case "Linux":
				foundLinux = true
				assert.Equal(t, "https://github.com/open-policy-agent/conftest/releases/download/v"+plugins.ConftestVersion+"/conftest_"+plugins.ConftestVersion+"_Linux_x86_64.tar.gz", b.URL, "URL for mac binary")
				t.Logf("found linux binary URL %s", b.URL)
			case "Windows":
				foundWindows = true
				assert.Equal(t, "https://github.com/open-policy-agent/conftest/releases/download/v"+plugins.ConftestVersion+"/conftest_"+plugins.ConftestVersion+"_Windows_x86_64.zip", b.URL, "URL for windows binary")
				t.Logf("found windows binary URL %s", b.URL)
			}
		}
	}
	assert.True(t, foundArm, "did not find an arm linux binary in the plugin %#v", plugin)
	assert.True(t, foundLinux, "did not find a linux binary in the plugin %#v", plugin)
	assert.True(t, foundMac, "did not find a mac binary in the plugin %#v", plugin)
	assert.True(t, foundWindows, "did not find a windows binary in the plugin %#v", plugin)
}

func TestKubeScorePlugin(t *testing.T) {
	t.Parallel()

	plugin := plugins.CreateKubeScorePlugin(plugins.KubeScoreVersion)

	assert.Equal(t, plugins.KubeScorePluginName, plugin.Name, "plugin.Name")
	assert.Equal(t, plugins.KubeScorePluginName, plugin.Spec.Name, "plugin.Spec.Name")

	foundLinux := false
	foundMac := false
	foundWindows := false
	foundArm := false
	for _, b := range plugin.Spec.Binaries {
		switch b.Goarch {
		case "arm64":
			switch b.Goos {
			case "Linux":
				foundArm = true
				assert.Equal(t, "https://github.com/zegl/kube-score/releases/download/v"+plugins.KubeScoreVersion+"/kube-score_"+plugins.KubeScoreVersion+"_linux_arm64.tar.gz", b.URL, "URL for linux arm binary")
				t.Logf("found linux binary URL %s", b.URL)
			}

		case "amd64":
			switch b.Goos {
			case "Darwin":
				foundMac = true
				assert.Equal(t, "https://github.com/zegl/kube-score/releases/download/v"+plugins.KubeScoreVersion+"/kube-score_"+plugins.KubeScoreVersion+"_darwin_amd64.tar.gz", b.URL, "URL for mac binary")
				t.Logf("found mac binary URL %s", b.URL)
			case "Linux":
				foundLinux = true
				assert.Equal(t, "https://github.com/zegl/kube-score/releases/download/v"+plugins.KubeScoreVersion+"/kube-score_"+plugins.KubeScoreVersion+"_linux_amd64.tar.gz", b.URL, "URL for mac binary")
				t.Logf("found linux binary URL %s", b.URL)
			case "Windows":
				foundWindows = true
				assert.Equal(t, "https://github.com/zegl/kube-score/releases/download/v"+plugins.KubeScoreVersion+"/kube-score_"+plugins.KubeScoreVersion+"_windows_amd64.tar.gz", b.URL, "URL for windows binary")
				t.Logf("found windows binary URL %s", b.URL)
			}
		}
	}
	assert.True(t, foundArm, "did not find an arm linux binary in the plugin %#v", plugin)
	assert.True(t, foundLinux, "did not find a linux binary in the plugin %#v", plugin)
	assert.True(t, foundMac, "did not find a mac binary in the plugin %#v", plugin)
	assert.True(t, foundWindows, "did not find a windows binary in the plugin %#v", plugin)
}

func TestKubevalPlugin(t *testing.T) {
	t.Parallel()

	plugin := plugins.CreateKubevalPlugin(plugins.KubevalVersion)

	assert.Equal(t, plugins.KubevalPluginName, plugin.Name, "plugin.Name")
	assert.Equal(t, plugins.KubevalPluginName, plugin.Spec.Name, "plugin.Spec.Name")

	foundLinux := false
	foundMac := false
	foundWindows := false
	foundArm := false
	for _, b := range plugin.Spec.Binaries {
		switch b.Goarch {
		case "arm64":
			switch b.Goos {
			case "Linux":
				foundArm = true
				//assert.Equal(t, "https://github.com/instrumenta/kubeval/releases/download/v"+plugins.KubevalVersion+"/kubeval-linux-arm64.tar.gz", b.URL, "URL for linux arm binary")
				assert.Equal(t, "https://github.com/jenkins-x-plugins/kubeval/releases/download/v"+plugins.KubevalVersion+"/kubeval_"+plugins.KubevalVersion+"_linux_arm64.tar.gz", b.URL, "URL for linux arm binary")
				t.Logf("found linux binary URL %s", b.URL)
			}

		case "amd64":
			switch b.Goos {
			case "Darwin":
				foundMac = true
				//assert.Equal(t, "https://github.com/instrumenta/kubeval/releases/download/v"+plugins.KubevalVersion+"/kubeval-darwin-amd64.tar.gz", b.URL, "URL for mac binary")
				assert.Equal(t, "https://github.com/jenkins-x-plugins/kubeval/releases/download/v"+plugins.KubevalVersion+"/kubeval_"+plugins.KubevalVersion+"_darwin_amd64.tar.gz", b.URL, "URL for mac binary")
				t.Logf("found mac binary URL %s", b.URL)
			case "Linux":
				foundLinux = true
				//assert.Equal(t, "https://github.com/instrumenta/kubeval/releases/download/v"+plugins.KubevalVersion+"/kubeval-linux-amd64.tar.gz", b.URL, "URL for mac binary")
				assert.Equal(t, "https://github.com/jenkins-x-plugins/kubeval/releases/download/v"+plugins.KubevalVersion+"/kubeval_"+plugins.KubevalVersion+"_linux_amd64.tar.gz", b.URL, "URL for mac binary")
				t.Logf("found linux binary URL %s", b.URL)
			case "Windows":
				foundWindows = true
				//assert.Equal(t, "https://github.com/instrumenta/kubeval/releases/download/v"+plugins.KubevalVersion+"/kubeval-windows-amd64.tar.gz", b.URL, "URL for windows binary")
				assert.Equal(t, "https://github.com/jenkins-x-plugins/kubeval/releases/download/v"+plugins.KubevalVersion+"/kubeval_"+plugins.KubevalVersion+"_windows_amd64.tar.gz", b.URL, "URL for windows binary")
				t.Logf("found windows binary URL %s", b.URL)
			}
		}
	}
	assert.True(t, foundArm, "did not find an arm linux binary in the plugin %#v", plugin)
	assert.True(t, foundLinux, "did not find a linux binary in the plugin %#v", plugin)
	assert.True(t, foundMac, "did not find a mac binary in the plugin %#v", plugin)
	assert.True(t, foundWindows, "did not find a windows binary in the plugin %#v", plugin)
}

func TestPolarisPlugin(t *testing.T) {
	t.Parallel()

	plugin := plugins.CreatePolarisPlugin(plugins.PolarisVersion)

	assert.Equal(t, plugins.PolarisPluginName, plugin.Name, "plugin.Name")
	assert.Equal(t, plugins.PolarisPluginName, plugin.Spec.Name, "plugin.Spec.Name")

	foundLinux := false
	foundMac := false
	foundWindows := false
	foundArm := false
	for _, b := range plugin.Spec.Binaries {
		switch b.Goarch {
		case "arm64":
			switch b.Goos {
			case "Linux":
				foundArm = true
				assert.Equal(t, "https://github.com/FairwindsOps/polaris/releases/download/"+plugins.PolarisVersion+"/polaris_"+plugins.PolarisVersion+"_linux_arm64.tar.gz", b.URL, "URL for linux arm binary")
				t.Logf("found linux binary URL %s", b.URL)
			}

		case "amd64":
			switch b.Goos {
			case "Darwin":
				foundMac = true
				assert.Equal(t, "https://github.com/FairwindsOps/polaris/releases/download/"+plugins.PolarisVersion+"/polaris_"+plugins.PolarisVersion+"_darwin_amd64.tar.gz", b.URL, "URL for mac binary")
				t.Logf("found mac binary URL %s", b.URL)
			case "Linux":
				foundLinux = true
				assert.Equal(t, "https://github.com/FairwindsOps/polaris/releases/download/"+plugins.PolarisVersion+"/polaris_"+plugins.PolarisVersion+"_linux_amd64.tar.gz", b.URL, "URL for mac binary")
				t.Logf("found linux binary URL %s", b.URL)
			case "Windows":
				foundWindows = true
				assert.Equal(t, "https://github.com/FairwindsOps/polaris/releases/download/"+plugins.PolarisVersion+"/polaris_"+plugins.PolarisVersion+"_windows_amd64.tar.gz", b.URL, "URL for windows binary")
				t.Logf("found windows binary URL %s", b.URL)
			}
		}
	}
	assert.True(t, foundArm, "did not find an arm linux binary in the plugin %#v", plugin)
	assert.True(t, foundLinux, "did not find a linux binary in the plugin %#v", plugin)
	assert.True(t, foundMac, "did not find a mac binary in the plugin %#v", plugin)
	assert.True(t, foundWindows, "did not find a windows binary in the plugin %#v", plugin)
}
