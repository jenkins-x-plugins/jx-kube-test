package run

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type BinaryPlugin struct {
	Name       string
	DownloadFn func(string) (string, error)
	Binary     string
	Version    string
	Args       []string
}

func (o *BinaryPlugin) GetBinary() (string, error) {
	if o.Binary != "" {
		return o.Binary, nil
	}
	var err error
	o.Binary, err = o.DownloadFn(o.Version)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download %s plugin", o.Name)
	}
	return o.Binary, nil
}

func (o *BinaryPlugin) AddFlags(cmd *cobra.Command, name string, version string, fn func(version string) (string, error)) {
	o.Name = name
	o.Version = version
	if fn != nil {
		o.DownloadFn = fn
	}

	cmd.Flags().StringVarP(&o.Binary, name+"-binary", "", "", fmt.Sprintf("specifies the %s binary location to use. If not specified we download the plugin", name))
	cmd.Flags().StringVarP(&o.Version, name+"-version", "", version, fmt.Sprintf("specifies the %s version to use. If not specified we download the plugin", name))
	cmd.Flags().StringArrayVarP(&o.Args, name+"-args", "", nil, fmt.Sprintf("specifies any optional %s command line arguments to pass", name))
}
