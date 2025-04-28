export default {
  // The root directory that Jest should scan for tests and modules
  rootDir: ".",

  // The test environment that will be used for testing
  testEnvironment: "jsdom",

  // The glob patterns Jest uses to detect test files
  testMatch: ["**/__tests__/**/*.js", "**/*.test.js", "**/*.spec.js"],

  // An array of file extensions your modules use
  moduleFileExtensions: ["js", "jsx", "ts", "tsx", "json", "node"],

  // A list of paths to directories that Jest should use to search for files in
  roots: ["<rootDir>/src"],

  // A map from regular expressions to module names or to arrays of module names
  // that allow to stub out resources with a single module
  moduleNameMapper: {
    // Handle CSS imports (with CSS modules)
    "\\.module\\.(css|sass|scss)$": "identity-obj-proxy",

    // Handle CSS imports (without CSS modules)
    "\\.(css|sass|scss)$": "<rootDir>/src/__mocks__/styleMock.js",

    // Handle image imports
    "\\.(jpg|jpeg|png|gif|webp|svg)$": "<rootDir>/src/__mocks__/fileMock.js",
  },

  // A list of paths to modules that run some code to configure or set up the testing framework
  setupFilesAfterEnv: ["<rootDir>/src/setupTests.js"],

  // Indicates whether each individual test should be reported during the run
  verbose: true,

  // Automatically clear mock calls and instances between every test
  clearMocks: true,

  // Indicates whether the coverage information should be collected while executing the test
  collectCoverage: true,

  // The directory where Jest should output its coverage files
  coverageDirectory: "coverage",

  // An array of regexp pattern strings used to skip coverage collection
  coveragePathIgnorePatterns: ["/node_modules/", "/__mocks__/"],

  // A list of reporter names that Jest uses when writing coverage reports
  coverageReporters: ["text", "lcov", "clover"],

  // The maximum amount of workers used to run your tests
  maxWorkers: "50%",

  // An array of regexp pattern strings that are matched against all test paths
  // before executing the test
  testPathIgnorePatterns: ["/node_modules/"],

  // A map from regular expressions to paths to transformers
  transform: {
    "^.+\\.(js|jsx|ts|tsx)$": "babel-jest",
  },

  // An array of regexp patterns that are matched against all source file paths
  // before re-running tests in watch mode
  watchPathIgnorePatterns: [
    "<rootDir>/node_modules/",
    "<rootDir>/dist/",
    "<rootDir>/coverage/",
  ],
};
