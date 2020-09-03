package calc

// LinBuffer has fMRI bold signal of 13362 voxels * 600 seconds
type LinBuffer [13362][600]float32

// MatBuffer hold 13362 by 13362 size matrix
type MatBuffer [13362][13362]float32

// LinStatEle represents the statistic of a timeseries of one voxel point
type LinStatEle struct {
	avg    float32
	stddev float32
}

// LinStat holds statistics for 13362 voxel points
type LinStat [13362]LinStatEle
