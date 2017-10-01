#!/usr/bin/make -f

.DEFAULT: test
.PHONY: test reset

test:
	./dev/test.sh

reset:
	./dev/reset.sh
