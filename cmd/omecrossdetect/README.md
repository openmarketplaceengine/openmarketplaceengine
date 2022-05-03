# Command line detect tollgate crossings tool

Two args required
```shell

-l string
    	path to locations csv file
-t string
    	path to tollgates yaml file
```

Run command example 
```shell
go run omecrossdetect.go -t ../../pkg/tollgate/yaml/tollgates.yaml -l ../../pkg/location/csv/testdata/coopdrive-gps-pings-2022.04.06.csv
```

# Output example
```shell
detected crossing [0] {"tollgate_id":"holland_tunnel","worker_id":"4026b7e2b949a25312a6ba7a038e852b1a927ff6","movement":{"from":{"lon":-74.008265,"lat":40.725799},"to":{"lon":-74.024297,"lat":40.728271}},"direction":"NW","alg":"vector"}
detected crossing [1] {"tollgate_id":"holland_tunnel","worker_id":"4026b7e2b949a25312a6ba7a038e852b1a927ff6","movement":{"from":{"lon":-74.024297,"lat":40.728271},"to":{"lon":-74.024297,"lat":40.728271}},"direction":"","alg":"bbox"}
detected crossing [2] {"tollgate_id":"holland_tunnel","worker_id":"4026b7e2b949a25312a6ba7a038e852b1a927ff6","movement":{"from":{"lon":-74.024297,"lat":40.728271},"to":{"lon":-74.029229,"lat":40.729248}},"direction":"NW","alg":"bbox"}
```