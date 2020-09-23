github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "4c40fa674130e6d92bcdb4ef9bd17954fdbf3fab",
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
        "PATH=\\\$(pwd)/\\\$(dirname $(out_location :mga)):\\\$PATH go generate -x ./...",
        "$(out_location :mga) generate kit endpoint ./...",
        "$(out_location :mga) generate event handler ./...",
        "$(out_location :mga) generate event handler --output subpkg:suffix=gen ./...",
        "$(out_location :mga) generate event dispatcher ./...",
        "$(out_location :mga) generate event dispatcher --output subpkg:suffix=gen ./...",
        "$(out_location :mga) generate testify mock ./...",
        "$(out_location :mga) generate testify mock --output subpkg:suffix=mocks ./...",
        "$(out_location :mga) create service --force internal/scaffold/service/test",
    ],
)

tarball(
    name = "package",
    srcs = ["README.md", "LICENSE", ":mga"],
    out = f"mga_{CONFIG.OS}_{CONFIG.ARCH}.tar.gz",
    gzip = True,
    labels = ["dist"],
)

subinclude("///pleasings2//misc")

sha256sum(
    name = "checksums.txt",
    srcs = [
        "@linux_amd64//:package",
        "@darwin_amd64//:package",
    ],
    out = "checksums.txt",
    labels = ["dist"],
)

subinclude("///pleasings2//github")

github_release(
    name = "publish",
    assets = [
        "@linux_amd64//:package",
        "@darwin_amd64//:package",
        ":checksums.txt",
    ],
    labels = ["dist"],
)
