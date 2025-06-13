import type { UserConfig } from 'cz-git'

/** @type {import('cz-git').UserConfig} */
const config: UserConfig = {
    extends: ['@commitlint/config-conventional'],
    prompt: {
        alias: {
            fd: 'docs: fix typos'
        },
        messages: {
            type: 'Select the type of change you are committing:',
            scope: 'Select a scope (optional):',
            customScope: 'Enter a custom scope:',
            subject: 'Write a short, concise description of the change:\n',
            body: 'Provide a more detailed description (optional). Use "|" for line breaks:\n',
            breaking: 'List any breaking changes (optional). Use "|" for line breaks:\n',
            footerPrefixesSelect: 'Select an Issue prefix (optional):',
            customFooterPrefix: 'Enter a custom Issue prefix:',
            footer: 'List related Issues (optional), e.g., #1, #2:\n',
            confirmCommit: 'Submit or modify the commit?'
        },
        types: [
            {
                value: 'feat',
                name: 'feat:      ✨ A new feature.',
                emoji: ':sparkles:'
            },
            {
                value: 'fix',
                name: 'fix:       🐛 A bug fix.',
                emoji: ':bug:'
            },
            {
                value: 'docs',
                name: 'docs:      📝 Documentation only changes.',
                emoji: ':memo:'
            },
            {
                value: 'style',
                name: 'style:     💄 Changes that do not affect the meaning of the code.',
                emoji: ':lipstick:'
            },
            {
                value: 'refactor',
                name: 'refactor:  ♻️  A code change that neither fixes a bug nor adds a feature.',
                emoji: ':recycle:'
            },
            {
                value: 'perf',
                name: 'perf:      ⚡️ A code change that improves performance.',
                emoji: ':zap:'
            },
            {
                value: 'test',
                name: 'test:      ✅ Adding missing tests or correcting existing tests.',
                emoji: ':white_check_mark:'
            },
            {
                value: 'build',
                name: 'build:     📦️ Changes that affect the build system or external dependencies.',
                emoji: ':package:'
            },
            {
                value: 'ci',
                name: 'ci:        🎡 Changes to our CI configuration files and scripts.',
                emoji: ':ferris_wheel:'
            },
            {
                value: 'revert',
                name: 'revert:    ⏪️ Revert to a commit.',
                emoji: ':rewind:'
            },
            {
                value: 'chore',
                name: 'chore:     🔨 Other changes that do not modify src or test files.',
                emoji: ':hammer:'
            }
        ],
        useEmoji: true,
        emojiAlign: 'center',
        useAI: false,
        aiNumber: 1,
        themeColorCode: '',
        scopes: [],
        allowCustomScopes: true,
        allowEmptyScopes: true,
        customScopesAlign: 'bottom',
        customScopesAlias: 'custom',
        emptyScopesAlias: 'empty',
        upperCaseSubject: false,
        markBreakingChangeMode: false,
        allowBreakingChanges: ['feat', 'fix'],
        breaklineNumber: 100,
        breaklineChar: '|',
        skipQuestions: [],
        issuePrefixes: [
            { value: 'link', name: 'link:     Link to ongoing ISSUES' },
            { value: 'closed', name: 'closed:   Mark ISSUES as completed' }
        ],
        customIssuePrefixAlign: 'top',
        emptyIssuePrefixAlias: 'skip',
        customIssuePrefixAlias: 'custom',
        allowCustomIssuePrefix: true,
        allowEmptyIssuePrefix: true,
        confirmColorize: true,
        maxHeaderLength: Infinity,
        maxSubjectLength: Infinity,
        minSubjectLength: 0,
        scopeOverrides: undefined,
        defaultBody: '',
        defaultIssues: '',
        defaultScope: '',
        defaultSubject: ''
    }
}

export default config
