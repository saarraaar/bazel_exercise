def _pkg_tar_impl(ctx):
    """Implements a simple pkg_tar rule."""

    tar_file = ctx.actions.declare_file("{}.tar".format(ctx.attr.name))

    args = ["--output", tar_file.path]
    if ctx.attr.package_dir:
        args.extend(["--root", ctx.attr.package_dir])

    args.extend([scr.path for scr in ctx.files.srcs])

    ctx.actions.run(
        inputs = ctx.files.srcs,
        outputs = [tar_file],
        arguments = args,
        progress_message = "Writing {}".format(tar_file.short_path),
        executable = ctx.executable._tar_tool,
    )

    return [DefaultInfo(files = depset([tar_file]))]

pkg_tar = rule(
    implementation = _pkg_tar_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "Input files to put into the tarball.",
        ),
        "package_dir": attr.string(),
        "_tar_tool": attr.label(
            default = Label("//utils:tar"),
            cfg = "exec",
            executable = True,
            allow_single_file = True,
        ),
    },
    doc = "A simplified version of bazel_tools' pkg_tar",
)

# A simple macro wrapping around native genrule.
# Alternatively, this could have been implemented as a rule following this example:
# https://github.com/bazelbuild/examples/blob/main/rules/shell_command/rules.bzl#L44
# All the issues with the shell-based appoach stated in the example above apply here as well.
def file_size(name, file):
    native.genrule(
        name = name,
        srcs = [file],
        outs = ["{}.sh".format(name)],
        cmd = """echo -n 'echo $$(basename $<) file size is ' > $@ && (ls -dnL $< | awk '{print $$5}') >> $@""",
        executable = True,
    )
