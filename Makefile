# Binary name
BINARY=bw
DATE=`date +%Y%m%d`
# Builds the project
build:
		go build -o ./bw/${BINARY} ./bw/main.go
release:
		go build -o ./bw/${BINARY} ./bw/main.go
		# Build image
		docker build -t "r.xiaozhou.net/projects/bing:${DATE}" -f Dockerfile .
		docker push "r.xiaozhou.net/projects/bing:${DATE}"
		go clean
# Cleans our projects: deletes binaries
clean:
		rm -rf ./bw/bw

.PHONY:  clean build
