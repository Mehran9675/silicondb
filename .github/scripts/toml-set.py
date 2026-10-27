#!/usr/bin/env python3
"""Sets key = value (value as TOML, e.g. '"text"' or 3) in a section of
giode.toml, replacing the key if the section has it and adding it after
the section header otherwise."""
import re
import sys

section, key, value = sys.argv[1:4]
path = sys.argv[4] if len(sys.argv) > 4 else "giode.toml"
lines = open(path).read().split("\n")
start = lines.index(f"[{section}]")
end = next((i for i in range(start + 1, len(lines)) if lines[i].startswith("[")), len(lines))
for i in range(start + 1, end):
    if re.match(rf"^{re.escape(key)}\s*=", lines[i]):
        lines[i] = f"{key} = {value}"
        break
else:
    lines.insert(start + 1, f"{key} = {value}")
open(path, "w").write("\n".join(lines))
