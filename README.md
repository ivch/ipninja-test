# Test task for ipninja

## Notes
- To save some time i have intentionally skipped writing test. I wrote couple of test just to show
that i can write test and i remember of their importance and usually i look for 90+% test coverage.
- Client part can look awfull from ther professional frontend developer point of view. But since i'm a 
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
- 'docker build -t ipninjatest .'
- 'docker run -p 8080:8080 ipninjatest'
- `http://localhost:8080`