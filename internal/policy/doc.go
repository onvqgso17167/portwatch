// Package policy provides rule-based evaluation of port change events.
//
// A Policy is a prioritised list of Rules. Each Rule specifies:
//   - A set of ports it applies to (empty means all ports).
//   - An Action: "alert", "ignore", or "log".
//   - An optional time window (TimeStart / TimeEnd in HH:MM 24-hour format)
//     during which the rule is active. Windows that span midnight are supported.
//
// Rules are evaluated in order; the first matching rule wins. If no rule
// matches, the default action is ActionAlert.
//
// Policies can be loaded from a JSON file via Load, or constructed
// programmatically with New.
package policy
