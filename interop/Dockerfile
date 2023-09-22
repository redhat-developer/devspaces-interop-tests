FROM registry.ci.openshift.org/openshift/release:golang-1.18

SHELL ["/bin/bash", "-c"]

# Install yq, kubectl, oc tools.
RUN yum install --assumeyes -d1 python3-pip  httpd-tools && \
    pip3 install --upgrade setuptools && \
    pip3 install yq && \
    curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin && \
    curl https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest-4.12/openshift-client-linux.tar.gz | tar xvzf - -C /usr/local/bin/ oc && \
    chmod ug+x /usr/local/bin/oc

# Set required permissions for OpenShift usage
RUN mkdir -p /.config && \
    chgrp -R 0 /.config && \
    chmod -R g=u /.config && \
    chgrp -R 0 /tmp && \
    chmod -R g=u /tmp
