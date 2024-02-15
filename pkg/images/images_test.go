package images

import (
	"testing"

	"github.com/rancher/rke/types"
)

func Test_toKeep(t *testing.T) {
	type args struct {
		info types.K8sVersionInfo
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "invalid semVer",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "2.6c.0",
				MaxRancherVersion:       "",
				DeprecateRancherVersion: "",
			}},
			want:    false,
			wantErr: true,
		},
		{
			name: "deprecated < 2.7.0",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "",
				MaxRancherVersion:       "",
				DeprecateRancherVersion: "2.6.99",
			}},
			want:    false,
			wantErr: false,
		},
		{
			name: "deprecated = 2.7.0",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "",
				MaxRancherVersion:       "",
				DeprecateRancherVersion: "2.7.0",
			}},
			want:    false,
			wantErr: false,
		},
		{
			name: "2.7.0 < deprecated < 2.9.0 ",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "",
				MaxRancherVersion:       "",
				DeprecateRancherVersion: "2.8.2",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "2.7.0 < deprecated < 2.9.0 - 2",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "",
				MaxRancherVersion:       "",
				DeprecateRancherVersion: "2.8.5-alpha1",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "deprecated = 2.9.0",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "",
				MaxRancherVersion:       "",
				DeprecateRancherVersion: "2.9.0",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "min < 2.7.0 < 2.9.0 < max",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "2.6.12-alpha2",
				MaxRancherVersion:       "2.9.2",
				DeprecateRancherVersion: "",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "min < 2.7.0 < max = 2.9.0",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "2.6.12-alpha2",
				MaxRancherVersion:       "2.9.0",
				DeprecateRancherVersion: "",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "min < 2.7.0 < max < 2.9.0",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "2.6.12-alpha2",
				MaxRancherVersion:       "2.8.12-patch1",
				DeprecateRancherVersion: "",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "min = 2.7.0 < max < 2.9.0",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "2.7.0",
				MaxRancherVersion:       "2.8.5",
				DeprecateRancherVersion: "",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "min = 2.7.0 < max = 2.9.0",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "2.7.0",
				MaxRancherVersion:       "2.9.0",
				DeprecateRancherVersion: "",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "2.7.0 < min < max < 2.9.0",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "2.7.1-patch-1",
				MaxRancherVersion:       "2.8.5",
				DeprecateRancherVersion: "",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "2.9.0 < min < 2.9.99 < max",
			args: args{info: types.K8sVersionInfo{
				MinRancherVersion:       "2.9.1-patch1",
				MaxRancherVersion:       "2.10.0-alpha1",
				DeprecateRancherVersion: "",
			}},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toKeep(tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("toKeep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toKeep() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_releaseToKeep(t *testing.T) {
	type args struct {
		release map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "invalid Server Version",
			args: args{release: map[string]interface{}{
				"minChannelServerVersion": "",
				"maxChannelServerVersion": "v2.6c.7-alpha2.3",
				"version":                 ""},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "missing version",
			args: args{release: map[string]interface{}{
				"minChannelServerVersion": "v2.6.0-alpha1",
				"maxChannelServerVersion": "v2.6.99",
				"version":                 ""},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "min < 2.7.0 < 2.9.0 < max",
			args: args{release: map[string]interface{}{
				"minChannelServerVersion": "v2.6.0-alpha1",
				"maxChannelServerVersion": "v2.9.1",
				"version":                 "v1.20.15+rke2r2"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "min < 2.7.0 < max = 2.9.0",
			args: args{release: map[string]interface{}{
				"minChannelServerVersion": "v2.6.99",
				"maxChannelServerVersion": "v2.9.0",
				"version":                 "v1.20.15+rke2r2"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "min < 2.7.0 < max < 2.9.0",
			args: args{release: map[string]interface{}{
				"minChannelServerVersion": "v2.6.12-alpha1",
				"maxChannelServerVersion": "v2.8.99",
				"version":                 "v1.20.15+rke2r2"},
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "min = 2.7.0 < max < 2.9.0",
			args: args{release: map[string]interface{}{
				"minChannelServerVersion": "v2.7.0",
				"maxChannelServerVersion": "v2.8.99",
				"version":                 "v1.20.15+rke2r2"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "min = 2.7.0 < max = 2.9.0",
			args: args{release: map[string]interface{}{
				"minChannelServerVersion": "v2.7.0",
				"maxChannelServerVersion": "v2.9.0",
				"version":                 "v1.20.15+rke2r2"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "2.7.0 < min < max < 2.9.0",
			args: args{release: map[string]interface{}{
				"minChannelServerVersion": "v2.7.2-alpha1",
				"maxChannelServerVersion": "v2.8.14-patch1",
				"version":                 "v1.20.15+rke2r2"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "2.9.0 < min < 2.9.99 < max",
			args: args{release: map[string]interface{}{
				"minChannelServerVersion": "v2.9.1-alpha1",
				"maxChannelServerVersion": "v2.10.0-patch1",
				"version":                 "v1.20.15+rke2r2"},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := releaseToKeep(tt.args.release)
			if (err != nil) != tt.wantErr {
				t.Errorf("releaseToKeep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("releaseToKeep() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getImageTagMap() map[string]map[string]bool {
	return map[string]map[string]bool{
		"rancher/mirrored-coredns-coredns": {
			"1.9.0": true,
			"1.9.3": true,
			"1.9.4": true,
		},
		"rancher/mirrored-calico-cni": {
			"v3.22.0": true,
		},
	}
}

func getImageTagSlice() []string {
	return []string{
		"rancher/mirrored-coredns-coredns:1.9.0",
		"rancher/mirrored-coredns-coredns:1.9.4",
		"rancher/mirrored-coredns-coredns:1.9.3",
		"rancher/mirrored-calico-cni:v3.22.0",
	}
}

func Test_unique(t *testing.T) {
	type args struct {
		imageTag map[string]map[string]bool
		images   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				imageTag: getImageTagMap(),
				images:   getImageTagSlice(),
			},
			wantErr: false,
		},
		{
			name: "more than 2 parts",
			args: args{
				imageTag: getImageTagMap(),
				images: append(getImageTagSlice(),
					"rancher/mirrored-calico-cni:tag1:tag2"),
			},
			wantErr: true,
		},
		{
			name: "missing tags",
			args: args{
				imageTag: getImageTagMap(),
				images: append(getImageTagSlice(),
					"rancher/mirrored-calico-cni"),
			},
			wantErr: true,
		},
		{
			name: "missing rancher prefix",
			args: args{
				imageTag: getImageTagMap(),
				images: append(getImageTagSlice(),
					"mirrored-calico-cni"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unique(tt.args.imageTag, tt.args.images)
			if (err != nil) != tt.wantErr {
				t.Errorf("unique() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
