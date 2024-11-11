TAR=test
BUILD_CMD=go build

APP_MAIN_PATH=./src/app/app.go
DEV_MAIN_PATH=./src/dev_test_core/main.go.dev
SRC=./src/main.go

ENTRY=APP_Entry

ENV_FILE_PRODUCT=./src/public/env.go.product
ENV_FILE_DEV=./src/public/env.go.dev
ENV_FILE=./src/public/env.go

all:
	cp ${ENV_FILE_PRODUCT} ${ENV_FILE}

	cp ${APP_MAIN_PATH} ${SRC}

	echo "`sed 's/${ENTRY}/main/g' ${SRC}`" > ${SRC}
	
	${BUILD_CMD} -o ${TAR} ${SRC}
	
	rm -rf ${SRC}
dev:
	cp ${ENV_FILE_DEV} ${ENV_FILE}
	
	cp $(DEV_MAIN_PATH) ${SRC}
	
	echo "`sed 's/target_op_file/${TAR}/g' ${SRC}`" > ${SRC}
	
	${BUILD_CMD} -o ${TAR} ${SRC}
	
	rm -rf ${SRC}

clean:
	rm -rf ${TAR}
