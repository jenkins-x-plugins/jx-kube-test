package run_test

import (
	"github.com/jenkins-x-plugins/jx-kube-test/pkg/cmd/run"
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
