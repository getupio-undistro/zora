package discovery

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

func NewForConfig(c *rest.Config) (ClusterDiscoverer, error) {
	kclient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	mclient, err := versioned.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	return &clusterDiscovery{kubernetes: kclient, metrics: mclient}, nil
}

func NewResources(available, usage resource.Quantity) Resources {
	fraction := float64(usage.MilliValue()) / float64(available.MilliValue()) * 100
	return Resources{Available: available, Usage: usage, UsagePercentage: int32(fraction)}
}

type clusterDiscovery struct {
	kubernetes *kubernetes.Clientset
	metrics    *versioned.Clientset
}

func (r *clusterDiscovery) Discover(ctx context.Context) (*ClusterInfo, error) {
	nodes, err := r.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, errors.New("cluster has no nodes")
	}

	prov := r.Provider(nodes[0])
	reg, err := r.Region(nodes)
	if err != nil {
		return nil, err
	}

	return &ClusterInfo{
		Nodes:             nodes,
		Resources:         avgNodeResources(nodes),
		CreationTimestamp: oldestNodeTimestamp(nodes),
		Provider:          prov,
		Region:            reg,
	}, nil
}

// Provider finds the cluster source by matching against provider specific
// labels on a node, returning the provider if the match succeeds and
// "unknown" if it fails.
func (r *clusterDiscovery) Provider(node NodeInfo) string {
	for l := range node.Labels {
		for pref, p := range ClusterSourcePrefixes {
			if strings.HasPrefix(l, pref) {
				return p
			}
		}
	}
	return "unknown"
}

// Region returns "multi-region" if the cluster nodes belong to distinct
// locations, otherwise it returns the region itself.
func (r *clusterDiscovery) Region(nodes []NodeInfo) (string, error) {
	regs := map[string]bool{}
	haslabel := false
	for c := 0; c < len(nodes); c++ {
		for l, v := range nodes[c].Labels {
			if l == RegionLabel {
				regs[v] = true
				if haslabel && len(regs) > 1 {
					return "multi-region", nil
				} else {
					haslabel = true
				}
			}
		}
	}
	if !haslabel {
		return "", fmt.Errorf("unable to discover region: %w",
			fmt.Errorf("no node has the label <%s>", RegionLabel))
	}
	reg := ""
	for reg = range regs {
		continue
	}
	return reg, nil
}

func (r *clusterDiscovery) Version() (string, error) {
	v, err := r.kubernetes.Discovery().ServerVersion()
	if err != nil {
		return "", fmt.Errorf("failed to discover server version: %w", err)
	}
	return v.String(), nil
}

func (r *clusterDiscovery) Nodes(ctx context.Context) ([]NodeInfo, error) {
	if err := r.checkMetricsAPI(); err != nil {
		return nil, err
	}
	metricsList, err := r.metrics.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list NodeMetrics: %w", err)
	}
	nodeList, err := r.kubernetes.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Nodes: %w", err)
	}

	return nodeResources(nodeList.Items, metricsList.Items), nil
}

func (r *clusterDiscovery) checkMetricsAPI() error {
	apiGroups, err := r.kubernetes.Discovery().ServerGroups()
	if err != nil {
		return err
	}
	for _, group := range apiGroups.Groups {
		if group.Name != v1beta1.GroupName {
			continue
		}
		for _, version := range group.Versions {
			if version.Version == v1beta1.SchemeGroupVersion.Version {
				return nil
			}
		}
	}
	return errors.New("metrics API not available")
}

func nodeResources(nodes []corev1.Node, nodeMetrics []v1beta1.NodeMetrics) []NodeInfo {
	infos := make([]NodeInfo, 0, len(nodes))
	metrics := make(map[string]corev1.ResourceList)
	for _, m := range nodeMetrics {
		metrics[m.Name] = m.Usage
	}
	for _, n := range nodes {
		usage := metrics[n.Name]
		info := NodeInfo{
			Name:              n.Name,
			Labels:            n.Labels,
			Resources:         make(map[corev1.ResourceName]Resources),
			Ready:             nodeIsReady(n),
			CreationTimestamp: n.CreationTimestamp,
		}
		for _, res := range MeasuredResources {
			info.Resources[res] = NewResources(n.Status.Allocatable[res], usage[res])
		}
		infos = append(infos, info)
	}
	return infos
}

func avgNodeResources(nodes []NodeInfo) map[corev1.ResourceName]Resources {
	totalAvailable := make(map[corev1.ResourceName]*resource.Quantity)
	totalUsage := make(map[corev1.ResourceName]*resource.Quantity)

	for _, node := range nodes {
		for _, res := range MeasuredResources {
			if r, found := node.Resources[res]; found {
				if _, ok := totalAvailable[res]; ok {
					totalAvailable[res].Add(r.Available)
					totalUsage[res].Add(r.Usage)
				} else {
					totalAvailable[res] = &r.Available
					totalUsage[res] = &r.Usage
				}
			}
		}
	}
	result := make(map[corev1.ResourceName]Resources)
	for _, res := range MeasuredResources {
		result[res] = NewResources(*totalAvailable[res], *totalUsage[res])
	}
	return result
}

func nodeIsReady(node corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func oldestNodeTimestamp(nodes []NodeInfo) metav1.Time {
	oldest := metav1.NewTime(time.Now().UTC())
	for _, node := range nodes {
		if node.CreationTimestamp.Before(&oldest) {
			oldest = node.CreationTimestamp
		}
	}
	return oldest
}
