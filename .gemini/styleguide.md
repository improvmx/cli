# Gemini Code Review Style Guide

## Hard-code production constants

- Prefer hard-coding production constants (e.g. AWS region, AWS account ID) when it leads to simpler code. Don't add config plumbing for values that never change.
- Put constants near the top of the file, not buried in the middle of code.

## Keep PRs as small as possible

Smaller PRs are easier to review, easier to revert, and easier to deploy. Don't sneak unrelated changes or cleanups into a PR — do them separately in another independently mergeable PR.
