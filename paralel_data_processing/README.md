# Paralel Data Processing

This directory contains files that I used to handle large amounts of data and perform any desired computation. This code can be adatped to do perform and computation on large amounts of data. 

I use load bearing mechanisms, go routines, and channels to accompish this task. 

## Getting Started
1) Go [1.20.3](https://go.dev/doc/install)

2) Internet Connection

## List of FIles
- `misc/generate_data.py`: generates senteces to be counted by `main.go` 
- `misc/sentences.csv`: csv file that conatins sentences.
- `main.go`: main code to be ran


## Setup and Run

#### Setup
In order to generate the csv file neccessary to run `main.go`. If you want to add more sentences to the csv file edit the following line of code

```
length = random.randint(1, number_of_sentences)
```

replace `number_of_sentences` with the amount of senteces you want.


#### Run
Run the command `go run main.go` to view the output of the code

The output of the word counts will be diaplyed in console at the end of the codes completion. 