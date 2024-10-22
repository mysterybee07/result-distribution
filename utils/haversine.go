package utils

import (
	"math"
)

const EarthRadiusKm = 6371.0 // Earth's radius in kilometers

// Haversine function to calculate the distance between two points on the Earth
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadiusKm * c
}

// degreesToRadians converts degrees to radians
// func degreesToRadians(degrees float64) float64 {
// 	return degrees * math.Pi / 180
// }

// Example findNearestCenter function
