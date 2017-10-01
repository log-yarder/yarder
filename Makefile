#!/usr/bin/make -f

.DEFAULT: test
.PHONY: test reset

test:
	./dev/test.sh

clean:
	./dev/clean.sh
