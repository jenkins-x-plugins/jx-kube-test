package run

import (
	"fmt"
	"github.com/jenkins-x-plugins/jx-gitops/pkg/plugins"
	"github.com/jenkins-x-plugins/jx-kube-test/pkg/apis/kubetest/v1alpha1"
	ktplugins "github.com/jenkins-x-plugins/jx-kube-test/pkg/plugins"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/jenkins-x/jx-helpers/v3/pkg/yamls"
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
	SettingsFile    string
	WorkDir         string
	OutFile         string
	ChartsDir       string
	RecurseCharts   bool
	Helm            BinaryPlugin
	ConftestPlugin  BinaryPlugin
	HelmPlugin      BinaryPlugin
	KubeScorePlugin BinaryPlugin
	KubevalPlugin   BinaryPlugin
	PolarisPlugin   BinaryPlugin
	CommandRunner   cmdrunner.CommandRunner
	Settings        *v1alpha1.KubeTest
}

type ResourceLocation struct {
	Description string
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

	o.ConftestPlugin.AddFlags(cmd, "conftest", ktplugins.ConftestVersion, ktplugins.GetConftestBinary)
	o.Helm.AddFlags(cmd, "helm", plugins.HelmVersion, plugins.GetHelmBinary)
	o.KubeScorePlugin.AddFlags(cmd, "kubescore", ktplugins.KubeScoreVersion, ktplugins.GetKubeScoreBinary)
	o.KubevalPlugin.AddFlags(cmd, "kubeval", ktplugins.KubevalVersion, ktplugins.GetKubevalBinary)
	o.PolarisPlugin.AddFlags(cmd, "polaris", ktplugins.PolarisVersion, ktplugins.GetPolarisBinary)

	cmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "the directory to look for helm, helmfile or kustomize files")
	cmd.Flags().StringVarP(&o.ChartsDir, "chart-dir", "", "charts", "the directory to look for helm charts if no .jx/kube-test/settings.yaml file is found")
	cmd.Flags().BoolVarP(&o.RecurseCharts, "recurse", "r", true, "should we recurse through the chart dir to find charts if no .jx/kube-test/settings.yaml file is found")
	cmd.Flags().StringVarP(&o.SettingsFile, "settings", "s", "", "the settings file to use. If not specified will look in .jx/kube-test/settings.yaml in the directory")
	cmd.Flags().StringVarP(&o.WorkDir, "work-dir", "w", "", "the work directory used to generate the output. If not specified a new temporary dir is created")
	cmd.Flags().StringVarP(&o.OutFile, "output", "o", "", "the file to generate")
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

	if o.WorkDir == "" {
		o.WorkDir, err = ioutil.TempDir("", "")
		if err != nil {
			return errors.Wrapf(err, "failed to create temp dir")
		}
	}

	if o.SettingsFile == "" {
		o.SettingsFile = filepath.Join(o.Dir, ".jx", "kube-test", "settings.yaml")
	}

	if o.Settings == nil {
		exists, err := files.FileExists(o.SettingsFile)
		if err != nil {
			return errors.Wrapf(err, "failed to check if file exists %s", o.SettingsFile)
		}
		if exists {
			o.Settings = &v1alpha1.KubeTest{}
			err = yamls.LoadFile(o.SettingsFile, o.Settings)
			if err != nil {
				return errors.Wrapf(err, "failed to load file %s", o.SettingsFile)
			}
			log.Logger().Debugf("loaded settings file %s", info(o.SettingsFile))
		} else {
			o.Settings = o.createDefaultSettings()
		}
	}
	if o.Settings == nil {
		return errors.Errorf("failed to discover or generate settings")
	}
	if o.Settings.Spec.OutputDir != "" {
		if o.Settings.Spec.Format == "" {
			log.Logger().Warnf("no spec.format is specified in %s so defaulting to 'tap'", o.SettingsFile)
			o.Settings.Spec.Format = "tap"
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

	for i := range o.Settings.Spec.Rules {
		rule := &o.Settings.Spec.Rules[i]

		if rule.Charts != nil {
			err = o.TestCharts(rule, rule.Charts)
			if err != nil {
				return errors.Wrapf(err, "failed to test charts at %s", rule.Charts.Dir)
			}
			continue
		}
		if rule.Resources != nil {
			err = o.TestResources(rule, rule.Resources)
			if err != nil {
				return errors.Wrapf(err, "failed to test resources at %s", rule.Resources.Dir)
			}
			continue
		}
		return errors.Errorf("invalid rule %#v has neither charts or resources", rule)
	}
	return nil
}

// TestResources tests the resources
func (o *Options) TestResources(rule *v1alpha1.Rule, resources *v1alpha1.Source) error {
	dir := resources.Dir
	exists, err := files.DirExists(dir)
	if err != nil {
		return errors.Wrapf(err, "failed to check if dir exists %s", dir)
	}
	if !exists {
		return errors.Errorf("the resource dir %s does not exist", dir)
	}

	co := &ResourceLocation{
		Description: fmt.Sprintf("resources %s", dir),
		OutputDir:   dir,
	}
	err = o.verifyResources(co, &rule.Tests)
	if err != nil {
		return errors.Wrapf(err, "failed to verify resources in dir  %s", dir)
	}
	return nil
}

// TestCharts tests the charts
func (o *Options) TestCharts(rule *v1alpha1.Rule, charts *v1alpha1.Charts) error {
	dir := charts.Dir
	exists, err := files.DirExists(dir)
	if err != nil {
		return errors.Wrapf(err, "failed to check if dir exists %s", dir)
	}
	if !exists {
		return errors.Errorf("the charts dir %s does not exist", dir)
	}

	helmbin, err := o.Helm.GetBinary(nil)
	if err != nil {
		return errors.Wrapf(err, "failed to find helm binary")
	}

	if !charts.Recurse {
		path := filepath.Join(dir, "Chart.yaml")
		exists, err = files.FileExists(path)
		if err != nil {
			return errors.Wrapf(err, "failed to check if file exists %s", path)
		}
		if !exists {
			return errors.Errorf("the charts dir %s does not contain a Chart.yaml file. You can enable 'recurse: true' to find charts inside the directory", dir)
		}

		err = o.helmTemplateAndVerify(rule, helmbin, dir)
		if err != nil {
			return errors.Wrapf(err, "failed to test chart %s", dir)
		}
		return nil
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
		err = o.helmTemplateAndVerify(rule, helmbin, d)
		if err != nil {
			return errors.Wrapf(err, "failed to test chart %s", d)
		}
	}

	return nil
}

func (o *Options) helmTemplateAndVerify(rule *v1alpha1.Rule, helmbin string, d string) error {
	rel, err := filepath.Rel(o.Dir, d)
	if err != nil {
		log.Logger().Warnf("failed to find relative chart dir from %s to %s", o.Dir, d)
		rel = d
	}

	// TODO support iterating through available values to pass in
	i := 1
	releaseName := "rel" + strconv.Itoa(i)
	outDir := filepath.Join(o.WorkDir, rel, releaseName)

	err = os.MkdirAll(outDir, files.DefaultDirWritePermissions)
	if err != nil {
		return errors.Wrapf(err, "failed to create output dir %s", outDir)
	}

	args := []string{"template", "--output-dir", outDir, releaseName, d}
	c := &cmdrunner.Command{
		Name: helmbin,
		Args: args,
		Out:  os.Stdout,
		Err:  os.Stderr,
	}

	_, err = o.CommandRunner(c)
	if err != nil {
		return errors.Wrapf(err, "failed to run %s", c.CLI())
	}

	co := &ResourceLocation{
		Description: fmt.Sprintf("chart %s release %s", d, releaseName),
		OutputDir:   outDir,
	}
	err = o.verifyResources(co, &rule.Tests)
	if err != nil {
		return errors.Wrapf(err, "failed to verify chart output for %s", d)
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

func (o *Options) verifyResources(co *ResourceLocation, tests *v1alpha1.Tests) error {
	log.Logger().Infof("verifying %s output at %s", co.Description, co.OutputDir)

	if tests.Kubeval != nil {
		err := o.kubeval(co, tests.Kubeval)
		if err != nil {
			return errors.Wrapf(err, "failed to run kubeval on %s", co.Description)
		}
	}
	if tests.Conftest != nil {
		err := o.conftest(co, tests.Conftest)
		if err != nil {
			return errors.Wrapf(err, "failed to run conftest on %s", co.Description)
		}
	}
	if tests.Kubescore != nil {
		err := o.kubescore(co, tests.Kubescore)
		if err != nil {
			return errors.Wrapf(err, "failed to run kube-score on %s", co.Description)
		}
	}
	if tests.Polaris != nil {
		err := o.polaris(co, tests.Polaris)
		if err != nil {
			return errors.Wrapf(err, "failed to run polaris on %s", co.Description)
		}
	}
	return nil
}

func (o *Options) kubeval(co *ResourceLocation, t *v1alpha1.Test) error {
	if len(t.Args) > 0 {
		o.KubevalPlugin.Args = t.Args
	}
	bin, err := o.KubevalPlugin.GetBinary(nil)
	if err != nil {
		return errors.Wrapf(err, "failed to get the kubeval binary")
	}

	args := []string{"-d", co.OutputDir}
	args = append(args, o.KubevalPlugin.Args...)
	args = append(args, t.Args...)
	args = AddFormatFlags(o.Settings, "o", "output", args)

	c := &cmdrunner.Command{
		Name: bin,
		Args: args,
	}
	return o.runTestCommand("kubeval", co, c)
}

func (o *Options) kubescore(co *ResourceLocation, t *v1alpha1.Test) error {
	bin, err := o.KubeScorePlugin.GetBinary(t)
	if err != nil {
		return errors.Wrapf(err, "failed to get the kube-score binary")
	}

	fileNames, err := o.findYAMLFiles(co.OutputDir)
	if err != nil {
		return errors.Wrapf(err, "failed to find YAML files in dir %s", co.OutputDir)
	}
	if len(fileNames) == 0 {
		log.Logger().Warnf("no YAML files found for %s in output dir %s", co.Description, co.OutputDir)
		return nil
	}

	args := []string{"score"}
	args = append(args, o.KubeScorePlugin.Args...)
	args = append(args, t.Args...)
	args = append(args, fileNames...)
	args = AddFormatFlags(o.Settings, "o", "output", args)
	c := &cmdrunner.Command{
		Name: bin,
		Args: args,
	}

	log.Logger().Infof("kube-score is verifying %s...", co.Description)

	text, err := o.CommandRunner(c)
	if err != nil {
		log.Logger().Infof("kube-score returned error %s", err.Error())
		//return errors.Wrapf(err, "failed to run %s", c.CLI())
	}

	log.Logger().Infof("result:\n%s", termcolor.ColorStatus(text))
	return nil
}

func (o *Options) conftest(co *ResourceLocation, t *v1alpha1.Test) error {
	bin, err := o.ConftestPlugin.GetBinary(t)
	if err != nil {
		return errors.Wrapf(err, "failed to get the polaris binary")
	}

	args := []string{"test", co.OutputDir}
	args = append(args, o.ConftestPlugin.Args...)
	args = append(args, t.Args...)
	args = AddFormatFlags(o.Settings, "o", "output", args)
	c := &cmdrunner.Command{
		Name: bin,
		Args: args,
	}

	log.Logger().Infof("conftest is verifying %s...", co.Description)

	text, err := o.CommandRunner(c)
	if err != nil {
		log.Logger().Infof("conftest returned error %s", err.Error())
		//return errors.Wrapf(err, "failed to run %s", c.CLI())
	}

	log.Logger().Infof("result:\n%s", termcolor.ColorStatus(text))
	return nil
}

func (o *Options) polaris(co *ResourceLocation, t *v1alpha1.Test) error {
	bin, err := o.PolarisPlugin.GetBinary(t)
	if err != nil {
		return errors.Wrapf(err, "failed to get the polaris binary")
	}

	args := []string{"audit", "--audit-path", co.OutputDir}
	args = append(args, o.PolarisPlugin.Args...)
	args = append(args, t.Args...)
	args = AddFormatFlags(o.Settings, "f", "format", args)
	c := &cmdrunner.Command{
		Name: bin,
		Args: args,
	}

	log.Logger().Infof("polaris is verifying %s...", co.Description)

	text, err := o.CommandRunner(c)
	if err != nil {
		log.Logger().Infof("polaris returned error %s", err.Error())
		//return errors.Wrapf(err, "failed to run %s", c.CLI())
	}

	log.Logger().Infof("result:\n%s", termcolor.ColorStatus(text))
	return nil
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

func (o *Options) createDefaultSettings() *v1alpha1.KubeTest {
	return &v1alpha1.KubeTest{
		Spec: v1alpha1.KubeTestSpec{
			Rules: []v1alpha1.Rule{
				{
					Charts: &v1alpha1.Charts{
						Dir:     o.ChartsDir,
						Recurse: o.RecurseCharts,
					},
					Tests: v1alpha1.Tests{
						Kubeval: &v1alpha1.Test{},
						Polaris: &v1alpha1.Test{},
					},
				},
			},
		},
	}
}

func (o *Options) runTestCommand(name string, co *ResourceLocation, c *cmdrunner.Command) error {
	outputDir := o.Settings.Spec.OutputDir
	format := o.Settings.Spec.Format

	if outputDir != "" {
		err := os.MkdirAll(outputDir, files.DefaultDirWritePermissions)
		if err != nil {
			return errors.Wrapf(err, "failed to create dir %s", outputDir)
		}
	}

	log.Logger().Infof("%s is verifying %s...", name, co.Description)

	text, err := o.CommandRunner(c)
	if err != nil {
		log.Logger().Infof("%s returned error", name)
		log.Logger().Debugf("%s returned error %s", name, err.Error())
	}

	if outputDir != "" {
		path := filepath.Join(outputDir, name+"."+format)
		err = ioutil.WriteFile(path, []byte(text), files.DefaultFileWritePermissions)
		if err != nil {
			return errors.Wrapf(err, "failed to save file %s", path)
		}
		log.Logger().Infof("saved %s results in %s", name, info(path))
		return nil
	}
	log.Logger().Infof("result:\n%s", termcolor.ColorStatus(text))
	return nil
}

// AddFormatFlags if the format is specified lets add it as a command line argument
func AddFormatFlags(settings *v1alpha1.KubeTest, flag string, optionName string, args []string) []string {
	if settings.Spec.Format == "" {
		return args
	}

	// lets check if we already have the format flag
	oldValue := options.ArgumentsOptionValue(args, flag, optionName)
	if oldValue != "" {
		return args
	}

	args = append(args, "--"+optionName, settings.Spec.Format)
	return args
}
