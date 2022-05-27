package report

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	inspectv1a1 "github.com/getupio-undistro/inspect/apis/inspect/v1alpha1"
	"github.com/getupio-undistro/inspect/worker/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// NewClusterIssue creates and returns a pointer to a <ClusterIssue> instance
// carrying issue metadata on its labels. The instance is set as a child of the
// Job whereby the plugin executed.
func NewClusterIssue(c *config.Config, cispec *inspectv1a1.ClusterIssueSpec, orefs []metav1.OwnerReference, jid *string) *inspectv1a1.ClusterIssue {
	cispec.Cluster = c.Cluster
	return &inspectv1a1.ClusterIssue{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterIssue",
			APIVersion: inspectv1a1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-%s-%s", c.Cluster, strings.ToLower(cispec.ID), *jid),
			Namespace:       c.ClusterIssuesNs,
			OwnerReferences: orefs,
			Labels: map[string]string{
				inspectv1a1.LabelScanID:   c.JobUid,
				inspectv1a1.LabelCluster:       c.Cluster,
				inspectv1a1.LabelSeverity:      string(cispec.Severity),
				inspectv1a1.LabelIssueID:       cispec.ID,
				inspectv1a1.LabelIssueCategory: cispec.Category,
			},
		},
		Spec: *cispec,
	}
}

// Parse receives a reader pointing to a plugin's report file, transforming
// such report into an array of <ClusterIssue> pointers according to the
// cluster name and issues namespace specified on the <Config> struct. The
// parsing for each plugin is left to dedicated functions which are called
// according to the plugin type.
func Parse(r io.Reader, c *config.Config) ([]*inspectv1a1.ClusterIssue, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid configuration: %w", err)
	}
	repby, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("Unable to read results of plugin <%s> from cluster <%s>: %w", c.Plugin, c.Cluster, err)
	}
	cispecs, err := config.PluginParsers[c.Plugin](repby)
	if err != nil {
		return nil, err
	}

	juid := c.JobUid[strings.LastIndex(c.JobUid, "-")+1:]
	orefs := []metav1.OwnerReference{{
		APIVersion: "batch/v1",
		Kind:       "Job",
		Name:       c.Job,
		UID:        types.UID(c.JobUid),
	}}
	ciarr := make([]*inspectv1a1.ClusterIssue, len(cispecs))
	for i := 0; i < len(cispecs); i++ {
		ciarr[i] = NewClusterIssue(c, cispecs[i], orefs, &juid)
	}
	return ciarr, nil
}
