.POSIX:
.PHONY: *
.EXPORT_ALL_VARIABLES:

KUBECONFIG = $(shell pwd)/kubeconfig
KUBE_CONFIG_PATH = $(KUBECONFIG)

default: ansible bootstrap smoke-test post-install clean

configure:
	./scripts/configure
	git status

ansible:
	make -C ansible

bootstrap:
	make -C bootstrap

smoke-test:
	make -C test filter=Smoke

post-install:
	@./scripts/hacks

test:
	make -C test

dev:
	make -C ansible cluster env=dev
	make -C bootstrap

git-hooks:
	pre-commit install
