# isro_algo_demo
Program to simulate ISRO satellite launch, using various golang sync constructs. 

## Instruction to install and run this program:
```bash
go get github.com/parthdesai/isro_algo_demo
go install github.com/parthdesai/isro_algo_demo
cd $GOBIN
./isro_algo_demo
```

## Command line options supported:
```bash
./isro_algo_demo -help
Usage of ./isro_algo_demo:
  -numberOfLS int
    	Number of launch station (default 2)
  -perLSLaunchCapacity int
    	Per Launch station launch capacity (default 4)
  -satToLaunch int
    	Total number of satellite to launch (default 500)
```

## To run tests:
```bash
cd $GOPATH/src/github.com/parthdesai/isro_algo_demo/lib
go test --v
```

