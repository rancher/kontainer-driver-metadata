package templates

const WeaveTemplate = `
---
# This ConfigMap can be used to configure a self-hosted Weave Net installation.
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: weave-net
      namespace: kube-system
  - apiVersion: extensions/v1beta1
    kind: DaemonSet
    metadata:
      name: weave-net
      labels:
        name: weave-net
      namespace: kube-system
    spec:
      template:
        metadata:
          annotations:
            scheduler.alpha.kubernetes.io/critical-pod: ''
            scheduler.alpha.kubernetes.io/tolerations: >-
              [{"key":"dedicated","operator":"Equal","value":"master","effect":"NoSchedule"}]
          labels:
            name: weave-net
        spec:
          affinity:
            nodeAffinity:
              requiredDuringSchedulingIgnoredDuringExecution:
                nodeSelectorTerms:
                  - matchExpressions:
                    - key: beta.kubernetes.io/os
                      operator: NotIn
                      values:
                        - windows
{{if .NodeSelector}}
          nodeSelector:
            {{ range $k, $v := .NodeSelector }}
              {{ $k }}: "{{ $v }}"
            {{ end }}
{{end}}
          containers:
            - name: weave
              command:
                - /home/weave/launch.sh
              env:
                - name: HOSTNAME
                  valueFrom:
                    fieldRef:
                      apiVersion: v1
                      fieldPath: spec.nodeName
                - name: IPALLOC_RANGE
                  value: "{{.ClusterCIDR}}"
                {{- if .WeavePassword}}
                - name: WEAVE_PASSWORD
                  value: "{{.WeavePassword}}"
                {{- end}}
                {{- if .MTU }}
                {{- if ne .MTU 0 }}
                - name: WEAVE_MTU
                  value: "{{.MTU}}"
                {{- end }}
                {{- end }}
              image: {{.Image}}
              readinessProbe:
                httpGet:
                  host: 127.0.0.1
                  path: /status
                  port: 6784
                initialDelaySeconds: 30
              resources:
                requests:
                  cpu: 10m
              securityContext:
                privileged: true
              volumeMounts:
                - name: weavedb
                  mountPath: /weavedb
                - name: cni-bin
                  mountPath: /host/opt
                - name: cni-bin2
                  mountPath: /host/home
                - name: cni-conf
                  mountPath: /host/etc
                - name: dbus
                  mountPath: /host/var/lib/dbus
                - name: lib-modules
                  mountPath: /lib/modules
                - name: xtables-lock
                  mountPath: /run/xtables.lock
            - name: weave-npc
              env:
                - name: HOSTNAME
                  valueFrom:
                    fieldRef:
                      apiVersion: v1
                      fieldPath: spec.nodeName
                - name: EXTRA_ARGS
                  value: --log-level=info
              image: {{.CNIImage}}
              resources:
                requests:
                  cpu: 10m
              securityContext:
                privileged: true
              volumeMounts:
                - name: xtables-lock
                  mountPath: /run/xtables.lock
            - name: weave-plugins
              command:
                - /opt/rke-tools/weave-plugins-cni.sh
              image: {{.WeaveLoopbackImage}}
              securityContext:
                privileged: true
              volumeMounts:
                - name: cni-bin
                  mountPath: /opt
          hostNetwork: true
          hostPID: true
          restartPolicy: Always
          securityContext:
            seLinuxOptions: {}
          serviceAccountName: weave-net
          tolerations:
          - operator: Exists
            effect: NoSchedule
          - operator: Exists
            effect: NoExecute
          volumes:
            - name: weavedb
              hostPath:
                path: /var/lib/weave
            - name: cni-bin
              hostPath:
                path: /opt
            - name: cni-bin2
              hostPath:
                path: /home
            - name: cni-conf
              hostPath:
                path: /etc
            - name: dbus
              hostPath:
                path: /var/lib/dbus
            - name: lib-modules
              hostPath:
                path: /lib/modules
            - name: xtables-lock
              hostPath:
                path: /run/xtables.lock
      updateStrategy:
{{if .UpdateStrategy}}
{{ toYaml .UpdateStrategy | indent 8}}
{{end}}
        type: RollingUpdate
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: weave-net
  labels:
    name: weave-net
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: weave-net
  labels:
    name: weave-net
rules:
  - apiGroups:
      - ''
    resources:
      - pods
      - namespaces
      - nodes
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - networking.k8s.io
    resources:
      - networkpolicies
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ''
    resources:
      - nodes/status
    verbs:
      - patch
      - update
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: weave-net
  labels:
    name: weave-net
roleRef:
  kind: ClusterRole
  name: weave-net
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name: weave-net
    namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: weave-net
  labels:
    name: weave-net
  namespace: kube-system
rules:
  - apiGroups:
      - ''
    resourceNames:
      - weave-net
    resources:
      - configmaps
    verbs:
      - get
      - update
  - apiGroups:
      - ''
    resources:
      - configmaps
    verbs:
      - create
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: weave-net
  labels:
    name: weave-net
  namespace: kube-system
roleRef:
  kind: Role
  name: weave-net
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name: weave-net
    namespace: kube-system
{{- end}}
`

