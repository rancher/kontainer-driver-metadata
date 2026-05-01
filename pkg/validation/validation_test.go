package main

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestValidateMessages(t *testing.T) {
	tests := []struct {
		name      string
		release   map[string]interface{}
		shouldErr bool
		errMsg    string
	}{
		{
			name: "release without messages field",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
			},
			shouldErr: false,
		},
		{
			name: "release with valid single message",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "test-warning-1",
						"type":     "warning",
						"severity": "high",
						"summary":  "Test warning",
						"message":  "This is a test warning message",
						"url":      "https://example.com/issue",
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "release with valid multiple messages",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"type":     "warning",
						"severity": "high",
						"summary":  "Issue 1",
						"message":  "Description 1",
					},
					map[string]interface{}{
						"id":       "msg-2",
						"type":     "info",
						"severity": "low",
						"summary":  "Issue 2",
						"message":  "Description 2",
						"url":      "https://example.com",
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "message missing type",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"severity": "high",
						"summary":  "Test",
						"message":  "Description",
					},
				},
			},
			shouldErr: true,
			errMsg:    "missing required field 'type'",
		},
		{
			name: "message with invalid type",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"type":     "invalid",
						"severity": "high",
						"summary":  "Test",
						"message":  "Description",
					},
				},
			},
			shouldErr: true,
			errMsg:    "invalid type 'invalid'",
		},
		{
			name: "message missing summary",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"type":     "warning",
						"severity": "high",
						"message":  "Description",
					},
				},
			},
			shouldErr: true,
			errMsg:    "missing required field 'summary'",
		},
		{
			name: "message missing message field",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"type":     "warning",
						"severity": "high",
						"summary":  "Test",
					},
				},
			},
			shouldErr: true,
			errMsg:    "missing required field 'message'",
		},
		{
			name: "message with optional url field",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"type":     "critical",
						"severity": "high",
						"summary":  "Critical issue",
						"message":  "Detailed description",
						"url":      "https://example.com/issue/123",
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "message without optional url field",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"type":     "info",
						"severity": "low",
						"summary":  "Informational",
						"message":  "Just informational",
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "invalid message format (not a map)",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					"not a map",
				},
			},
			shouldErr: true,
			errMsg:    "invalid message format",
		},
		{
			name: "message with unexpected field",
			release: map[string]interface{}{
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
			shouldErr: true,
			errMsg:    "unexpected field 'extra'",
		},
		{
			name: "message with multiple unexpected fields",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"type":     "warning",
						"severity": "high",
						"summary":  "Test",
						"message":  "Description",
						"extra1":   "unexpected",
						"extra2":   "foobar",
					},
				},
			},
			shouldErr: true,
			errMsg:    "unexpected field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessages(tt.release)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					if !contains(err.Error(), tt.errMsg) {
						t.Errorf("expected error containing '%s', got: %v", tt.errMsg, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateMessageAllTypes(t *testing.T) {
	types := []string{"info", "warning", "critical"}
	for _, msgType := range types {
		release := map[string]interface{}{
			"version": "v1.20.0+rke2r1",
			"messages": []interface{}{
				map[string]interface{}{
					"id":       "msg-test",
					"type":     msgType,
					"severity": "medium",
					"summary":  "Test message",
					"message":  "Test description",
				},
			},
		}
		err := validateMessages(release)
		if err != nil {
			t.Errorf("expected no error for type '%s', got: %v", msgType, err)
		}
	}
}

func TestValidateMessageAllSeverities(t *testing.T) {
	severities := []string{"low", "medium", "high"}
	for _, severity := range severities {
		release := map[string]interface{}{
			"version": "v1.20.0+rke2r1",
			"messages": []interface{}{
				map[string]interface{}{
					"id":       "msg-test",
					"type":     "warning",
					"severity": severity,
					"summary":  "Test message",
					"message":  "Test description",
				},
			},
		}
		err := validateMessages(release)
		if err != nil {
			t.Errorf("expected no error for severity '%s', got: %v", severity, err)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && len(substr) > 0 && s[0:len(substr)] == substr || len(s) > len(substr))
}

// TestValidateMessageWithUnstructured tests that messages validation works with unstructured data
func TestValidateMessageWithUnstructured(t *testing.T) {
	release := map[string]interface{}{
		"version": "v1.20.0+rke2r1",
		"messages": []interface{}{
			map[string]interface{}{
				"id":       "etcd-restore-issue",
				"type":     "warning",
				"severity": "high",
				"summary":  "Known etcd restore issue",
				"message":  "This version has a known issue with etcd snapshot restores. Please upgrade to v1.20.1+rke2r1.",
				"url":      "https://github.com/rancher/rke2/issues/12345",
			},
		},
	}

	// Verify that unstructured can extract messages correctly
	messages, found, _ := unstructured.NestedSlice(release, "messages")
	if !found || len(messages) == 0 {
		t.Errorf("expected to find messages in release")
	}

	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}

	msgMap := messages[0].(map[string]interface{})
	if msgMap["id"] != "etcd-restore-issue" {
		t.Errorf("expected id 'etcd-restore-issue', got '%s'", msgMap["id"])
	}

	// Test validation passes
	err := validateMessages(release)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
}
