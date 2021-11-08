package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-android/cache"
	"github.com/bitrise-io/go-android/gradle"
	utilscache "github.com/bitrise-io/go-steputils/cache"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-steputils/tools"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/env"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	shellquote "github.com/kballard/go-shellquote"
)

const (
	apkEnvKey     = "BITRISE_APK_PATH"
	testApkEnvKey = "BITRISE_TEST_APK_PATH"
	testSuffix    = "AndroidTest"
)

// Configs ...
type Configs struct {
	ProjectLocation string `env:"project_location,dir"`
	APKPathPattern  string `env:"apk_path_pattern"`
	Variant         string `env:"variant,required"`
	Module          string `env:"module,required"`
	Arguments       string `env:"arguments"`
	CacheLevel      string `env:"cache_level,opt[none,only_deps,all]"`
	DeployDir       string `env:"BITRISE_DEPLOY_DIR,dir"`
}

var cmdFactory = command.NewFactory(env.NewRepository())
var logger = log.NewLogger(false)

func getArtifacts(gradleProject gradle.Project, started time.Time, pattern string, includeModule bool) (artifacts []gradle.Artifact, err error) {
	artifacts, err = gradleProject.FindArtifacts(started, pattern, includeModule)
	if err != nil {
		return
	}
	if len(artifacts) == 0 {
		if !started.IsZero() {
			logger.Warnf("No artifacts found with pattern: %s that has modification time after: %s", pattern, started)
			logger.Warnf("Retrying without modtime check....")
			fmt.Println()
			return getArtifacts(gradleProject, time.Time{}, pattern, includeModule)
		}
		logger.Warnf("No artifacts found with pattern: %s without modtime check", pattern)
	}

	return
}

func exportArtifacts(artifacts []gradle.Artifact, deployDir string) ([]string, error) {
	var paths []string
	for _, artifact := range artifacts {
		exists, err := pathutil.IsPathExists(filepath.Join(deployDir, artifact.Name))
		if err != nil {
			return nil, fmt.Errorf("failed to check path, error: %v", err)
		}

		artifactName := filepath.Base(artifact.Path)

		if exists {
			timestamp := time.Now().Format("20060102150405")
			ext := filepath.Ext(artifact.Name)
			name := strings.TrimSuffix(filepath.Base(artifact.Name), ext)
			artifact.Name = fmt.Sprintf("%s-%s%s", name, timestamp, ext)
		}

		logger.Printf("  Export [ %s => $BITRISE_DEPLOY_DIR/%s ]", artifactName, artifact.Name)

		if err := artifact.Export(deployDir); err != nil {
			logger.Warnf("failed to export artifact (%s), error: %v", artifact.Path, err)
			continue
		}

		paths = append(paths, filepath.Join(deployDir, artifact.Name))
	}
	return paths, nil
}

func filterVariants(module, variant string, variantsMap gradle.Variants) (gradle.Variants, error) {
	filteredVariants := gradle.Variants{}
	var testVariant string
	var appVariant string
	for _, v := range variantsMap[module] {
		if strings.EqualFold(v, variant) {
			appVariant = v
		} else if strings.EqualFold(v, variant+testSuffix) {
			testVariant = v
		}
	}

	if appVariant == "" {
		return nil, fmt.Errorf("variant: %s not found in %s module", variant, module)
	}

	if testVariant == "" {
		return nil, fmt.Errorf("variant: %s not found in %s module", variant+testSuffix, module)
	}

	filteredVariants[module] = []string{appVariant, testVariant}
	return filteredVariants, nil
}

// androidTestVariantPairs returns (build - AndroidTest) variant pairs
func androidTestVariantPairs(module string, variantsMap gradle.Variants) (gradle.Variants, error) {
	appVariants := gradle.Variants{}
	testVariants := gradle.Variants{}
	for m, variants := range variantsMap {
		for _, v := range variants {
			if strings.HasSuffix(strings.ToLower(v), strings.ToLower(testSuffix)) {
				testVariants[m] = append(testVariants[m], v)
			} else {
				appVariants[m] = append(appVariants[m], v)
			}
		}
	}

	variantPairs := gradle.Variants{}
	for m, appVariant := range appVariants {
		for _, variant := range appVariant {
			if sliceutil.IsStringInSlice(variant+testSuffix, testVariants[m]) {
				variantPairs[m] = append(variantPairs[m], []string{variant, variant + testSuffix}...)
			}
		}
	}

	return variantPairs, nil
}

