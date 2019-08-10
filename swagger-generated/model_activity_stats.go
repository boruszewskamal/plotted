/*
 * Strava API v3
 *
 * Strava API
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

// A set of rolled-up statistics and totals for an athlete
type ActivityStats struct {
	// The longest distance ridden by the athlete.
	BiggestRideDistance float64 `json:"biggest_ride_distance,omitempty"`
	// The highest climb ridden by the athlete.
	BiggestClimbElevationGain float64 `json:"biggest_climb_elevation_gain,omitempty"`
	// The recent (last 4 weeks) ride stats for the athlete.
	RecentRideTotals *ActivityTotal `json:"recent_ride_totals,omitempty"`
	// The recent (last 4 weeks) run stats for the athlete.
	RecentRunTotals *ActivityTotal `json:"recent_run_totals,omitempty"`
	// The recent (last 4 weeks) swim stats for the athlete.
	RecentSwimTotals *ActivityTotal `json:"recent_swim_totals,omitempty"`
	// The year to date ride stats for the athlete.
	YtdRideTotals *ActivityTotal `json:"ytd_ride_totals,omitempty"`
	// The year to date run stats for the athlete.
	YtdRunTotals *ActivityTotal `json:"ytd_run_totals,omitempty"`
	// The year to date swim stats for the athlete.
	YtdSwimTotals *ActivityTotal `json:"ytd_swim_totals,omitempty"`
	// The all time ride stats for the athlete.
	AllRideTotals *ActivityTotal `json:"all_ride_totals,omitempty"`
	// The all time run stats for the athlete.
	AllRunTotals *ActivityTotal `json:"all_run_totals,omitempty"`
	// The all time swim stats for the athlete.
	AllSwimTotals *ActivityTotal `json:"all_swim_totals,omitempty"`
}
