package demand

var driverLocationFarWest = LatLon{
	Lat: 40.734462,
	Lon: -73.927502,
}

var driverLocationMiddle = LatLon{
	Lat: 40.769987,
	Lon: -73.591087,
}

var longBeachJobs = []*Job{
	{
		ID: "long-beach-2-river-head",
		PickUp: LatLon{
			Lat: 40.588773,
			Lon: -73.664954,
		},
		DropOff: LatLon{
			Lat: 40.904901,
			Lon: -72.743725,
		},
	},
	{
		ID: "long-beach-2-west-babylon",
		PickUp: LatLon{
			Lat: 40.588773,
			Lon: -73.664954,
		},
		DropOff: LatLon{
			Lat: 40.731495,
			Lon: -73.354590,
		},
	},
	{
		ID: "long-beach-2-garden-city",
		PickUp: LatLon{
			Lat: 40.588773,
			Lon: -73.664954,
		},
		DropOff: LatLon{
			Lat: 40.733576,
			Lon: -73.631995,
		},
	},
	{
		ID: "long-beach-2-valley-stream",
		PickUp: LatLon{
			Lat: 40.588773,
			Lon: -73.664954,
		},
		DropOff: LatLon{
			Lat: 40.674234,
			Lon: -73.725378,
		},
	},
}

var hicksvilleJobs = []*Job{
	{
		ID: "hicksville-2-river-head",
		PickUp: LatLon{
			Lat: 40.790631,
			Lon: -73.534741,
		},
		DropOff: LatLon{
			Lat: 40.904901,
			Lon: -72.743725,
		},
	},
	{
		ID: "hicksville-2-west-babylon",
		PickUp: LatLon{
			Lat: 40.790631,
			Lon: -73.534741,
		},
		DropOff: LatLon{
			Lat: 40.731495,
			Lon: -73.354590,
		},
	},
}

var greatNeckJobs = []*Job{
	{
		ID: "greatNeck-2-river-head",
		PickUp: LatLon{
			Lat: 40.786626,
			Lon: -73.718512,
		},
		DropOff: LatLon{
			Lat: 40.904901,
			Lon: -72.743725,
		},
	},
	{
		ID: "greatNeck-2-west-babylon",
		PickUp: LatLon{
			Lat: 40.786626,
			Lon: -73.718512,
		},
		DropOff: LatLon{
			Lat: 40.731495,
			Lon: -73.354590,
		},
	},
}
