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
             "ipam":{
                "type":"host-local"
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
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      priorityClassName: system-node-critical
      serviceAccountName: kube-router
      serviceAccount: kube-router
      containers:
      - name: kube-router
        image: {{.CNIImage}}
        imagePullPolicy: Always
        args:
        - "--run-router=true"
        - "--run-firewall=true"
        - "--run-service-proxy={{.RunServiceProxy}}"
        - "--kubeconfig=/var/lib/kube-router/kubeconfig"
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
      initContainers:
      - name: install-cni
        image: alpine
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
            mv -f ${TMP} /var/lib/kube-router/kubeconfig;
          fi;
          mkdir -p /opt/cni/bin;
          wget https://github.com/containernetworking/plugins/releases/download/v0.8.1/cni-plugins-linux-amd64-v0.8.1.tgz -O /tmp/cni-plugins-linux-amd64-v0.8.1.tgz && tar -xf /tmp/cni-plugins-linux-amd64-v0.8.1.tgz -C /opt/cni/bin/
        volumeMounts:
        - mountPath: /etc/cni/net.d
          name: cni-conf-dir
        - mountPath: /opt/cni/bin
          name: cni-bin-dir
        - mountPath: /etc/kube-router
          name: kube-router-cfg
        - name: kubeconfig
          mountPath: /var/lib/kube-router
      hostNetwork: true
      dnsPolicy: ClusterFirstWithHostNet
      tolerations:
      - key: CriticalAddonsOnly
        operator: Exists
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
        operator: Exists
      - effect: NoSchedule
        key: node.kubernetes.io/not-ready
        operator: Exists
      volumes:
      - name: lib-modules
        hostPath:
          path: /lib/modules
      - name: cni-conf-dir
        hostPath:
          path: /etc/cni/net.d
      - name: cni-bin-dir
        hostPath:
          path: /opt/cni/bin
      - name: kube-router-cfg
        configMap:
          name: kube-router-cfg
      - name: kubeconfig
        hostPath:
          path: /var/lib/kube-router
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-router
  namespace: kube-system
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
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
apiVersion: rbac.authorization.k8s.io/v1beta1
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
             "ipam":{
                "type":"host-local"
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
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      priorityClassName: system-node-critical
      serviceAccountName: kube-router
      serviceAccount: kube-router
      containers:
      - name: kube-router
        image: {{.CNIImage}}
        imagePullPolicy: Always
        args:
        - "--run-router=true"
        - "--run-firewall=true"
        - "--run-service-proxy={{.RunServiceProxy}}"
        - "--kubeconfig=/var/lib/kube-router/kubeconfig"
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
      initContainers:
      - name: install-cni
        image: alpine
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
            mv -f ${TMP} /var/lib/kube-router/kubeconfig;
          fi;
          mkdir -p /opt/cni/bin;
          wget https://github.com/containernetworking/plugins/releases/download/v0.8.1/cni-plugins-linux-amd64-v0.8.1.tgz -O /tmp/cni-plugins-linux-amd64-v0.8.1.tgz && tar -xf /tmp/cni-plugins-linux-amd64-v0.8.1.tgz -C /opt/cni/bin/
        volumeMounts:
        - mountPath: /etc/cni/net.d
          name: cni-conf-dir
        - mountPath: /opt/cni/bin
          name: cni-bin-dir
        - mountPath: /etc/kube-router
          name: kube-router-cfg
        - name: kubeconfig
          mountPath: /var/lib/kube-router
      hostNetwork: true
      dnsPolicy: ClusterFirstWithHostNet
      tolerations:
      - key: CriticalAddonsOnly
        operator: Exists
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
        operator: Exists
      - effect: NoSchedule
        key: node.kubernetes.io/not-ready
        operator: Exists
      volumes:
      - name: lib-modules
        hostPath:
          path: /lib/modules
      - name: cni-conf-dir
        hostPath:
          path: /etc/cni/net.d
      - name: cni-bin-dir
        hostPath:
          path: /opt/cni/bin
      - name: kube-router-cfg
        configMap:
          name: kube-router-cfg
      - name: kubeconfig
        hostPath:
          path: /var/lib/kube-router
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-router
  namespace: kube-system
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
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
apiVersion: rbac.authorization.k8s.io/v1beta1
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
