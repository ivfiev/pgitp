
all:
	mkdir -p ./bin/android
	mkdir -p ./bin/linux
	CGO_ENABLED=0 GOARCH=arm64 GOOS=android go build -ldflags="-s -w" -o ./bin/android/pgitp src/*
	go build -ldflags="-s -w" -o ./bin/linux/pgitp src/*
