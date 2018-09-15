# curlsh - better than `curl <url> | sh`

There is a common installation method that is quite controversial. Run
`curl https://url | sh`.

This installation method is dangerous and insecure.

This installation method is not going away.

This project's aim is to make things a bit better.

## Usage

```
Usage of ./curlsh:
  -hash value
        SRI hash
  -pager string
        select pager (CURLSH_PAGER, PAGER) (default "less -R")
  -sudo
        run the script with sudo
  -trusted
        whenver the script is trusted
  -url value
        URL to fetch
```

## Example

```
$ ./curlsh -url https://zimbatm.github.io/curlsh/sri_test.js \
  -hash "sha256-ySadHRVML1LfcwlPIxXx4CQpk64arq0Yv32cBpu9CFQ="
```

## Features

### No timing attachs

Because the script is fully fetched before being executed.

TODO: add reference

### Secure by default

Nudges the user towards the right things: read the script and check the hashes

## ChangeLog

* [CHANGELOG](CHANGELOG.md)

## Research

<https://github.com/chrisgreg/sri-gen>
