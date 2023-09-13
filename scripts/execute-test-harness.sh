#!/usr/bin/env bash

set -e

USER_NAME=$1
DS_VERSION=$2
DEVSPACES_NAMESPACE="openshift-devspaces"
OPERATORS_NAMESPACE="openshift-operators"
USER_NAMESPACE="${USER_NAME}-devspaces"
REPORT_DIR="ds-interop-report"

if [[ ${USER_NAME} == ""  ]]; then
  echo "Please specify user-name of your OCP cluster as the first argument to run test harness."
  echo "into a distinct namespace '<user-name>-devspaces' where will run our test workspace"
  echo "execute-test-harness.sh <user-name> !!ds-version"
  exit 1
fi

# Ensure there are no already existed projects
oc delete namespace ${DEVSPACES_NAMESPACE} --wait=true --ignore-not-found
oc delete namespace ${USER_NAMESPACE} --wait=true --ignore-not-found

oc create namespace ${DEVSPACES_NAMESPACE}
oc create namespace ${USER_NAMESPACE}
oc project ${USER_NAMESPACE}

ID=$(date +%s)
OPENSHIFT_API_URL=$(oc config view --minify -o jsonpath='{.clusters[*].cluster.server}')
OPENSHIFT_API_TOKEN=$(oc whoami -t)

TMP_POD_YML=$(mktemp)
TMP_KUBECONFIG_YML=$(mktemp)

cat kubeconfig.template.yml |
    sed -e "s#__OPENSHIFT_API_URL__#${OPENSHIFT_API_URL}#g" |
    sed -e "s#__OPENSHIFT_API_TOKEN__#${OPENSHIFT_API_TOKEN}#g" |
    cat >${TMP_KUBECONFIG_YML}

cat ${TMP_KUBECONFIG_YML}

oc delete configmap -n ${OPERATORS_NAMESPACE} ds-testsuite-kubeconfig || true
oc create configmap -n ${OPERATORS_NAMESPACE} ds-testsuite-kubeconfig \
    --from-file=config=${TMP_KUBECONFIG_YML}

if [[ "${DS_VERSION}" == "" ]]; then
    export DS_VERSION=$(oc get packagemanifest devspaces -o json | jq -r '.status.channels[] | select(.name == "stable") | .currentCSV')
fi

cat test-harness.pod.template.yml |
    sed -e "s#__ID__#${ID}#g" |
    sed -e "s#__OPERATORS_NS__#${OPERATORS_NAMESPACE}#g" |
    sed -e "s#__DEVSPACES_NS__#${DEVSPACES_NAMESPACE}#g" |
    sed -e "s#__USER_NS__#${USER_NAMESPACE}#g" |
    sed -e "s#__DEV_SPACES_VERSION__#${DS_VERSION}#g" |
    cat >${TMP_POD_YML}

cat ${TMP_POD_YML}

# start the test
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
    ds-testsuite-${ID}:/test-run-results ${REPORT_DIR}/${ID} -c download

oc exec -n ${OPERATORS_NAMESPACE} ds-testsuite-${ID} -c download \
    -- touch /tmp/done
