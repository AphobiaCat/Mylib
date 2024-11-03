SRC_=`find ./src -type f -name "*.go"`
SRC=./src/main.go
TAR=test
BUILD_CMD=go build


all:
	${BUILD_CMD} -o ${TAR} ${SRC}
clean:
	rm -rf ${TAR}
