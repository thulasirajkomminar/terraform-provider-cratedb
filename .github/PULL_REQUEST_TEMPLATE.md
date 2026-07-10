### Description

<!---
Please provide a helpful description of what change this pull request will introduce.
--->


### Relations

<!---
If your pull request fully resolves and should automatically close the linked issue, use Closes. Otherwise, use Relates.

For Example:

Relates #0000
or
Closes #0000
--->

Closes #0000

### References

<!---
Optionally, provide any helpful references that may help the reviewer(s), e.g. links to the
CrateDB Cloud API documentation (https://console.cratedb.cloud/api/docs).
--->


### Checklist

- [ ] Acceptance tests added/updated for the changed resources or data sources (`make testacc`)
- [ ] Documentation regenerated with `make docs` (edit `templates/` and `examples/`, not `docs/` directly)
- [ ] `CHANGELOG.md` updated under `[Unreleased]`
- [ ] `go build ./...`, `go vet ./...`, `gofmt`, and `make lint` pass

### Output from Acceptance Testing

<!--
Replace TestAccXXX with a pattern that matches the tests affected by this PR.

For more information on the `-run` flag, see the `go test` documentation at https://tip.golang.org/cmd/go/#hdr-Testing_flags.
-->

```console
% make testacc TESTARGS='-run=TestAccXXX'

...
```
