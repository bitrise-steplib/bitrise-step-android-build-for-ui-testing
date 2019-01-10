package main

import (
	"reflect"
	"testing"

	"github.com/bitrise-tools/go-android/gradle"
)

func Test_androidTestVariantPairs(t *testing.T) {
	tests := []struct {
		name        string
		module      string
		variantsMap gradle.Variants
		want        gradle.Variants
		wantErr     bool
	}{
		{
			name:        "one AndroidTest for app module",
			module:      "app",
			variantsMap: oneAndroidTestVariantsMap(),
			want:        wantOneAndroidTestForApp(),
			wantErr:     false,
		},
		{
			name:        "one AndroidTest for another_app module",
			module:      "another_app",
			variantsMap: oneAndroidTestVariantsMap(),
			want:        wantOneAndroidTestForAnotherApp(),
			wantErr:     false,
		},
		{
			name:        "no AndroidTest for another_app module",
			module:      "another_app",
			variantsMap: noAndroidTestVariantsMap(),
			want:        map[string][]string{},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := androidTestVariantPairs(tt.module, tt.variantsMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("androidTestVariantPairs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("androidTestVariantPairs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func oneAndroidTestVariantsMap() gradle.Variants {
	variantsMap := gradle.Variants{}
	variantsMap["app"] = []string{"AndroidTest", "Debug", "Demo", "Release", "DemoDebug", "DemoDebugAndroidTest", "DemoDebugUnitTest", "DemoRelease", "DemoReleaseUnitTest"}
	variantsMap["another_app"] = []string{"AndroidTest", "AnotherDemo", "Debug", "Demo", "AnotherDemoDebug", "AnotherDemoDebugAndroidTest", "AnotherDemoDebugUnitTest", "AnotherDemoRelease"}
	return variantsMap
}

func noAndroidTestVariantsMap() gradle.Variants {
	variantsMap := gradle.Variants{}
	variantsMap["another_app"] = []string{"AndroidTest", "AnotherDemo", "Debug", "Demo", "AnotherDemoDebug", "AnotherDemoDebugUnitTest", "AnotherDemoRelease"}
	return variantsMap
}

func wantOneAndroidTestForApp() gradle.Variants {
	want := gradle.Variants{}
	want["app"] = []string{"DemoDebug", "DemoDebugAndroidTest"}
	return want
}

func wantOneAndroidTestForAnotherApp() gradle.Variants {
	want := gradle.Variants{}
	want["another_app"] = []string{"AnotherDemoDebug", "AnotherDemoDebugAndroidTest"}
	return want
}

func Test_filterVariants(t *testing.T) {
	tests := []struct {
		name        string
		module      string
		variant     string
		variantsMap gradle.Variants
		want        gradle.Variants
		wantErr     bool
	}{
		{
			name:        "variant filter for app module",
			module:      "app",
			variant:     "DemoDebug",
			variantsMap: oneAndroidTestVariantsMap(),
			want:        wantVariantFilterForApp(),
			wantErr:     false,
		},
		{
			name:        "variant filter for app module",
			module:      "another_app",
			variant:     "AnotherDemoDebug",
			variantsMap: oneAndroidTestVariantsMap(),
			want:        wantVariantFilterForAnotherApp(),
			wantErr:     false,
		},
		{
			name:        "variant filter for app module",
			module:      "app",
			variant:     "DemoDebug",
			variantsMap: appVariantNotFoundVariantsMap(),
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "variant filter for app module",
			module:      "app",
			variant:     "DemoDebug",
			variantsMap: testVariantNotFoundVariantsMap(),
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterVariants(tt.module, tt.variant, tt.variantsMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterVariants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterVariants() = %v, want %v", got, tt.want)
			}
		})
	}
}

func appVariantNotFoundVariantsMap() gradle.Variants {
	variantsMap := gradle.Variants{}
	variantsMap["app"] = []string{"AndroidTest", "Debug", "Demo", "Release", "DemoDebugAndroidTest", "DemoDebugUnitTest", "DemoRelease", "DemoReleaseUnitTest"}
	variantsMap["another_app"] = []string{"AndroidTest", "AnotherDemo", "Debug", "Demo", "AnotherDemoDebug", "AnotherDemoDebugAndroidTest", "AnotherDemoDebugUnitTest", "AnotherDemoRelease"}
	return variantsMap
}

func testVariantNotFoundVariantsMap() gradle.Variants {
	variantsMap := gradle.Variants{}
	variantsMap["app"] = []string{"AndroidTest", "Debug", "Demo", "Release", "DemoDebug", "DemoDebugUnitTest", "DemoRelease", "DemoReleaseUnitTest"}
	variantsMap["another_app"] = []string{"AndroidTest", "AnotherDemo", "Debug", "Demo", "AnotherDemoDebug", "AnotherDemoDebugAndroidTest", "AnotherDemoDebugUnitTest", "AnotherDemoRelease"}
	return variantsMap
}
func wantVariantFilterForApp() gradle.Variants {
	want := gradle.Variants{}
	want["app"] = []string{"DemoDebug", "DemoDebugAndroidTest"}
	return want
}

func wantVariantFilterForAnotherApp() gradle.Variants {
	want := gradle.Variants{}
	want["another_app"] = []string{"AnotherDemoDebug", "AnotherDemoDebugAndroidTest"}
	return want
}
