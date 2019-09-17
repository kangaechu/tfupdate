package tfupdate

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
)

func TestNewTerraformUpdater(t *testing.T) {
	cases := []struct {
		version string
		want    Updater
		ok      bool
	}{
		{
			version: "0.12.7",
			want: &TerraformUpdater{
				version: "0.12.7",
			},
			ok: true,
		},
		{
			version: "",
			want:    nil,
			ok:      false,
		},
	}

	for _, tc := range cases {
		got, err := NewTerraformUpdater(tc.version)
		if tc.ok && err != nil {
			t.Errorf("NewTerraformUpdater() with version = %s returns unexpected err: %+v", tc.version, err)
		}

		if !tc.ok && err == nil {
			t.Errorf("NewTerraformUpdater() with version = %s expects to return an error, but no error", tc.version)
		}

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("NewTerraformUpdater() with version = %s returns %#v, but want = %#v", tc.version, got, tc.want)
		}
	}
}

func TestUpdateTerraform(t *testing.T) {
	cases := []struct {
		src     string
		version string
		want    string
		ok      bool
	}{
		{
			src: `
terraform {
  required_version = "0.12.6"
}
`,
			version: "0.12.7",
			want: `
terraform {
  required_version = "0.12.7"
}
`,
			ok: true,
		},
		{
			src: `
terraform {
  required_providers {
    null = "2.1.1"
  }
}
`,
			version: "0.12.7",
			want: `
terraform {
  required_providers {
    null = "2.1.1"
  }
}
`,
			ok: true,
		},
		{
			src: `
provider "aws" {
  version = "2.11.0"
  region  = "ap-northeast-1"
}
`,
			version: "0.12.7",
			want: `
provider "aws" {
  version = "2.11.0"
  region  = "ap-northeast-1"
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		u := &TerraformUpdater{
			version: tc.version,
		}
		f, diags := hclwrite.ParseConfig([]byte(tc.src), "", hcl.Pos{Line: 1, Column: 1})
		if diags.HasErrors() {
			t.Fatalf("unexpected diagnostics: %s", diags)
		}

		err := u.Update(f)
		if tc.ok && err != nil {
			t.Errorf("Update() with src = %s, version = %s returns unexpected err: %+v", tc.src, tc.version, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("Update() with src = %s, version = %s expects to return an error, but no error", tc.src, tc.version)
		}

		got := string(hclwrite.Format(f.BuildTokens(nil).Bytes()))
		if got != tc.want {
			t.Errorf("Update() with src = %s, version = %s returns %s, but want = %s", tc.src, tc.version, got, tc.want)
		}
	}
}
