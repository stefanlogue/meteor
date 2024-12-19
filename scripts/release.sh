#!/bin/bash

set -e

export TAG=$(svu next --force-patch-increment --strip-prefix)

read -p "Creating new release for v$TAG. Do you want to continue? [Y/n] " prompt

if [[ $prompt == "y" || $prompt == "Y" || $prompt == "yes" || $prompt == "Yes" ]]; then
  python scripts/prepare_changelog.py
  git add CHANGELOG.md
  git commit -m "chore: bump version to $TAG for release" || true && git push
  echo "Creating new git tag $TAG"
  git tag "v$TAG" -m "v$TAG"
  git push --tags
else
  echo "Cancelled"
  exit 1
fi
