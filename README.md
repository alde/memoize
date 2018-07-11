# Memoize

## Purpose
Memoize slow requests, and return the cached response instead.

Useful when trying to assemble a `jq` query, or a `sed` nightmare string.

## Usage
```
$ time ./memoize curl "https://artifactory.internal/artifactory/api/search/artifact?name=many-whelps" | jq '.results[-1].uri' | sed "s|.*/\([0-9\.]*\)/.*|\1|"
13.1.105

real    0m2.156s
user    0m0.050s
sys     0m0.034s

$
# Run again
$ time ./memoize curl "https://artifactory.internal/artifactory/api/search/artifact?name=many-whelps" | jq '.results[-1].uri' | sed "s|.*/\([0-9\.]*\)/.*|\1|"
13.1.105

real    0m0.015s
user    0m0.011s
sys     0m0.014s
$
```

To clear cache, run `memoize clear`.

## Todo
* Clear only a specific cache, not all.
* TTL caches, make new fresh request if some time has passed.

