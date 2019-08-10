/*
 * Strava API v3
 *
 * Strava API
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ActivityZone struct {
	Score int32 `json:"score,omitempty"`
	DistributionBuckets *TimedZoneDistribution `json:"distribution_buckets,omitempty"`
	Type_ string `json:"type,omitempty"`
	SensorBased bool `json:"sensor_based,omitempty"`
	Points int32 `json:"points,omitempty"`
	CustomZones bool `json:"custom_zones,omitempty"`
	Max int32 `json:"max,omitempty"`
}