const WeaveTemplateV116 = `
---
# This ConfigMap can be used to configure a self-hosted Weave Net installation.
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: weave-net
      namespace: kube-system
  - apiVersion: apps/v1
    kind: DaemonSet
    metadata:
      name: weave-net
      labels:
        name: weave-net
      namespace: kube-system
    spec:
      selector:
        matchLabels:
          name: weave-net
      template:
        metadata:
          annotations:
            scheduler.alpha.kubernetes.io/critical-pod: ''
            scheduler.alpha.kubernetes.io/tolerations: >-
              [{"key":"dedicated","operator":"Equal","value":"master","effect":"NoSchedule"}]
          labels:
            name: weave-net
        spec:
          affinity:
            nodeAffinity:
              requiredDuringSchedulingIgnoredDuringExecution:
                nodeSelectorTerms:
                  - matchExpressions:
                    - key: beta.kubernetes.io/os
                      operator: NotIn
                      values:
                        - windows
# Rancher specific change
{{- if .WeaveNetPriorityClassName }}
          priorityClassName: {{ .WeaveNetPriorityClassName }}
{{- end }}
{{if .NodeSelector}}
          nodeSelector:
            {{ range $k, $v := .NodeSelector }}
              {{ $k }}: "{{ $v }}"
            {{ end }}
{{end}}
          containers:
            - name: weave
              command:
                - /home/weave/launch.sh
              env:
                - name: HOSTNAME
                  valueFrom:
                    fieldRef:
                      apiVersion: v1
                      fieldPath: spec.nodeName
                - name: IPALLOC_RANGE
                  value: "{{.ClusterCIDR}}"
                {{- if .WeavePassword}}
                - name: WEAVE_PASSWORD
                  value: "{{.WeavePassword}}"
                {{- end}}
                {{- if .MTU }}
                {{- if ne .MTU 0 }}
                - name: WEAVE_MTU
                  value: "{{.MTU}}"
                {{- end }}
                {{- end }}
              image: {{.Image}}
              readinessProbe:
                httpGet:
                  host: 127.0.0.1
                  path: /status
                  port: 6784
                initialDelaySeconds: 30
              resources:
                requests:
                  cpu: 10m
              securityContext:
                privileged: true
              volumeMounts:
                - name: weavedb
                  mountPath: /weavedb
                - name: cni-bin
                  mountPath: /host/opt
                - name: cni-bin2
                  mountPath: /host/home
                - name: cni-conf
                  mountPath: /host/etc
                - name: dbus
                  mountPath: /host/var/lib/dbus
                - name: lib-modules
                  mountPath: /lib/modules
                - name: xtables-lock
                  mountPath: /run/xtables.lock
            - name: weave-npc
              env:
                - name: HOSTNAME
                  valueFrom:
                    fieldRef:
                      apiVersion: v1
                      fieldPath: spec.nodeName
                - name: EXTRA_ARGS
                  value: --log-level=info
              image: {{.CNIImage}}
              resources:
                requests:
                  cpu: 10m
              securityContext:
                privileged: true
              volumeMounts:
                - name: xtables-lock
                  mountPath: /run/xtables.lock
            - name: weave-plugins
              command:
                - /opt/rke-tools/weave-plugins-cni.sh
              image: {{.WeaveLoopbackImage}}
              securityContext:
                privileged: true
              volumeMounts:
                - name: cni-bin
                  mountPath: /opt
          hostNetwork: true
          hostPID: true
          restartPolicy: Always
          securityContext:
            seLinuxOptions: {}
          serviceAccountName: weave-net
          tolerations:
          - operator: Exists
            effect: NoSchedule
          - operator: Exists
            effect: NoExecute
          volumes:
            - name: weavedb
              hostPath:
                path: /var/lib/weave
            - name: cni-bin
              hostPath:
                path: /opt
            - name: cni-bin2
              hostPath:
                path: /home
            - name: cni-conf
              hostPath:
                path: /etc
            - name: dbus
              hostPath:
                path: /var/lib/dbus
            - name: lib-modules
              hostPath:
                path: /lib/modules
            - name: xtables-lock
              hostPath:
                path: /run/xtables.lock
      updateStrategy:
{{if .UpdateStrategy}}
{{ toYaml .UpdateStrategy | indent 8}}
{{end}}
        type: RollingUpdate
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: weave-net
  labels:
    name: weave-net
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: weave-net
  labels:
    name: weave-net
rules:
  - apiGroups:
      - ''
    resources:
      - pods
      - namespaces
      - nodes
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - networking.k8s.io
    resources:
      - networkpolicies
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ''
    resources:
      - nodes/status
    verbs:
      - patch
      - update
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: weave-net
  labels:
    name: weave-net
roleRef:
  kind: ClusterRole
  name: weave-net
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name: weave-net
    namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: weave-net
  labels:
    name: weave-net
  namespace: kube-system
rules:
  - apiGroups:
      - ''
    resourceNames:
      - weave-net
    resources:
      - configmaps
    verbs:
      - get
      - update
  - apiGroups:
      - ''
    resources:
      - configmaps
    verbs:
      - create
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: weave-net
  labels:
    name: weave-net
  namespace: kube-system
roleRef:
  kind: Role
  name: weave-net
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name: weave-net
    namespace: kube-system
{{- end}}
`

