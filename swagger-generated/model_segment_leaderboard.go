/*
 * Strava API v3
 *
 * Strava API
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

// A
type SegmentLeaderboard struct {
	// The total number of entries for this leaderboard
	EntryCount int32 `json:"entry_count,omitempty"`
	// Deprecated, use entry_count
	EffortCount int32 `json:"effort_count,omitempty"`
	KomType string `json:"kom_type,omitempty"`
	Entries []SegmentLeaderboardEntry `json:"entries,omitempty"`
}
