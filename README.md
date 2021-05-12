# jx-kube-test

[![Documentation](https://godoc.org/github.com/jenkins-x-plugins/jx-kube-test?status.svg)](https://pkg.go.dev/mod/github.com/jenkins-x-plugins/jx-kube-test)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x-plugins/jx-kube-test)](https://goreportcard.com/report/github.com/jenkins-x-plugins/jx-kube-test)
[![Releases](https://img.shields.io/github/release-pre/jenkins-x/jx-kube-test.svg)](https://github.com/jenkins-x-plugins/jx-kube-test/releases)
[![LICENSE](https://img.shields.io/github/license/jenkins-x/jx-kube-test.svg)](https://github.com/jenkins-x-plugins/jx-kube-test/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://slack.k8s.io/)

jx-kube-test is a command line tool for testing Kubernetes resources generated via helm, helmfile or kustomize.
  

## Configuration

You can configure the tests to run by creating a `.jx/kube-test/settings.yaml` file using the 

* [KubeTest Configuration Reference](docs/config.md#kubetest.jenkins-x.io/v1alpha1.KubeTest)
         

## Using in a GitOps repository

If you use a GitOps repository layout like the [Jenkins X GitOps Layout Conventions](https://github.com/jenkins-x-plugins/jx-gitops/blob/main/docs/git_layout.md) then you'll have a root folder like `config-root`  in which case you can perform the default kubernetes validation on your resources via:


```bash 
jx kube test run --source-dir config-root
```

To see the available arguments run `jx kube test run --help` or [browse the CLI reference](docs/cmd/jx-kube-test_run.md#options)

## Commands

See the [jx-kube-test command reference](docs/cmd/jx-kube-test.md#see-also)

