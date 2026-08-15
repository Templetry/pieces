# Dependency updates

This project uses [Renovate](https://docs.renovatebot.com) to keep dependencies current.

## How it behaves

| Update kind | What happens |
|---|---|
| minor / patch | Grouped into a single scheduled PR |
| major | One PR each, labelled, always read by a human |
| security advisory | Opened immediately, ignoring the schedule |
| lock files | Refreshed monthly |

The schedule and automerge policy live in `renovate.json`. A dependency dashboard issue lists everything pending.

## Enabling it

Install the [Renovate GitHub App](https://github.com/apps/renovate) on the repository (or run the self-hosted CLI in CI). Nothing else is needed — the configuration is already here.

## Why this is not optional work

A project's dependencies rot whether or not anyone is looking. The failure mode is not "an update is missing": it is discovering, on the day of an urgent fix, that upgrading requires crossing three major versions at once.
