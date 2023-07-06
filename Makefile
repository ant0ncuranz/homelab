.POSIX:
.PHONY: *
.EXPORT_ALL_VARIABLES:

KUBECONFIG = $(shell pwd)/kubeconfig
KUBE_CONFIG_PATH = $(KUBECONFIG)

default: metal bootstrap external # smoke-test post-install clean

configure:
	./scripts/configure
	git status

metal:
	make -C metal

bootstrap:
	make -C bootstrap

external:
	make -C external

smoke-test:
	make -C test filter=Smoke

post-install:
	@./scripts/hacks

test:
	make -C test

dev:
	make -C metal cluster env=dev
	make -C bootstrap

git-hooks:
	pre-commit install
