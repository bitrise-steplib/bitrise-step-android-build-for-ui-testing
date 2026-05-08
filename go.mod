module github.com/bitrise-steplib/bitrise-step-android-build-for-ui-testing

go 1.24.0

require (
	github.com/bitrise-io/bitrise-build-cache-cli/v2 v2.6.0
	github.com/bitrise-io/go-android v0.0.0-20211122121320-da1559f13057
	github.com/bitrise-io/go-steputils v0.0.0-20210819160244-b3962254d553
	github.com/bitrise-io/go-utils v1.0.13
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
)

require (
	github.com/bitrise-io/go-utils/v2 v2.0.0-alpha.34 // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-version v1.3.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/shirou/gopsutil/v4 v4.25.7 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/sys v0.41.0 // indirect
)

// Pin go-utils to the v0 pseudo-version that still ships the `env` package
// (removed in v1.x). Required because the CLI dep transitively requests
// v1.0.13, and MVS would otherwise upgrade us past the version the legacy
// imports rely on. The CLI itself uses go-utils/v2 (separate module path)
// so this pin only affects the legacy go-utils import line.
replace github.com/bitrise-io/go-utils => github.com/bitrise-io/go-utils v0.0.0-20210819143908-bbd923881fab
