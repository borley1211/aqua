package controller

import (
	"errors"
	"testing"
)

func Test_validateConfig(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		cfg   *Config
	}{
		{
			title: "normal",
			cfg: &Config{
				Packages: []*Package{
					{
						Name:     "suzuki-shunsuke/ci-info",
						Registry: "standard",
						Version:  "v1.0.0",
					},
				},
				Registries: Registries{
					"standard": &Registry{
						Name:      "standard",
						RepoOwner: "aquaproj",
						RepoName:  "aqua-registry",
						Ref:       "v0.8.0",
						Path:      "registry.yaml",
						Type:      "github_content",
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			if err := validateConfig(d.cfg); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func Test_validateRegistries(t *testing.T) {
	t.Parallel()
	data := []struct {
		title      string
		registries Registries
		isErr      bool
	}{
		{
			title: "normal",
			registries: Registries{
				"ci-info": &Registry{
					Name:      "ci-info",
					RepoOwner: "suzuki-shunsuke",
					RepoName:  "ci-info",
					Ref:       "v1.0.0",
					Path:      "registry.yaml",
					Type:      "github_content",
				},
				"standard": &Registry{
					Name:      "standard",
					RepoOwner: "aquaproj",
					RepoName:  "aqua-registry",
					Ref:       "v0.8.0",
					Path:      "registry.yaml",
					Type:      "github_content",
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			err := validateRegistries(d.registries)
			if d.isErr {
				if !errors.Is(err, errRegistryNameIsDuplicated) {
					t.Fatal(err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func Test_validatePackages(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		pkgs  []*Package
		isErr bool
	}{
		{
			title: "normal",
			pkgs: []*Package{
				{
					Name:     "suzuki-shunsuke/cmdx",
					Registry: "standard",
				},
			},
		},
		{
			title: "duplicated",
			pkgs: []*Package{
				{
					Name:     "suzuki-shunsuke/cmdx",
					Registry: "standard",
				},
				{
					Name:     "suzuki-shunsuke/cmdx",
					Registry: "standard",
				},
			},
			isErr: false,
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			err := validatePackages(d.pkgs)
			if d.isErr {
				t.Fatal(err)
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func Test_validatePackageInfos(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		title    string
		pkgInfos PackageInfos
		isErr    bool
	}{
		{
			title: "normal",
			pkgInfos: PackageInfos{
				&PackageInfo{
					Name: "foo",
					Files: []*File{
						{
							Name: "foo",
						},
					},
				},
				&PackageInfo{
					Name: "bar",
					Files: []*File{
						{
							Name: "bar",
						},
					},
				},
			},
		},
		{
			title: "duplicated",
			isErr: true,
			pkgInfos: PackageInfos{
				&PackageInfo{
					Name: "foo",
					Files: []*File{
						{
							Name: "foo",
						},
					},
				},
				&PackageInfo{
					Name: "foo",
					Files: []*File{
						{
							Name: "foo",
						},
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			err := validatePackageInfos(d.pkgInfos)
			if d.isErr {
				if !errors.Is(err, errPkgNameMustBeUniqueInRegistry) {
					t.Fatal(err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
