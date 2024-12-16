code tree  
src  
|-- app //user code here  
|-- dev_test_core //dev env code  
|-- example //include all module use case  
|-- module //all module tool  
|-- public //include public func and config file

how to use:  
1: ./init_project.sh test "this is a test"  
2: make  
3: ./install.sh

make type   
all: build release version  
dev: build dev version - auto rebuild and run when dev version is working  
debug: build debug version, No compilation optimizations  

change config in  
./src/public/env.go.dev  
./src/public/env.go.product  

Import the required modules from src/module/  
and modify your code in src/app/app.go  
