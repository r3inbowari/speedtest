package location

type Location struct {
	Lat float64
	Lon float64
}

// lat lon
var (
	Beijing      = &Location{Lat: 39.5600, Lon: 116.2000}
	Paris        = &Location{Lat: 48.8600, Lon: 2.3390}
	SanFrancisco = &Location{Lat: 37.7687, Lon: -122.4754}
	London       = &Location{Lat: 51.5063, Lon: -0.1201}
	Moscow       = &Location{Lat: 55.7520, Lon: 37.6175}
	HongKong     = &Location{Lat: 22.3207, Lon: 114.1689}
)
