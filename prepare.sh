#!/bin/bash

go mod tidy

bazel run //:gazelle --ui_event_filters=-DEBUG,+INFO
bazel run //:gazelle --ui_event_filters=-DEBUG,+INFO -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
bazel run //:gazelle --ui_event_filters=-DEBUG,+INFO 
