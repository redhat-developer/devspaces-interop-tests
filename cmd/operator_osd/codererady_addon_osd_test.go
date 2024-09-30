package operator_tests

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
	"github.com/redhat-developer/devspaces-interop-tests/internal/hlog"
	reporter "github.com/redhat-developer/devspaces-interop-tests/internal/reporters"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/client"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/deploy"
	testContext "github.com/redhat-developer/devspaces-interop-tests/pkg/deploy/context"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/deploy/olm"
	_ "github.com/redhat-developer/devspaces-interop-tests/tests"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Create Constant file
const (
	testResultsDirectory = "/test-run-results"
	jUnitOutputFilename  = "junit-dev-spaces.xml"
	addonMetadataName    = "addon-metadata.json"
	DebugSummaryOutput   = "debug_tests.json"
	RedHatDeveloperGHOrg = "redhat-developer"
	DSCtlRepoName        = "devspaces-chectl"
)

// Start to register flags
func init() {
	registerCheFlags(flag.CommandLine)
}

// SynchronizedBeforeSuite blocks are primarily meant to solve the problem of setting up the custom resources for
// Dev Spaces
var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	// Initialize Dev Spaces Kubernetes client to create resources in a giving namespace
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}

	hlog.Log.Info("Installing Dev Spaces objects before running Test Harness...")

	if testContext.Config.CSVName == "" {
		hlog.Log.Fatal("Failed to get Dev Spaces version from packagemanifest 'devspaces'.", zap.String("csv-name", testContext.Config.CSVName))
	}

	hlog.Log.Infof("Using CSV: %s", testContext.Config.CSVName)

	olmClient, err := olm.NewOLMK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create olm client go", zap.Error(err))
	}
	hlog.Log.Info("Start to create OLM Kubernetes Objects for Dev Spaces")

	controller := olm.NewOLMController(olmClient)
	controller.InstallOLMOperator()

	// Initialize Dev Spaces Kubernetes client to create resources in a giving namespace
	deploy := deploy.NewTestHarnessController(k8sClient)

	if !deploy.DeployDevSpaces() {
		hlog.Log.Panic("Failed to deploy Dev Spaces", zap.Error(err))
	}

	return nil
}, func(data []byte) {})

var _ = ginkgo.SynchronizedAfterSuite(func() {
	// Initialize Dev Spaces Kubernetes client to create resources in a giving namespace
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}
	//Delete all objects after pass all test suites.
	hlog.Log.Info("Clean up all created objects by Test Harness.")
	// Initialize Dev Spaces Kubernetes client to create resources in a giving namespace
	deploy := deploy.NewTestHarnessController(k8sClient)

	if err := deploy.DeleteCustomResource(); err != nil {
		hlog.Log.Panic("Failed to delete custom resources in cluster")
	}

}, func() {})

func TestHarnessCodeReadyWorkspaces(t *testing.T) {
	// configure zap logging for Dev Spaces addon, Zap Logger create a file <*.log> where is possible
	//to find information about addon execution.
	gomega.RegisterFailHandler(ginkgo.Fail)

	var r []ginkgo.Reporter
	r = append(r, reporters.NewJUnitReporter(filepath.Join(testResultsDirectory, jUnitOutputFilename)))
	r = append(r, reporter.NewDetailsReporterFile(filepath.Join(testResultsDirectory, DebugSummaryOutput)))

	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "Dev Spaces Operator Test Harness", r)

	err := testContext.Instance.WriteToJSON(filepath.Join(testResultsDirectory, addonMetadataName))
	if err != nil {
		hlog.Log.Panic("error while writing metadata")
	}
}

// Get namespace from cluster with given name
func GetNamespace(namespace string) (*v1.Namespace, error) {
	// Initialize Dev Spaces Kubernetes client to create resources in a giving namespace
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}

	return k8sClient.Kube().CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
}

func registerCheFlags(flags *flag.FlagSet) {
	flags.StringVar(&testContext.Config.DevSpacesNamespace, "devspaces-namespace", "openshift-devspaces", "Indicate where to install and deploy Dev Spaces application. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.OperatorsNamespace, "operators-namespace", "openshift-operators", "Indicate where to install Dev Spaces Operator. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.UserNamespace, "user-namespace", "admin-devspaces", "Indicate where to run a test workspace.")
	flags.StringVar(&testContext.Config.SubscriptionName, "subscription-name", "devspaces", "Indicate the name of your subscription. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.OLMChannel, "channel", "stable", "Indicate the channel for the subscription. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.SourceNamespace, "source-ns", "openshift-marketplace", "Indicate namespace where catalog source it is installed. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.CatalogSourceName, "catalog-name", "redhat-operators", "Indicate the name for the catalog source where you have exposed the bundles. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.OLMPackage, "package-name", "devspaces", "Indicate the name of Dev Spaces package. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.CSVName, "csv-name", "devspacesoperator.v3.16.0", "Indicates csv version to install. If 'osd-provider' is true this flag it is ignored .")

	if testContext.Config.CSVName == "" {
		hlog.Log.Panic("Please specify csv name in order to install Dev Spaces via olm. Eg. --csv-name=devspacesoperator.v3.16.1'")
	}
}
