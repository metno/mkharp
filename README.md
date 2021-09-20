# Create observation files for harp

This application creates [harp](https://github.com/harphub/harp) observation files from [csv](https://en.wikipedia.org/wiki/Comma-separated_values) files. It is meant for making harp more accessible when you have observation data in a non-standard format, by providing a conversion from a very simple csv format to the more complex file format expected by harp.

## Installation

The best way to install this is by downloading a binary release from the releases page. Detailed instructions for how to do this is available [here](docs/installation.md).

## CSV file format

The application works on csv files, containing observation data. They look like this:

```csv
time, T2m, AccPcp24h
2021-09-17T06:00:00Z, 12.4, 0.2
2021-09-17T07:00:00Z, 12.7,
```

The first line must contain a header, describing the data to come. The first element _must_ be `time`, but the rest can be any valid harp parameter name. The following lines contains the data - with the first element containing the observation time, formatted using the [RFC3339](https://duckduckgo.com/?t=ffab&q=RFC3339&ia=web) format. The following values are floats, and can be empty to signal missing data.

Note that most spreadsheets have a feature for exporting data as csv, so working with data in spreadsheets should be fairly easy.


## Invocation

A more complete example of how to use this application can be found [here](docs/usage.md).

The application takes several command-line arguments, many of which are mandatory. Run the application with `--help` to see the full list of arguments.

An example invocation may look like this:

```bash
$ mkharp -create -sid 1492 -elevation 94 -lat 10.72 -lon 59.9423 -obstype synop -out OBSTABLE_2021.sqlite < data.csv
```

This creates a file, called `OBSTABLE_2021.sqlite`, populating with data for a station with id `1492`, and the given latitude, longitude and elevation. Observation type is set to synop, and the actual observation data is read from a file, called `data.csv`. The format of the data file is described above.

Note that this command fails if the output file already exists with the given obstype. Drop the `-create` option if you wish to add data to an existing file.
