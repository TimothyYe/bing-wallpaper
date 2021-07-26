# Binary name
BINARY=bw
DATE=`date +%Y%m%d`
# Builds the project
build:
		go build -o ./bw/${BINARY} ./bw/main.go
release:
		# Build image
		docker build -t "timothyye/bing:${DATE}" -f Dockerfile .
		docker push "timothyye/bing:${DATE}"
test:
		go test .
# Cleans our projects: deletes binaries
clean:
		rm -rf ./bw/bw

.PHONY:  clean build
