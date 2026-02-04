# Copyright © 2026 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# Dell Technologies, Dell and other trademarks are trademarks of Dell Inc.
# or its subsidiaries. Other trademarks may be trademarks of their respective
# owners.

download-csm-common:
	git clone --depth 1 git@github.com:dell/csm.git temp-repo
	cp temp-repo/config/csm-common.mk .
	rm -rf temp-repo

vendor: clean
	go generate ./...
	GOPRIVATE=github.com go mod vendor

go-code-tester:
	git clone --depth 1 git@github.com:dell/actions.git temp-repo
	cp temp-repo/go-code-tester/entrypoint.sh ./go-code-tester
	chmod +x go-code-tester
	rm -rf temp-repo
