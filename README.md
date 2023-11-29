Pod dependency Services
This container can be used as init container to specify dependency of other pod. It will check for exiting pod with success status. If any pod with given label selector is found running in current namespace it will exit with success else exit with failure after timeout.

## Environment Variable:
POD_LABELS - This is comma (,) separated string of labels of dependency pods which will be checked for Running phase.

MAX_RETRY -	Maximum number of times for which init container will try to check if dependency pods are Running.

RETRY_TIME_OUT_SECOND -	Number of seconds init container will pause between each retry.


## Example usage:
spec:
containers:
...
```
serviceAccountName: {{ .Values.serviceAccount }} #optional
initContainers:
- name: pod-dependency
  image: {{ .Values.pod-dependency-init-container.image }}
  env:
  - name: POD_LABELS
    value: "usermanagement"
  - name: MAX_RETRY
    value: "10"
  - name: RETRY_TIME_OUT_SECOND
    value: "60"
  - name: NAMESPACE
    value: {{ .Values.Release.Namespace }}
     
```
___
## RBAC
In case of RBAC this container requires pods resource get, list, watch access. Which can be provided by below yaml
```---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Values.serviceAccount }}
  namespace: {{ .Values.Release.Namespace }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pods-dependency-role
  namespace: {{ .Values.Release.Namespace }}
rules:
  - apiGroups:
      - ""
      - apps
      - autoscaling
      - batch
      - extensions
      - policy
      - rbac.authorization.k8s.io
    resources:
      - pods
      - componentstatuses
      - configmaps
      - daemonsets
      - deployments
      - events
      - endpoints
      - horizontalpodautoscalers
      - ingress
      - jobs
      - cronjobs
      - limitranges
      - namespaces
      - nodes
      - pods
      - persistentvolumes
      - persistentvolumeclaims
      - resourcequotas
      - replicasets
      - replicationcontrollers
      - serviceaccounts
      - secrets
      - services
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: pods-dependency-role-binding
  namespace: {{ .Values.Release.Namespace }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: pods-dependency-role
subjects:
  - namespace: {{ .Values.Release.Namespace }}
    kind: ServiceAccount
    name: {{ .Values.serviceAccount }}
---
apiVersion: v1
kind: Pod
metadata:
  name: pods-dependency
  namespace: {{ .Values.Release.Namespace }}
  annotations:
    kubectl.kubernetes.io/default-container: "busybox-container"
spec:
  serviceAccountName: {{ .Values.serviceAccount }}
  initContainers:
    - name: pod-dependency
      image: docker-sandbox.infra.cloudera.com/manish.paneri/pod-dependency-init-container:v11
      env:
        - name: POD_LABELS
          value: "usermanagement,enviromentservice"
        - name: MAX_RETRY
          value: "5"
        - name: RETRY_TIME_OUT_SECOND
          value: "60"
        - name: NAMESPACE
          value: {{ .Values.Release.Namespace }}
  containers:
    - name: busybox-container
      image: busybox
      command: ["/bin/sh", "-ec", "sleep 1000"]```