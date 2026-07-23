GHCR_REPOSITORY = ghcr.io/noaled-lab/mediamtx

ghcr:
	$(eval VERSION := $(shell git describe --tags | tr -d v))

	echo "$(GITHUB_TOKEN)" | docker login ghcr.io -u $(GITHUB_USER) --password-stdin

	docker buildx rm builder 2>/dev/null || true
	docker buildx create --name=builder

#	docker build --builder=builder \
#	-f docker/ffmpeg-rpi.Dockerfile . \
#	--platform=linux/arm/v6,linux/arm/v7,linux/arm64 \
#	-t $(GHCR_REPOSITORY):$(VERSION)-ffmpeg-rpi \
#	-t $(GHCR_REPOSITORY):1-ffmpeg-rpi \
#	-t $(GHCR_REPOSITORY):latest-ffmpeg-rpi \
#	--push
#
#	docker build --builder=builder \
#	-f docker/rpi.Dockerfile . \
#	--platform=linux/arm/v6,linux/arm/v7,linux/arm64 \
#	-t $(GHCR_REPOSITORY):$(VERSION)-rpi \
#	-t $(GHCR_REPOSITORY):1-rpi \
#	-t $(GHCR_REPOSITORY):latest-rpi \
#	--push
#
#	docker build --builder=builder \
#	-f docker/ffmpeg.Dockerfile . \
#	--platform=linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm64 \
#	-t $(GHCR_REPOSITORY):$(VERSION)-ffmpeg \
#	-t $(GHCR_REPOSITORY):1-ffmpeg \
#	-t $(GHCR_REPOSITORY):latest-ffmpeg \
#	--push

	docker build --builder=builder \
	-f docker/standard.Dockerfile . \
	--platform=linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm64 \
	-t $(GHCR_REPOSITORY):$(VERSION) \
	-t $(GHCR_REPOSITORY):1 \
	-t $(GHCR_REPOSITORY):latest \
	--push

	docker buildx rm builder
