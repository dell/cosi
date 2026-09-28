# Copyright © 2026 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# Dell Technologies, Dell and other trademarks are trademarks of Dell Inc.
# or its subsidiaries. Other trademarks may be trademarks of their respective
# owners.

.PHONY: copy-csm-common vendor go-code-tester

copy-csm-common:
	cp ../csm/config/csm-common.mk .

vendor: clean
	go generate ./...
	rm -rf vendor
	GOPRIVATE=github.com go mod vendor

go-code-tester:
	git clone --depth 1 git@github.com:dell/actions.git temp-repo
	cp temp-repo/go-code-tester/entrypoint.sh ./go-code-tester
	chmod +x go-code-tester
	rm -rf temp-repo
