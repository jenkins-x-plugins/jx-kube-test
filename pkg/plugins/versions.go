package plugins

import jenkinsv1 "github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"

const (
	// ConftestPluginName the default name of the conftest plugin
	ConftestPluginName = "conftest"

	// ConftestVersion the default version of conftest to use
	ConftestVersion = "0.24.0"

	// KubeScorePluginName the default name of the kube-score plugin
	KubeScorePluginName = "kube-score"

	// KubeScoreVersion the default version of kube-score to use
	KubeScoreVersion = "1.11.0"

	// KubevalPluginName the default name of the kubeval plugin
	KubevalPluginName = "kubeval"

	// KubevalVersion the default version of kubeval to use
	KubevalVersion = "0.16.1"

	// PolarisPluginName the default name of the polaris plugin
	PolarisPluginName = "polaris"

	// PolarisVersion the default version of polaris to use
	PolarisVersion = "3.2.1"
)

var (
	// Plugins default plugins
	Plugins = []jenkinsv1.Plugin{
		CreateConftestPlugin(ConftestVersion),
		CreatePolarisPlugin(KubeScoreVersion),
		CreateKubeScorePlugin(KubeScoreVersion),
		CreateKubevalPlugin(KubevalVersion),
	}
)
