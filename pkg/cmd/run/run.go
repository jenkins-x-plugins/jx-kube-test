package run

import (
	"github.com/jenkins-x-plugins/jx-gitops/pkg/plugins"
	ktplugins "github.com/jenkins-x-plugins/jx-kube-test/pkg/plugins"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	info = termcolor.ColorInfo

	cmdLong = templates.LongDesc(`
		Runs all of the kubernetes tests
`)

	cmdExample = templates.Examples(`
		# runs the kubernetes tests 
		jx kube test run
	`)
)

// Options the options for the command
type Options struct {
	options.BaseOptions

	Dir             string
	WorkDir         string
	OutFile         string
	HelmBinary      string
	PolarisBinary   string
	ConftestBinary  string
	KubeScoreBinary string
	KubevalBinary   string
	ConftestArgs    []string
	KubeScoreArgs   []string
	KubevalArgs     []string
	PolarisArgs     []string
	UseHelmPlugin   bool
	CommandRunner   cmdrunner.CommandRunner
}

type ChartOutput struct {
	ChartDir    string
	ReleaseName string
	OutputDir   string
}

// NewCmdRun creates a command object for the command
func NewCmdRun() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Runs all of the kubernetes tests",
		Long:    cmdLong,
		Example: cmdExample,
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}
	o.BaseOptions.AddBaseFlags(cmd)
	cmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "the directory to look for helm, helmfile or kustomize files")
	cmd.Flags().StringVarP(&o.WorkDir, "work-dir", "w", "", "the work directory used to generate the output. If not specified a new temporary dir is created")
	cmd.Flags().StringVarP(&o.OutFile, "output", "o", "", "the file to generate")
	cmd.Flags().StringVarP(&o.HelmBinary, "helm-binary", "", "", "specifies the helm binary location to use. If not specified we download the plugin")
	cmd.Flags().StringVarP(&o.ConftestBinary, "conftest-binary", "", "", "specifies the conftest binary location to use. If not specified we download the plugin")
	cmd.Flags().StringVarP(&o.KubeScoreBinary, "kube-score-binary", "", "", "specifies the kube-score binary location to use. If not specified we download the plugin")
	cmd.Flags().StringVarP(&o.KubevalBinary, "kubeval-binary", "", "", "specifies the kubeval binary location to use. If not specified we download the plugin")
	cmd.Flags().StringVarP(&o.PolarisBinary, "polaris-binary", "n", "", "specifies the polaris binary location to use. If not specified we download the plugin")
	cmd.Flags().BoolVarP(&o.UseHelmPlugin, "use-helm-plugin", "", false, "uses the jx binary plugin for helm rather than whatever helm is on the $PATH")
	return cmd, o
}

// Run implements the command
func (o *Options) Validate() error {
	err := o.BaseOptions.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate options")
	}
	if o.CommandRunner == nil {
		o.CommandRunner = cmdrunner.QuietCommandRunner
	}
	if o.HelmBinary == "" {
		if o.UseHelmPlugin {
			o.HelmBinary, err = plugins.GetHelmBinary(plugins.HelmVersion)
			if err != nil {
				return err
			}
		}
		if o.HelmBinary == "" {
			o.HelmBinary = "helm"
		}
	}
	if o.WorkDir == "" {
		o.WorkDir, err = ioutil.TempDir("", "")
		if err != nil {
			return errors.Wrapf(err, "failed to create temp dir")
		}
	}
	return nil
}

