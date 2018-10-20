# Test task for ipninja

## Notes
- To save some time i have intentionally skipped writing test. I wrote couple of them just to show
that i know how to write tests and i know about tests importance. Usually i strive for 90+% test coverage.
- Client part can look awful from the professional frontend developer point of view. But since i'm a 
backend developer please do not judge strictly :) 

## Godep
- `go get github.com/tools/godep`
- Save deps: `godep save -t ./...`
- Restore deps: `godep restore`

## Run tests
- `go test ./...`

## Run application
- restore dependecies
- `go build -o app`
- `./app`
- `http://localhost:8080`

OR

- `docker build -t ipninjatest .`
- `docker run -p 8080:8080 ipninjatest`
- `http://localhost:8080`