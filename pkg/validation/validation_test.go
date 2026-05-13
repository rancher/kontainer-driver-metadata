package main

import (
	"strings"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

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
				ID:       strPtr("test-warning-1"),
				Type:     strPtr("warning"),
				Severity: strPtr("high"),
				Summary:  strPtr("Test warning"),
				Message:  strPtr("This is a test warning message"),
			}),
			shouldErr: false,
		},
		{
			name: "release with valid multiple messages",
			release: releaseWithMessages(
				Message{
					ID:       strPtr("msg-1"),
					Type:     strPtr("warning"),
					Severity: strPtr("high"),
					Summary:  strPtr("Issue 1"),
					Message:  strPtr("Description 1"),
				},
				Message{
					ID:       strPtr("msg-2"),
					Type:     strPtr("info"),
					Severity: strPtr("low"),
					Summary:  strPtr("Issue 2"),
					Message:  strPtr("Description 2"),
				},
			),
			shouldErr: false,
		},
		{
			name: "message missing id",
			release: releaseWithMessages(Message{
				Type:     strPtr("warning"),
				Severity: strPtr("high"),
				Summary:  strPtr("Test"),
				Message:  strPtr("Description"),
			}),
			shouldErr: true,
			errMsg:    "missing required field 'id'",
		},
		{
			name: "message with empty id",
			release: releaseWithMessages(Message{
				ID:       strPtr(""),
				Type:     strPtr("warning"),
				Severity: strPtr("high"),
				Summary:  strPtr("Test"),
				Message:  strPtr("Description"),
			}),
			shouldErr: true,
			errMsg:    "missing required field 'id'",
		},
		{
			name: "message missing type",
			release: releaseWithMessages(Message{
				ID:       strPtr("msg-1"),
				Severity: strPtr("high"),
				Summary:  strPtr("Test"),
				Message:  strPtr("Description"),
			}),
			shouldErr: true,
			errMsg:    "missing required field 'type'",
		},
		{
			name: "message with invalid type",
			release: releaseWithMessages(Message{
				ID:       strPtr("msg-1"),
				Type:     strPtr("invalid"),
				Severity: strPtr("high"),
				Summary:  strPtr("Test"),
				Message:  strPtr("Description"),
			}),
			shouldErr: true,
			errMsg:    "invalid type 'invalid'",
		},
		{
			name: "message with invalid severity",
			release: releaseWithMessages(Message{
				ID:       strPtr("msg-1"),
				Type:     strPtr("warning"),
				Severity: strPtr("critical"),
				Summary:  strPtr("Test"),
				Message:  strPtr("Description"),
			}),
			shouldErr: true,
			errMsg:    "invalid severity 'critical'",
		},
		{
			name: "message missing summary",
			release: releaseWithMessages(Message{
				ID:       strPtr("msg-1"),
				Type:     strPtr("warning"),
				Severity: strPtr("high"),
				Message:  strPtr("Description"),
			}),
			shouldErr: true,
			errMsg:    "missing required field 'summary'",
		},
		{
			name: "message missing message field",
			release: releaseWithMessages(Message{
				ID:       strPtr("msg-1"),
				Type:     strPtr("warning"),
				Severity: strPtr("high"),
				Summary:  strPtr("Test"),
			}),
			shouldErr: true,
			errMsg:    "missing required field 'message'",
		},
		{
			name: "message without optional severity field",
			release: releaseWithMessages(Message{
				ID:      strPtr("msg-1"),
				Type:    strPtr("info"),
				Summary: strPtr("Informational"),
				Message: strPtr("Just informational"),
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
					ID:       strPtr("msg-1"),
					Type:     strPtr("warning"),
					Severity: strPtr("high"),
					Summary:  strPtr("Issue 1"),
					Message:  strPtr("Description 1"),
				}),
				releaseWithMessages(Message{
					ID:       strPtr("msg-2"),
					Type:     strPtr("info"),
					Severity: strPtr("low"),
					Summary:  strPtr("Issue 2"),
					Message:  strPtr("Description 2"),
				}),
			},
			shouldErr: false,
		},
		{
			name: "duplicate IDs across releases are rejected",
			releases: []Releases{
				releaseWithMessages(Message{
					ID:       strPtr("msg-dup"),
					Type:     strPtr("warning"),
					Severity: strPtr("high"),
					Summary:  strPtr("Issue"),
					Message:  strPtr("Description"),
				}),
				releaseWithMessages(Message{
					ID:       strPtr("msg-dup"),
					Type:     strPtr("info"),
					Severity: strPtr("low"),
					Summary:  strPtr("Same issue"),
					Message:  strPtr("Another description"),
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
						ID:       strPtr("msg-dup"),
						Type:     strPtr("warning"),
						Severity: strPtr("high"),
						Summary:  strPtr("Issue 1"),
						Message:  strPtr("Description 1"),
					},
					Message{
						ID:       strPtr("msg-dup"),
						Type:     strPtr("info"),
						Severity: strPtr("low"),
						Summary:  strPtr("Issue 2"),
						Message:  strPtr("Description 2"),
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
			ID:       strPtr("msg-test"),
			Type:     strPtr(string(msgType)),
			Severity: strPtr("medium"),
			Summary:  strPtr("Test message"),
			Message:  strPtr("Test description"),
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
			ID:       strPtr("msg-test-" + string(severity)),
			Type:     strPtr("warning"),
			Severity: strPtr(string(severity)),
			Summary:  strPtr("Test message"),
			Message:  strPtr("Test description"),
		})
		seenIDs := make(map[string]bool)
		if err := validateMessages(release, seenIDs); err != nil {
			t.Fatalf("expected no error for severity %q, got: %v", severity, err)
		}
	}
}

