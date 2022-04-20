module.exports = {
    branches: ['main', 'semver'],
    plugins: [
        [
            "@semantic-release-plus/docker",
            {
                name: {
                    registry: "ghcr.io",
                    namespace: "openmarketplaceengine",
                    repository: "openmarketplaceengine",
                    tag: undefined,
                    sha: undefined,
                },
                skipLogin: true,
            },
        ],
    ],
};

// <registry>/<namespace>/<repo>:<tag> or <registry>/<namespace>/<repo>@<sha>