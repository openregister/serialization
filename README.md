# serialization

This project contains a script which will convert a TSV file in the format supplied by the RDA into a Serialized Register.

- Timestamps are added with the time the SR is created.
- Item JSON is canonicalized.
- fields JSON is used to format fields with Cardinality 'n'
- quotes all allowed in the body of a field and will be escaped in JSON
- if the first or last character of a field is a quote then the whole field should be surrounded by quotes e.g ""Go" Cafe"

### Building

Assuming you checkout the project into a Go workspace i.e. $GOPATH is defined and $GOPATH/bin is on your PATH.

One extra package is required for YAML parsing.

>go get gopkg.in/yaml.v2

>cd $GOPATH/src

TODO check into github

>git clone git@bitbucket.org:john_ollier/register-serializer.git

>go install register.gov.uk/register-serializer

### Usage

Pass the argument 'tsv' or 'yaml'

Pass the path to fields JSON and TSV files/ YAML directory to be loaded as arguments.

e.g.

>cd openregister

>register-serializer tsv field-records.json address/address.tsv

>register-serializer yaml field-records.json registry-data/data/beta/register
