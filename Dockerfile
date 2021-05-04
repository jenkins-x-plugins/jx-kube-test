FROM gcr.io/jenkinsxio/jx-cli-base:latest

ENTRYPOINT ["jx-kube-test"]

COPY ./build/linux/jx-kube-test /usr/bin/jx-kube-test