const WeaveTemplateV120 = `
---
# This ConfigMap can be used to configure a self-hosted Weave Net installation.
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: weave-net
      namespace: kube-system
  - apiVersion: apps/v1
    kind: DaemonSet
    metadata:
      name: weave-net
      labels:
        name: weave-net
      namespace: kube-system
    spec:
      # Wait 5 seconds to let pod connect before rolling next pod
      selector:
        matchLabels:
          name: weave-net
      minReadySeconds: 5
      template:
        metadata:
          annotations:
            scheduler.alpha.kubernetes.io/critical-pod: ''
            scheduler.alpha.kubernetes.io/tolerations: >-
              [{"key":"dedicated","operator":"Equal","value":"master","effect":"NoSchedule"}]
          labels:
            name: weave-net
        spec:
          affinity:
            nodeAffinity:
              requiredDuringSchedulingIgnoredDuringExecution:
                nodeSelectorTerms:
                  - matchExpressions:
                    - key: beta.kubernetes.io/os
                      operator: NotIn
                      values:
                        - windows
# Rancher specific change
{{- if .WeaveNetPriorityClassName }}
          priorityClassName: {{ .WeaveNetPriorityClassName }}
{{- end }}
{{if .NodeSelector}}
          nodeSelector:
            {{ range $k, $v := .NodeSelector }}
              {{ $k }}: "{{ $v }}"
            {{ end }}
{{end}}
          initContainers:
            - name: weave-init
              image: {{.Image}}
              command:
                - /home/weave/init.sh
              env:
              securityContext:
                privileged: true
              volumeMounts:
                - name: cni-bin
                  mountPath: /host/opt
                - name: cni-bin2
                  mountPath: /host/home
                - name: cni-conf
                  mountPath: /host/etc
                - name: lib-modules
                  mountPath: /lib/modules
                - name: xtables-lock
                  mountPath: /run/xtables.lock
                  readOnly: false
          containers:
            - name: weave
              command:
                - /home/weave/launch.sh
              env:
                - name: INIT_CONTAINER
                  value: "true"
                - name: HOSTNAME
                  valueFrom:
                    fieldRef:
                      apiVersion: v1
                      fieldPath: spec.nodeName
                - name: IPALLOC_RANGE
                  value: "{{.ClusterCIDR}}"
                {{- if .WeavePassword}}
                - name: WEAVE_PASSWORD
                  value: "{{.WeavePassword}}"
                {{- end}}
                {{- if .MTU }}
                {{- if ne .MTU 0 }}
                - name: WEAVE_MTU
                  value: "{{.MTU}}"
                {{- end }}
                {{- end }}
              image: {{.Image}}
              readinessProbe:
                httpGet:
                  host: 127.0.0.1
                  path: /status
                  port: 6784
                initialDelaySeconds: 30
              resources:
                requests:
                  cpu: 50m
              securityContext:
                privileged: true
              volumeMounts:
                - name: weavedb
                  mountPath: /weavedb
                - name: dbus
                  mountPath: /host/var/lib/dbus
                  readOnly: true
                - mountPath: /host/etc/machine-id
                  name: cni-machine-id
                  readOnly: true
                - name: xtables-lock
                  mountPath: /run/xtables.lock
                  readOnly: false
            - name: weave-npc
              env:
                - name: HOSTNAME
                  valueFrom:
                    fieldRef:
                      apiVersion: v1
                      fieldPath: spec.nodeName
                - name: EXTRA_ARGS
                  value: --log-level=info
              image: {{.CNIImage}}
              resources:
                requests:
                  cpu: 50m
              securityContext:
                privileged: true
              volumeMounts:
                - name: xtables-lock
                  mountPath: /run/xtables.lock
            - name: weave-plugins
              command:
                - /opt/rke-tools/weave-plugins-cni.sh
              image: {{.WeaveLoopbackImage}}
              securityContext:
                privileged: true
              volumeMounts:
                - name: cni-bin
                  mountPath: /opt
          hostNetwork: true
          dnsPolicy: ClusterFirstWithHostNet
          hostPID: false
          restartPolicy: Always
          securityContext:
            seLinuxOptions: {}
          serviceAccountName: weave-net
          tolerations:
          - operator: Exists
            effect: NoSchedule
          - operator: Exists
            effect: NoExecute
          volumes:
            - name: weavedb
              hostPath:
                path: /var/lib/weave
            - name: cni-bin
              hostPath:
                path: /opt
            - name: cni-bin2
              hostPath:
                path: /home
            - name: cni-conf
              hostPath:
                path: /etc
            - name: cni-machine-id
              hostPath:
                path: /etc/machine-id
            - name: dbus
              hostPath:
                path: /var/lib/dbus
            - name: lib-modules
              hostPath:
                path: /lib/modules
            - name: xtables-lock
              hostPath:
                path: /run/xtables.lock
                type: FileOrCreate
      updateStrategy:
{{if .UpdateStrategy}}
{{ toYaml .UpdateStrategy | indent 8}}
{{end}}
        type: RollingUpdate
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: weave-net
  labels:
    name: weave-net
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: weave-net
  labels:
    name: weave-net
rules:
  - apiGroups:
      - ''
    resources:
      - pods
      - namespaces
      - nodes
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - networking.k8s.io
    resources:
      - networkpolicies
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ''
    resources:
      - nodes/status
    verbs:
      - patch
      - update
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: weave-net
  labels:
    name: weave-net
roleRef:
  kind: ClusterRole
  name: weave-net
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name: weave-net
    namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: weave-net
  labels:
    name: weave-net
  namespace: kube-system
rules:
  - apiGroups:
      - ''
    resourceNames:
      - weave-net
    resources:
      - configmaps
    verbs:
      - get
      - update
  - apiGroups:
      - ''
    resources:
      - configmaps
    verbs:
      - create
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: weave-net
  labels:
    name: weave-net
  namespace: kube-system
roleRef:
  kind: Role
  name: weave-net
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name: weave-net
    namespace: kube-system
{{- end}}
`
