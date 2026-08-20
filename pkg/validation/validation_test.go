package validation

import (
	"strings"
	"testing"
)

func releaseWithMessages(messages ...Message) Releases {
	return Releases{
		Version:  "v1.20.0+rke2r1",
		Messages: messages,
	}
}

func TestValidateMessages(t *testing.T) {
	tests := []struct {
		name      string
		release   Releases
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "release without messages field",
			release:   Releases{Version: "v1.20.0+rke2r1"},
			shouldErr: false,
		},
		{
			name: "release with valid single message",
			release: releaseWithMessages(Message{
				ID:       "test-warning-1",
				Type:     MessageTypeWarning,
				Severity: MessageSeverityHigh,
				Summary:  "Test warning",
				Message:  "This is a test warning message",
			}),
			shouldErr: false,
		},
		{
			name: "release with valid multiple messages",
			release: releaseWithMessages(
				Message{
					ID:       "msg-1",
					Type:     MessageTypeWarning,
					Severity: MessageSeverityHigh,
					Summary:  "Issue 1",
					Message:  "Description 1",
				},
				Message{
					ID:       "msg-2",
					Type:     MessageTypeInfo,
					Severity: MessageSeverityLow,
					Summary:  "Issue 2",
					Message:  "Description 2",
				},
			),
			shouldErr: false,
		},
		{
			name: "message missing id",
			release: releaseWithMessages(Message{
				Type:     MessageTypeWarning,
				Severity: MessageSeverityHigh,
				Summary:  "Test",
				Message:  "Description",
			}),
			shouldErr: true,
			errMsg:    "missing required field 'id'",
		},
		{
			name: "message with empty id",
			release: releaseWithMessages(Message{
				ID:       "",
				Type:     MessageTypeWarning,
				Severity: MessageSeverityHigh,
				Summary:  "Test",
				Message:  "Description",
			}),
			shouldErr: true,
			errMsg:    "missing required field 'id'",
		},
		{
			name: "message missing type",
			release: releaseWithMessages(Message{
				ID:       "msg-1",
				Severity: MessageSeverityHigh,
				Summary:  "Test",
				Message:  "Description",
			}),
			shouldErr: true,
			errMsg:    "missing required field 'type'",
		},
		{
			name: "message with invalid type",
			release: releaseWithMessages(Message{
				ID:       "msg-1",
				Type:     MessageType("invalid"),
				Severity: MessageSeverityHigh,
				Summary:  "Test",
				Message:  "Description",
			}),
			shouldErr: true,
			errMsg:    "invalid type 'invalid'",
		},
		{
			name: "message with invalid severity",
			release: releaseWithMessages(Message{
				ID:       "msg-1",
				Type:     MessageTypeWarning,
				Severity: MessageSeverity("critical"),
				Summary:  "Test",
				Message:  "Description",
			}),
			shouldErr: true,
			errMsg:    "invalid severity 'critical'",
		},
		{
			name: "message missing summary",
			release: releaseWithMessages(Message{
				ID:       "msg-1",
				Type:     MessageTypeWarning,
				Severity: MessageSeverityHigh,
				Message:  "Description",
			}),
			shouldErr: true,
			errMsg:    "missing required field 'summary'",
		},
		{
			name: "message missing message field",
			release: releaseWithMessages(Message{
				ID:       "msg-1",
				Type:     MessageTypeWarning,
				Severity: MessageSeverityHigh,
				Summary:  "Test",
			}),
			shouldErr: true,
			errMsg:    "missing required field 'message'",
		},
		{
			name: "message without optional severity field",
			release: releaseWithMessages(Message{
				ID:      "msg-1",
				Type:    MessageTypeInfo,
				Summary: "Informational",
				Message: "Just informational",
			}),
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seenIDs := make(map[string]bool)
			err := validateMessages(tt.release, seenIDs)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Fatalf("expected error containing %q, got: %v", tt.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateMessageDuplicateIDs(t *testing.T) {
	tests := []struct {
		name      string
		releases  []Releases
		shouldErr bool
		errMsg    string
	}{
		{
			name: "no duplicate IDs across releases",
			releases: []Releases{
				releaseWithMessages(Message{
					ID:       "msg-1",
					Type:     MessageTypeWarning,
					Severity: MessageSeverityHigh,
					Summary:  "Issue 1",
					Message:  "Description 1",
				}),
				releaseWithMessages(Message{
					ID:       "msg-2",
					Type:     MessageTypeInfo,
					Severity: MessageSeverityLow,
					Summary:  "Issue 2",
					Message:  "Description 2",
				}),
			},
			shouldErr: false,
		},
		{
			name: "duplicate IDs across releases are rejected",
			releases: []Releases{
				releaseWithMessages(Message{
					ID:       "msg-dup",
					Type:     MessageTypeWarning,
					Severity: MessageSeverityHigh,
					Summary:  "Issue",
					Message:  "Description",
				}),
				releaseWithMessages(Message{
					ID:       "msg-dup",
					Type:     MessageTypeInfo,
					Severity: MessageSeverityLow,
					Summary:  "Same issue",
					Message:  "Another description",
				}),
			},
			shouldErr: true,
			errMsg:    "duplicate id 'msg-dup'",
		},
		{
			name: "duplicate IDs within same release are rejected",
			releases: []Releases{
				releaseWithMessages(
					Message{
						ID:       "msg-dup",
						Type:     MessageTypeWarning,
						Severity: MessageSeverityHigh,
						Summary:  "Issue 1",
						Message:  "Description 1",
					},
					Message{
						ID:       "msg-dup",
						Type:     MessageTypeInfo,
						Severity: MessageSeverityLow,
						Summary:  "Issue 2",
						Message:  "Description 2",
					},
				),
			},
			shouldErr: true,
			errMsg:    "duplicate id 'msg-dup'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seenIDs := make(map[string]bool)
			var err error
			for _, release := range tt.releases {
				err = validateMessages(release, seenIDs)
				if err != nil {
					break
				}
			}
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Fatalf("expected error containing %q, got: %v", tt.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateMessageAllTypes(t *testing.T) {
	types := []MessageType{MessageTypeInfo, MessageTypeWarning, MessageTypeCritical}
	for _, msgType := range types {
		release := releaseWithMessages(Message{
			ID:       "msg-test",
			Type:     msgType,
			Severity: MessageSeverityMedium,
			Summary:  "Test message",
			Message:  "Test description",
		})
		seenIDs := make(map[string]bool)
		if err := validateMessages(release, seenIDs); err != nil {
			t.Fatalf("expected no error for type %q, got: %v", msgType, err)
		}
	}
}

func TestValidateMessageAllSeverities(t *testing.T) {
	severities := []MessageSeverity{MessageSeverityLow, MessageSeverityMedium, MessageSeverityHigh}
	for _, severity := range severities {
		release := releaseWithMessages(Message{
			ID:       "msg-test-" + string(severity),
			Type:     MessageTypeWarning,
			Severity: severity,
			Summary:  "Test message",
			Message:  "Test description",
		})
		seenIDs := make(map[string]bool)
		if err := validateMessages(release, seenIDs); err != nil {
			t.Fatalf("expected no error for severity %q, got: %v", severity, err)
		}
	}
}

func TestLoadReleasesStrictValidation(t *testing.T) {
	t.Run("loads valid typed releases", func(t *testing.T) {
		loaded, err := loadReleases(map[string]any{
			"releases": []any{
				map[string]any{
					"version": "v1.25.11",
					"charts": map[string]any{
						"rke2-cilium": map[string]any{
							"repo":    "rancher-rke2-charts",
							"version": "1.0.0",
						},
					},
					"featureVersions": map[string]any{
						"encryption-key-rotation": "v1.25.11",
					},
					"messages": []any{
						map[string]any{
							"id":       "msg-1",
							"type":     "warning",
							"severity": "high",
							"summary":  "Test",
							"message":  "Description",
						},
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(loaded.Releases) != 1 {
			t.Fatalf("expected 1 release, got %d", len(loaded.Releases))
		}
		if loaded.Releases[0].Charts == nil || loaded.Releases[0].Charts["rke2-cilium"] == nil {
			t.Fatalf("expected typed chart data to load")
		}
		if len(loaded.Releases[0].Messages) != 1 {
			t.Fatalf("expected 1 message, got %d", len(loaded.Releases[0].Messages))
		}
	})

	t.Run("loads k3s args keys used in data.json", func(t *testing.T) {
		loaded, err := loadReleases(map[string]any{
			"releases": []any{
				map[string]any{
					"version": "v1.33.0+k3s1",
					"serverArgs": map[string]any{
						"default-local-storage-path": map[string]any{"type": "string", "default": "/var/lib/rancher/k3s/storage"},
						"disable-apiserver":          map[string]any{"type": "boolean", "default": false},
						"disable-controller-manager": map[string]any{"type": "boolean", "default": false},
						"disable-etcd":               map[string]any{"type": "boolean", "default": false},
						"disable-network-policy":     map[string]any{"type": "boolean", "default": false},
						"etcd-s3-bucket-lookup-type": map[string]any{"type": "enum", "default": "auto", "options": []any{"auto", "dns", "path"}},
						"flannel-backend":            map[string]any{"type": "enum", "options": []any{"none", "vxlan", "ipsec", "host-gw", "wireguard-native"}},
						"flannel-ipv6-masq":          map[string]any{"type": "boolean"},
						"helm-controller-arg":        map[string]any{"type": "string"},
						"kine-tls":                   map[string]any{"type": "boolean"},
						"prime":                      map[string]any{"type": "boolean"},
						"secrets-encryption":         map[string]any{"type": "boolean", "default": false},
						"secrets-encryption-provider": map[string]any{
							"type":    "enum",
							"default": "aescbc",
							"options": []any{"aescbc", "secretbox"},
						},
						"write-kubeconfig":      map[string]any{"type": "string"},
						"write-kubeconfig-mode": map[string]any{"type": "string"},
					},
					"agentArgs": map[string]any{
						"disable-apiserver-lb": map[string]any{"type": "boolean"},
						"docker":               map[string]any{"type": "boolean", "default": false},
						"flannel-cni-conf":     map[string]any{"type": "string"},
						"flannel-conf":         map[string]any{"type": "string"},
						"flannel-iface":        map[string]any{"type": "string"},
						"image-service-endpoint": map[string]any{
							"type": "string",
						},
						"node-external-dns": map[string]any{"type": "array"},
						"node-internal-dns": map[string]any{"type": "array"},
						"pause-image":       map[string]any{"type": "string"},
						"prefer-bundled-bin": map[string]any{
							"type": "boolean",
						},
						"snapshotter":   map[string]any{"type": "string"},
						"vpn-auth":      map[string]any{"type": "string"},
						"vpn-auth-file": map[string]any{"type": "string"},
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(loaded.Releases) != 1 {
			t.Fatalf("expected 1 release, got %d", len(loaded.Releases))
		}
		if loaded.Releases[0].ServerArgs == nil || loaded.Releases[0].AgentArgs == nil {
			t.Fatalf("expected serverArgs and agentArgs to decode")
		}
		if loaded.Releases[0].ServerArgs.DisableApiserver.Type != "boolean" {
			t.Fatalf("expected disable-apiserver type to decode")
		}
		if loaded.Releases[0].ServerArgs.EtcdS3BucketLookupType.Default != "auto" {
			t.Fatalf("expected etcd-s3-bucket-lookup-type default to decode")
		}
		if loaded.Releases[0].ServerArgs.DefaultLocalStoragePath.Default != "/var/lib/rancher/k3s/storage" {
			t.Fatalf("expected default-local-storage-path default to decode")
		}
		if loaded.Releases[0].ServerArgs.Prime.Type != "boolean" {
			t.Fatalf("expected prime type to decode")
		}
		if loaded.Releases[0].ServerArgs.HelmControllerArg.Type != "string" {
			t.Fatalf("expected helm-controller-arg type to decode")
		}
		if loaded.Releases[0].ServerArgs.WriteKubeconfig.Type != "string" {
			t.Fatalf("expected write-kubeconfig type to decode")
		}
		if loaded.Releases[0].ServerArgs.WriteKubeconfigMode.Type != "string" {
			t.Fatalf("expected write-kubeconfig-mode type to decode")
		}
		if loaded.Releases[0].AgentArgs.Docker.Type != "boolean" {
			t.Fatalf("expected docker type to decode")
		}
		if loaded.Releases[0].AgentArgs.NodeExternalDNS.Type != "array" {
			t.Fatalf("expected node-external-dns type to decode")
		}
	})

	t.Run("rejects unexpected message field", func(t *testing.T) {
		_, err := loadReleases(map[string]any{
			"releases": []any{
				map[string]any{
					"version": "v1.20.0+rke2r1",
					"messages": []any{
						map[string]any{
							"id":       "msg-1",
							"type":     "warning",
							"severity": "high",
							"summary":  "Test",
							"message":  "Description",
							"extra":    "unexpected field",
						},
					},
				},
			},
		})
		if err == nil {
			t.Fatalf("expected error but got none")
		}
		if !strings.Contains(err.Error(), "extra") {
			t.Fatalf("expected error to mention unexpected field, got: %v", err)
		}
	})

	t.Run("rejects invalid message format", func(t *testing.T) {
		_, err := loadReleases(map[string]any{
			"releases": []any{
				map[string]any{
					"version": "v1.20.0+rke2r1",
					"messages": []any{
						"not a map",
					},
				},
			},
		})
		if err == nil {
			t.Fatalf("expected error but got none")
		}
	})
}
