package lib

import (
	"sync"
)

// LS represents launch site.
type LS struct {
	id                int
	next              *LS
	countDownFrom     int
	reporter          Reporter
	maxLaunchCapacity int

	completionWaitGroup *sync.WaitGroup
}

// Init initializes launch site.
//
//Parameters:
//	id: id of launch station
//	next: Next launch station to inform
//	maxLaunchCapacity: maximum satellites this launch station can launch in one go
//	reporter: Reporter interface implementaton
func (l *LS) Init(id int,
	next *LS,
	countDownFrom int,
	maxLaunchCapacity int,
	reporter Reporter,
	completionWaitGroup *sync.WaitGroup) {
	l.id = id
	l.next = next
	l.countDownFrom = countDownFrom
	l.reporter = reporter
	l.maxLaunchCapacity = maxLaunchCapacity
	l.completionWaitGroup = completionWaitGroup
}

// LaunchBatch launches next batch of satellites
// Parameters:
//	completedUpto: number of satellites launched before
//	remainingSatellitesToLaunch: number of satellite remains to be launched
func (l *LS) LaunchBatch(completedUpto, remainingSatellitesToLaunch int) {
	defer l.completionWaitGroup.Done()

	startWaitGroup := sync.WaitGroup{}
	currentCount := l.countDownFrom
	numberOfSatelliteToLaunch := l.maxLaunchCapacity
	notifyNextLS := true

	if numberOfSatelliteToLaunch >= remainingSatellitesToLaunch {
		numberOfSatelliteToLaunch = remainingSatellitesToLaunch
		notifyNextLS = false
	}

	endChannel := make(chan int, numberOfSatelliteToLaunch)

	l.setupLaunch(numberOfSatelliteToLaunch, completedUpto, &currentCount, &startWaitGroup, endChannel)

	for currentCount >= 0 {
		startWaitGroup.Wait()
		l.reporter.ReportCountDown(l.id, currentCount)
		currentCount--
		startWaitGroup.Add(numberOfSatelliteToLaunch)

		for i := 0; i < numberOfSatelliteToLaunch; i++ {
			endChannel <- 0
		}
	}

	if notifyNextLS {
		l.reporter.ReportControlTransfer(l.id, l.next.id, numberOfSatelliteToLaunch)
		l.completionWaitGroup.Add(1)
		go l.next.LaunchBatch(completedUpto+numberOfSatelliteToLaunch, remainingSatellitesToLaunch-numberOfSatelliteToLaunch)
	} else {
		l.reporter.ReportCompletion(l.id, completedUpto+numberOfSatelliteToLaunch, numberOfSatelliteToLaunch)
	}
}

func (l *LS) setupLaunch(numberOfSatelliteToLaunch int,
	completedUpto int,
	currentCountPtr *int,
	startWaitGroup *sync.WaitGroup,
	endChannel chan int) {
	startWaitGroup.Add(numberOfSatelliteToLaunch)
	for i := 0; i < numberOfSatelliteToLaunch; i++ {
		go l.launchSatellite(completedUpto+i, currentCountPtr, startWaitGroup, endChannel)
	}
}

// launchSatellite utilizes WaitGroup and Channel to sync.
// LS controls WorkGroup and channel such that all go routines are at same count any given time
func (l *LS) launchSatellite(id int,
	currentCount *int,
	startWaitGroup *sync.WaitGroup,
	endChannel chan int) {
	for *currentCount >= 0 {
		l.reporter.ReportSatteliteCountDown(l.id, id, *currentCount)
		startWaitGroup.Done()
		<-endChannel
	}
}
