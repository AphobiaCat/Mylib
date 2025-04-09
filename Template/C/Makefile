# Compiler setup
SUB_MODULES_CXXFLAGS :=
SUB_MODULES_LDFLAGS :=
include src/modules/Makefile

CC = g++
CXX = g++
SRC_DIR = src
BUILD_DIR = build/obj
BIN_DIR = build
APP_DIR = $(SRC_DIR)/app

# Collect all include directories automatically
INC_DIRS := $(shell find $(SRC_DIR) -type d -name inc)
INCLUDES := $(addprefix -I, $(INC_DIRS))

CFLAGS = -std=c++20 -Wall -Wextra $(INCLUDES)
CXXFLAGS = $(CFLAGS)

CFLAGS    += $(SUB_MODULES_CXXFLAGS)
CXXFLAGS  += $(SUB_MODULES_CXXFLAGS)
LDFLAGS   += $(SUB_MODULES_LDFLAGS)

# Collect shared C and C++ source files from modules and public
MODULE_C_SOURCES := $(shell find $(SRC_DIR)/modules -name "*.c" ! -path "$(SRC_DIR)/modules/Examples/*") \
                    $(shell find $(SRC_DIR)/public -name "*.c")
MODULE_CPP_SOURCES := $(shell find $(SRC_DIR)/modules -name "*.cpp" ! -path "$(SRC_DIR)/modules/Examples/*") \
                      $(shell find $(SRC_DIR)/public -name "*.cpp")

MODULE_OBJECTS := $(patsubst $(SRC_DIR)/%.c, $(BUILD_DIR)/%.o, $(MODULE_C_SOURCES)) \
                  $(patsubst $(SRC_DIR)/%.cpp, $(BUILD_DIR)/%.o, $(MODULE_CPP_SOURCES))

# List all app directories automatically (e.g., header, worker)
APPS := $(notdir $(wildcard $(APP_DIR)/*))

# Default target: build all apps
all: $(addprefix $(BIN_DIR)/, $(APPS))

# Per-app build rules using dynamic Make function
define MAKE_APP_RULES
APP_$1_C_SOURCES := $(shell find $(APP_DIR)/$1/src -name "*.c")
APP_$1_CPP_SOURCES := $(shell find $(APP_DIR)/$1/src -name "*.cpp")
APP_$1_OBJECTS := $$(patsubst $(SRC_DIR)/%.c,$(BUILD_DIR)/%.o,$$(APP_$1_C_SOURCES)) \
                  $$(patsubst $(SRC_DIR)/%.cpp,$(BUILD_DIR)/%.o,$$(APP_$1_CPP_SOURCES))

$(BIN_DIR)/$1: $$(APP_$1_OBJECTS) $(MODULE_OBJECTS)
	@echo "Linking $$@ ..."
	$$(CXX) $$^ $(LDFLAGS) -o $$@
endef

# Generate rules for each app
$(foreach app,$(APPS),$(eval $(call MAKE_APP_RULES,$(app))))

# Compilation rule for C files
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c
	@mkdir -p $(dir $@)
	@echo "Compiling C: $< ..."
	$(CC) $(CFLAGS) -c $< -o $@

# Compilation rule for C++ files
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.cpp
	@mkdir -p $(dir $@)
	@echo "Compiling C++: $< ..."
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Clean generated files
clean:
	rm -rf $(BUILD_DIR) $(BIN_DIR)/*

.PHONY: all clean

