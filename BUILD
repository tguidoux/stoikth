load("@gazelle//:def.bzl", "gazelle")
load("@npm//:defs.bzl", "npm_link_all_packages")
load("@rules_go//go:def.bzl", "TOOLS_NOGO", "nogo")
load("//tools/multirun:defs.bzl", "multi_run")

# gazelle:map_kind go_binary stoikth_go_binary //tools/rules/golang:defs.bzl
# gazelle:map_kind go_library stoikth_go_library //tools/rules/golang:defs.bzl
# gazelle:map_kind go_test stoikth_go_test //tools/rules/golang:defs.bzl
gazelle(name = "gazelle")

nogo(
    name = "nogo",
    visibility = ["//visibility:public"],
    deps = TOOLS_NOGO,
)

npm_link_all_packages(name = "node_modules")

multi_run(
    name = "clean",
    targets = [
        "//:gazelle",
    ],
)
