export CED_VERSION="0.2.0"
# Create git tags
if git rev-parse "v$CED_VERSION" >/dev/null 2>&1; then
    echo "Tag already present"
else
    echo "Tag not present. Creating"
    git tag -a "v$CED_VERSION" &&
    git push origin "v$CED_VERSION"
fi
# Run go releaser for binary releases
goreleaser release --rm-dist &&
# Run docker build push
docker buildx build --platform linux/arm64,linux/arm/v7,linux/arm/v6,linux/amd64 --build-arg CED_VERSION=$CED_VERSION --push -t blmhemu/ced:$CED_VERSION -t blmhemu/ced:latest .