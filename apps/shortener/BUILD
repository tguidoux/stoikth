load("//tools/rules/golang:defs.bzl", "stoikth_go_binary", "stoikth_go_image", "stoikth_go_library", "stoikth_go_test")

stoikth_go_library(
    name = "shortener_lib",
    srcs = [
        "db.go",
        "main.go",
        "types.go",
        "utils.go",
    ],
    importpath = "github.com/tguidoux/stoith/apps/shortener",
    visibility = ["//visibility:private"],
    deps = ["@com_github_mattn_go_sqlite3//:go-sqlite3"],
)

stoikth_go_binary(
    name = "main",
    embed = [":shortener_lib"],
    visibility = ["//visibility:public"],
)

# Build and run this image with
#     bazel run //apps/shortener:shortener_img_load_docker
#     docker run --rm stoikth/shortener:latest
stoikth_go_image(
    name = "shortener_img",
    binary = ":main",
    image_tags = ["stoikth/shortener:latest"],
)

stoikth_go_test(
    name = "shortener_test",
    srcs = [
        "db_test.go",
        "main_test.go",
        "utils_test.go",
    ],
    embed = [":shortener_lib"],
    deps = ["@com_github_stretchr_testify//assert"],
)
