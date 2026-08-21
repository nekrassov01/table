.PHONY: example

TARGET = $(if $(target),$(target),all)

example:
ifdef data
ifndef mode
	$(error mode is undefined)
endif
endif
	@go run ./cmd '$(TARGET)'$(if $(mode), '$(mode)')$(if $(data), '$(data)')
