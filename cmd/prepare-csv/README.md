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

`prepare-csv` requires two input files, one is configuration file (see below) and another is the data file that is to be processed. Data file should be in CSV-format, default separator is ';', but can be specified in the configuration file. The file should not contain a header file, i.e. all lines are processed. Processed CSV is written to stdout-stream an can be redirected to a file, if necessary, e.g.

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
	{ "from": 0, "title": "solecris.id"},
	{ "from": 1, "title": "dc.title"},
	{ "from": 2, "title": "dc.author", "split-by": ";"},
    { "from": 3, "title": "dc.language.iso", "filters": ["uef.isolang"]},
    { "from": 9, "title": "dc.identifier.issn"},
    { "from": 12, "discard": true, "title": "dc.identifier.issue"},
    ]
}
```
 * `input-separator` is a string (only first character is relevant) specifying the CSV separator in the input file (default: ";")
 * `output-sepator` is a string used to separate fields in outputted CSV, if the field itself contains this character, its content is escaped with double-quotes (default: ";")
 * `split-separator` is a string used to separate items within fields that have `split-by` definied (default: ";")
 * `columns` is a list containing column-maps
 * A column map is a map containing following keys:
   * `from` an integer (starting from zero) that specifies input column
   * `discard` an boolean, if true discards that column (useful in temporarily disabling column, as JSON doesn't have comments)
   * `title` a string specifying the title of this column in output
   * `split-by` a string specifying a string used to separate items within fields in the input file
   * `filters` a list of strings specifying names of filters columns are filtered with. Filtering takes places after replacing the splitter string (`split-by`) and are applied in the order they are in the list. The up-to-date names can be found in the source code file `filter.go` and they are listed also below (hopefully up-to-date as well):
     * `uef.isolang` replace language string with its ISO-639-1 code, eg. "suomi" -> "FI". Source languages are primary those found in UEF's SoleCRIS system.
     * `uef.peerreview` peer review status (eprint.status), map 0/1 to either http://purl.org/eprint/status/PeerReviewed or http://purl.org/eprint/status/NonPeerReviewed
	 * `uef.type` tries to map document types used in UEF's SoleCRIS system into ePrintTypes.
	 * `uef.doi` tries to format dois into urls (10.1111/etc -> http://doi.org/doi:10.1111/etc) if it seems likely to succeed.
Columns are output in the order they are in the `columns` list.

# Author & License

The program was written during 2016 in SURIMA (Suomi rinnakkaistallennuksen mallimaaksi - Finland for a model country in parallel publishing) -project, in the University of Eastern Finland by Ilja Sidoroff <ilja.sidoroff@uef.fi>. It is licensed with a MIT License.
