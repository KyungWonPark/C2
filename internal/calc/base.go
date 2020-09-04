package calc

// LinStatEle represents the statistic of a timeseries of one voxel point
type LinStatEle struct {
	avg    float32
	stddev float32
}

// Config holds wokrder related configuration
type Config struct {
	NumLoader     int
	NumComputer   int
	NumBufferSize int
}
