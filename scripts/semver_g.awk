#!/usr/bin/awk -f
# Exit 0 if A>B by major.minor.patch, ignoring pre-release suffix; else 1.
# Usage: scripts/semver_g.awk <A> <B>
function norm(s,  n,i,arr) {
    gsub(/^v/,"",s); sub(/-.*/,"",s)
    n = split(s, arr, ".")
    for (i=1; i<=3; i++) if (i>n) arr[i]=0
    return sprintf("%d.%d.%d", arr[1], arr[2], arr[3])
}
BEGIN {
    if (ARGC < 3) exit 2
    as = norm(ARGV[1]); bs = norm(ARGV[2])
    split(as,a,"."); split(bs,b,".")
    if (a[1]>b[1] || (a[1]==b[1] && (a[2]>b[2] || (a[2]==b[2] && a[3]>b[3])))) exit 0
    exit 1
}
