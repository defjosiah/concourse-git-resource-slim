# concourse-git-resource-slim

This uses the github api to do the concourse git resource
This uses git to keep track of versions and ordering, but it is limited so that
it is very quick.

The /check command uses the github list commits api
The /in command uses the github tarball api

## Running

Put your GITHUB_TOKEN into the environment

```bash
export GITHUB_TOKEN=$(gh auth token)
```

Run any of the `./inputs/*.json` files, make sure to put the `jq` command after
it runs, so you can ensure that you are properly outputting to stdout vs. stderr
for logs.

```bash
GITHUB_TOKEN=$(gh auth token); cat ./inputs/check-input-version.json | jq --arg token "$GITHUB_TOKEN" '.source["auth-token"] = $token' | make run-check | jq
```
