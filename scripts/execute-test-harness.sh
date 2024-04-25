#!/usr/bin/env bash

set -e

USER_NAME=$1

if [[ ${USER_NAME} == ""  ]]; then
  export USER_NAME="admin"
  echo "[INFO] Into a distinct namespace '<user-name>-devspaces' will run a test workspace.
By default will be used the 'admin' user-name You can specify an user-name of your OCP cluster 
as a parameter of this script, e.g.: 'execute-test-harness.sh <user-name>'"
fi

DEVSPACES_NAMESPACE="openshift-devspaces"
OPERATORS_NAMESPACE="openshift-operators"
USER_NAMESPACE="${USER_NAME}-devspaces"
REPORT_DIR="test-run-results"

oc create namespace ${DEVSPACES_NAMESPACE}
oc create namespace ${USER_NAMESPACE}
oc project ${USER_NAMESPACE}

ID=$(date +%s)
OPENSHIFT_API_URL=$(oc config view --minify -o jsonpath='{.clusters[*].cluster.server}')
OPENSHIFT_API_TOKEN=$(oc whoami -t)

echo "USER: $(oc whoami)"
echo "TOKEN: $(oc whoami -t)"
export OPENSHIFT_API_USER=$(oc whoami)

TMP_POD_YML=$(mktemp)
TMP_KUBECONFIG_YML=$(mktemp)

cat kubeconfig.template.yml |
    sed -e "s#__OPENSHIFT_API_URL__#${OPENSHIFT_API_URL}#g" |
    sed -e "s#__OPENSHIFT_API_TOKEN__#${OPENSHIFT_API_TOKEN}#g" |
    sed -e "s#__OPENSHIFT_API_USER__#${OPENSHIFT_API_USER}#g" |
    cat >${TMP_KUBECONFIG_YML}


oc delete configmap -n ${OPERATORS_NAMESPACE} ds-testsuite-kubeconfig || true

echo "[INFO] Creating configmap 'ds-testsuite-kubeconfig' in namespace '${OPERATORS_NAMESPACE}'"
oc create configmap -n ${OPERATORS_NAMESPACE} ds-testsuite-kubeconfig \
    --from-file=config=${TMP_KUBECONFIG_YML}

echo "[INFO] Dev Spaces will use the latest production version"
export DS_VERSION=$(oc get packagemanifest devspaces -o json | jq -r '.status.channels[] | select(.name == "stable") | .currentCSV')
echo "[INFO] The version of Dev Spaces: ${DS_VERSION}"


cat test-harness.pod.template.yml |
    sed -e "s#__ID__#${ID}#g" |
    sed -e "s#__OPERATORS_NS__#${OPERATORS_NAMESPACE}#g" |
    sed -e "s#__DEVSPACES_NS__#${DEVSPACES_NAMESPACE}#g" |
    sed -e "s#__USER_NS__#${USER_NAMESPACE}#g" |
    sed -e "s#__DEV_SPACES_VERSION__#${DS_VERSION}#g" |
    cat >${TMP_POD_YML}


# start the test
echo "[INFO] Creating pod 'ds-testsuite-${ID}' in namespace '${OPERATORS_NAMESPACE}'"
oc create -f ${TMP_POD_YML}

# wait for the pod to start
while true; do
    sleep 3
    PHASE=$(oc get pod -n ${OPERATORS_NAMESPACE} ds-testsuite-${ID} \
        --template='{{ .status.phase }}')
    if [[ ${PHASE} == "Running" ]]; then
        break
    fi
done

# wait for the test to finish
oc logs -n ${OPERATORS_NAMESPACE} ds-testsuite-${ID} -c test -f

# just to sleep
sleep 3

# download the test results
mkdir -p ${REPORT_DIR}/${ID}

oc rsync -n ${OPERATORS_NAMESPACE} \
    ds-testsuite-${ID}:/test-run-results ${REPORT_DIR}/${ID} -c download || true

oc exec -n ${OPERATORS_NAMESPACE} ds-testsuite-${ID} -c download \
    -- touch /tmp/done || true

echo "Retrieve test results"
oc cp ds-testsuite-${ID}:/test-run-results ${REPORT_DIR}
