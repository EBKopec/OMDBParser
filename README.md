## OMDB Parser

This is the OMDB Parser.

## Before you start

1. Download the file: `title.basics.tsv.gz` from [IMDB](https://datasets.imdbws.com/).
2. Get an API key at [omdbapi](https://www.omdbapi.com/)

What is it doing?

* Read OMDB database file and apply some filters passed as program flags;
* With the filters results a query to OMDB API to get the plot of each movie;
* Once the movie is found a result if showed as below:
```
IMDB_ID     |   Title               |   Plot
tt0000005   |   Blacksmith Scene    |   Three men hammer on an anvil and pass a bottle of beer around. 
```
* if the program is still running and the `maxRunTime` threshold is reached, 
the program will print any results and gracefully exit;
* the program will gracefully exit with no output when stopped with `SIGTERM`;
* If the API Key is not be informed an error message will be shown:
```
You must provide an API Key in the "apiKey" Flag
```
* If the request API limit were reached the following message will be shown:
```
URL: [http://www.omdbapi.com/?apikey=<APIKEY>&&t=&i=tt0002406&plot=short&type=&y=]
Request limit reached! - Unauthorized
```


### Requirements
```
- apiKey string
Should not there be an API key here?
- endYear string
Should not there be an end year here?
- filePath string
Should not there be file here?
- genres string
Should not there be genres here?
- maxApiRequests int
Should not there be max API requests numbers here? (default 1)
- maxRequests int
Should not there be max requests numbers here? (default 1)
- plot string
Should not there be a plot here? (default "short")
- primaryTitle string
Should not there be a primary title here?
- runtimeMinutes string
Should not there be run time minutes here?
- startYear string
Should not there be a start year here?
- timeOut string
Should not there be timeout here? (default "1s")
- titleType string
Should not there be a title type here?
- workers int
Should not there be workers here?
```
### Usage mode
```
go run cmd/main.go -filePath="data.tsv" -timeOut=10s -workers=10 -genres="Drama" -maxApiRequests=10 -plot="short" -apiKey="<APIKey>"
``` 
### Results
``` 
IMDB_ID       |  Title                  |  Plot  
tt0000591     |  The Prodigal Son       |  It's a play in three parts. This film is supposed to be the first long feature film released in Europe  
IMDB_ID       |  Title                    |  Plot  
tt0000615     |  Robbery Under Arms       |  N/A  
IMDB_ID       |  Title        |  Plot  
tt0000630     |  Hamlet       |  Hamlet suspects his uncle has murdered his father to claim the throne of Denmark and the hand of Hamlet's mother, but the prince cannot decide whether or not he should take vengeance.
``` 
