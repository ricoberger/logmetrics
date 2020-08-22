BRANCH      ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILDTIME   ?= $(shell date '+%Y-%m-%d@%H:%M:%S')
BUILDUSER   ?= $(shell id -un)
REPO        ?= github.com/ricoberger/logmetrics
REVISION    ?= $(shell git rev-parse HEAD)
VERSION     ?= $(shell git describe --tags)

.PHONY: build release-major release-minor release-patch

build:
	go build -ldflags "-X ${REPO}/pkg/version.Version=${VERSION} \
		-X ${REPO}/pkg/version.Revision=${REVISION} \
		-X ${REPO}/pkg/version.Branch=${BRANCH} \
		-X ${REPO}/pkg/version.BuildUser=${BUILDUSER} \
		-X ${REPO}/pkg/version.BuildDate=${BUILDTIME}" \
		-o ./bin/logmetrics ./cmd/logmetrics;

release-major:
	$(eval OLD_VERSION=$(shell git describe --tags --abbrev=0))
	$(eval MAJORVERSION=$(shell git describe --tags --abbrev=0 | sed s/v// | awk -F. '{print $$1+1".0.0"}'))
	git checkout master
	git pull
	sed -i'.backup' 's/${OLD_VERSION}/${MAJORVERSION}/g' charts/README.md
	sed -i'.backup' 's/${OLD_VERSION}/${MAJORVERSION}/g' charts/logmetrics/Chart.yaml
	sed -i'.backup' 's/${OLD_VERSION}/${MAJORVERSION}/g' charts/logmetrics/values.yaml
	rm charts/README.md.backup
	rm charts/logmetrics/Chart.yaml.backup
	rm charts/logmetrics/values.yaml.backup
	git add .
	git commit -m 'Prepare release $(MAJORVERSION)'
	git push
	git tag -a $(MAJORVERSION) -m 'Release $(MAJORVERSION)'
	git push origin --tags

release-minor:
	$(eval OLD_VERSION=$(shell git describe --tags --abbrev=0))
	$(eval MINORVERSION=$(shell git describe --tags --abbrev=0 | sed s/v// | awk -F. '{print $$1"."$$2+1".0"}'))
	git checkout master
	git pull
	sed -i'.backup' 's/${OLD_VERSION}/${MINORVERSION}/g' charts/README.md
	sed -i'.backup' 's/${OLD_VERSION}/${MINORVERSION}/g' charts/logmetrics/Chart.yaml
	sed -i'.backup' 's/${OLD_VERSION}/${MINORVERSION}/g' charts/logmetrics/values.yaml
	rm charts/README.md.backup
	rm charts/logmetrics/Chart.yaml.backup
	rm charts/logmetrics/values.yaml.backup
	git add .
	git commit -m 'Prepare release $(MINORVERSION)'
	git push
	git tag -a $(MINORVERSION) -m 'Release $(MINORVERSION)'
	git push origin --tags

release-patch:
	$(eval OLD_VERSION=$(shell git describe --tags --abbrev=0))
	$(eval PATCHVERSION=$(shell git describe --tags --abbrev=0 | sed s/v// | awk -F. '{print $$1"."$$2"."$$3+1}'))
	git checkout master
	git pull
	sed -i'.backup' 's/${OLD_VERSION}/${PATCHVERSION}/g' charts/README.md
	sed -i'.backup' 's/${OLD_VERSION}/${PATCHVERSION}/g' charts/logmetrics/Chart.yaml
	sed -i'.backup' 's/${OLD_VERSION}/${PATCHVERSION}/g' charts/logmetrics/values.yaml
	rm charts/README.md.backup
	rm charts/logmetrics/Chart.yaml.backup
	rm charts/logmetrics/values.yaml.backup
	git add .
	git commit -m 'Prepare release $(PATCHVERSION)'
	git push
	git tag -a $(PATCHVERSION) -m 'Release $(PATCHVERSION)'
	git push origin --tags
