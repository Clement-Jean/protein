#!/bin/bash
set -eo pipefail

rm -rf "$(bazel info output_path)/_coverage/genhtml"
genhtml --output "$(bazel info output_path)/_coverage/genhtml" "$(bazel info output_path)/_coverage/_coverage_report.dat"
open "$(bazel info output_path)/_coverage/genhtml/index.html"