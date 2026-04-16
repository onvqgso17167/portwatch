// Package tag provides a thread-safe registry for associating
// human-readable labels with port numbers.
//
// Tags are loaded from a JSON file or set programmatically and are
// used by the reporter and notifier to enrich alert messages with
// context such as "http", "postgres", or "internal-service".
package tag