// Run implements the command
func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate")
	}

	dir := o.Dir
	exists, err := files.DirExists(dir)
	if err != nil {
		return errors.Wrapf(err, "failed to check if dir exists %s", dir)
	}
	if !exists {
		return options.InvalidOptionf("dir", o.Dir, "dir does not exist")
	}

	chartDirs, err := o.findChartDirs(err, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to find charts in dir %s", dir)
	}

	if len(chartDirs) == 0 {
		log.Logger().Infof("no charts found in dir %s", dir)
		return nil
	}

	for _, d := range chartDirs {

		rel, err := filepath.Rel(dir, d)
		if err != nil {
			return errors.Wrapf(err, "failed to find relative chart dir from %s to %s", rel, d)
		}

		// TODO have a loop for all the different value sets to pass in....
		i := 1
		releaseName := "rel" + strconv.Itoa(i)
		outDir := filepath.Join(o.WorkDir, rel, releaseName)

		err = os.MkdirAll(outDir, files.DefaultDirWritePermissions)
		if err != nil {
			return errors.Wrapf(err, "failed to create output dir %s", outDir)
		}

		args := []string{"template", "--output-dir", outDir, releaseName, d}
		c := &cmdrunner.Command{
			Name: o.HelmBinary,
			Args: args,
			Out:  os.Stdout,
			Err:  os.Stderr,
		}

		_, err = o.CommandRunner(c)
		if err != nil {
			return errors.Wrapf(err, "failed to run %s", c.CLI())
		}

		co := &ChartOutput{
			ChartDir:    d,
			ReleaseName: releaseName,
			OutputDir:   outDir,
		}
		err = o.verifyChart(co)
		if err != nil {
			return errors.Wrapf(err, "failed to verify chart output for %s", d)
		}
	}

	return nil
}

// findChartDirs lets find all the charts in the given dir
func (o *Options) findChartDirs(err error, dir string) ([]string, error) {
	var chartDirs []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		name := info.Name()
		if info.IsDir() || name != "Chart.yaml" {
			return nil
		}

		chartDir := filepath.Dir(path)
		log.Logger().Infof("found chart in dir %s", chartDir)
		chartDirs = append(chartDirs, chartDir)
		return nil
	})
	return chartDirs, err
}

func (o *Options) verifyChart(co *ChartOutput) error {
	log.Logger().Infof("verifying chart %s output at %s", co.ChartDir, co.OutputDir)

	err := o.kubeval(co)
	if err != nil {
		return errors.Wrapf(err, "failed to run kubeval on chart %s", co.ChartDir)
	}
	err = o.conftest(co)
	if err != nil {
		return errors.Wrapf(err, "failed to run conftest on chart %s", co.ChartDir)
	}
	err = o.polaris(co)
	if err != nil {
		return errors.Wrapf(err, "failed to run polaris on chart %s", co.ChartDir)
	}
	err = o.kubeScore(co)
	if err != nil {
		return errors.Wrapf(err, "failed to run kube-score on chart %s", co.ChartDir)
	}
	return nil
}

func (o *Options) kubeval(co *ChartOutput) error {
	bin, err := o.GetKubevalBinary()
	if err != nil {
		return errors.Wrapf(err, "failed to get the kubeval binary")
	}

	args := []string{"-o", "json", "-d", co.OutputDir}
	args = append(args, o.KubevalArgs...)
	c := &cmdrunner.Command{
		Name: bin,
		Args: args,
	}

	log.Logger().Infof("kubeval is verifying chart %s...", co.ChartDir)

	text, err := o.CommandRunner(c)
	if err != nil {
		return errors.Wrapf(err, "failed to run %s", c.CLI())
	}

	log.Logger().Infof("result:\n%s", termcolor.ColorStatus(text))
	return nil
}

func (o *Options) kubeScore(co *ChartOutput) error {
	bin, err := o.GetKubeScoreBinary()
	if err != nil {
		return errors.Wrapf(err, "failed to get the kube-score binary")
	}

	fileNames, err := o.findYAMLFiles(co.OutputDir)
	if err != nil {
		return errors.Wrapf(err, "failed to find YAML files in dir %s", co.OutputDir)
	}
	if len(fileNames) == 0 {
		log.Logger().Warnf("no YAML files found for chart %s in output dir %s", co.ChartDir, co.OutputDir)
		return nil
	}

	args := []string{"score", "-o", "json"}
	args = append(args, o.KubeScoreArgs...)
	args = append(args, fileNames...)
	c := &cmdrunner.Command{
		Name: bin,
		Args: args,
	}

	log.Logger().Infof("kube-score is verifying chart %s...", co.ChartDir)

	text, err := o.CommandRunner(c)
	if err != nil {
		log.Logger().Infof("kube-score returned error %s", err.Error())
		//return errors.Wrapf(err, "failed to run %s", c.CLI())
	}

	log.Logger().Infof("result:\n%s", termcolor.ColorStatus(text))
	return nil
}

