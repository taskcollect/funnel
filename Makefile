distroless:
	docker build . -f Dockerfile.distroless -t ghcr.io/taskcollect/funnel

alpine:
	docker build . -f Dockerfile.alpine -t ghcr.io/taskcollect/funnel:alpine