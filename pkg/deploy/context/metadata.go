package context

import (
	"encoding/json"
	"io/ioutil"
	"os"

	v1 "k8s.io/api/core/v1"
)

// metadata houses metadata to be written out to the additional-metadata.json
type metadata struct {
	// Whether the CRD was found. Typically Spyglass seems to have issues displaying non-strings, so
	// this will be written out as a string despite the native JSON boolean type.
	FoundCRD bool `json:"found-crd,string"`

	// ClusterTimeUp return how many time take Dev Spaces to be up
	ClusterTimeUp float64 `json:"cluster-time-up,int"`

	// DevSpacesPodsInfo returns all Dev Spaces containers infos
	CodeReadyPodsInfo []CodeReadyPods `json:"pods-info,string"`

	//CRWPodTime returns all Dev Spaces up time from pulled image to start container
	CRWPodTime PodTimes `json:"pods-up-times int"`

	//DevSpacesServerIsUp Returns true or false depending if code ready server is UP
	DevSpacesServerIsUp bool `json:"devspaces_apiserver_is_up,bool"`
}

type PodTimes struct {
	// StartUp time of dashboard pod
	DashboardUpTime float64 `json:"devspaces-dashboard-up-time,float64"`

	// StartUp time of devfile pod
	DevFileUpTime float64 `json:"devfile-up-time,float64"`

	// StartUp time of plugin-registry pod
	PluginRegUpTime float64 `json:"plugins-up-time,float64"`

	// StartUp time of devspaces server pod
	CodeReadyUpTime float64 `json:"devspaces-up-time,float64"`
}

type CodeReadyPods struct {
	// Name show the name of the pod
	Name string `json:"name,string"`

	// DockerImage show the image used by a container
	DockerImage string `json:"docker_image,string"`

	// Indicate the status of a pod
	Status v1.PodPhase `json:"status,string"`

	// Labels defined to all Codeready Pods
	Labels map[string]string `json:"labels,string"`
}

type TestConfig struct {
	// Namespace where to install Dev Spaces and DevWorkspace operators
	OperatorsNamespace string

	// Namespace where to install Dev Spaces components
	DevSpacesNamespace string

	// Namespace where to run a test workspace
	UserNamespace string

	// Subscription name metadata for crw installation
	SubscriptionName string

	// Indicate the channel from where to install Dev Spaces
	OLMChannel string

	// Package name of Dev Spaces. By default it is devspaces
	OLMPackage string

	// Catalog Source name where Dev Spaces bundles are
	CatalogSourceName string

	// CSV name it is the version of Dev Spaces to install
	CSVName string

	// Source where catalog source is
	SourceNamespace string

	// Indicate if the tests are osd or not
	IS_OSD bool

	// check if test harness are up and working
	UP bool
}

// Config instance
var Config = TestConfig{}

// Metadata instance
var Instance = metadata{}

// WriteToJSON will marshall the metadata struct and write it into the given file.
func (m *metadata) WriteToJSON(outputFilename string) (err error) {
	var data []byte
	if data, err = json.Marshal(m); err != nil {
		return err
	}

	if err = ioutil.WriteFile(outputFilename, data, os.FileMode(0644)); err != nil {
		return err
	}

	return nil
}
