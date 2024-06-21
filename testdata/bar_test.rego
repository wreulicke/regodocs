# this file will be ignored
package bar_test

import rego.v1

deny_test_xxx contains msg if {
	not input.test

	msg := "bar"
}
