package lib

import "sync"

// ISRO struct represents ISRO organization, it manages launch sites
// completionWaitGroup is used to wait for entire launch to complete
type ISRO struct {
	ls                  []*LS
	completionWaitGroup *sync.WaitGroup
}

// Init initializes ISRO instance
//
// Parameters:
//	numberOfLS: Number of Launch stations
//  perLSSatellites: Max number of satellites which can be launched per Launch station
//	reporter: struct instance which implemets reporter interface
func (isro *ISRO) Init(numberOfLS int, perLSSatellites int, reporter Reporter) {
	isro.ls = make([]*LS, numberOfLS)
	isro.completionWaitGroup = &sync.WaitGroup{}

	for i := 0; i < len(isro.ls); i++ {
		isro.ls[i] = &LS{}
	}

	// perLSSatellites is separate because of flexibility that, each LS can have
	// different capabilities to launch satellites.
	for i, satCount := 0, 0; i < (len(isro.ls) - 1); i, satCount = i+1, satCount+perLSSatellites {
		isro.ls[i].Init(i+1, isro.ls[i+1], 10, perLSSatellites, reporter, isro.completionWaitGroup)
	}

	last := numberOfLS - 1
	isro.ls[last].Init(last+1, isro.ls[0], 10, perLSSatellites, reporter, isro.completionWaitGroup)
}

// StartLaunch invokes entire chain of launching satellites
// Parameters:
//	numberOfSatelliteToLaunch: Number of satellites to launch
func (isro *ISRO) StartLaunch(numberOfSatellitesToLaunch int) {
	isro.completionWaitGroup.Add(1)
	isro.ls[0].LaunchBatch(0, numberOfSatellitesToLaunch)
	isro.completionWaitGroup.Wait()
}
