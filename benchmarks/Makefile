benchtime  ?= 10000x
count      ?= 5
cpuprofile ?= cpu.prof
memprofile ?= mem.prof

PROFILE = $(if $(cpuprofile),-cpuprofile $(cpuprofile)) $(if $(memprofile),-memprofile $(memprofile))
CMD     = go test -benchmem -count $(count) -benchtime $(benchtime) $(PROFILE) .
TARGETS = all text html markdown backlog csv comparison cold
TARGET  = $(if $(target),$(target),all)

.PHONY: bench clean

bench:
ifeq ($(filter $(TARGET),$(TARGETS)),)
	$(error unknown benchmark target: $(TARGET))
else ifeq ($(TARGET),all)
	$(CMD) -bench .
else ifeq ($(TARGET),text)
	$(CMD) -bench '^BenchmarkText'
else ifeq ($(TARGET),html)
	$(CMD) -bench '^BenchmarkHTML'
else ifeq ($(TARGET),markdown)
	$(CMD) -bench '^BenchmarkMarkdown'
else ifeq ($(TARGET),backlog)
	$(CMD) -bench '^BenchmarkBacklog'
else ifeq ($(TARGET),csv)
	$(CMD) -bench '^BenchmarkCSV'
else ifeq ($(TARGET),comparison)
	$(CMD) -bench '^BenchmarkComparison'
else ifeq ($(TARGET),cold)
	go test -tags cold -benchmem -count 1 -benchtime 100x -bench 'Cold$$' .
endif

clean:
	go clean
	rm -f $(cpuprofile) $(memprofile) benchmarks.test
