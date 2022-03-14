package dom

type WorkerLocation struct {
	Recnum    int64 // auto-increasing record number
	Worker    SUID
	Stamp     Time
	Longitude float64
	Latitude  float64
	Speed     int
}
