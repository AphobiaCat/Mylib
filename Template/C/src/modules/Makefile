SUB_MODULES_CXXFLAGS ?=
SUB_MODULES_LDFLAGS ?=

TINYCC_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))TinyCC_Manager/
SUB_MODULES_CXXFLAGS += -I$(TINYCC_DIR)inc
SUB_MODULES_LDFLAGS += $(TINYCC_DIR)lib/libtcc.a

export SUB_MODULES_CXXFLAGS
export SUB_MODULES_LDFLAGS