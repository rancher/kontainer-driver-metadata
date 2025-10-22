#!/usr/bin/awk -f
# Usage: channels-minchan-for-version.awk <channels-file> <exact-version>
# Prints the minChannelServerVersion for that release (if present), else nothing.

function trim(s){ sub(/^[ \t\r\n]+/,"",s); sub(/[ \t\r\n]+$/,"",s); return s }

BEGIN {
    if (ARGC < 3) exit 2
    target = ARGV[2]
    ARGV[2] = ""             # prevent awk from treating target as a file
    inBlock = 0
    curVer  = ""
}

# Start of a release block: "- version: <ver>"
/^[[:space:]]*-[[:space:]]*version:[[:space:]]*/ {
    inBlock = 1
    line = $0
    sub(/^[[:space:]]*-[[:space:]]*version:[[:space:]]*/, "", line)
    curVer = trim(line)
    next
}

# Capture minChannelServerVersion only if we are in the target block
inBlock && curVer == target && /^[[:space:]]*minChannelServerVersion:[[:space:]]*/ {
    line = $0
    sub(/^[[:space:]]*minChannelServerVersion:[[:space:]]*/, "", line)
    print trim(line)
    exit 0
}