package templates

const SystemComponentsTemplate = `
---
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: controlplane-proxy
  namespace: kube-system
spec:
  template:
    metadata:
      labels:
        app: controlplane-proxy
    spec:
      nodeSelector:
        node-role.kubernetes.io/controlplane: "true"
      tolerations:
      - key: node-role.kubernetes.io/controlplane
        operator: "Equal"
        value: "true"
      - key: node-role.kubernetes.io/etcd
        operator: "Equal"
        value: "true"
      hostNetwork: true
      containers:
      - name: controlplane-proxy
        command: ["sh","-c","while :; do sleep 100; done"]
        image: alpine

---
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: etcd-proxy
  namespace: kube-system
spec:
  template:
    metadata:
      labels:
        app: etcd-proxy
    spec:
      nodeSelector:
        node-role.kubernetes.io/etcd: "true"
      tolerations:
      - key: node-role.kubernetes.io/controlplane
        operator: "Equal"
        value: "true"
      - key: node-role.kubernetes.io/etcd
        operator: "Equal"
        value: "true"
      hostNetwork: true
      containers:
      - name: etcd-proxy
        command: ["sh","-c","while :; do sleep 100; done"]
        image: alpine
---
apiVersion: v1
kind: Service
metadata:
  name: controller-manager
  namespace: kube-system
spec:
  ports:
  - port: 10257
    targetPort: 10257
  selector:
    app: controlplane-proxy
---
apiVersion: v1
kind: Service
metadata:
  name: scheduler
  namespace: kube-system
spec:
  ports:
  - port: 10251
    targetPort: 10251
  selector:
    app: controlplane-proxy
---
apiVersion: v1
kind: Service
metadata:
  name: etcd
  namespace: kube-system
spec:
  ports:
  - port: 2379
    targetPort: 2379
  selector:
    app: etcd-proxy
`