func (o *Options) conftest(co *ChartOutput) error {
	bin, err := o.GetConftestBinary()
	if err != nil {
		return errors.Wrapf(err, "failed to get the polaris binary")
	}

	args := []string{"test", "-o", "json", co.OutputDir}
	args = append(args, o.ConftestArgs...)
	c := &cmdrunner.Command{
		Name: bin,
		Args: args,
	}

	log.Logger().Infof("conftest is verifying chart %s...", co.ChartDir)

	text, err := o.CommandRunner(c)
	if err != nil {
		log.Logger().Infof("conftest returned error %s", err.Error())
		//return errors.Wrapf(err, "failed to run %s", c.CLI())
	}

	log.Logger().Infof("result:\n%s", termcolor.ColorStatus(text))
	return nil
}

func (o *Options) polaris(co *ChartOutput) error {
	bin, err := o.GetPolarisBinary()
	if err != nil {
		return errors.Wrapf(err, "failed to get the polaris binary")
	}

	args := []string{"audit", "--audit-path", co.OutputDir, "--format", "json"}
	args = append(args, o.PolarisArgs...)
	c := &cmdrunner.Command{
		Name: bin,
		Args: args,
	}

	log.Logger().Infof("polaris is verifying chart %s...", co.ChartDir)

	text, err := o.CommandRunner(c)
	if err != nil {
		log.Logger().Infof("polaris returned error %s", err.Error())
		//return errors.Wrapf(err, "failed to run %s", c.CLI())
	}

	log.Logger().Infof("result:\n%s", termcolor.ColorStatus(text))
	return nil
}

func (o *Options) GetConftestBinary() (string, error) {
	if o.ConftestBinary != "" {
		return o.ConftestBinary, nil
	}
	var err error
	o.ConftestBinary, err = ktplugins.GetConftestBinary(ktplugins.ConftestVersion)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download conftest plugin")
	}
	return o.ConftestBinary, nil
}

func (o *Options) GetPolarisBinary() (string, error) {
	if o.PolarisBinary != "" {
		return o.PolarisBinary, nil
	}
	var err error
	o.PolarisBinary, err = ktplugins.GetPolarisBinary(ktplugins.PolarisVersion)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download polaris plugin")
	}
	return o.PolarisBinary, nil
}

func (o *Options) GetKubeScoreBinary() (string, error) {
	if o.KubeScoreBinary != "" {
		return o.KubeScoreBinary, nil
	}
	var err error
	o.KubeScoreBinary, err = ktplugins.GetKubeScoreBinary(ktplugins.KubeScoreVersion)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download kube-score plugin")
	}
	return o.KubeScoreBinary, nil
}

func (o *Options) GetKubevalBinary() (string, error) {
	if o.KubevalBinary != "" {
		return o.KubevalBinary, nil
	}
	var err error
	o.KubevalBinary, err = ktplugins.GetKubevalBinary(ktplugins.KubevalVersion)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download kubeval plugin")
	}
	return o.KubevalBinary, nil
}

func (o *Options) findYAMLFiles(dir string) ([]string, error) {
	var answer []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		name := info.Name()
		if info.IsDir() || !strings.HasSuffix(name, ".yaml") {
			return nil
		}

		answer = append(answer, path)
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find YAML files in dir %s", dir)
	}
	return answer, nil
}
