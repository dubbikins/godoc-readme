package template_functions

import (
	"bytes"
	"fmt"
	"go/doc"

	"golang.org/x/tools/go/packages"
)

// NOTE(Alert): Use this alert to provide additional information

// TIP(Alert): Use this alert to provide helpful tips to users

// IMPORTANT(Alert): Use this alert to highlight important information

// WARNING(Alert): Use this alert to warn users about potential issues

// CAUTION(Alert): Use this alert to caution users about serious issues

// Alert returns a function that, given the name of a target, returns a string representing the alerts for that target
// Can be used in a template by calling `{{ Alert "target_name" }}` where `target_name` is the name of the package, a Type, Func, Var, or Const in the package.
// Alerts are rendered AFTER the doc comment for the target by default. Provide your own templates to modify this behavior.
func Alert(pkg *packages.Package, notes map[string][]*doc.Note) func(string) string {

	_alert_types := []string{"NOTE", "WARNING", "IMPORTANT", "CAUTION", "TIP"}
	var package_alerts = map[string]map[string][]*doc.Note{}
	for _, alert_type := range _alert_types {
		// Iterate over all of the alerts for a given type (NOTE, WARNING, IMPORTANT, CAUTION, TIP) and add them to the alerts for a given key
		// That way in the function that we return can access all of the given alerts for the key without having to
		// iterate over all of the alerts for the entire package
		for _, _alert_for_type := range notes[alert_type] {
			//if a key doesn't exist for the alert's UID in the alerts map, create a new entry in the map
			if _, found := package_alerts[_alert_for_type.UID]; !found {
				package_alerts[_alert_for_type.UID] = map[string][]*doc.Note{}
			}
			if _, found := package_alerts[_alert_for_type.UID][alert_type]; !found {
				package_alerts[_alert_for_type.UID][alert_type] = []*doc.Note{}
			}
			package_alerts[_alert_for_type.UID][alert_type] = append(package_alerts[_alert_for_type.UID][alert_type], _alert_for_type)
		}
	}
	return func(key string) string {

		var buf = bytes.NewBuffer(nil)
		alert_types := []string{"NOTE", "WARNING", "IMPORTANT", "CAUTION", "TIP"}
		for i, alert_type := range alert_types {
			//Check if notes exist for any of the github markdown alert types
			alerts, found := package_alerts[key][alert_type]
			if !found {
				continue
			}
			var header_written bool
			for _, alert := range alerts {
				if alert.UID == key {
					if !header_written {
						buf.WriteString("\n") // Alerts should always start on a new line
						buf.WriteString(fmt.Sprintf(">[!%s]\n", alert_type))
						header_written = true
					}
					buf.WriteString(fmt.Sprintf(">%s", alert.Body))
				}
			}
			buf.WriteString("\n")
			if i == len(alert_types)-1 {
				buf.WriteString("\n") // Alerts should always be followed by a blank line
			}
		}
		return buf.String()
	}
}


