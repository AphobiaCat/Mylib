code tree  
src  
|-- app //user code here  
|-- dev_test_core //dev env code  
|-- example //include all module use case  
|-- module //all module tool  
|-- public //include public func

how to use:  
1: ./init_project.sh test "this is a test"  
2: make  
3: ./install.sh  
4: or run it by ./test --config ./config_debug.yaml

make type   
all: build release version  
dev: build dev version - auto rebuild and run when dev version is working  
debug: build debug version, No compilation optimizations  

change config in  
./config_debug.yaml
./config_prod.yaml

Import the required modules from src/module/  
and modify your code in src/app/app.go  
