package operator_tests

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
	"github.com/redhat-developer/devspaces-interop-tests/internal/hlog"
	reporter "github.com/redhat-developer/devspaces-interop-tests/internal/reporters"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/api/github"
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
	OSDQeNamespace       = "codeready-workspaces-operator-qe"
	OSDCrwNamespace      = "codeready-workspaces-operator"
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
	// In case if --osd-provider=false DS will start to be installed from a specific catalog source
	if !testContext.Config.IS_OSD {
		github := github.NewGitubClient("redhat-developer", "devspaces-chectl")
		crwVersion, err := github.GetLatestCodeReadyWorkspacesTag()
		if err != nil {
			hlog.Log.Fatal("Failed to get version from github.", zap.Error(err))
		}
		/*if crwVersion == "" {
			hlog.Log.Fatal("Failed to get Dev Spaces version from github.", zap.String("crwVersion", crwVersion))
		}*/

		if testContext.Config.CSVName == "" {
			hlog.Log.Info("Flag `--csv-name` is not defined. Getting latest stable version of Dev Spaces from github...")
			testContext.Config.CSVName = "devspacesoperator.v" + crwVersion
		}

		/*if testContext.Config.CSVName != "devspacesoperator.v"+crwVersion {
			hlog.Log.Fatalf("Failed to define csv. You specify an old version. Latest Dev Spaces stable version is: %s", crwVersion)
		}*/
		hlog.Log.Infof("Using CSV: %s", testContext.Config.CSVName)

		olmClient, err := olm.NewOLMK8sClient()
		if err != nil {
			hlog.Log.Panic("Failed to create olm client go", zap.Error(err))
		}
		hlog.Log.Info("Test Harness will run outside of OSD. Start to create OLM Kubernetes Objects for Dev Spaces")

		controller := olm.NewOLMController(olmClient)
		controller.InstallOLMOperator()
	} else {
		hlog.Log.Info("Test Harness detect Dev Spaces in OSD cluster")
		// Check if Dev Spaces operator is installed on OSD namespace
		start := CheckOSDNamespace()
		if !start {
			// In case if Dev Spaces operator not found in any namespace specified the software will crush
			os.Exit(1)
		}
	}

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

// Check where devspaces operator it is installed
func CheckOSDNamespace() bool {
	hlog.Log.Info("Start to detect OSD namespace where operator it is installed...")
	OsdNamespaces := []string{OSDCrwNamespace, OSDQeNamespace}

	for _, namespace := range OsdNamespaces {
		_, err := GetNamespace(namespace)
		if err == nil {
			hlog.Log.Info("Dev Spaces operator detected on namespace: " + namespace)
			testContext.Config.OperatorsNamespace = namespace

			return true
		}
	}

	hlog.Log.Error("Error on start Dev Spaces Test Harness. Please check provided namespace")

	return false
}

func registerCheFlags(flags *flag.FlagSet) {
	flags.BoolVar(&testContext.Config.IS_OSD, "osd-provider", true, "Indicates if `test-harness` run in osd or not.")
	flags.StringVar(&testContext.Config.DevSpacesNamespace, "devspaces-namespace", "openshift-devspaces", "Indicate where to install and deploy Dev Spaces application. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.OperatorsNamespace, "operators-namespace", "openshift-operators", "Indicate where to install Dev Spaces Operator. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.UserNamespace, "user-namespace", "admin-devspaces", "Indicate where to run a test workspace.")
	flags.StringVar(&testContext.Config.SubscriptionName, "subscription-name", "devspaces", "Indicate the name of your subscription. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.OLMChannel, "channel", "stable", "Indicate the channel for the subscription. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.SourceNamespace, "source-ns", "openshift-marketplace", "Indicate namespace where catalog source it is installed. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.CatalogSourceName, "catalog-name", "redhat-operators", "Indicate the name for the catalog source where you have exposed the bundles. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.OLMPackage, "package-name", "devspaces", "Indicate the name of Dev Spaces package. If 'osd-provider' is true this flag it is ignored .")
	flags.StringVar(&testContext.Config.CSVName, "csv-name", "devspacesoperator.v3.1.0", "Indicates csv version to install. If 'osd-provider' is true this flag it is ignored .")

	if testContext.Config.CSVName == "" {
		hlog.Log.Panic("Please specify csv name in order to install Dev Spaces via olm. Eg. 'che-operator-test-harness --osd-provider=false --csv-name=devspacesoperator.v3.1.0'")
	}
}
