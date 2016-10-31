## Serializer

This project contains a script which will convert a TSV file or set of YAML files in the formats supplied by the RDA into a Serialized Register. Output is sent to STDOUT.

- timestamps are added with the time the SR is created.
- item JSON is canonicalized.
- fields JSON is used to format fields with Cardinality 'n'
- quotes all allowed in the body of a field and will be escaped in JSON
- if the first or last character of a field is a quote then the whole field should be surrounded by quotes e.g ""Go" Cafe"

### Building

Assuming you checkout the project into a Go workspace i.e. $GOPATH is defined with sub directories *bin, src, pkg* and $GOPATH/bin is on your PATH.

One extra package is required for YAML parsing.

    >go get gopkg.in/yaml.v2
    >cd $GOPATH/src

If this directory does not exist, create it:

    >mkdir github.com/openregister/
    >git clone git@github.com:openregister/serializer.git
    >go install github.com/openregister/serializer

Will build the script and put an executable file for your architecture in $GOPATH/bin

Alternatively, from any directory, but note this will clone the repo over https :

    >go get github.com/openregister/serializer

### Tests

Tests are in **.../openregister/serializer_test.go**. To run

    >cd [path to ]/openregister/serializer/
    >go test

### Usage

Pass the argument **tsv** or **yaml**

Pass the path to the **fields** JSON; and the TSV files/ YAML directory to be loaded as arguments.

e.g.

    >cd openregister (where register data is located)
    >serializer tsv field-records.json address/address.tsv
    >serializer yaml field-records.json registry-data/data/beta/register
