module.exports = {
  root: true,
  env: {
    node: true,
    es2021: true,
  },
  extends: ["eslint:recommended", "plugin:json/recommended"],
  parserOptions: {
    ecmaVersion: "latest",
    sourceType: "module",
  },
  plugins: ["json"],
  rules: {
    "no-unused-vars": ["error", { argsIgnorePattern: "^_" }],
  },
  ignorePatterns: ["node_modules/", "build/", "dist/", "coverage/", "data/", "ui/", "**/*.d.ts"],
};