func mainE(config Configs) error {
	gradleProject, err := gradle.NewProject(config.ProjectLocation, cmdFactory)
	if err != nil {
		return fmt.Errorf("Failed to open project, error: %s", err)
	}

	buildTask := gradleProject.GetTask("assemble")

	args, err := shellquote.Split(config.Arguments)
	if err != nil {
		return fmt.Errorf("Failed to parse arguments, error: %s", err)
	}

	logger.Infof("Variants:")

	variants, err := buildTask.GetVariants(args...)
	if err != nil {
		return fmt.Errorf("Failed to fetch variants, error: %s", err)
	}

	variantPairs, err := androidTestVariantPairs(config.Module, variants)
	if err != nil {
		return fmt.Errorf("Failed to find variant pairs (build and AndroidTest variant), error: %s", err)
	}

	filteredVariants, err := filterVariants(config.Module, config.Variant, variants)
	if err != nil {
		// List all the variants if there is an error
		for module, variants := range variants {
			logger.Printf("%s:", module)
			for _, variant := range variants {
				logger.Printf("- %s", variant)
			}
		}
		fmt.Println()

		return fmt.Errorf("Failed to find buildable variants, error: %s", err)
	}

	// List the variants only which has (Build - AndroidTest) variant pair
	for module, variants := range variantPairs {
		logger.Printf("%s:", module)
		for _, variant := range variants {
			if sliceutil.IsStringInSlice(variant, filteredVariants[module]) {
				logger.Donef("✓ %s", variant)
				continue
			}
			logger.Printf("- %s", variant)
		}
		fmt.Println()
	}

	started := time.Now()

	logger.Infof("Run build:")
	buildCommand := buildTask.GetCommand(filteredVariants, args...)

	logger.Donef("$ " + buildCommand.PrintableCommandArgs())
	fmt.Println()

	if err := buildCommand.Run(); err != nil {
		return fmt.Errorf("Build task failed, error: %v", err)
	}

	fmt.Println()

	logger.Infof("Export APKs:")
	fmt.Println()

	apks, err := getArtifacts(gradleProject, started, config.APKPathPattern, false)
	if err != nil {
		return fmt.Errorf("failed to find apks, error: %v", err)
	}

	exportedArtifactPaths, err := exportArtifacts(apks, config.DeployDir)
	if err != nil {
		return fmt.Errorf("Failed to export artifact: %v", err)
	}

	var exportedAppArtifact string
	var exportedTestArtifact string
	for _, pth := range exportedArtifactPaths {
		if strings.HasSuffix(strings.ToLower(path.Base(pth)), strings.ToLower("AndroidTest.apk")) {
			exportedTestArtifact = pth
		} else {
			exportedAppArtifact = pth
		}
	}

	if exportedAppArtifact == "" {
		return fmt.Errorf("Could not find the exported app APK")
	}

	if exportedTestArtifact == "" {
		return fmt.Errorf("Could not find the exported test APK")
	}

	fmt.Println()
	if err := tools.ExportEnvironmentWithEnvman(apkEnvKey, exportedAppArtifact); err != nil {
		return fmt.Errorf("Failed to export environment variable: %s", apkEnvKey)
	}
	logger.Printf("  Env    [ $%s = $BITRISE_DEPLOY_DIR/%s ]", apkEnvKey, filepath.Base(exportedAppArtifact))

	if err := tools.ExportEnvironmentWithEnvman(testApkEnvKey, exportedTestArtifact); err != nil {
		return fmt.Errorf("Failed to export environment variable: %s", apkEnvKey)
	}
	logger.Printf("  Env    [ $%s = $BITRISE_DEPLOY_DIR/%s ]", testApkEnvKey, filepath.Base(exportedTestArtifact))

	var paths, sep string
	for _, path := range exportedArtifactPaths {
		paths += sep + "$BITRISE_DEPLOY_DIR/" + filepath.Base(path)
		sep = "| \\\n" + strings.Repeat(" ", 11)
	}
	fmt.Println()

	return nil
}

func failf(s string, args ...interface{}) {
	logger.Errorf(s, args...)
	os.Exit(1)
}

func main() {
	var config Configs

	if err := stepconf.Parse(&config); err != nil {
		failf("Couldn't create step config: %v", err)
	}

	stepconf.Print(config)

	fmt.Println()

	if err := mainE(config); err != nil {
		failf("%s", err)
	}

	fmt.Println()
	logger.Infof("Collecting cache:")
	if warning := cache.Collect(config.ProjectLocation, utilscache.Level(config.CacheLevel), cmdFactory); warning != nil {
		logger.Warnf("%s", warning)
	}

	logger.Donef("  Done")
}
