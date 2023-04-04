package currentversions

import (
	"context"
	"errors"
	"strings"

	"github.com/catenax-ng/maintenance-dashboard/internal/data"
	"github.com/catenax-ng/maintenance-dashboard/internal/helpers"
	"github.com/catenax-ng/maintenance-dashboard/internal/parseversion"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Returns the versions of nodes and labeled services
func GetCurrentVersions(ctx context.Context) []*data.AppVersionInfo {
	clientSet := newClientSet()
	var result []*data.AppVersionInfo

	log.Infoln("Getting version info about nodes.")
	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Errorf("Unable to get node info: %v", err)
	}

	log.Infof("Found %v nodes to scan.", len(nodes.Items))
	for _, node := range nodes.Items {
		semverVersion, err := parseversion.ToSemver(node.Status.NodeInfo.KubeletVersion)
		if err != nil {
			log.Warnf("Skipping invalid version: %v", node.Status.NodeInfo.KubeletVersion)
		} else {
			result = append(result, &data.AppVersionInfo{
				CurrentVersion:  semverVersion,
				NewReleasesName: "kubernetes/kubernetes",
				ResourceName:    node.Name,
			})
		}
	}

	log.Infoln("Getting version info about services.")
	apps := getAppsToScan(ctx, clientSet)
	result = append(result, apps...)

	log.Infoln("Resources in the cluster to be scanned with their current version:")
	for _, res := range result {
		log.Infof("%v: %v", res.NewReleasesName, res.CurrentVersion.String())
	}
	return result
}

// Initializes new ClientSet either based on kubeconfig or in-cluster
func newClientSet() *kubernetes.Clientset {
	kubeconfig := helpers.GetEnv("KUBE_CONFIG", "")
	incluster := helpers.GetEnv("IN_CLUSTER", "false")
	if incluster != "true" {

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Panic(err.Error())
		}
		cs, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Panic(err.Error())
		}

		return cs
	}

	// Creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic(err.Error())
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	return cs
}

// Get services annotated with maintenance/scan=true
func getAppsToScan(ctx context.Context, clientSet *kubernetes.Clientset) []*data.AppVersionInfo {
	var result []*data.AppVersionInfo

	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"maintenance/scan": "true"}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}

	deployments, err := clientSet.AppsV1().Deployments("").List(ctx, listOptions)
	if err != nil {
		log.Panicf("Unable to get deployments to scan: %v", err.Error())
	} else {
		for _, deployment := range deployments.Items {
			appVersionInfo, err := createAppVersionInfo(deployment.ObjectMeta, deployment.Spec.Template)
			if err != nil {

			} else {
				result = append(result, appVersionInfo)
			}
		}
	}

	statefulsets, err := clientSet.AppsV1().StatefulSets("").List(ctx, listOptions)
	if err != nil {
		log.Panicf("Unable to get statefulsets to scan: %v", err.Error())
	} else {
		for _, statefulset := range statefulsets.Items {
			appVersionInfo, err := createAppVersionInfo(statefulset.ObjectMeta, statefulset.Spec.Template)
			if err != nil {

			} else {
				result = append(result, appVersionInfo)
			}
		}
	}

	daemonsets, err := clientSet.AppsV1().DaemonSets("").List(ctx, listOptions)
	if err != nil {
		log.Panicf("Unable to get daemonsets to scan: %v", err.Error())
	} else {
		for _, daemonset := range daemonsets.Items {
			appVersionInfo, err := createAppVersionInfo(daemonset.ObjectMeta, daemonset.Spec.Template)
			if err != nil {

			} else {
				result = append(result, appVersionInfo)
			}
		}
	}

	log.Infof("Found %v apps to scan.", len(result))
	return result
}

// Parse image tag and return a data.AppVersionInfo struct. If version can't be parsed, returns error instead.
func createAppVersionInfo(objectMeta metav1.ObjectMeta, podTemplate v1.PodTemplateSpec) (*data.AppVersionInfo, error) {
	versionLabel, keyExists := objectMeta.Labels["app.kubernetes.io/version"]

	if !keyExists {
		image := podTemplate.Spec.Containers[0].Image
		versionLabel = strings.Split(image, ":")[1]
		log.Warnf("Label app.kubernetes.io/version missing for app %v. Using the tag of the %v image instead.", objectMeta.Name, image)
	}

	semverVersion, err := parseversion.ToSemver(versionLabel)

	if err == nil {
		return &data.AppVersionInfo{
			CurrentVersion:  semverVersion,
			NewReleasesName: objectMeta.Annotations["maintenance/releasename"],
			ResourceName:    objectMeta.Name,
		}, nil
	}
	return &data.AppVersionInfo{}, errors.New("Unable to parse version label " + versionLabel)
}
