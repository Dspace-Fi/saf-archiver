# Prepare CSV

This is a simple program that transforms a CSV file to CSV form that is required by the `saf-packager` program. Transformation means basically rearranging columns and changing the separator tokens. It is probably useful mainly in the context of University of Eastern
Finland's SoleCRIS to DSpace import -process, but provided here just in case it might provide useful for others.

# Installation

`prepare-csv` is a go-program with no external dependencies. To build it, go to source code folder and type

```
$ go build
```

which should result in executable `prepare-csv`.

# Usage

```
$ prepare-csv config.json input-filename.csv
```

`prepare-csv` requires two input files, one is configuration file (see below) and another is the data file that is to be processed. Data file should be in CSV-format, default separator is ';', but can be specified in the configuration file. Processed CSV is written to stdout-stream an can be redirected to a file, if necessary, e.g.

```
$ prepare-csv config.json input-filename.csv > output-file.csv
```

# Configuration file

Configuration file is a JSON map with following contents:

```
{
    "input-separator": ";",
    "output-separator": ";",
    "split-separator": "||",
    "columns":   [
	{ "from": 0, "to": 0, "title": "solecris.id"},
	{ "from": 1, "to": 1, "title": "dc.title"},
	{ "from": 2, "to": 2, "title": "dc.author", "split-by": ";"},
	{ "from": 9, "to": 3, "title": "dc.identifier.issn"}
    ]
}
```
 * `input-separator` is a string (only first character is relevant) specifying the CSV separator in the input file
 * `output-sepator` is a string used to separate fields in outputted CSV, if the field itself contains this character, its content is escaped with double-quotes
 * `split-separator` is a string used to separate items within fields that have `split-by` definied
 * `columns` is a list containing column-maps
 * A column map is a map containing following keys:
   * `from` an integer (starting from zero) that specifies input column
   * `to` an integer (starting from zero) that specifies to which column input is mapped in output
   * `title` a string specifying the title of this column in output
   * `split-by` a string specifying a string used to separate items within fields in the input file
  
TODO document error conditions.

# Author & License

The program was written during 2016 in SURIMA (Suomi rinnakkaistallennuksen mallimaaksi - Finland for a model country in parallel publishing) -project, in the University of Eastern Finland by Ilja Sidoroff <ilja.sidoroff@uef.fi>. It is licensed with a MIT License.
