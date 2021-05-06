package run_test

import (
	"github.com/jenkins-x-plugins/jx-kube-test/pkg/apis/kubetest/v1alpha1"
	"github.com/jenkins-x-plugins/jx-kube-test/pkg/cmd/run"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jenkins-x/jx-helpers/v3/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

var (
	// generateTestOutput enable to regenerate the expected output
	generateTestOutput = true
)

func TestCmdRun(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err, "failed to create temp dir")

	testDir := filepath.Join("test_data")
	fs, err := ioutil.ReadDir(testDir)
	require.NoError(t, err, "failed to read test dir %s", testDir)
	for _, f := range fs {
		if f == nil || !f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		srcDir := filepath.Join(testDir, name)
		runDir := filepath.Join(tmpDir, name)

		t.Logf("running test %s in dir %s\n", name, runDir)

		_, o := run.NewCmdRun()
		o.Dir = srcDir
		o.ChartsDir = srcDir
		o.CommandRunner = cmdrunner.DefaultCommandRunner
		err = o.Run()
		require.NoError(t, err, "failed to run the command")

		f := o.OutFile
		if f == "" {
			t.Logf("no output file")
			continue
		}
		require.FileExists(t, f, "should have generated file")
		t.Logf("generated file %s\n", f)

		if generateTestOutput {
			generatedFile := f
			expectedPath := filepath.Join(srcDir, "expected.sh")
			data, err := ioutil.ReadFile(generatedFile)
			require.NoError(t, err, "failed to load %s", generatedFile)

			err = ioutil.WriteFile(expectedPath, data, 0666)
			require.NoError(t, err, "failed to save file %s", expectedPath)

			t.Logf("saved file %s\n", expectedPath)
			continue
		}

		testhelpers.AssertTextFilesEqual(t, filepath.Join(runDir, "expected.sh"), f, "generated file")
	}

}

func TestAddFormatFlags(t *testing.T) {
	testCases := []struct {
		args     []string
		format   string
		expected []string
	}{
		{
			args:     []string{"-o", "json"},
			format:   "tap",
			expected: []string{"-o", "json"},
		},
		{
			args:     []string{"-o", "tap"},
			format:   "tap",
			expected: []string{"-o", "tap"},
		},
		{
			format:   "tap",
			expected: []string{"--output", "tap"},
		},
		{
			args:     []string{"-c", "cheese"},
			format:   "tap",
			expected: []string{"-c", "cheese", "--output", "tap"},
		},
	}

	flag := "o"
	optionName := "output"

	for _, tc := range testCases {
		settings := &v1alpha1.KubeTest{}
		settings.Spec.Format = tc.format
		got := run.AddFormatFlags(settings, flag, optionName, tc.args)

		assert.Equal(t, tc.expected, got, "for format %s", tc.format)

		t.Logf("got args: %v for format %s\n", got, tc.format)
	}
}
