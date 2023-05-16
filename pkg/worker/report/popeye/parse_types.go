// Copyright 2022 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package popeye

import zorav1a1 "github.com/undistro/zora/api/zora/v1alpha1"

var (
	// LevelToIssueSeverity maps Popeye's <Level> type to Zora's
	// <ClusterIssueSeverity>.
	LevelToIssueSeverity = [4]zorav1a1.ClusterIssueSeverity{
		zorav1a1.SeverityUnknown,
		zorav1a1.SeverityLow,
		zorav1a1.SeverityMedium,
		zorav1a1.SeverityHigh,
	}

	// IssueIDtoGenericMsg maps Popeye's issue codes to generic versions of the
	// issue description. The original issues can be found on Popeye's source
	// file <internal/issues/assets/codes.yml>.
	IssueIDtoGenericMsg = map[string]string{
		// Container
		"POP-105": "Unnamed probe port in use",
		"POP-108": "Unnamed port",
		"POP-109": "CPU reached request threshold",
		"POP-110": "Memory reached request threshold",
		"POP-111": "CPU reached user threshold",
		"POP-112": "Memory reached user threshold",
		"POP-113": "Container image not hosted on an allowed docker registry",

		// Pod
		"POP-200": "Pod is terminating",
		"POP-201": "Pod is terminating a process",
		"POP-202": "Pod is waiting",
		"POP-203": "Pod is waiting a process",
		"POP-204": "Pod is not ready",
		"POP-205": "Pod was restarted",
		"POP-207": "Pod is in an unhappy phase",

		// Security
		"POP-304": "ServiceAccount references a secret which does not exist",
		"POP-305": "ServiceAccount references a docker-image pull secret which does not exist",

		// General
		"POP-401": "Unable to locate key reference",
		"POP-402": "No metrics-server detected",
		"POP-403": "Deprecated API group",
		"POP-404": "Deprecation check failed",

		// Deployment and StatefulSet
		"POP-501": "Unhealthy, mismatch between desired and available state",
		"POP-503": "At current load, CPU under allocated",
		"POP-504": "At current load, CPU over allocated",
		"POP-505": "At current load, Memory under allocated",
		"POP-506": "At current load, Memory over allocated",
		"POP-507": "Deployment references ServiceAccount which does not exist",

		// HPA
		"POP-600": "HPA references a Deployment which does not exist",
		"POP-601": "HPA references a StatefulSet which does not exist",
		"POP-602": "Replicas at burst will match or exceed cluster CPU capacity",
		"POP-603": "Replicas at burst will match or exceed cluster memory capacity",
		"POP-604": "If ALL HPAs are triggered, cluster CPU capacity will match or exceed threshold",
		"POP-605": "If ALL HPAs are triggered, cluster memory capacity will match or exceed threshold",

		// Node
		"POP-700": "Found taint that no pod can tolerate",
		"POP-704": "Insufficient memory on Node (MemoryPressure condition)",
		"POP-705": "Insufficient disk space on Node (DiskPressure condition)",
		"POP-706": "Insufficient PIDs on Node (PIDPressure condition)",
		"POP-707": "No network configured on Node (NetworkUnavailable condition)",
		"POP-709": "Node CPU threshold reached",
		"POP-710": "Node Memory threshold reached",

		// PodDisruptionBudget
		"POP-901": "MinAvailable is greater than the number of pods currently running",

		// Service
		"POP-1101": "Skip ports check. No explicit ports detected on pod",
		"POP-1102": "Unnamed service port in use",
		"POP-1106": "No target ports match service port",

		// ReplicaSet
		"POP-1120": "Unhealthy ReplicaSet",

		// NetworkPolicies
		"POP-1200": "No pods match pod selector",
		"POP-1201": "No namespaces match namespace selector",

		// RBAC
		"POP-1300": "References a role which does not exist",
	}

	// IssueIDtoUrl maps Popeye's issue codes to urls for wiki pages, blog
	// posts and other sources documenting the issue.
	IssueIDtoUrl = map[string]string{
		// Container
		"POP-100": "https://kubernetes.io/docs/concepts/containers/images/#image-names",
		"POP-101": "https://kubernetes.io/docs/concepts/containers/images/#image-names",
		"POP-102": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/",
		"POP-103": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/",
		"POP-104": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-readiness-probes",
		"POP-105": "",
		"POP-106": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-107": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-108": "",
		"POP-109": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-110": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-111": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-112": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-113": "",

		// Pod
		"POP-200": "https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle",
		"POP-201": "https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle",
		"POP-202": "https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle",
		"POP-203": "https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle",
		"POP-204": "https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle",
		"POP-205": "https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle",
		"POP-206": "https://kubernetes.io/docs/concepts/workloads/pods/disruptions",
		"POP-207": "https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle",
		"POP-208": "https://kubernetes.io/docs/concepts/configuration/overview/#naked-pods-vs-replicasets-deployments-and-jobs",

		// Security
		"POP-300": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/",
		"POP-301": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/",
		"POP-302": "https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted",
		"POP-303": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/",
		"POP-304": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account",
		"POP-305": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account",
		"POP-306": "https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted",

		// General
		"POP-400": "",
		"POP-401": "",
		"POP-402": "https://kubernetes.io/docs/tasks/debug/debug-cluster/resource-metrics-pipeline/#metrics-server",
		"POP-403": "https://kubernetes.io/docs/reference/using-api/deprecation-guide",
		"POP-404": "https://kubernetes.io/docs/reference/using-api/deprecation-guide",
		"POP-405": "https://kubernetes.io/docs/tasks/administer-cluster/cluster-upgrade/",
		"POP-406": "https://kubernetes.io/releases/",

		// Deployment and StatefulSet
		"POP-500": "https://kubernetes.io/docs/concepts/workloads/",
		"POP-501": "https://kubernetes.io/docs/concepts/workloads/",
		"POP-503": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-504": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-505": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-506": "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
		"POP-507": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/",

		// HPA
		"POP-600": "https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/",
		"POP-601": "https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/",
		"POP-602": "https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/",
		"POP-603": "https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/",
		"POP-604": "https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/",
		"POP-605": "https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/",

		// Node
		"POP-700": "https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/",
		"POP-701": "https://kubernetes.io/docs/concepts/architecture/nodes/#node-status",
		"POP-702": "https://kubernetes.io/docs/concepts/architecture/nodes/#node-status",
		"POP-703": "https://kubernetes.io/docs/concepts/architecture/nodes/#node-status",
		"POP-704": "https://kubernetes.io/docs/concepts/architecture/nodes/#node-status",
		"POP-705": "https://kubernetes.io/docs/concepts/architecture/nodes/#node-status",
		"POP-706": "https://kubernetes.io/docs/concepts/architecture/nodes/#node-status",
		"POP-707": "https://kubernetes.io/docs/concepts/architecture/nodes/#node-status",
		"POP-708": "https://kubernetes.io/docs/tasks/debug/debug-cluster/resource-metrics-pipeline/",
		"POP-709": "https://kubernetes.io/docs/concepts/architecture/nodes/",
		"POP-710": "https://kubernetes.io/docs/concepts/architecture/nodes/",
		"POP-711": "https://kubernetes.io/docs/concepts/architecture/nodes/#manual-node-administration",
		"POP-712": "https://kubernetes.io/docs/concepts/overview/components/",

		// Namespace
		"POP-800": "https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/",

		// PodDisruptionBudget
		"POP-900": "https://kubernetes.io/docs/concepts/workloads/pods/disruptions/",
		"POP-901": "https://kubernetes.io/docs/concepts/workloads/pods/disruptions/",

		// PV and PVC
		"POP-1000": "https://kubernetes.io/docs/concepts/storage/persistent-volumes/",
		"POP-1001": "https://kubernetes.io/docs/concepts/storage/persistent-volumes/",
		"POP-1002": "https://kubernetes.io/docs/concepts/storage/persistent-volumes/",
		"POP-1003": "https://kubernetes.io/docs/concepts/storage/persistent-volumes/",
		"POP-1004": "https://kubernetes.io/docs/concepts/storage/persistent-volumes/",

		// Service
		"POP-1100": "https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service",
		"POP-1101": "https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service",
		"POP-1102": "https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service",
		"POP-1103": "https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer",
		"POP-1104": "https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport",
		"POP-1105": "https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors",
		"POP-1106": "https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service",
		"POP-1107": "https://kubernetes.io/docs/concepts/services-networking/service/#external-traffic-policy",
		"POP-1108": "https://kubernetes.io/docs/concepts/services-networking/service/#external-traffic-policy",
		"POP-1109": "https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service",

		// ReplicaSet
		"POP-1120": "https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/",

		// NetworkPolicies
		"POP-1200": "https://kubernetes.io/docs/concepts/services-networking/network-policies/#networkpolicy-resource",
		"POP-1201": "https://kubernetes.io/docs/concepts/services-networking/network-policies/#networkpolicy-resource",

		// RBAC
		"POP-1300": "https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding",
	}

	// IssueIDtoCategory maps Popeye's issue codes to Category as described
	IssueIDtoCategory = map[string]string{
		"POP-100": "Container",
		"POP-101": "Container",
		"POP-102": "Container",
		"POP-103": "Container",
		"POP-104": "Container",
		"POP-105": "Container",
		"POP-106": "Container",
		"POP-107": "Container",
		"POP-108": "Container",
		"POP-109": "Container",
		"POP-110": "Container",
		"POP-111": "Container",
		"POP-112": "Container",
		"POP-113": "Container",

		"POP-200": "Pod",
		"POP-201": "Pod",
		"POP-202": "Pod",
		"POP-203": "Pod",
		"POP-204": "Pod",
		"POP-205": "Pod",
		"POP-206": "Pod",
		"POP-207": "Pod",
		"POP-208": "Pod",

		"POP-300": "Security",
		"POP-301": "Security",
		"POP-302": "Security",
		"POP-303": "Security",
		"POP-304": "Security",
		"POP-305": "Security",
		"POP-306": "Security",

		"POP-400": "General",
		"POP-401": "General",
		"POP-402": "General",
		"POP-403": "General",
		"POP-404": "General",
		"POP-405": "General",
		"POP-406": "General",

		"POP-500": "Workloads",
		"POP-501": "Workloads",
		"POP-503": "Workloads",
		"POP-504": "Workloads",
		"POP-505": "Workloads",
		"POP-506": "Workloads",
		"POP-507": "Workloads",

		"POP-600": "HorizontalPodAutoscaler",
		"POP-601": "HorizontalPodAutoscaler",
		"POP-602": "HorizontalPodAutoscaler",
		"POP-603": "HorizontalPodAutoscaler",
		"POP-604": "HorizontalPodAutoscaler",
		"POP-605": "HorizontalPodAutoscaler",

		"POP-700": "Node",
		"POP-701": "Node",
		"POP-702": "Node",
		"POP-703": "Node",
		"POP-704": "Node",
		"POP-705": "Node",
		"POP-706": "Node",
		"POP-707": "Node",
		"POP-708": "Node",
		"POP-709": "Node",
		"POP-710": "Node",
		"POP-711": "Node",
		"POP-712": "Node",

		"POP-800": "Namespace",

		"POP-900": "PodDisruptionBudget",
		"POP-901": "PodDisruptionBudget",

		"POP-1000": "Volumes",
		"POP-1001": "Volumes",
		"POP-1002": "Volumes",
		"POP-1003": "Volumes",
		"POP-1004": "Volumes",

		"POP-1100": "Service",
		"POP-1101": "Service",
		"POP-1102": "Service",
		"POP-1103": "Service",
		"POP-1104": "Service",
		"POP-1105": "Service",
		"POP-1106": "Service",
		"POP-1107": "Service",
		"POP-1108": "Service",
		"POP-1109": "Service",

		"POP-1120": "ReplicaSet",

		"POP-1200": "NetworkPolicies",
		"POP-1201": "NetworkPolicies",

		"POP-1300": "RBAC",
	}
)
