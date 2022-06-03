def _push_to_ecr_impl(ctx):
    ctx.actions.expand_template(
        template = ctx.file._pusher_template,
        output = ctx.outputs.executable,
        substitutions = {
            "%{image_name}": ctx.attr.image_name,
            "%{registry}": ctx.attr.registry,
            "%{repo}": ctx.attr.repo,
            "%{name}":ctx.label.name,
            "%{aws_profile}":ctx.attr.aws_profile,
            "%{log_verbosity}": ctx.attr.log_verbosity,
            "%{extra_tag}": ctx.attr.extra_tag,
        },
        is_executable = True,
    )

    runfiles = ctx.runfiles()
    return [DefaultInfo(runfiles = runfiles)]

push_to_ecr =  rule(
    attrs = {
        "image_name": attr.string(
            mandatory = True,
        ),

        "registry": attr.string(
            mandatory = True,
        ),

        "repo": attr.string(
            mandatory = True,
        ),

        "aws_profile": attr.string(
            mandatory = True,
        ),

        "log_verbosity": attr.string(
            default = "ERROR",
        ),

        "extra_tag": attr.string(
            mandatory = True,
        ),

        "_pusher_template": attr.label(
            default = Label("@//bazel/ecr/templates:pusher.py.tpl"),
            allow_single_file = True,
        ),
    },
    implementation = _push_to_ecr_impl,
    executable =  True,
)
