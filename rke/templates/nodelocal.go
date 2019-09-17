package templates

const NodelocalTemplateV116 = `
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: node-local-dns
  namespace: kube-system
  labels:
    kubernetes.io/cluster-service: "true"
    addonmanager.kubernetes.io/mode: Reconcile
{{- end }}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: node-local-dns
  namespace: kube-system
  labels:
    addonmanager.kubernetes.io/mode: Reconcile
data:
  Corefile: |
    {{.ClusterDomain}}:53 {
        errors
        cache {
                success 9984 30
                denial 9984 5
        }
        reload
        loop
        bind {{.IPAddress}}
        forward . {{.ClusterDNSServer}} {
                force_tcp
        }
        prometheus :9253
        health {{.IPAddress}}:8080
        }
    in-addr.arpa:53 {
        errors
        cache 30
        reload
        loop
        bind {{.IPAddress}}
        forward . {{.ClusterDNSServer}} {
                force_tcp
        }
        prometheus :9253
        }
    ip6.arpa:53 {
        errors
        cache 30
        reload
        loop
        bind {{.IPAddress}}
        forward . {{.ClusterDNSServer}} {
                force_tcp
        }
        prometheus :9253
        }
    .:53 {
        errors
        cache 30
        reload
        loop
        bind {{.IPAddress}}
        forward . /etc/resolv.conf {
                force_tcp
        }
        prometheus :9253
        }
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-local-dns
  namespace: kube-system
  labels:
    k8s-app: node-local-dns
    kubernetes.io/cluster-service: "true"
    addonmanager.kubernetes.io/mode: Reconcile
spec:
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 10%
  selector:
    matchLabels:
      k8s-app: node-local-dns
  template:
    metadata:
       labels:
          k8s-app: node-local-dns
    spec:
      priorityClassName: system-node-critical
{{- if eq .RBACConfig "rbac"}}
      serviceAccountName: node-local-dns
{{- end }}
      hostNetwork: true
      dnsPolicy: Default  # Don't use cluster DNS.
      tolerations:
      - key: "CriticalAddonsOnly"
        operator: "Exists"
      containers:
      - name: node-cache
        image: {{.NodelocalImage}}
        resources:
          limits:
            memory: 30Mi
          requests:
            cpu: 25m
            memory: 5Mi
        args: [ "-localip", "{{.IPAddress}}", "-conf", "/etc/coredns/Corefile" ]
        securityContext:
          privileged: true
        ports:
        - containerPort: 53
          name: dns
          protocol: UDP
        - containerPort: 53
          name: dns-tcp
          protocol: TCP
        - containerPort: 9253
          name: metrics
          protocol: TCP
        livenessProbe:
          httpGet:
            host: {{.IPAddress}}
            path: /health
            port: 8080
          initialDelaySeconds: 60
          timeoutSeconds: 5
        volumeMounts:
        - mountPath: /run/xtables.lock
          name: xtables-lock
          readOnly: false
        - name: config-volume
          mountPath: /etc/coredns
      volumes:
      - name: xtables-lock
        hostPath:
          path: /run/xtables.lock
          type: FileOrCreate
      - name: config-volume
        configMap:
          name: node-local-dns
          items:
            - key: Corefile
              path: Corefile
`
