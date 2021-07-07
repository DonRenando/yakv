# Contributing Guide

Thank you for considering contributing to yakv. 😀

yakv is an open source project and we love to receive contributions from our community! ❤

<h2>Table of Contents:</h2>

- [Contributing Guide](#contributing-guide)
  - [How can I contribute?](#how-can-i-contribute)
  - [Code of Conduct](#code-of-conduct)
  - [How to make your first contribution!](#how-to-make-your-first-contribution)
    - [Make a Pull Request](#make-a-pull-request)
  - [Reporting Issues](#reporting-issues)
  - [Benchmarking](#benchmarking)

## How can I contribute?

**There are many ways to contribute:**
- [Issues](https://github.com/burntcarrot/yakv/issues)
- Documentation
- Tutorials/Blog Posts
- [Feature Requests](https://github.com/burntcarrot/yakv/discussions/categories/feature-requests)

## Code of Conduct

Read the [Code of Conduct](https://github.com/burntcarrot/yakv/blob/main/CODE_OF_CONDUCT.md).

Following these guidelines helps to communicate that you respect the time of the developers managing and developing this open source project.

## How to make your first contribution!

Uncertain about where to start contributing to yakv?

You can always look at **[open issues](https://github.com/burntcarrot/yakv/issues?q=is%3Aissue+is%3Aopen+)**, but if you are a complete beginner, check out these issues:

- [Good First Issues](https://github.com/burntcarrot/yakv/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22): These issues are suitable for beginners.
- [Up for Grabs](https://github.com/burntcarrot/yakv/issues?q=is%3Aissue+is%3Aopen+label%3A%22up+for+grabs%22): These issues are often low-effort, intended for beginners.
- [Help Wanted](https://github.com/burntcarrot/yakv/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22): These issues require more attention and are crucial for a smooth development of yakv.

> ### New to Github? 🤔
> Here are some tutorials which might help you in making a contribution:
> - [Make a Pull Request](http://makeapullrequest.com/)
> - [First Timers Only](http://www.firsttimersonly.com/)


### Make a Pull Request

- Fork yakv.
- Create a new branch for the fix. (for example, `fix-api`)
- Make required changes to the code (fixing bugs, typos, etc.)
- **(Recommended) Run tests using `go test`.**
- **(Recommended) Install golangci-lint. Run `golangci-lint run` locally before making a pull request.**
- Commit changes.
- Make a pull request with a suitable title and description about the fix.
- Request for review (if your PR is ready to merge).
- Wait for your PR to get merged. ⌚

At this point, you're ready to make your first contribution to yakv! Feel free to ask for help, have fun! 🥳


## Reporting Issues

When filing an issue, make sure to answer these questions:

1. What version of Go are you using (go version)?
2. What operating system and processor architecture are you using?
3. What did you do?
4. What did you expect to see?

Attach the output of the error/bug enclosed in code blocks.

## Benchmarking

Benchmarks are done using [vegeta](https://github.com/tsenart/vegeta).

If you are submitting a benchmark, please specify:

- Device Specifications:
    - Operating System
    - Processor Architecture
    - Processor
    - RAM
    - Other device specifications
- Benchmark title
- Method (for example, `GET`, `PUT`, etc.)
- Rate (for example, `50 requests/seconds`)
- Attack Report
- Attack Plot
