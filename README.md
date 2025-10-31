# Android Build for UI Testing

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-android-build-for-ui-testing?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-android-build-for-ui-testing/releases)

Builds your Android project with Gradle with the belonging AndroidTest variant.

<details>
<summary>Description</summary>

[This Step](https://github.com/bitrise-steplib/bitrise-step-android-build-for-ui-testing) generates all the APKs you need to run instrumentation tests for your Android app: both an APK from your app and the belonging test APK, for example,  `:app:assembleDemoDebug`, `:app:assembleDemoDebugAndroidTest`

 ### Configuring the Step
 1. Add the **Project Location** which is the root directory of your Android project.
 2. Set the **Module** you want to build. To see your available modules, open your project in Android Studio and go to **Project Structure** and see the list on the left.
 3. Set the **Variant** you want to build. To see your available variants, open your project in Android Studio and go to **Project Structure** and then the **variants** section.
 Under **Options**:
 4. Set the **APK location pattern**: Once the build has run, the Step finds the APK files with the given pattern.
 5. **Set the level of cache** where `all` caches build cache and dependencies, `only_deps` caches dependencies only, `none` does not cache anything.
 6. If you wish to pass any extra Gradle arguments to the gradle task, you can add those in the **Additional Gradle Arguments** input.

 ### Useful links
- [Testing with Bitrise](https://devcenter.bitrise.io/testing/testing-index/)
- [Deploying an Android app with Bitrise](https://devcenter.bitrise.io/deploy/android-deploy/android-deployment-index/)

### Related Steps
- [Android Build](https://www.bitrise.io/integrations/steps/android-build)
- [Gradle Runner](https://www.bitrise.io/integrations/steps/gradle-runner)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `project_location` | The root directory of your android project, for example, where your root build gradle file exist (also gradlew, settings.gradle, etc...) | required | `$BITRISE_SOURCE_DIR` |
| `module` | Set the module to build. Valid syntax examples: `app`, `feature:nested-module`  To see your available modules please open your project in Android Studio and go in [Project Structure] and see the list on the left.  | required |  |
| `variant` | Set the variant that you want to build. To see your available variants please open your project in Android Studio and go in [Project Structure] -> variants section.  | required |  |
| `apk_path_pattern` | Will find the APK files with the given pattern. | required | `*/build/outputs/apk/*.apk` |
| `cache_level` | `all` - will cache build cache and dependencies `only_deps` - will cache dependencies only `none` - will not cache anything | required | `only_deps` |
| `arguments` | Extra arguments passed to the gradle task |  |  |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `BITRISE_APK_PATH` | This output will include the path of the generated APK after filtering based on the filter inputs. |
| `BITRISE_TEST_APK_PATH` | This output will include the path of the generated test APK after filtering based on the filter inputs. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-android-build-for-ui-testing/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-android-build-for-ui-testing/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
