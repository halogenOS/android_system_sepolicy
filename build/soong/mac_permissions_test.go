// SPDX-FileCopyrightText: The halogenOS Project
// SPDX-License-Identifier: Apache-2.0

package selinux

import (
	"testing"

	"github.com/google/blueprint/proptools"

	"android/soong/android"
)

func TestMacPermissionsKeyDirs(t *testing.T) {
	t.Parallel()

	result := android.GroupFixturePreparers(
		android.PrepareForTestWithAndroidBuildComponents,
		android.PrepareForTestWithHostTools("insertkeys"),
		android.FixtureRegisterWithContext(func(ctx android.RegistrationContext) {
			ctx.RegisterModuleType("mac_permissions", macPermissionsFactory)
		}),
		android.FixtureModifyProductVariables(func(variables android.FixtureProductVariables) {
			variables.MainlineSepolicyDevCertificates = proptools.StringPtr("vendor/mainline")
		}),
		android.FixtureMergeMockFs(android.MockFS{
			"build/make/target/product/security/platform.x509.pem": nil,
			"vendor/mainline/networkstack.x509.pem":                nil,
		}),
		android.FixtureWithRootAndroidBp(`
			mac_permissions {
				name: "plat_mac_permissions.xml",
				keys: ["keys.conf"],
				srcs: ["mac_permissions.xml"],
			}
		`),
	).RunTest(t)

	macPermissions := result.ModuleForTests(t, "plat_mac_permissions.xml", "android_common").
		Output("plat_mac_permissions.xml")
	android.AssertStringDoesContain(t, "command", macPermissions.RuleParams.Command,
		"DEFAULT_SYSTEM_DEV_CERTIFICATE=${keys_dir} MAINLINE_SEPOLICY_DEV_CERTIFICATES=vendor/mainline "+
			"MAINLINE_BLUETOOTH_SEPOLICY_DEV_CERTIFICATES=${keys_dir} MAINLINE_NFC_SEPOLICY_DEV_CERTIFICATES=${keys_dir} ")
	implicits := android.PathsRelativeToTop(macPermissions.Implicits)
	android.AssertStringListContains(t, "implicits", implicits, "out/soong/test-keys/platform.x509.pem")
	android.AssertStringListContains(t, "implicits", implicits, "vendor/mainline/networkstack.x509.pem")
}
