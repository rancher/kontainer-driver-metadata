package templates

const KubeRouterTemplate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: kube-router-cfg
  namespace: kube-system
  labels:
    tier: node
    k8s-app: kube-router
data:
  cni-conf.json: |
    {
      "cniVersion":"0.3.0",
      "name":"mynet",
      "plugins":[
        {
            "name":"kubernetes",
            "type":"bridge",
            "bridge":"kube-bridge",
            "isDefaultGateway":true,
            "forceAddress": true,
            "ipam":{
              "type":"host-local"
            }
        },
        {
            "type":"portmap",
            "capabilities":{
              "snat":true,
              "portMappings":true
            }
        }
      ]
    }
  kubeconfig: |
    apiVersion: v1
    kind: Config
    clusterCIDR: "{{.ClusterCIDR}}"
    clusters:
    - name: cluster
      cluster:
        certificate-authority: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        server: "{{.APIRoot}}"
    users:
    - name: kube-router
      user:
        tokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    contexts:
    - context:
        cluster: cluster
        user: kube-router
      name: kube-router-context
    current-context: kube-router-context

---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: kube-router
    tier: node
  name: kube-router
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: kube-router
      tier: node
  template:
    metadata:
      labels:
        k8s-app: kube-router
        tier: node
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "20243"
    spec:
      priorityClassName: system-node-critical
      serviceAccountName: kube-router
      containers:
      - name: kube-router
        image: {{.CNIImage}}
        imagePullPolicy: IfNotPresent
        args:
        - "--run-router=true"
        - "--run-firewall={{.RunFirewall}}"
        - "--run-service-proxy={{.RunServiceProxy}}"
        - "--kubeconfig=/var/lib/kube-router/kubeconfig"
        - "--bgp-graceful-restart" 
        - "--metrics-port=20243"
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: KUBE_ROUTER_CNI_CONF_FILE
          value: /etc/cni/net.d/10-kuberouter.conflist
        livenessProbe:
          httpGet:
            path: /healthz
            port: 20244
          initialDelaySeconds: 10
          periodSeconds: 3
        resources:
          requests:
            cpu: 250m
            memory: 250Mi
        securityContext:
          privileged: true
        volumeMounts:
        - name: lib-modules
          mountPath: /lib/modules
          readOnly: true
        - name: cni-conf-dir
          mountPath: /etc/cni/net.d
        - name: kubeconfig
          mountPath: /var/lib/kube-router
          readOnly: true
        - name: xtables-lock
          mountPath: /run/xtables.lock
          readOnly: false
      initContainers:
      - name: install-cni
        image: {{.CNIImage}}
        imagePullPolicy: Always
        command:
        - /bin/sh
        - -c
        - set -e -x;
          if [ ! -f /etc/cni/net.d/10-kuberouter.conflist ]; then
            if [ -f /etc/cni/net.d/*.conf ]; then
              rm -f /etc/cni/net.d/*.conf;
            fi;
            TMP=/etc/cni/net.d/.tmp-kuberouter-cfg;
            cp /etc/kube-router/cni-conf.json ${TMP};
            mv ${TMP} /etc/cni/net.d/10-kuberouter.conflist;
          fi;
          if [ ! -f /var/lib/kube-router/kubeconfig ]; then
            TMP=/var/lib/kube-router/.tmp-kubeconfig;
            cp /etc/kube-router/kubeconfig ${TMP};
            mv ${TMP} /var/lib/kube-router/kubeconfig;
          fi
        volumeMounts:
        - mountPath: /etc/cni/net.d
          name: cni-conf-dir
        - mountPath: /etc/kube-router
          name: kube-router-cfg
        - name: kubeconfig
          mountPath: /var/lib/kube-router
      hostNetwork: true
      tolerations:
      - effect: NoSchedule
        operator: Exists
      - key: CriticalAddonsOnly
        operator: Exists
      - effect: NoExecute
        operator: Exists
      volumes:
      - name: lib-modules
        hostPath:
          path: /lib/modules
      - name: cni-conf-dir
        hostPath:
          path: /etc/cni/net.d
      - name: kube-router-cfg
        configMap:
          name: kube-router-cfg
      - name: kubeconfig
        hostPath:
          path: /var/lib/kube-router
      - name: xtables-lock
        hostPath:
          path: /run/xtables.lock
          type: FileOrCreate

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-router
  namespace: kube-system

---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kube-router
  namespace: kube-system
rules:
  - apiGroups:
    - ""
    resources:
      - namespaces
      - pods
      - services
      - nodes
      - endpoints
    verbs:
      - list
      - get
      - watch
  - apiGroups:
    - "networking.k8s.io"
    resources:
      - networkpolicies
    verbs:
      - list
      - get
      - watch
  - apiGroups:
    - extensions
    resources:
      - networkpolicies
    verbs:
      - get
      - list
      - watch
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kube-router
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kube-router
subjects:
- kind: ServiceAccount
  name: kube-router
  namespace: kube-system
`

const KubeRouterTemplateV116 = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: kube-router-cfg
  namespace: kube-system
  labels:
    tier: node
    k8s-app: kube-router
data:
  cni-conf.json: |
    {
      "cniVersion":"0.3.0",
      "name":"mynet",
      "plugins":[
        {
            "name":"kubernetes",
            "type":"bridge",
            "bridge":"kube-bridge",
            "isDefaultGateway":true,
            "forceAddress": true,
            "ipam":{
              "type":"host-local"
            }
        },
        {
            "type":"portmap",
            "capabilities":{
              "snat":true,
              "portMappings":true
            }
        }
      ]
    }
  kubeconfig: |
    apiVersion: v1
    kind: Config
    clusterCIDR: "{{.ClusterCIDR}}"
    clusters:
    - name: cluster
      cluster:
        certificate-authority: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        server: "{{.APIRoot}}"
    users:
    - name: kube-router
      user:
        tokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    contexts:
    - context:
        cluster: cluster
        user: kube-router
      name: kube-router-context
    current-context: kube-router-context

---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: kube-router
    tier: node
  name: kube-router
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: kube-router
      tier: node
  template:
    metadata:
      labels:
        k8s-app: kube-router
        tier: node
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "20243"
    spec:
      priorityClassName: system-node-critical
      serviceAccountName: kube-router
      containers:
      - name: kube-router
        image: {{.CNIImage}}
        imagePullPolicy: IfNotPresent
        args:
        - "--run-router=true"
        - "--run-firewall={{.RunFirewall}}"
        - "--run-service-proxy={{.RunServiceProxy}}"
        - "--kubeconfig=/var/lib/kube-router/kubeconfig"
        - "--bgp-graceful-restart" 
        - "--metrics-port=20243"
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: KUBE_ROUTER_CNI_CONF_FILE
          value: /etc/cni/net.d/10-kuberouter.conflist
        livenessProbe:
          httpGet:
            path: /healthz
            port: 20244
          initialDelaySeconds: 10
          periodSeconds: 3
        resources:
          requests:
            cpu: 250m
            memory: 250Mi
        securityContext:
          privileged: true
        volumeMounts:
        - name: lib-modules
          mountPath: /lib/modules
          readOnly: true
        - name: cni-conf-dir
          mountPath: /etc/cni/net.d
        - name: kubeconfig
          mountPath: /var/lib/kube-router
          readOnly: true
        - name: xtables-lock
          mountPath: /run/xtables.lock
          readOnly: false
      initContainers:
      - name: install-cni
        image: {{.CNIImage}}
        imagePullPolicy: Always
        command:
        - /bin/sh
        - -c
        - set -e -x;
          if [ ! -f /etc/cni/net.d/10-kuberouter.conflist ]; then
            if [ -f /etc/cni/net.d/*.conf ]; then
              rm -f /etc/cni/net.d/*.conf;
            fi;
            TMP=/etc/cni/net.d/.tmp-kuberouter-cfg;
            cp /etc/kube-router/cni-conf.json ${TMP};
            mv ${TMP} /etc/cni/net.d/10-kuberouter.conflist;
          fi;
          if [ ! -f /var/lib/kube-router/kubeconfig ]; then
            TMP=/var/lib/kube-router/.tmp-kubeconfig;
            cp /etc/kube-router/kubeconfig ${TMP};
            mv ${TMP} /var/lib/kube-router/kubeconfig;
          fi
        volumeMounts:
        - mountPath: /etc/cni/net.d
          name: cni-conf-dir
        - mountPath: /etc/kube-router
          name: kube-router-cfg
        - name: kubeconfig
          mountPath: /var/lib/kube-router
      hostNetwork: true
      tolerations:
      - effect: NoSchedule
        operator: Exists
      - key: CriticalAddonsOnly
        operator: Exists
      - effect: NoExecute
        operator: Exists
      volumes:
      - name: lib-modules
        hostPath:
          path: /lib/modules
      - name: cni-conf-dir
        hostPath:
          path: /etc/cni/net.d
      - name: kube-router-cfg
        configMap:
          name: kube-router-cfg
      - name: kubeconfig
        hostPath:
          path: /var/lib/kube-router
      - name: xtables-lock
        hostPath:
          path: /run/xtables.lock
          type: FileOrCreate

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-router
  namespace: kube-system

---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kube-router
  namespace: kube-system
rules:
  - apiGroups:
    - ""
    resources:
      - namespaces
      - pods
      - services
      - nodes
      - endpoints
    verbs:
      - list
      - get
      - watch
  - apiGroups:
    - "networking.k8s.io"
    resources:
      - networkpolicies
    verbs:
      - list
      - get
      - watch
  - apiGroups:
    - extensions
    resources:
      - networkpolicies
    verbs:
      - get
      - list
      - watch
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kube-router
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kube-router
subjects:
- kind: ServiceAccount
  name: kube-router
  namespace: kube-system
`
