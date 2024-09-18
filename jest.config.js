module.exports = {
  testEnvironment: 'node',
  roots: ['<rootDir>/test', '<rootDir>/integ-tests'],
  testMatch: ['**/*.test.ts'],
  transform: {
    '^.+\\.tsx?$': 'ts-jest'
  }
};
