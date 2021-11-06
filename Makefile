distroless:
	docker build . -f Dockerfile.distroless -t gobase:distroless

alpine:
	docker build . -f Dockerfile.alpine -t gobase:alpine