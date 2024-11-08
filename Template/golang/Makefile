TAR=test
BUILD_CMD=go build

APP_MAIN_PATH=./src/app/app.go
DEV_MAIN_PATH=./src/dev_test_core/main.go.dev
SRC=./src/main.go

ENTRY=APP_Entry

all:
	cp ${APP_MAIN_PATH} ${SRC}

	echo "`sed 's/${ENTRY}/main/g' ${SRC}`" > ${SRC}
	
	${BUILD_CMD} -o ${TAR} ${SRC}
	
	rm -rf ${SRC}
dev:
	cp $(DEV_MAIN_PATH) ${SRC}
	
	echo "`sed 's/target_op_file/${TAR}/g' ${SRC}`" > ${SRC}
	
	${BUILD_CMD} -o ${TAR} ${SRC}

clean:
	rm -rf ${TAR}
