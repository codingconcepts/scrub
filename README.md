# scrub
A simple CLI to securely wipe a file from existence

## Installation

``` bash
$ go get -u github.com/codingconcepts/scrub
```

## Usage

**-d** "debug" overwrites files with random data but does not delete them

**-s** "sweeps" the number of overwrite sweeps (default 10) to perform

### Example

The following example will securly delete `fileOne.txt`, `fileTwo.txt` and all files in `dirOne` and `dirTwo`:

``` bash
$ scrub fileOne.txt fileTwo.txt dirOne dirTwo
```
