#!/usr/bin/env bash
# Sets the app version from the release tag ($TAG, e.g. v1.2.0) and the
# build number from the workflow run, in giode.toml and the app's config.
# Tags with more or fewer than three numbers are normalized: v1.2 is
# 1.2.0 and v1.0.0.0 is 1.0.0.
set -euo pipefail
version="$(python3 - "$TAG" <<'PY'
import re, sys
parts = re.findall(r"\d+", sys.argv[1])[:3]
print(".".join((parts + ["0", "0", "0"])[:3]))
PY
)"
.github/scripts/toml-set.py app version "\"$version\""
.github/scripts/toml-set.py app build "$GITHUB_RUN_NUMBER"
python3 - "$version" <<'PY'
import re, sys
p = "internal/config/config.go"
s = open(p).read()
s = re.sub(r'^var Version = ".*"$', f'var Version = "{sys.argv[1]}"', s, flags=re.M)
open(p, "w").write(s)
PY
echo "Version $version, build $GITHUB_RUN_NUMBER"
