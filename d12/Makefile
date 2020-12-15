first-example:
	@if [ "`./first.awk example`" != "25" ]; then \
        echo "Expected 25, got:"; \
		./first.awk example; \
        exit 1; \
    fi

first:
	./first.awk input

second-example:
	@if [ "`./second.awk example`" != "286" ]; then \
        echo "Expected 286, got:"; \
		./second.awk example; \
        exit 1; \
    fi

second:
	./second.awk input