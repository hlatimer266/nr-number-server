# New Relice Take Home Challenge

### Build & Run
* `make build` (runs units tests)
* `make run`
  
### Testing
* Unit tests can be run using `make test`
* To run performance test:
  * `make build-pvs`
  * start application `make run` in one terminal
  * in separate terminal set to the root of this project, run `make test-pvs`
* Simple testing can be done using ```nc localhost 4000``` and manually entering number with application running
  
### Design Assumptions + Approach
* See `instructions.txt` for requirements
* Leading zeros required for write and are written to file
* Server specific new line is `$` i.e `123456789$` is valid
* Up to five connections are able to write and more are able to connect but are unable to write
* `server` handles kicking off all concurrent running routines and shuts down all rountines on `terminate$` sequence
* status report, file writing, and connection management are all executed on their own goroutine.
* file writing is performed on a batch basis every second while server is accepting new content and queueing for the next write
* console report runs every 10 seconds during reads / writes to cache + number.log 

### TODO
* report doesn't sync very well with incoming writes because I didn't manage to implement report aggregation in time :(
* unit test for server is incomplete. would have liked to added tests for: 
  * confirm bad input closes connection
  * confirm `terminate$` sequence closes all connections
  * both of the above can be observed when running actual program
