# Gemini Code Review Style Guide

## Hard-code production constants

- Prefer hard-coding production constants (e.g. AWS region, AWS account ID) when it leads to simpler code. Don't add config plumbing for values that never change.
- Put constants that will be changed often or which it's important to know the value of near the top of the file, not buried in the middle of code.
  - But reduce indirection whenever possible. Avoiding defining single-use and/or implementation specific constants that will never change at the top of the file. They should be defined where they're used. Similarly, avoid re-defining single-use variables for no reason, just inline them where they're used.

## Keep PRs as small as possible

Smaller PRs are easier to review, easier to revert, and easier to deploy. Don't sneak unrelated changes or cleanups into a PR — do them separately in another independently mergeable PR.
