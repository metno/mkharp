# Usage example

This is an example of how to use the `mkharp` application. Is assumes that you have [installed](installation.md) the latest version of `mkharp.

## Step one: Create input file

Open your favourite text editor, and create a file, called `data.csv`. Fill in some data, like this:

```
time, T2m, AccPcp24h
2021-09-17T06:00:00Z, 12.4, 0.2
2021-09-17T07:00:00Z, 12.7,
```

This contains two observations, with air temperature (T2m), and precipitation, accumulated over 24 hours (AccPcp24h). Note that it is not mandatory to give values for all parameters for all timesteps, but the number of commas `,` must be the same on every line.

## Step two: Create a harp observation file.

Run `mkharp` like this:

```bash
mkharp \
    -create \
    -sid 1492 \
    -elevation 94 \
    -lat 10.72 \
    -lon 59.9423 \
    -obstype synop \
    -out OBSTABLE_2021.sqlite \
    < data.csv
```

This will create a harp observation file, called `OBSTABLE_2021.sqlite`, and populate it with data from the previously created file, `data.csv`. The command line arguments specify other metadata for the observations. In this case, we say that this is `synop` data for a station with id `1492`, and the given latitude, longitude and elevation.

## Step three: Add more data to the same file (optional)

If you create another file with observations, either for other times, or other stations, you can run the same command as before, while dropping the `-create` command line argument. For example, if you have another observation file, called `more_data.csv`, you can add that to an existing harp file, using the following command:

```bash
mkharp \
    -sid 1492 \
    -elevation 94 \
    -lat 10.72 \
    -lon 59.9423 \
    -obstype synop \
    -out OBSTABLE_2021.sqlite \
    < more_data.csv
```

## Step four: Copy observation file into harp's expected directory structure

Harp expects files to be placed in a particular file structure, with observation files placed under `OBSTABLE/OBSTABLE_{YYYY}.sqlite`. Copy the files into the correct directory and you should be ready to verify your forecasts. Refer to the [harp tutorial](https://harphub.github.io/harp_tutorial/index.html) for details about how to do this.
