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
  targets = "docker-bake"
}
target "moduleA" {
  dockerfile = "Dockerfile"
  context    = "./modules/message_monitoring"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = ["$${DOCKER_USERNAME}/$${DOCKER_REGISTRY_PREFIX}-message_monitoring:$${TAG}"]
}
target "moduleA_debug_arm64" {
  dockerfile = "Dockerfile.debug.arm64"
  context    = "./modules/message_monitoring"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = ["$${DOCKER_USERNAME}/$${DOCKER_REGISTRY_PREFIX}-message_monitoring:$${TAG}"]
}
target "moduleB" {
  dockerfile = "Dockerfile"
  context    = "./modules/message_monitoring"
  platforms  = ["linux/amd64", "linux/arm64/v8"]
  tags       = ["$${DOCKER_USERNAME}/$${DOCKER_REGISTRY_PREFIX}-message_monitoring:$${TAG}"]
}
