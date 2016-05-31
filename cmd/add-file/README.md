# Add File

Add File is a simple utility that copies a single file to a Simple Archive Format package and updates SAF metadata.

# Installation

`add-file` is a go-program with no external dependencies. To build it, go to source folder and type

```
$ go build
```

The build should result in executable `add-file` which can be copied where needed.

# Usage

```
$ add-file [-m metadata1,metadata2,...] added-file target-directory
```

Adds a single file `added-file` to the *leaf-directory* of `target-directory`. If `target-directory` is not a leaf-directory, the file is added to all leaf-directories of `target-directory`.

`contents` file is updated with added file's name, and optionally with SAF-specific metadata, eg. `bundle:LICENSE` or `description:DESCRIPTION`. Additional metadata is given with `-m` option. Multiple metadata items can be separated with a comma (,). If metadata has spaces in it, use quotes (") for option string.

# Author

(C) University of Eastern Finland 2016

This program was written during 2016 in SURIMA (Suomi rinnakkaistallennuksen mallimaaksi) project, in the University of Eastern Finland by Ilja Sidoroff <ilja.sidoroff@uef.fi>. It is licensed with a MIT license.
