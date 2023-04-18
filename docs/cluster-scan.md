# Configure a cluster scan

Since your clusters are connected the next and last step is configure a scan for them
by creating a `ClusterScan` in the same namespace as `Cluster` resource.

The `ClusterScan` will be responsible for reporting issues and vulnerabilities of your clusters.

Failure to perform this step implies that the scan will not be performed, and therefore the health of your cluster will be unknown.


## Create a `ClusterScan`

The `ClusterScan` scans the `Cluster` referenced in `clusterRef.name` field periodically on a given schedule, 
written in [Cron](https://en.wikipedia.org/wiki/Cron) format.

Here is a sample configuration that scan `mycluster` once an hour.
You can modify putting your desired periodicity.

```yaml
cat << EOF | kubectl apply -f -
apiVersion: zora.undistro.io/v1alpha1
kind: ClusterScan
metadata:
  name: mycluster
  namespace: zora-system
spec:
  clusterRef:
    name: mycluster
  schedule: "0 * * * *"  # at minute 0 past every hour
EOF
```

### Cron schedule syntax

Cron expression has five fields separated by a space, and each field represents a time unit.


```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of the month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday;
│ │ │ │ │                                   7 is also Sunday on some systems)
│ │ │ │ │                                   OR sun, mon, tue, wed, thu, fri, sat
│ │ │ │ │
* * * * *
```

| Operator | Descriptor           | Example                                                                                            |
|----------|----------------------|----------------------------------------------------------------------------------------------------|
| *        | Any value            | `15 * * * *` runs at every minute 15 of every hour of every day.                                   |
| ,        | Value list separator | `2,10 4,5 * * *` runs at minute 2 and 10 of the 4th and 5th hour of every day.                     |
| -        | Range of values      | `30 4-6 * * *` runs at minute 30 of the 4th, 5th, and 6th hour.                                    |
| /        | Step values          | `20/15 * * * *` runs every 15 minutes starting from minute 20 through 59 (minutes 20, 35, and 50). |


Now Zora is ready to help you to identify potential issues and vulnerabilities in your kubernetes clusters.

You can check the scans status and the reported issues by the following steps:

## List cluster scans

Listing the `ClusterScans`, the information of the last scans are available:

```shell
kubectl get clusterscan -o wide
```
```
NAME        CLUSTER     SCHEDULE    SUSPEND   PLUGINS         LAST STATUS   LAST SCHEDULE   LAST SUCCESSFUL   ISSUES   READY   SAAS   AGE   NEXT SCHEDULE
mycluster   mycluster   0 * * * *   false     marvin,popeye   Complete      13s             1s                34       True    OK     39s   2023-04-18T14:00:00Z
```

The `LAST STATUS` column represents the status (`Active`, `Complete` or `Failed`) of the last **scan** 
that was scheduled at the time represented by `LAST SCHEDULE` column.

## List cluster issues

Once the cluster is successfully scanned,
the reported issues are available in `ClusterIssue` resources:

```shell
kubectl get clusterissues -l cluster=mycluster
```
```
NAME                              CLUSTER     ID         MESSAGE                                                                             SEVERITY   CATEGORY                  AGE
mycluster-m-102-18e887d99ccb      mycluster   M-102      Privileged container                                                                High       Security                  100s
mycluster-m-103-18e887d99ccb      mycluster   M-103      Insecure capabilities                                                               High       Security                  100s
mycluster-m-104-18e887d99ccb      mycluster   M-104      HostPath volume                                                                     High       Security                  100s
mycluster-m-105-18e887d99ccb      mycluster   M-105      Not allowed hostPort                                                                High       Security                  100s
mycluster-m-111-18e887d99ccb      mycluster   M-111      Not allowed volume type                                                             Low        Security                  100s
mycluster-m-112-18e887d99ccb      mycluster   M-112      Allowed privilege escalation                                                        Medium     Security                  100s
mycluster-m-113-18e887d99ccb      mycluster   M-113      Container could be running as root user                                             Medium     Security                  100s
mycluster-m-115-18e887d99ccb      mycluster   M-115      Not allowed seccomp profile                                                         Low        Security                  100s
mycluster-m-201-18e887d99ccb      mycluster   M-201      Application credentials stored in configuration files                               High       Security                  100s
mycluster-m-300-18e887d99ccb      mycluster   M-300      Root filesystem write allowed                                                       Low        Security                  100s
mycluster-pop-102-c6d6b0eefab4    mycluster   POP-102    No probes defined                                                                   Medium     Container                 103s
mycluster-pop-106-c6d6b0eefab4    mycluster   POP-106    No resources requests/limits defined                                                Medium     Container                 103s
mycluster-pop-605-c6d6b0eefab4    mycluster   POP-605    If ALL HPAs are triggered, cluster memory capacity will match or exceed threshold   Medium     HorizontalPodAutoscaler   103s
mycluster-pop-710-c6d6b0eefab4    mycluster   POP-710    Node Memory threshold reached                                                       Medium     Node                      103s
```

It's possible filter issues by cluster, issue ID, severity and category 
using [label selector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/):

```shell
# issues from mycluster
kubectl get clusterissues -l cluster=mycluster

# clusters with issue POP-106
kubectl get clusterissues -l id=POP-106

# issues from mycluster with high severity
kubectl get clusterissues -l cluster=mycluster,severity=High

# only issues reported by the last scan from mycluster
kubectl get clusterissues -l cluster=mycluster,scanID=fa4e63cc-5236-40f3-aa7f-599e1c83208b

# issues reported from marvin plugin
kubectl get clusterissues -l plugin=marvin
```

!!! tip "Why is it an issue?"

    The field `url` in `ClusterIssue` spec represents a link for a documentation about this issue.
    It is displayed in the UI and you can see by `kubectl` with the `-o=yaml` flag or the command below.
    
    ```shell
    kubectl get clusterissues -o=custom-columns="NAME:.metadata.name,MESSAGE:.spec.message,URL:.spec.url"
    ```
    ```
    NAME                          MESSAGE                                                                        URL
    mycluster-pop-102-27557035    No probes defined                                                              https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
    mycluster-pop-105-27557035    Liveness probe uses a port#, prefer a named port                               <none>
    mycluster-pop-106-27557035    No resources requests/limits defined                                           https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
    mycluster-pop-1100-27557035   No pods match service selector                                                 https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service
    mycluster-pop-306-27557035    Container could be running as root user. Check SecurityContext/Image           https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted
    mycluster-pop-500-27557035    Zero scale detected                                                            https://kubernetes.io/docs/concepts/workloads/
    ```
    
    These docs should help you understand why it's an issue and how to fix it.
    
    All URLs are available [here](https://github.com/undistro/zora/blob/main/pkg/worker/report/popeye/parse_types.go#L109) 
    and you can contribute to Zora adding new links. See our [contribution guidelines](https://github.com/undistro/zora/blob/main/CONTRIBUTING.md).
