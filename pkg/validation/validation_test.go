package main

import (
	"testing"
	"strings"

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
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "message missing id",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"type":     "warning",
						"severity": "high",
						"summary":  "Test",
						"message":  "Description",
					},
				},
			},
			shouldErr: true,
			errMsg:    "missing required field 'id'",
		},
		{
			name: "message with empty id",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "",
						"type":     "warning",
						"severity": "high",
						"summary":  "Test",
						"message":  "Description",
					},
				},
			},
			shouldErr: true,
			errMsg:    "missing required field 'id'",
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
			name: "message with invalid severity",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"type":     "warning",
						"severity": "critical",
						"summary":  "Test",
						"message":  "Description",
					},
				},
			},
			shouldErr: true,
			errMsg:    "invalid severity 'critical'",
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
			name: "message without optional severity field",
			release: map[string]interface{}{
				"version": "v1.20.0+rke2r1",
				"messages": []interface{}{
					map[string]interface{}{
						"id":       "msg-1",
						"type":     "info",
						"summary":  "Informational",
						"message":  "Just informational",
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "message with url field is rejected",
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
			shouldErr: true,
			errMsg:    "unexpected field 'url'",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seenIDs := make(map[string]bool)
			err := validateMessages(tt.release, seenIDs)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing '%s', got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateMessageDuplicateIDs(t *testing.T) {
	tests := []struct {
		name      string
		releases  []map[string]interface{}
		shouldErr bool
		errMsg    string
	}{
		{
			name: "no duplicate IDs across releases",
			releases: []map[string]interface{}{
				{
					"version": "v1.20.0+rke2r1",
					"messages": []interface{}{
						map[string]interface{}{
							"id":       "msg-1",
							"type":     "warning",
							"severity": "high",
							"summary":  "Issue 1",
							"message":  "Description 1",
						},
					},
				},
				{
					"version": "v1.21.0+rke2r1",
					"messages": []interface{}{
						map[string]interface{}{
							"id":       "msg-2",
							"type":     "info",
							"severity": "low",
							"summary":  "Issue 2",
							"message":  "Description 2",
						},
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "duplicate IDs across releases are rejected",
			releases: []map[string]interface{}{
				{
					"version": "v1.20.0+rke2r1",
					"messages": []interface{}{
						map[string]interface{}{
							"id":       "msg-dup",
							"type":     "warning",
							"severity": "high",
							"summary":  "Issue",
							"message":  "Description",
						},
					},
				},
				{
					"version": "v1.21.0+rke2r1",
					"messages": []interface{}{
						map[string]interface{}{
							"id":       "msg-dup",
							"type":     "info",
							"severity": "low",
							"summary":  "Same issue",
							"message":  "Another description",
						},
					},
				},
			},
			shouldErr: true,
			errMsg:    "duplicate id 'msg-dup'",
		},
		{
			name: "duplicate IDs within same release are rejected",
			releases: []map[string]interface{}{
				{
					"version": "v1.20.0+rke2r1",
					"messages": []interface{}{
						map[string]interface{}{
							"id":       "msg-dup",
							"type":     "warning",
							"severity": "high",
							"summary":  "Issue 1",
							"message":  "Description 1",
						},
						map[string]interface{}{
							"id":       "msg-dup",
							"type":     "info",
							"severity": "low",
							"summary":  "Issue 2",
							"message":  "Description 2",
						},
					},
				},
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
					t.Errorf("expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing '%s', got: %v", tt.errMsg, err)
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
	types := []MessageType{MessageTypeInfo, MessageTypeWarning, MessageTypeCritical}
	for _, msgType := range types {
		release := map[string]interface{}{
			"version": "v1.20.0+rke2r1",
			"messages": []interface{}{
				map[string]interface{}{
					"id":       "msg-test",
					"type":     string(msgType),
					"severity": "medium",
					"summary":  "Test message",
					"message":  "Test description",
				},
			},
		}
		seenIDs := make(map[string]bool)
		err := validateMessages(release, seenIDs)
		if err != nil {
			t.Errorf("expected no error for type '%s', got: %v", msgType, err)
		}
	}
}

func TestValidateMessageAllSeverities(t *testing.T) {
	severities := []MessageSeverity{MessageSeverityLow, MessageSeverityMedium, MessageSeverityHigh}
	for _, severity := range severities {
		release := map[string]interface{}{
			"version": "v1.20.0+rke2r1",
			"messages": []interface{}{
				map[string]interface{}{
					"id":       "msg-test-" + string(severity),
					"type":     "warning",
					"severity": string(severity),
					"summary":  "Test message",
					"message":  "Test description",
				},
			},
		}
		seenIDs := make(map[string]bool)
		err := validateMessages(release, seenIDs)
		if err != nil {
			t.Errorf("expected no error for severity '%s', got: %v", severity, err)
		}
	}
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
	seenIDs := make(map[string]bool)
	err := validateMessages(release, seenIDs)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
}
