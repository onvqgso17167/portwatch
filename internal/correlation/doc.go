// Package correlation groups scanner diff results into named incidents.
//
// When multiple ports open or close within a short time window on the same
// network, they are assigned a shared incident ID. This allows downstream
// alerting and audit components to distinguish isolated single-port changes
// from coordinated or cascading activity.
//
// Usage:
//
//	c := correlation.New(10 * time.Second)
//	event := c.Add("tcp", openedResults)
//	fmt.Println(event.ID, event.Correlated)
package correlation
