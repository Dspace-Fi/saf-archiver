# SAF Archiver

SAF Archiver is a program for creating Simple Archive Format archives
(see e.g. [https://wiki.duraspace.org/display/DSDOC6x/Importing+and+Exporting+Items+via+Simple+Archive+Format]). It creates SAF archives from CSV-files and was developed for University of Eastern Finland's DSpace-instance, but is released here, in case it might prove useful for others.

There is also another tool for the same purpose: https://wiki.duraspace.org/display/DSPACE/Simple+Archive+Format+Packager
which might be more suitable for generic purposes.

# Installation

`saf-archiver` is a go-program with no external dependencies. To build it, go to source folder and type

```
$ go build
```
which should result in executable `saf-archiver`

# Usage

```
$ saf-archiver input-filename.csv output-directory
```

`saf-archiver` requires one input file, csv-datafile containing the imported information. You can use `prepare-csv` -program to format the input file. Another required parameter is output directory *that should not exist* before program is run. Created directory should contain a simple archive formatted package created from the inputted data, which can be optionally zipped and imported to DSpace with `[dspace]/bin/dspace import` command. `saf-archiver` handles only metadata - you can add files to generated archive with `add-file` command.

CSV.input file format is as follows. Input separator for fields is ';', if the field itself contains ';', the field can be put in double quotes. If field contains multiple values, e.g. dc.contributor.author -field can contain many authors; they can be separated with "||" (two bars).

# TODO

 * Configurable input and value separators
 * Support for metadata language qualifier

# Author

(C) University of Eastern Finland 2016

The program was written during 2016 in SURIMA (Suomi rinnakkaistallennuksen mallimaaksi) project, in the University of Eastern Finland by Ilja Sidoroff <ilja.sidoroff@uef.fi>. It is licensed with a MIT License.

