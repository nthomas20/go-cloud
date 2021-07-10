# This how we want to name the binary output
BINARY_LINUX=go-cloud
BINARY_WINDOWS=go-cloud.exe
BINARY_MACOS=go-cloud-darwin
DESCRIPTOR_LINUX=Linux-amd64
DESCRIPTOR_WINDOWS=Windows-amd64
DESCRIPTOR_MACOS=MacOS

RELEASE_TITLE=GO Cloud
RELEASE_MESSAGE=Official Release

# These are the values we want to pass for VERSION and BUILD ( Semantic Versioning Recommended: https://semver.org/ )
VERSION=`git describe --tags --abbrev=0 \`git rev-list --tags --max-count=1)\``
BUILD=`date +%FT%T%z`

# Configure build and coverage flags
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD}"
LDFLAGS_HERE=-ldflags "-w -s -X main.version=${VERSION}-local -X main.buildDate=${BUILD}"

COVERAGEFLAGS=-race -coverprofile=coverage.txt -covermode=atomic -v
COVERAGE_TOKEN = ${CODECOV_TOKEN}

# Current git branch
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# .PHONY: all test coverage changelog clean build release reqcheck
.PHONY: all test clean build

# A plain-old `make` command will run the build process for local executable
all: build

# Tests the project
test:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go test ./... ${LDFLAGS}
	# env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go test ./... ${LDFLAGS}
	# env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go test ./... ${LDFLAGS}

# Tests the project linux-only and upload coverage report
coverage:
	@[ "${CODECOV_TOKEN}" ] && echo "all good" || ( echo "CODECOV_TOKEN is not set"; exit 1 )
	if [ -f coverage.txt ] ; then rm coverage.txt ; fi
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go test ./... ${LDFLAGS} ${COVERAGEFLAGS}
	curl -s https://codecov.io/bash | bash

# Generate changelog
# changelog: reqcheck
# 	@[ "${MVERSION}" ] && echo "Tagging Version ${MVERSION}" || ( echo "MVERSION not specified. 'make build MVERSION=v#.#.#'"; exit 1 )
# 	git-chglog -o CHANGELOG.md --next-tag ${MVERSION}
# 	git add CHANGELOG.md
# 	git commit -m "update changelog"
# 	git push

# Install golang depenencies
install-dev-dependencies: clean
	go mod tidy

# Cleans our project: deletes binaries
clean:
	if [ -f coverage.txt ] ; then rm coverage.txt ; fi
	if [ -f bin/${BINARY_LINUX} ] ; then rm bin/${BINARY_LINUX} ; fi
	# if [ -f bin/${BINARY_WINDOWS} ] ; then rm bin/${BINARY_WINDOWS} ; fi
	# if [ -f bin/${BINARY_MACOS} ] ; then rm bin/${BINARY_MACOS} ; fi

# Generate and push changelog ( https://github.com/git-chglog/git-chglog )
changelog:
	git pull
	git-chglog -o CHANGELOG.md
	git commit -a -m "update changelog"
	git push
	git push origin ${VERSION}

# Builds the project ( https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04 )
build: clean
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build ${LDFLAGS} -o bin/${BINARY_LINUX}
	# env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build ${LDFLAGS} -o bin/${BINARY_WINDOWS}
	# env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build ${LDFLAGS} -o bin/${BINARY_MACOS}

# Check our requirements
reqcheck:
	# hub version
	git version
	git-chglog -v

# Pushes release ( Requires https://github.com/github/hub )
# release: build
# 	@[ "${GITHUB_TOKEN}" ] || ( echo "GITHUB_TOKEN not specified"; exit 1 )
# 	@[ "${MVERSION}" ] && echo "Releasing Version ${MVERSION}" || ( echo "MVERSION not specified. 'make release MVERSION=v#.#.#'"; exit 1 )
# 	git pull
# 	# Create the release branch
# 	git branch release/${MVERSION}
# 	git checkout release/${MVERSION}
# 	git push --set-upstream origin release/${MVERSION}
# 	# Generate the changelog file
# 	make changelog
# 	hub release create \
# 		-a "bin/${BUILD_EXECUTABLE}#${BUILD_EXECUTABLE} (${BUILD_DESCRIPTOR})" \
# 		-m "${RELEASE_TITLE} ${MVERSION}" \
# 		-m "${RELEASE_MESSAGE}" \
# 		-t release/${MVERSION} \
# 		${MVERSION}
# 	git pull
# 	# Copy changelog
# 	cp CHANGELOG.md CHANGELOG.md.t
# 	# Return to originating branch
# 	git checkout ${BRANCH}
# 	# Update the changelog from the release branch
# 	cat CHANGELOG.md.t > CHANGELOG.md
# 	rm CHANGELOG.md.t
# 	git add CHANGELOG.md
# 	git commit -m "update changelog"
# 	git push -f
