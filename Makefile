#!/usr/bin/make -f

.DEFAULT: test
.PHONY: test clean

test:
	./dev/test.sh

clean:
	./dev/clean.sh
