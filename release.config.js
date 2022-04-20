module.exports = {
    branches: ['main', 'semver'],
    plugins: [
        [
            "@semantic-release-plus/docker",
            {
                name: "ome:ci-release",
                skipLogin: true,
            },
        ],
    ],
};