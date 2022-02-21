group "default" {
  targets = ["gnssaggr"]
}

target "gnssaggr" {
  dockerfile = "Dockerfile"
  target = "production"
  platforms = ["linux/amd64", "linux/arm64"]
  tags = ["jonikahara/gnssaggr:latest"]
}
