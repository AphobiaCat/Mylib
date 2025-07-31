TAR=test
BUILD_CMD=go build

APP_MAIN_PATH=./src/app/app.go
APP_SRC=./src/app/main.go
BUILD_SRC=./src/app

DEV_MAIN_PATH=./src/dev_test_core/main.go.dev
DEV_SRC=./src/main.go

ENTRY=APP_Entry

all:
	cp ${APP_MAIN_PATH} ${APP_SRC}

	sed -i 's/!prod/prod/g; s/${ENTRY}/main/g' ${APP_SRC}

	${BUILD_CMD} -tags=prod -o ${TAR} ${BUILD_SRC}

	rm -f ${APP_SRC}

debug:
	cp ${APP_MAIN_PATH} ${APP_SRC}

	sed -i 's/!prod/prod/g; s/${ENTRY}/main/g' ${APP_SRC}
	
	${BUILD_CMD} -gcflags="all=-N -l" -tags=prod -o ${TAR} ${BUILD_SRC}

	rm -f ${APP_SRC}

dev:	
	cp $(DEV_MAIN_PATH) ${DEV_SRC}
	
	echo "`sed 's/target_op_file/${TAR}/g' ${DEV_SRC}`" > ${DEV_SRC}
	
	${BUILD_CMD} -o ${TAR} ${DEV_SRC}
	
	rm -rf ${DEV_SRC}

clean:
	rm -rf ${TAR}

