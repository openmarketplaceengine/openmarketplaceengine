module.exports = {
    branches: ['main'],
    plugins: [
        [
            "@semantic-release-plus/docker",
            {
                name: "ghcr.io/openmarketplaceengine/openmarketplaceengine:main",
            },
        ],
        "@semantic-release/github",
    ],
};