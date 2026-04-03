// Package mlbapi provides a Go client for the public MLB Stats API.
//
// The package offers two levels of access:
//
//   - Typed helper functions for common workflows such as schedules, standings,
//     lookups, highlights, linescores, box scores, rosters, and player summaries.
//   - Raw endpoint access through Client, Get, GetForce, and generated Call...
//     methods that map closely to the underlying Stats API endpoints.
//
// The package also includes stable TeamID constants and current team metadata so
// callers can work with well-known MLB clubs without making a lookup request.
package mlbapi
