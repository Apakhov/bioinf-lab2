# bioinf-lab2
bioinf lab2

## Build

```bash
go build -o bld/amino *.go
```

## Usage

```bash
./bld/amino {-flag [val]} file [file2]
```

## Flags

```
-g -gap float
    gap value (default -2)
-no-color
    disables colored output in cosole
-no-connections
    disables connections in output
-o -out string
    output file
-t -type string
    table type, one of Blosum64, DNA, Default (default "Default")
-oa -outalignment uint
    alignment of result sequences, if 0 no alignment used
-log-time
    print time of processing in log
-threads int
    amount of threads for computing, for optimal speed use available amount of cpu (default 8)
```
