#!/bin/bash
set -eo pipefail

find ./bazel-testlogs/ -type f -name "*.dat" -print0 | xargs -0 rm -f
bazel coverage --combined_report=lcov --nocache_test_results ...:all