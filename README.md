<br/>

<p align="center">
  <img src="https://github.com/user-attachments/assets/2ceb1e71-6b8e-4a4a-bf91-b0d303baea28" width="158" height="158" alt="Gopher"/>
</p>

<div id="user-content-toc" align="center">
  <ul>
    <summary><h1 style="display: inline-block;">Go Modkit</h1></summary>
  </ul>
</div>

<p align="center">
  <a aria-label="License" href="https://github.com/akalanka47000/go-modkit/blob/main/LICENSE">
    <img alt="" src="https://img.shields.io/badge/License-MIT-yellow.svg">
  </a>
   <a aria-label="License" href="https://turborepo.com">
    <img alt="" src="https://img.shields.io/badge/Maintained%20with-Turbo-f700ff.svg?style=flat-square">
  </a>
  <a aria-label="CI Release" href="https://github.com/akalanka47000/go-modkit/actions/workflows/tests.yml">
    <img alt="" src="https://github.com/akalanka47000/go-modkit/actions/workflows/tests.yml/badge.svg">
  </a>
</p>

<br/>

A monorepository designed to house go modules which make developer life a bit easier

Go Modkit strictly follows [Semantic versioning](https://semver.org). All packages are versioned independently and their stable versions start from `<module>/v1.0.0`. **No breaking changes will be introduced before any major version upgrades**. All packages are available under the `latest` tag, and for some packages, the `blizzard` tag if present indicates it's a pre-release version.


## Getting started

### Prerequisites

- [Go](https://golang.org/doc/install) version 1.23 or higher (recommended)
- [Bun](https://bun.sh/docs/installation) or [Node.js](https://nodejs.org/en/download/) if you want to utilize the power of [Turborepo](https://turbo.build/repo/docs) during development (Optional)

### Common commands

- Run `bun install` to install all dependencies // or alternatively `go get` for each module
- Run `bun run test` to run all test suites
- Run `bun format` to format all packages
- Run `bun lint` to run all linters

## Commit messages

- We follow conventional commits during our development workflow as a good practice. More information can be found at their official [documentation](https://www.conventionalcommits.org/en/v1.0.0-beta.4/#examples)
- Refer the [commitlint.config.js](https://github.com/akalanka47000/go-modkit/blob/main/commitlint.config.js) file for a full list of supported commit message prefixes
