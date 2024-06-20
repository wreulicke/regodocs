package test

import rego.v1

# test

# This rule is ignored because it has no metadata
deny_foo contains msg if {
	input.path == "/etc/shadow"

	msg := "deny_test"
}

allow if {
	input.path == "/etc/hosts"
}
