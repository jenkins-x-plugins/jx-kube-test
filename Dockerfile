FROM ghcr.io/jenkins-x/jx-boot:latest

ENTRYPOINT ["jx-kube-test"]

COPY ./build/linux/jx-kube-test /usr/bin/jx-kube-test