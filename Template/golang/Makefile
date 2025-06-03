TAR=test
BUILD_CMD=go build

APP_MAIN_PATH=./src/app/app.go
APP_SRC=./src/app/main.go
ALL_APP_SRC=`find ./src/app -type f -name "*.go" ! -name "app.go" ! -name "main.go"`

DEV_MAIN_PATH=./src/dev_test_core/main.go.dev
DEV_SRC=./src/main.go

ENTRY=APP_Entry

all:
	cp ${APP_MAIN_PATH} ${APP_SRC}

	echo "`sed 's/${ENTRY}/main/g' ${APP_SRC}`" > ${APP_SRC}

	${BUILD_CMD} -o ${TAR} ${APP_SRC} ${ALL_APP_SRC}

	rm -f ${APP_SRC}

debug:
	cp ${APP_MAIN_PATH} ${APP_SRC}

	echo "`sed 's/${ENTRY}/main/g' ${APP_SRC}`" > ${APP_SRC}
	
	${BUILD_CMD} -gcflags="all=-N -l" -o ${TAR} ${APP_SRC} ${ALL_APP_SRC}

	rm -f ${APP_SRC}

dev:	
	cp $(DEV_MAIN_PATH) ${DEV_SRC}
	
	echo "`sed 's/target_op_file/${TAR}/g' ${DEV_SRC}`" > ${DEV_SRC}
	
	${BUILD_CMD} -o ${TAR} ${DEV_SRC}
	
	rm -rf ${DEV_SRC}

clean:
	rm -rf ${TAR}

