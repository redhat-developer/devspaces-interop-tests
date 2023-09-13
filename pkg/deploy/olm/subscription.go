package olm

import (
	"context"
	"time"

	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	testContext "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy/context"

	v1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	gherr "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// InstallSubscription desc
func (o *OlmController) InstallSubscription() (err error) {
	hlog.Log.Infof("Creating new subscription '%s' in namespace '%s'", testContext.Config.SubscriptionName, testContext.Config.OperatorsNamespace)

	if testContext.Config.CSVName == "" {
		hlog.Log.Panic("Please specify csv name in order to install Dev Spaces via olm. Eg. 'che-operator-test-harness --osd-provider=false --csv-name=devspacesoperator.v3.1.0'")
	}

	if err := o.k8s.Create(context.TODO(), o.GetSubscriptionSpec()); err != nil {
		if errors.IsAlreadyExists(err) {
			return gherr.Wrapf(err, "Subscription '%s' already exist in namespace '%s'. Please remove it and try again ", testContext.Config.SubscriptionName, testContext.Config.OperatorsNamespace)
		} else {
			return gherr.Wrapf(err, "Failed to create subscription %v in namespace '%s'", err, testContext.Config.OperatorsNamespace)
		}
	}

	hlog.Log.Infof("Waiting subscription '%s' to be installed in namespace '%s'", testContext.Config.SubscriptionName, testContext.Config.OperatorsNamespace)

	if sub, err := o.WaitForSubscriptionState(IsSubscriptionInstalledCSVPresent); err != nil {
		return gherr.Wrapf(err, "Error to install subscription '%s' in namespace '%s", sub.Name, testContext.Config.OperatorsNamespace)
	}

	// Check if all csv compnents are created without errors
	// !Timeout is 5 minutes
	hlog.Log.Info("Waiting for CSV to be installed in cluster")
	o.WaitForClusterServiceVersionState(v1alpha1.CSVPhaseSucceeded)

	return nil
}

// GetSubscriptionSpec desc
func (o *OlmController) GetSubscriptionSpec() *v1alpha1.Subscription {
	return &v1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testContext.Config.SubscriptionName,
			Namespace: testContext.Config.OperatorsNamespace,
		},
		Spec: &v1alpha1.SubscriptionSpec{
			CatalogSource:          testContext.Config.CatalogSourceName,
			Package:                testContext.Config.OLMPackage,
			CatalogSourceNamespace: testContext.Config.SourceNamespace,
			Channel:                testContext.Config.OLMChannel,
			StartingCSV:            testContext.Config.CSVName,
			InstallPlanApproval:    "Automatic",
		},
		Status: v1alpha1.SubscriptionStatus{},
	}
}

// WaitForSubscriptionState desc
func (o *OlmController) WaitForSubscriptionState(inState func(s *v1alpha1.Subscription, err error) (bool, error)) (*v1alpha1.Subscription, error) {
	waitErr := wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		lastState, err := o.k8s.OLM.OperatorsV1alpha1().
			Subscriptions(testContext.Config.OperatorsNamespace).
			Get(context.TODO(), testContext.Config.SubscriptionName, metav1.GetOptions{})
		return inState(lastState, err)
	})

	return nil, waitErr
}

// IsSubscriptionInstalledCSVPresent desc
func IsSubscriptionInstalledCSVPresent(s *v1alpha1.Subscription, err error) (bool, error) {
	return s.Status.InstalledCSV != "" && s.Status.InstalledCSV != "<none>", err
}

// WaitForClusterServiceVersionState desc
func (o *OlmController) WaitForClusterServiceVersionState(status v1alpha1.ClusterServiceVersionPhase) (err error) {
	timeout := time.After(5 * time.Minute)
	tick := time.Tick(1 * time.Second)
	var csvInstallCompleted = time.Now()

	stopCh := make(chan struct{})
	defer close(stopCh)

	for {
		select {
		case <-timeout:
			return gherr.New("Error. CSV didn't install completed in 5 mins")
		case <-tick:
			csv, _ := o.GetClusterServiceVersion()
			if csv.Status.Phase == status {
				testContext.Instance.ClusterTimeUp = time.Since(csvInstallCompleted).Seconds()
				hlog.Log.Info("Successfully install CSV completed in ", testContext.Instance.ClusterTimeUp)

				return nil
			}
		}
	}
}
