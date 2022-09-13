GOOS=windows GOARCH=amd64 go build -o bin/WaterMarker_Windows.exe WaterMarker.go

GOOS=darwin GOARCH=amd64 go build -o bin/WaterMarker_MacOS WaterMarker.go

GOOS=linux GOARCH=amd64 go build -o bin/WaterMarker_Linux WaterMarker.go