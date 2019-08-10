/*
 * Strava API v3
 *
 * Strava API
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type SummaryGear struct {
	// The gear's unique identifier.
	Id string `json:"id,omitempty"`
	// Resource state, indicates level of detail. Possible values: 2 -> \"summary\", 3 -> \"detail\"
	ResourceState int32 `json:"resource_state,omitempty"`
	// Whether this gear's is the owner's default one.
	Primary bool `json:"primary,omitempty"`
	// The gear's name.
	Name string `json:"name,omitempty"`
	// The distance logged with this gear.
	Distance float32 `json:"distance,omitempty"`
}
