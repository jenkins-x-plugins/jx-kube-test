### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-kube-test/releases/download/v{{.Version}}/jx-kube-test-linux-amd64.tar.gz | tar xzv 
sudo mv jx-kube-test /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-kube-test/releases/download/v{{.Version}}/jx-kube-test-darwin-amd64.tar.gz | tar xzv
sudo mv jx-kube-test /usr/local/bin
```