func TestLoadReleasesStrictValidation(t *testing.T) {
	t.Run("loads valid typed releases", func(t *testing.T) {
		loaded, err := loadReleases(map[string]interface{}{
			"releases": []interface{}{
				map[string]interface{}{
					"version": "v1.25.11",
					"featureVersions": map[string]interface{}{
						"encryption-key-rotation": "enabled",
					},
					"charts": map[string]interface{}{
						"rke2-cilium": map[string]interface{}{
							"repo":    "rancher-rke2-charts",
							"version": "1.0.0",
						},
					},
					"messages": []interface{}{
						map[string]interface{}{
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
		if loaded.Releases[0].FeatureVersions == nil || loaded.Releases[0].FeatureVersions.EncryptionKeyRotation == nil {
			t.Fatalf("expected typed featureVersions data to load")
		}
		if len(loaded.Releases[0].Messages) != 1 {
			t.Fatalf("expected 1 message, got %d", len(loaded.Releases[0].Messages))
		}
	})

	t.Run("loads k3s args keys used in data.json", func(t *testing.T) {
		loaded, err := loadReleases(map[string]interface{}{
			"releases": []interface{}{
				map[string]interface{}{
					"version": "v1.33.0+k3s1",
					"serverArgs": map[string]interface{}{
						"default-local-storage-path": map[string]interface{}{"type": "string", "default": "/var/lib/rancher/k3s/storage"},
						"disable-apiserver":          map[string]interface{}{"type": "boolean", "default": false},
						"disable-controller-manager": map[string]interface{}{"type": "boolean", "default": false},
						"disable-etcd":               map[string]interface{}{"type": "boolean", "default": false},
						"disable-network-policy":     map[string]interface{}{"type": "boolean", "default": false},
						"etcd-s3-bucket-lookup-type": map[string]interface{}{"type": "enum", "default": "auto", "options": []interface{}{"auto", "dns", "path"}},
						"flannel-backend":            map[string]interface{}{"type": "enum", "options": []interface{}{"none", "vxlan", "ipsec", "host-gw", "wireguard-native"}},
						"flannel-ipv6-masq":          map[string]interface{}{"type": "boolean"},
						"kine-tls":                   map[string]interface{}{"type": "boolean"},
						"secrets-encryption":         map[string]interface{}{"type": "boolean", "default": false},
						"secrets-encryption-provider": map[string]interface{}{
							"type":    "enum",
							"default": "aescbc",
							"options": []interface{}{"aescbc", "secretbox"},
						},
					},
					"agentArgs": map[string]interface{}{
						"disable-apiserver-lb": map[string]interface{}{"type": "boolean"},
						"docker":               map[string]interface{}{"type": "boolean", "default": false},
						"flannel-cni-conf":     map[string]interface{}{"type": "string"},
						"flannel-conf":         map[string]interface{}{"type": "string"},
						"flannel-iface":        map[string]interface{}{"type": "string"},
						"image-service-endpoint": map[string]interface{}{
							"type": "string",
						},
						"node-external-dns": map[string]interface{}{"type": "array"},
						"node-internal-dns": map[string]interface{}{"type": "array"},
						"pause-image":       map[string]interface{}{"type": "string"},
						"prefer-bundled-bin": map[string]interface{}{
							"type": "boolean",
						},
						"snapshotter":   map[string]interface{}{"type": "string"},
						"vpn-auth":      map[string]interface{}{"type": "string"},
						"vpn-auth-file": map[string]interface{}{"type": "string"},
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
		if loaded.Releases[0].AgentArgs.Docker.Type != "boolean" {
			t.Fatalf("expected docker type to decode")
		}
		if loaded.Releases[0].AgentArgs.NodeExternalDNS.Type != "array" {
			t.Fatalf("expected node-external-dns type to decode")
		}
	})

	t.Run("rejects unexpected message field", func(t *testing.T) {
		_, err := loadReleases(map[string]interface{}{
			"releases": []interface{}{
				map[string]interface{}{
					"version": "v1.20.0+rke2r1",
					"messages": []interface{}{
						map[string]interface{}{
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
		_, err := loadReleases(map[string]interface{}{
			"releases": []interface{}{
				map[string]interface{}{
					"version": "v1.20.0+rke2r1",
					"messages": []interface{}{
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

func TestValidateEncryptedKeyRotation(t *testing.T) {
	t.Run("skips versions before threshold", func(t *testing.T) {
		release := Releases{Version: "v1.25.10"}
		if err := validateEncryptedKeyRotation(release); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("requires featureVersions at threshold", func(t *testing.T) {
		release := Releases{Version: "v1.25.11"}
		if err := validateEncryptedKeyRotation(release); err == nil {
			t.Fatalf("expected error but got none")
		}
	})

	t.Run("requires encryption-key-rotation field", func(t *testing.T) {
		release := Releases{
			Version:         "v1.25.11",
			FeatureVersions: &FeatureVersions{},
		}
		if err := validateEncryptedKeyRotation(release); err == nil {
			t.Fatalf("expected error but got none")
		}
	})

	t.Run("accepts encryption-key-rotation field", func(t *testing.T) {
		release := Releases{
			Version: "v1.25.11",
			FeatureVersions: &FeatureVersions{
				EncryptionKeyRotation: strPtr("enabled"),
			},
		}
		if err := validateEncryptedKeyRotation(release); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})
}
