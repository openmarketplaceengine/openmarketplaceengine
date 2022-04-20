module.exports = {
    branches: ['main', 'semver'],
    plugins: [
        [
            "@semantic-release-plus/docker",
            {
                name: "ghcr.io/openmarketplaceengine/openmarketplaceengine:ci",
                skipLogin: true,
            },
        ],
    ],
};