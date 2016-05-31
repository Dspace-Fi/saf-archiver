# Tools for building Simple Archive Format files

Utilities for building Simple Archive Format files suitable for DSpace import. See directory `cmd` and respective commands for more information.

Build e.g. with

```
$ go build cmd/prepare-csv/prepare-csv.go
$ go build cmd/saf-archiver/saf-archiver.go
$ go build cmd/add-file/add-file.go
```

`prepare-csv` is a program for manipulating CSV files which can then be transformed into Simple Archive Format package with `saf-archiver`. `saf-archiver` does not handle files, only metadata information, but files can be added to generated archive with `add-file` program.

# Licence

The programs are (C) 2016 University of Eastern Finland and are licensed with a MIT licence.
