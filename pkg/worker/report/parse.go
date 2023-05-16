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

package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-logr/logr"

	zorav1a1 "github.com/undistro/zora/api/zora/v1alpha1"
	"github.com/undistro/zora/pkg/worker/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// NewClusterIssue creates and returns a pointer to a <ClusterIssue> instance
// carrying issue metadata on its labels. The instance is set as a child of the
// Job whereby the plugin executed.
func NewClusterIssue(c *config.Config, cispec *zorav1a1.ClusterIssueSpec, orefs []metav1.OwnerReference, juid *string) *zorav1a1.ClusterIssue {
	cispec.Cluster = c.Cluster
	return &zorav1a1.ClusterIssue{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterIssue",
			APIVersion: zorav1a1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-%s-%s", c.Cluster, strings.ToLower(cispec.ID), *juid),
			Namespace:       c.ClusterIssuesNs,
			OwnerReferences: orefs,
			Labels: map[string]string{
				zorav1a1.LabelScanID:   c.JobUID,
				zorav1a1.LabelCluster:  c.Cluster,
				zorav1a1.LabelSeverity: string(cispec.Severity),
				zorav1a1.LabelIssueID:  cispec.ID,
				zorav1a1.LabelCategory: cispec.Category,
				zorav1a1.LabelPlugin:   c.Plugin,
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
func Parse(log logr.Logger, r io.Reader, c *config.Config) ([]*zorav1a1.ClusterIssue, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid configuration: %w", err)
	}
	repby, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("Unable to read results of plugin <%s> from cluster <%s>: %w", c.Plugin, c.Cluster, err)
	}
	cispecs, err := config.PluginParsers[c.Plugin](log, repby)
	if err != nil {
		return nil, err
	}

	juid := c.JobUID[strings.LastIndex(c.JobUID, "-")+1:]
	orefs := []metav1.OwnerReference{{
		APIVersion: "batch/v1",
		Kind:       "Job",
		Name:       c.Job,
		UID:        types.UID(c.JobUID),
	}}
	ciarr := make([]*zorav1a1.ClusterIssue, len(cispecs))
	for i := 0; i < len(cispecs); i++ {
		ciarr[i] = NewClusterIssue(c, cispecs[i], orefs, &juid)
	}
	return ciarr, nil
}
