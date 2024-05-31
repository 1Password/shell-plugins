package aiven

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessToken: "/TKopDBC2nHPi1TTa4kMakmTC3m+bHoC2UIUOvEjUTJeC/WxjnpsCqlrmR8VXgSx/hJUUQ8jnd+6gylz8RtnUYovkiiDq9pP/y54SNmqa1AMvR1AnYXevuvUWupZBDujRYkjyQvdu+QsUPtGOppmKc7ymZa1otRqGlFdVu5jhh3/7j8RcxsM4z0WdUSCRnBt3lL3nNRQE5diRE8xbkWBfUCZu1kY6XSzDQcrTko6AvkLY3wdvbfLfENL/l2pp6WmNVsftW5XjxihjL1O+9Klg1wuYxko40CcseL8W3up7HcCSVtgVZSMHjF9LQMLcmd0U9CVpngGY/fcM89MO0Wz4BzNXDyhhxx4ox+LdgriSm13p0o5hF6pVwSFRHxlVEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AIVEN_AUTH_TOKEN": "/TKopDBC2nHPi1TTa4kMakmTC3m+bHoC2UIUOvEjUTJeC/WxjnpsCqlrmR8VXgSx/hJUUQ8jnd+6gylz8RtnUYovkiiDq9pP/y54SNmqa1AMvR1AnYXevuvUWupZBDujRYkjyQvdu+QsUPtGOppmKc7ymZa1otRqGlFdVu5jhh3/7j8RcxsM4z0WdUSCRnBt3lL3nNRQE5diRE8xbkWBfUCZu1kY6XSzDQcrTko6AvkLY3wdvbfLfENL/l2pp6WmNVsftW5XjxihjL1O+9Klg1wuYxko40CcseL8W3up7HcCSVtgVZSMHjF9LQMLcmd0U9CVpngGY/fcM89MO0Wz4BzNXDyhhxx4ox+LdgriSm13p0o5hF6pVwSFRHxlVEXAMPLE",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"AIVEN_AUTH_TOKEN": "/TKopDBC2nHPi1TTa4kMakmTC3m+bHoC2UIUOvEjUTJeC/WxjnpsCqlrmR8VXgSx/hJUUQ8jnd+6gylz8RtnUYovkiiDq9pP/y54SNmqa1AMvR1AnYXevuvUWupZBDujRYkjyQvdu+QsUPtGOppmKc7ymZa1otRqGlFdVu5jhh3/7j8RcxsM4z0WdUSCRnBt3lL3nNRQE5diRE8xbkWBfUCZu1kY6XSzDQcrTko6AvkLY3wdvbfLfENL/l2pp6WmNVsftW5XjxihjL1O+9Klg1wuYxko40CcseL8W3up7HcCSVtgVZSMHjF9LQMLcmd0U9CVpngGY/fcM89MO0Wz4BzNXDyhhxx4ox+LdgriSm13p0o5hF6pVwSFRHxlVEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessToken: "/TKopDBC2nHPi1TTa4kMakmTC3m+bHoC2UIUOvEjUTJeC/WxjnpsCqlrmR8VXgSx/hJUUQ8jnd+6gylz8RtnUYovkiiDq9pP/y54SNmqa1AMvR1AnYXevuvUWupZBDujRYkjyQvdu+QsUPtGOppmKc7ymZa1otRqGlFdVu5jhh3/7j8RcxsM4z0WdUSCRnBt3lL3nNRQE5diRE8xbkWBfUCZu1kY6XSzDQcrTko6AvkLY3wdvbfLfENL/l2pp6WmNVsftW5XjxihjL1O+9Klg1wuYxko40CcseL8W3up7HcCSVtgVZSMHjF9LQMLcmd0U9CVpngGY/fcM89MO0Wz4BzNXDyhhxx4ox+LdgriSm13p0o5hF6pVwSFRHxlVEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.config/aiven/aiven-credentials.json": plugintest.LoadFixture(t, "aiven-credentials.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessToken: "/EXAMPLE/WxjnpsCqlrmR8VXgSx/hJUUQ8jnd+6gylz8RtnUYovkiiDhh3/7j8RcxsM4z0WdUSCRnBt3lL3nNRQE5diRE8xbkWBfUCZu1kY6XSzDQcrTko6AvkLY3wdvbfLfENL/l2pp6WmNVsftW5XjxihjL1O+9Klg1wuYxko40CcseL8W3up7EXAMPLE=",
					},
				},
			},
		},
	})
}
