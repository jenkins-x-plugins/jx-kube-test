package plugins

import jenkinsv1 "github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"

const (
	// KubeScorePluginName the default name of the kube-score plugin
	KubeScorePluginName = "kube-score"

	// KubeScoreVersion the default version of kube-score to use
	KubeScoreVersion = "1.11.0"

	// PolarisPluginName the default name of the kube-score plugin
	PolarisPluginName = "polaris"

	// PolarisVersion the default version of kube-score to use
	PolarisVersion = "3.2.1"
)

var (
	// Plugins default plugins
	Plugins = []jenkinsv1.Plugin{
		CreateKubeScorePlugin(KubeScoreVersion),
	}
)
