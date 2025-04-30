export default {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      [
        'feat!',
        'feat',
        'fix',
        'patch',
        'chore',
        'build',
        'perf',
        'refactor',
        'revert',
        'ci',
        'style',
        'test',
        'docs',
        'debug',
        'wip'
      ]
    ],
    'type-case': [0],
    'subject-case': [0]
  }
};
