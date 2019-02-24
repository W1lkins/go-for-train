NAME := go-for-train
PKG := github.com/W1lkins/$(NAME)

CGO_ENABLED := 0

BUILDTAGS :=

include basic.mk

.PHONY: prebuild
prebuild:
