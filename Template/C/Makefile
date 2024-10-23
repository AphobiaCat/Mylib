
GCC=g++

SRC=./src/main.c `find ./src -type f -name "*.cpp"` `find ./src -type f -name "*.h"`

TAR=test

FLAG=-g


all:
	${GCC} ${FLAG} ${SRC} -o ${TAR} -lpthread


clean:
	rm -rf ${TAR}
