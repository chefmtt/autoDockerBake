variable "DOCKER_USERNAME" {
  default = "usR"
}
variable "DOCKER_REGISTRY_PREFIX" {
  default = "foo"
}
variable "TAG" {
  default = "latest"
}
group "foo-modules" {
  targets = ["moduleC-lint", "moduleC-lint-debug-test", "moduleC-lint-debug", "moduleA", "moduleA-debug-arm64", "moduleB-debug"]
}
target "moduleC-lint" {
  dockerfile = "lint.Dockerfile"
  context    = "test/moduleC"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-moduleC:${TAG}-lint"
}
target "moduleC-lint-debug-test" {
  dockerfile = "lint.debug-test.Dockerfile"
  context    = "test/moduleC"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-moduleC:${TAG}-lint-debug-test"
}
target "moduleC-lint-debug" {
  dockerfile = "lint.debug.Dockerfile"
  context    = "test/moduleC"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-moduleC:${TAG}-lint-debug"
}
target "moduleA" {
  dockerfile = "Dockerfile"
  context    = "test/moduleA"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-moduleA:${TAG}"
}
target "moduleA-debug-arm64" {
  dockerfile = "Dockerfile.debug.arm64"
  context    = "test/moduleA"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-moduleA:${TAG}-debug-arm64"
}
target "moduleB-debug" {
  dockerfile = "Dockerfile.debug"
  context    = "test/moduleB"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-moduleB:${TAG}-debug"
}
