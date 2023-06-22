variable "DOCKER_USERNAME" {
  default = "user"
}
variable "DOCKER_REGISTRY_PREFIX" {
  default = "prefix"
}
variable "TAG" {
  default = "latest"
}
group "prefix-modules" {
  targets = ["moduleA", "moduleA-debug-arm64", "moduleB-debug"]
}
target "moduleA" {
  dockerfile = "Dockerfile"
  context    = "./modules/message_monitoring"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-moduleA:${TAG}"
}
target "moduleA-debug-arm64" {
  dockerfile = "Dockerfile.debug.arm64"
  context    = "./modules/message_monitoring"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-moduleA:${TAG}-debug-arm64"
}
target "moduleB-debug" {
  dockerfile = "Dockerfile.debug"
  context    = "./modules/message_monitoring"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-moduleB:${TAG}-debug"
}
