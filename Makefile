.PHONY: all updater

all: updater
	./updater/updater -no-pr

updater:
	$(MAKE) -C updater