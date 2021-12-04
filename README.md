# bazel_exercise

* `utils/tar.go` -- implements a cli tool that creates a tar file from a list of inputs.
* `rules/rules.bzl` -- contains a `pkg_tar` rule build using the cli tool above and the `file_size` executable rule/macro that prints out the size ouf the given input file to stdout.