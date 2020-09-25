github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "f8a12721c6f929db3e227e07c152d428ac47ab1b",
)

subinclude("///pleasings2//go")

timestamp = git_show("%ct")
date_fmt = "+%FT%T%z"

go_build(
    name = "mga",
    definitions = {
        "main.version": "${VERSION:-" + git_branch() + "}",
        "main.commitHash": git_commit()[0:8],
        "main.buildDate": f'$(date -u -d "@{timestamp}" "{date_fmt}" 2>/dev/null || date -u -r "{timestamp}" "{date_fmt}" 2>/dev/null || date -u "{date_fmt}")',
    },
    pass_env = ["VERSION"],
    stamp = True,
    trimpath = True,
    labels = ["binary"],
)

sh_cmd(
    name = "generate",
    deps = [":mga"],
    cmd = [
        "PATH=\\\$(pwd)/\\\$(dirname $(out_exe :mga)):\\\$PATH go generate -x ./...",
        "$(out_exe :mga) generate kit endpoint ./...",
        "$(out_exe :mga) generate event handler ./...",
        "$(out_exe :mga) generate event handler --output subpkg:suffix=gen ./...",
        "$(out_exe :mga) generate event dispatcher ./...",
        "$(out_exe :mga) generate event dispatcher --output subpkg:suffix=gen ./...",
        "$(out_exe :mga) generate testify mock ./...",
        "$(out_exe :mga) generate testify mock --output subpkg:suffix=mocks ./...",
        "$(out_exe :mga) create service --force internal/scaffold/service/test",
    ],
)

tarball(
    name = "artifact",
    srcs = ["README.md", "LICENSE", ":mga"],
    out = f"mga_{CONFIG.OS}_{CONFIG.ARCH}.tar.gz",
    gzip = True,
    labels = ["dist"],
)

subinclude("///pleasings2//misc")

build_artifacts(
    name = "artifacts",
    artifacts = [
        "@linux_amd64//:artifact",
        "@darwin_amd64//:artifact",
    ],
    labels = ["manual"],
)

subinclude("///pleasings2//github")

github_release(
    name = "publish",
    assets = [":artifacts"],
    labels = ["manual"],
)
