import superinvoke


class Tags(superinvoke.Tags):
    DEV = "dev"
    CI = "ci"
    OPS = "ops"


class Tools(superinvoke.Tools):
    Tester = superinvoke.Tool(
        name="gotestsum",
        version="1.7.0",
        tags=[Tags.DEV, Tags.CI],
        links={
            superinvoke.Platforms.LINUX: (
                "https://github.com/gotestyourself/gotestsum/releases/download/v1.7.0/gotestsum_1.7.0_linux_amd64.tar.gz",
                "gotestsum",
            ),
        },
    )

    CodeLinter = superinvoke.Tool(
        name="golangci-lint",
        version="1.42.0",
        tags=[Tags.DEV, Tags.CI],
        links={
            superinvoke.Platforms.LINUX: (
                "https://github.com/golangci/golangci-lint/releases/download/v1.42.0/golangci-lint-1.42.0-linux-amd64.tar.gz",
                "golangci-lint-1.42.0-linux-amd64/golangci-lint",
            ),
        },
    )

    Migrator = superinvoke.Tool(
        name="golang-migrate",
        version="4.15.1",
        tags=[Tags.DEV],
        links={
            superinvoke.Platforms.LINUX: (
                "https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz",
                "migrate",
            ),
        },
    )

    QueryLinter = superinvoke.Tool(
        name="sqlvet",
        version="1.1.1",
        tags=[Tags.DEV, Tags.CI],
        links={
            superinvoke.Platforms.LINUX: (
                "https://github.com/houqp/sqlvet/releases/download/v1.1.3/sqlvet-v1.1.3-linux-amd64.tar.gz",
                "sqlvet",
            ),
        },
    )

    MigrationLinter = superinvoke.Tool(
        name="squawk",
        version="0.8.0",
        tags=[Tags.DEV, Tags.CI],
        links={
            superinvoke.Platforms.LINUX: (
                "https://github.com/sbdchd/squawk/releases/download/v0.8.0/squawk-linux-x86_64",
                ".",
            ),
        },
    )
