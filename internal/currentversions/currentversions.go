package currentversions

import (
	"context"

	"github.com/catenax-ng/maintenance-dashboard/internal/data"
	"github.com/catenax-ng/maintenance-dashboard/internal/helpers"
	"github.com/catenax-ng/maintenance-dashboard/internal/parseversion"
	log "github.com/sirupsen/logrus"
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
	for _, deployment := range deployments.Items {
		versionLabel := deployment.ObjectMeta.Labels["app.kubernetes.io/version"]
		semverVersion, err := parseversion.ToSemver(versionLabel)
		if err != nil {
			log.Warnf("Skipping invalid version: %v", versionLabel)
		} else {
			result = append(result, &data.AppVersionInfo{
				CurrentVersion:  semverVersion,
				NewReleasesName: deployment.ObjectMeta.Annotations["maintenance/releasename"],
				ResourceName:    deployment.ObjectMeta.Name,
			})
		}
	}
	if err != nil {
		log.Panicf("Unable to get deployments to scan: %v", err.Error())
	}

	statefulsets, err := clientSet.AppsV1().StatefulSets("").List(ctx, listOptions)
	for _, statefulset := range statefulsets.Items {
		versionLabel := statefulset.ObjectMeta.Labels["app.kubernetes.io/version"]
		semverVersion, err := parseversion.ToSemver(versionLabel)
		if err != nil {
			log.Warnf("Skipping invalid version: %v", versionLabel)
		} else {
			result = append(result, &data.AppVersionInfo{
				CurrentVersion:  semverVersion,
				NewReleasesName: statefulset.ObjectMeta.Annotations["maintenance/releasename"],
				ResourceName:    statefulset.ObjectMeta.Name,
			})
		}
	}
	if err != nil {
		log.Panicf("Unable to get statefulsets to scan: %v", err.Error())
	}

	daemonsets, err := clientSet.AppsV1().DaemonSets("").List(ctx, listOptions)
	for _, daemonset := range daemonsets.Items {
		versionLabel := daemonset.ObjectMeta.Labels["app.kubernetes.io/version"]
		semverVersion, err := parseversion.ToSemver(versionLabel)
		if err != nil {
			log.Warnf("Skipping invalid version: %v", versionLabel)
		} else {
			result = append(result, &data.AppVersionInfo{
				CurrentVersion:  semverVersion,
				NewReleasesName: daemonset.ObjectMeta.Annotations["maintenance/releasename"],
				ResourceName:    daemonset.ObjectMeta.Name,
			})
		}
	}
	if err != nil {
		log.Panicf("Unable to get daemonsets to scan: %v", err.Error())
	}

	log.Infof("Found %v apps to scan.", len(result))
	return result
}
