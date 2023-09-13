# CodeReady-Test-Harness
Testing solution written in golang using ginkgo framework for CodeReady Workspaces. This tests runs in Openshift CI Platform. 

# Specifications
* Instrumented tests with ginkgo framework. Find more info: https://onsi.github.io/ginkgo/
* Structured logging with logrus.
* Use client-go to connect to Openshift Cluster.
* Deploy CodeReady Workspaces in OSD Cluster or outside of OSD.
* Defined events watcher oriented to CodeReady Workspaces Resources. Please look `pkg/monitors/watcher.go`
* Create, start CodeReady Workspaces Workspaces
* Writes out an `addon-metadata.json` file which will also be consumed by the osde2e test framework.
* Writes out a junit XML file with tests results to the /test-run-results directory as expected
  by the [https://github.com/openshift/osde2e](osde2e) test framework.
* Check CodeReady Workspaces pods health
* Check all kubernetes objects created by CodeReady Workspaces installation
* In case of no osd tests, CodeReady Test Harness creates olm related objects for installation

# Setup

Log into your openshift cluster, using `oc login -u <user> -p <password> <oc_api_url>.`

A properly setup Go workspace using **Go 1.13+ is required**.

Install dependencies:
```
# Install dependencies
$ go mod tidy
# Copy the dependencies to vendor folder
$ go mod vendor
# Create che-operator-test-harness binary in bin folder. Please add the binary to the path or just execute ./bin/che-operator-test-harness
$ make build
```

## The `che-operator-test-harness` command

The `che-operator-test-harness` command is the root command that executes all test harness functionality through a number of variables

### CodeReady Test Harness Arguments

Che Test Harness comes with a number of arguments that can be passed to the `che-operator-test-harness` command. Supported arguments:

| Argument | Usage | Default |
| -- | -- | -- |
| `--help` | Prints all available arguments | "" |
| `--osd-provider` | Indicates if `test-harness` run in osd or not | `true` |
| `--namespace` | Indicate where to install and deploy Dev Spaces Operator. If 'osd-provider' is true this flag it is ignored . | `openshift-devspaces` |
| `--subscription-name` | Indicate the name of your subscription. If 'osd-provider' is true this flag it is ignored . | `devspaces-sub` |
| `--channel` | Indicate the channel for the subscription. If 'osd-provider' is true this flag it is ignored . | `latest` |
| `--source-ns` | Indicate namespace where catalog source it is installed. If 'osd-provider' is true this flag it is ignored | `openshift-marketplace` |
| `--catalog-name` | Indicate the name for the catalog source where you have exposed the bundles. If 'osd-provider' is true this flag it is ignored | `redhat-operators` |
| `--package-name` | Indicate the name of codeready package. If 'osd-provider' is true this flag it is ignored  | `devspaces` |
| `--csv-name` | Indicates csv version to install. If 'osd-provider' is true this flag it is ignored | `devspacesoperator.v3.1.0` |
Also `che-operator-test-harness` command support all ``Ginkgo`` flags...

If you plan to execute test harness outside of OSD clusters please consider to run inside of folder `scripts`  execute-test-harness.sh. Please check docs in the `scripts` folder.
# Interoperability Tesing CodeReady Workspaces on Openshift

* It is used `Interop Hosted Service` model. Find more info: https://docs.engineering.redhat.com/display/PIT/Hosted+Service+Workflow
* `Interop QE team` launches `execute-test-harness.sh` in `script` folder on unreleased platform builds against the last production version of CodeReady Workspaces
* When `Interop QE team` performs run test scenario, subscribers get emails with results. Find more info: https://docs.engineering.redhat.com/display/PIT/CodeReady+Workspaces+on+Openshift%3A+Hosted+Project+Plan
# Openshift CI

Che-Test-Harness run as a part of Openshift CI every 24 hours. To visualize the jobs please go to [PROW](https://prow.ci.openshift.org/?job=*osde2e*addon*).
Openshift CI Job Configuration lives in [ci-operator](https://github.com/openshift/release/blob/master/ci-operator/jobs/openshift/osde2e/openshift-osde2e-main-periodics.yaml#L1572).
