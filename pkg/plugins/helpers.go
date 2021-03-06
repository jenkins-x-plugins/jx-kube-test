package plugins

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"

	jenkinsv1 "github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-helpers/v3/pkg/extensions"
	"github.com/jenkins-x/jx-helpers/v3/pkg/homedir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PluginBinDir returns the plugin dir
func PluginBinDir() (string, error) {
	return PluginBinDirFunc(os.Getenv)
}

// PluginBinDirFunc uses a function for looking up env vars for easier testing
func PluginBinDirFunc(fn func(string) string) (string, error) {
	for _, e := range []string{"JX_GITOPS_HOME", "JX3_HOME", "JX_HOME"} {
		v := fn(e)
		if v != "" {
			return homedir.PluginBinDir(v, ".jx")
		}
	}
	return homedir.PluginBinDir("", ".jx")
}

// GetConftestBinary returns the path to the locally installed kube-score extension
func GetConftestBinary(version string) (string, error) {
	if version == "" {
		version = ConftestVersion
	}
	pluginBinDir, err := PluginBinDir()
	if err != nil {
		return "", errors.Wrapf(err, "failed to find plugin home dir")
	}
	plugin := CreateConftestPlugin(version)
	return extensions.EnsurePluginInstalled(plugin, pluginBinDir)
}

// CreateConftestPlugin creates the kube-score plugin
func CreateConftestPlugin(version string) jenkinsv1.Plugin {
	binaries := extensions.CreateBinaries(func(p extensions.Platform) string {
		goos := p.Goos
		goarch := strings.ToLower(p.Goarch)
		if goarch == "amd64" {
			goarch = "x86_64"
		}
		ext := ".tar.gz"
		if p.IsWindows() {
			ext = ".zip"
		}
		return fmt.Sprintf("https://github.com/open-policy-agent/conftest/releases/download/v%s/conftest_%s_%s_%s%s", version, version, goos, goarch, ext)
	})

	plugin := jenkinsv1.Plugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: ConftestPluginName,
		},
		Spec: jenkinsv1.PluginSpec{
			SubCommand:  "kube-score",
			Binaries:    binaries,
			Description: "kube score binary",
			Name:        ConftestPluginName,
			Version:     version,
		},
	}
	return plugin
}

// GetKubeScoreBinary returns the path to the locally installed kube-score extension
func GetKubeScoreBinary(version string) (string, error) {
	if version == "" {
		version = KubeScoreVersion
	}
	pluginBinDir, err := PluginBinDir()
	if err != nil {
		return "", errors.Wrapf(err, "failed to find plugin home dir")
	}
	plugin := CreateKubeScorePlugin(version)
	return extensions.EnsurePluginInstalled(plugin, pluginBinDir)
}

// CreateKubeScorePlugin creates the kube-score plugin
func CreateKubeScorePlugin(version string) jenkinsv1.Plugin {
	binaries := extensions.CreateBinaries(func(p extensions.Platform) string {
		ext := ".tar.gz"
		return fmt.Sprintf("https://github.com/zegl/kube-score/releases/download/v%s/kube-score_%s_%s_%s%s", version, version, strings.ToLower(p.Goos), strings.ToLower(p.Goarch), ext)
	})

	plugin := jenkinsv1.Plugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: KubeScorePluginName,
		},
		Spec: jenkinsv1.PluginSpec{
			SubCommand:  "kube-score",
			Binaries:    binaries,
			Description: "kube score binary",
			Name:        KubeScorePluginName,
			Version:     version,
		},
	}
	return plugin
}

// GetKubevalBinary returns the path to the locally installed kube-score extension
func GetKubevalBinary(version string) (string, error) {
	if version == "" {
		version = KubevalVersion
	}
	pluginBinDir, err := PluginBinDir()
	if err != nil {
		return "", errors.Wrapf(err, "failed to find plugin home dir")
	}
	plugin := CreateKubevalPlugin(version)
	return extensions.EnsurePluginInstalled(plugin, pluginBinDir)
}

// CreateKubevalPlugin creates the kube-score plugin
func CreateKubevalPlugin(version string) jenkinsv1.Plugin {
	binaries := extensions.CreateBinaries(func(p extensions.Platform) string {
		ext := ".tar.gz"
		//return fmt.Sprintf("https://github.com/instrumenta/kubeval/releases/download/v%s/kubeval-%s-%s%s", version, strings.ToLower(p.Goos), strings.ToLower(p.Goarch), ext)
		return fmt.Sprintf("https://github.com/jenkins-x-plugins/kubeval/releases/download/v%s/kubeval_%s_%s_%s%s", version, version, strings.ToLower(p.Goos), strings.ToLower(p.Goarch), ext)
	})

	plugin := jenkinsv1.Plugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: KubevalPluginName,
		},
		Spec: jenkinsv1.PluginSpec{
			SubCommand:  "kube-score",
			Binaries:    binaries,
			Description: "kube score binary",
			Name:        KubevalPluginName,
			Version:     version,
		},
	}
	return plugin
}

// GetPolarisBinary returns the path to the locally installed kube-score extension
func GetPolarisBinary(version string) (string, error) {
	if version == "" {
		version = PolarisVersion
	}
	pluginBinDir, err := PluginBinDir()
	if err != nil {
		return "", errors.Wrapf(err, "failed to find plugin home dir")
	}
	plugin := CreatePolarisPlugin(version)
	return extensions.EnsurePluginInstalled(plugin, pluginBinDir)
}

// CreatePolarisPlugin creates the kube-score plugin
func CreatePolarisPlugin(version string) jenkinsv1.Plugin {
	binaries := extensions.CreateBinaries(func(p extensions.Platform) string {
		ext := ".tar.gz"
		return fmt.Sprintf("https://github.com/FairwindsOps/polaris/releases/download/%s/polaris_%s_%s_%s%s", version, version, strings.ToLower(p.Goos), strings.ToLower(p.Goarch), ext)
	})

	plugin := jenkinsv1.Plugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: PolarisPluginName,
		},
		Spec: jenkinsv1.PluginSpec{
			SubCommand:  "kube-score",
			Binaries:    binaries,
			Description: "kube score binary",
			Name:        PolarisPluginName,
			Version:     version,
		},
	}
	return plugin
}
