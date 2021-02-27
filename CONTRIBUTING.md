# Contributing

When contributing to this repository, please first discuss the change you wish
to make via issue, email, or any other method with the owners of this repository
before making a change.

Please note we have a code of conduct, please follow it in all your interactions
with the project.

## Issues

Issues should be used to report problems, request a new feature, or to discuss
potential changes before a PR is created. When you create a new Issue, a
template will be loaded that will guide you through collecting and providing the
information we need to investigate.

If you find an Issue that addresses the problem you're having, please add your
own reproduction information to the existing issue rather than creating a new
one. Adding a reaction can also help be indicating to our maintainers that a
particular problem is affecting more than just the reporter.

## Pull Requests

1. Ensure any intermediary artifacts (e.g. generated test reports, example
output, build archives) are not committed, and are properly ignored on both
`.gitignore` and `.dockerignore`.
2. For docker images, ensure any install or build dependencies are not present
on the final target.
3. Update the `README.md` and docs with details of changes to code components
such as interfaces, objects, variables, packages and logic changes.
4. For docker images, update the `README.md` and docs with details of changes to
build targets, environment variables, exposed ports, volumes, build arguments,
useful file locations and how to use the final image.
5. Increase the version numbers in any examples files and the `README.md` to the
new version that this Pull Request would represent. The versioning scheme we use
is [SemVer](http://semver.org/).
6. You may merge the Pull Request in once workflows are all passing, the target
branch can be rebased onto your branch and you have the sign-off of at least one
developer. If you do not have permission to merge, you may request a reviewer to
merge it for you.
