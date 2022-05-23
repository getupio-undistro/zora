package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	inspectv1a1 "github.com/getupio-undistro/inspect/apis/inspect/v1alpha1"
	"github.com/getupio-undistro/inspect/worker/report/popeye"
)

const (
	DefaultDoneDir        = "/tmp/undistro-inspect/results"
	DoneDirEnvVar         = "DONE_DIR"
	PluginEnvVar          = "PLUGIN_NAME"
	ClusterEnvVar         = "CLUSTER_NAME"
	ClusterIssuesNsEnvVar = "CLUSTER_ISSUES_NAMESPACE"
	JobEnvVar             = "JOB_NAME"
	JobUidEnvVar          = "JOB_UID"
)

// PluginParsers correlates plugins with their respective parsing functions.
var PluginParsers = map[string]func([]byte) ([]*inspectv1a1.ClusterIssueSpec, error){
	"popeye": popeye.Parse,
}

// Config stores information used by the worker to create a list of
// <ClusterIssue> instances, and to specify the "done" file path.
type Config struct {
	DonePath        string `json:"donePath"`
	Plugin          string `json:"plugin"`
	Cluster         string `json:"cluster"`
	ClusterIssuesNs string `json:"listClusterIssueNs"`
	Job             string `json:"job"`
	JobUid          string `json:"jobUid"`
}

// New instantiates a new <Config> struct, with the default path for the
// "done" file.
func New() *Config {
	return &Config{DonePath: fmt.Sprintf("%s/done", DefaultDoneDir)}
}

// FromEnv instantiates a new <Config> struct, with values taken from the
// environment. It'll return an error in case cluster related variables aren't
// found.
func FromEnv() (*Config, error) {
	c := New()
	if e := os.Getenv(PluginEnvVar); len(e) != 0 {
		c.Plugin = e
	} else {
		return nil, fmt.Errorf("Empty environment variable <%s>", PluginEnvVar)
	}
	if e := os.Getenv(ClusterEnvVar); len(e) != 0 {
		c.Cluster = e
	} else {
		return nil, fmt.Errorf("Empty environment variable <%s>", ClusterEnvVar)
	}
	if e := os.Getenv(ClusterIssuesNsEnvVar); len(e) != 0 {
		c.ClusterIssuesNs = e
	} else {
		return nil, fmt.Errorf("Empty environment variable <%s>", ClusterIssuesNsEnvVar)
	}

	if e := os.Getenv(JobEnvVar); len(e) != 0 {
		c.Job = e
	} else {
		return nil, fmt.Errorf("Empty environment variable <%s>", JobEnvVar)
	}
	if e := os.Getenv(JobUidEnvVar); len(e) != 0 {
		c.JobUid = e
	} else {
		return nil, fmt.Errorf("Empty environment variable <%s>", JobUidEnvVar)
	}

	if e := os.Getenv(DoneDirEnvVar); len(e) != 0 {
		c.DonePath = fmt.Sprintf("%s/done", e)
	}
	return c, nil
}

// Validate ensures a <Config> instance has all its fields populated, and the
// plugin specified is supported by the worker.
func (r *Config) Validate() error {
	if len(r.DonePath) == 0 {
		return errors.New("Config's <DonePath> field is empty")
	}
	if len(r.Cluster) == 0 {
		return errors.New("Config's <Cluster> field is empty")
	}
	if len(r.ClusterIssuesNs) == 0 {
		return errors.New("Config's <ClusterIssuesNs> field is empty")
	}
	if len(r.Plugin) == 0 {
		return errors.New("Config's <Plugin> field is empty")
	}

	if len(r.Job) == 0 {
		return errors.New("Config's <Job> field is empty")
	} else if i := strings.LastIndex(r.Job, "-"); i == -1 || i == len(r.Job)-1 {
		return errors.New("Config's <Job> field is invalid")
	}
	if len(r.JobUid) == 0 {
		return errors.New("Config's <JobUid> field is empty")
	}

	if _, ok := PluginParsers[r.Plugin]; !ok {
		return fmt.Errorf("Invalid plugin: <%s>", r.Plugin)
	}
	return nil
}

// HandleDonePath ensures the directory wherefrom the "done" file will be
// written exists.
func (r *Config) HandleDonePath() error {
	if len(r.DonePath) == 0 {
		return errors.New("Empty <DonePath>")

	}
	dir := path.Dir(r.DonePath)
	if _, err := os.Stat(dir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Unable to check existance of dir <%s>: %w", dir, err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("Unable to create results dir <%s>: %w", dir, err)
	}
	return nil
}

// HandleResultsPath returns an <io.Reader> pointing to the path inside the
// "done" file.
func (r *Config) HandleResultsPath() (io.Reader, error) {
	fiby, err := ioutil.ReadFile(r.DonePath)
	if err != nil {
		return nil, errors.New("Unable to read 'done' file")
	}
	if len(fiby) == 0 {
		return nil, errors.New("Empty 'done' file")
	}

	fid, err := os.Open(strings.TrimSpace(string(fiby)))
	if err != nil {
		return nil, fmt.Errorf("Unable to open 'done' file: %w", err)
	}
	finf, err := fid.Stat()
	if err != nil {
		return nil, fmt.Errorf("Invalid path in 'done' file: %w", err)
	}
	if finf.IsDir() {
		return nil, errors.New("Path in the 'done' file points to a directory")
	}
	return fid, nil
}
