## Serializer

This project contains a script which will convert a TSV file or YAML files in a directory in the format supplied by the RDA into a Serialized Register. Output is written to STDOUT.

- timestamps are added with the time the SR is created.
- item JSON is canonicalized.
- fields JSON is used to format fields with Cardinality 'n'
- quotes are allowed in the body of a field and will be escaped in JSON
- if the first or last character of a field is a quote then the whole field should be surrounded by quotes e.g ""Go" Cafe"

### Building

Install Go if necessary.

    >brew install go

Standard practice is to use one directory as your Go workspace and define an environment variable $GOPATH pointing to it. In that directory there should be sub directories *bin, src & pkg*. The $GOPATH/bin directory should be added to your $PATH.

Assuming you are going to check out the project into a Go workspace:

One extra package is required for YAML parsing.

    >go get gopkg.in/yaml.v2
    >cd $GOPATH/src

    If this directory does not exist, create it:

    >mkdir -p github.com/openregister/
    >cd github.com/openregister
    >git clone git@github.com:openregister/serializer.git
    >go install github.com/openregister/serializer

Will build the script and put an executable file for your architecture in $GOPATH/bin.

Alternatively, from any directory, but note this will clone the repo over https :

    >go get github.com/openregister/serializer

### Tests

Tests are in **.../openregister/serializer_test.go**. To run

    >cd [path to ]/openregister/serializer/
    >go test

### Usage

Pass the argument **tsv** or **yaml**

Pass the paths to the **fields** JSON; and the TSV files/ YAML directory to be loaded as arguments.

You need to have a copy of the fields JSON in a file e.g.

    >curl http://field.discovery.openregister.org/records.json > field-records.json

Then e.g.

    >cd openregister (where register data is located)
    >serializer tsv field-records.json address/address.tsv
    >serializer yaml field-records.json registry-data/data/beta/register
