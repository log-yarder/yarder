#!/usr/bin/make -f

.PHONY: test
.DEFAULT: test
test:
	./ci/test.sh
