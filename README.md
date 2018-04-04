# SCOIR Technical Interview for Back-End Engineers
This repo contains an exercise intended for Back-End Engineers.

## Instructions
1. Fork this repo.
1. Using technology of your choice, complete [the assignment](./Assignment.md).
1. Update this README with
    * a `How-To` section containing any instructions needed to execute your program.
    * an `Assumptions` section containing documentation on any assumptions made while interpreting the requirements.
1. Before the deadline, submit a pull request with your solution.

## Expectations
1. Please take no more than 8 hours to work on this exercise. Complete as much as possible and then submit your solution.
1. This exercise is meant to showcase how you work. With consideration to the time limit, do your best to treat it like a production system.


## How-to

1. [Install the go programming language on your OS of choice](https://golang.org/doc/install). 
2. Add go/bin to your $PATH
3. Install with `go install https://github.com/stvpalumbo/be_exam_candidate_el`
4. Run as follows:

```bash
    be_exam_candia <in_path> <out_path> <err_path>
```

## testing

1. Download the repository with `go get https://github.com/stvpalumbo/be_exam_candidate_el`
2. Navigate to the directory containing the source code (`$gopath/src/github.com/stvpalumbo/be_exam_candidate_el`)
3. Test with `go test --cover`

where in, out, and err are absolute or relative directories to where csv_to_json is

##Assumptions:

###Program Design

REQUIREMENT: csv files will be placed into input-directory
once the application starts it watches input-directory for any new files that need to be processed
files will be considered new if the file name has not been recorded as processed before.

INTERPRETATION: I am assuming this is a meant to run as a background process that continually monitors it's input folder, rather than a script. As such, it runs continuously for DEFAULT_TIMEOUT seconds, then exits. There is currently no proceedure to kill the process outside of the OS. This 

This application does not preserve state across runs. That is, during a single run, only a single file named "foo.csv" in the input-directory will be processed. However, if the application is restarted, it will process "foo.csv" again.

### Error Handling

The basic approach I have taken is that parsing errors are recoverable, but filesystem errors are not. The program is expected to run into improperly formatted .csvs

###Validation

#### INTERNAL_ID
REQUIREMENT: 8 digit positive integer. Cannot be empty.
INTERPRETATION: The term "8-digit positive integer" is a little vague. Is 00000001 the first 8-digit positive integer, or is 10000000? I have chosen to only accept integers in the half-open interval [10^7, 10^8).

####PHONE-NUM
REQUIREMENT: string that matches this pattern ###-###-####. Cannot be empty.
INTERPRETATION: matches this pattern ###-###-####, where # _is a digit_. Specifically, I use the regex `^[0-9]{3}-[0-9]{3}-[0-9]{4}$`