package main

import (
	"flag"
	"fmt"

	"github.com/parthdesai/isro_algo_demo/lib"
)

func main() {

	numberOfLSPtr := flag.Int("numberOfLS", 2, "Number of launch station")
	numberOfSatellitesToLaunchPtr := flag.Int("satToLaunch", 500, "Total number of satellite to launch")
	perLSLaunchCapacityPtr := flag.Int("perLSLaunchCapacity", 4, "Per Launch station launch capacity")
	flag.Parse()

	if *numberOfLSPtr <= 0 || *numberOfSatellitesToLaunchPtr <= 0 || *perLSLaunchCapacityPtr <= 0 {
		panic(fmt.Errorf("Parameters should be greater than zero"))
	}

	isro := lib.ISRO{}
	isro.Init(*numberOfLSPtr, *perLSLaunchCapacityPtr, &lib.PrintReporter{})
	isro.StartLaunch(*numberOfSatellitesToLaunchPtr)
}
