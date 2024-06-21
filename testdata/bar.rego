# METADATA
# description: |
#   This package contains test policies.
package bar

import rego.v1

# test

# This is a test policy that denies access to the /etc/shadow file.
# METADATA
# description: |
#   This is a test policy that denies access to the /etc/shadow file."
#   next line
#   next line
deny_test contains msg if {
	input.path == "/etc/shadow"

	msg := "deny_test"
}
