const path = require('path');
const webpack = require('webpack');
const NodePolyfillPlugin = require("node-polyfill-webpack-plugin");
const CopyPlugin = require("copy-webpack-plugin");
const Dotenv = require('dotenv-webpack');
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;


module.exports = {
    mode: process.env.NODE_ENV,
    entry: path.resolve(__dirname, 'src', 'index.js'),
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'bundle.js'
    },
    devServer: {
        contentBase: path.resolve(__dirname, 'public'),
        open: true,
        clientLogLevel: 'debug',
        port: 9000
    },
    resolve: {
        extensions: ['*', '.js', '.jsx'],
        alias: {
            modules: path.resolve(__dirname + '/node_modules/')
        },
    },
    plugins: [
        new NodePolyfillPlugin(),
        new CopyPlugin({
            patterns: [
                {from: "public", to: "."}
            ]
        }),
        new webpack.EnvironmentPlugin(['API_URL']),
       // new BundleAnalyzerPlugin()
    ],
    module: {
        rules: [
            {
                test: /\.(jsx|js)$/,
                include: path.resolve(__dirname, 'src'),
                exclude: /node_modules/,
                use: [{loader: 'babel-loader',}]
            },
            {
                test: /\.css$/i,
                use: ["style-loader", "css-loader"],
            },
        ]
    }
}