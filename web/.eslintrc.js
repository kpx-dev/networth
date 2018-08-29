module.exports = {
  parser: "babel-eslint",
  env: {
    es6: true,
    node: true,
    browser: true
  },
  parserOptions: {
    ecmaVersion: 6,
    sourceType: "module",
    ecmaFeatures: {
      jsx: true
    }
  },
  plugins: ["react"],
  extends: [
    "eslint:recommended",
    "plugin:react/recommended",
    // "plugin:prettier/recommended" // TODO: enable
  ],
  rules: {
    "no-unused-vars": "off", // TODO: delete
    "react/prop-types": "off", // TODO: delete
    "max-len": ["warn", { code: 200 }],
    "no-console": "off" // TODO: warn
  },
  settings: {
    react: {
      version: "16"
    }
  }
};
