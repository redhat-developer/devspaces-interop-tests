FROM registry.ci.openshift.org/openshift/release:golang-1.19

SHELL ["/bin/bash", "-c"]

# Temporary workaround since mirror.centos.org is down and can be replaced with vault.centos.org
RUN sed -i s/mirror.centos.org/vault.centos.org/g /etc/yum.repos.d/*.repo && \
    sed -i s/^#.*baseurl=http/baseurl=http/g /etc/yum.repos.d/*.repo && \
    sed -i s/^mirrorlist=http/#mirrorlist=http/g /etc/yum.repos.d/*.repo

# Install yq, kubectl, oc tools.
RUN yum install --assumeyes -d1 python3-pip  httpd-tools
RUN pip3 install --upgrade setuptools
RUN pip3 install yq
    
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin

# Download the oc.tar.gz file
RUN curl -L https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest-4.12/openshift-client-linux.tar.gz -o /tmp/oc.tar.gz

# Extract the oc binary from oc.tar.gz to /usr/local/bin/
RUN tar xvzf /tmp/oc.tar.gz -C /usr/local/bin/

# Make the oc binary executable
RUN chmod +x /usr/local/bin/oc

# Clean up the temporary file
RUN rm /tmp/oc.tar.gz

# Copy the devspaces-interop-tests repo into /tmp/devspaces folder
RUN mkdir /tmp/devspaces
WORKDIR /tmp/devspaces
COPY . .

# Set required permissions for OpenShift usage
RUN mkdir -p /.config && \
    chgrp -R 0 /.config && \
    chmod -R g=u /.config && \
    chgrp -R 0 /tmp && \
    chmod -R g=u /tmp
