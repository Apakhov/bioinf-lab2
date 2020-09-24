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
-g float
    gap value (default -2)
-gap float
    gap value (default -2)
-no-color
    disables colored output in cosole
-no-connections
    disables connections in output
-o string
    output file
-out string
    output file
-t string
    table type, one of Blosum64, DNA, Default (default "Default")
-type string
    table type, one of Blosum64, DNA, Default (default "Default")
```
