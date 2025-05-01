import js from "@eslint/js";
// jsonPlugin is not used since we're ignoring JSON files for now

export default [
  js.configs.recommended,
  {
    ignores: [
      "node_modules/",
      "out/",
      "bin/",
      "build/",
      "dist/",
      "internal/generated/",
      "data/",
      "ui/",
      "coverage/",
      "**/*.d.ts",
      "**/*.json", // Temporarily ignore JSON files
      ".eslintrc.js", // Ignore old config files
      ".prettierrc.js",
      ".vscode/",
      "package-lock.json",
    ],
  },
  {
    files: ["**/*.{js,mjs,cjs}"],
    languageOptions: {
      ecmaVersion: "latest",
      sourceType: "module",
      globals: {
        module: "readonly",
      },
    },
    rules: {
      "no-unused-vars": ["error", { argsIgnorePattern: "^_" }],
    },
  },
];
