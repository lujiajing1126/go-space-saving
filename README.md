# Spacing-Saving algorithm

> Spacing-Saving algorithm is an integrated approach for solving both problems of finding the most popular k elements,
> and finding frequent elements in a data stream.

## Usage

```go
// epsilon determines the total number of counters
ss, err := spaceSaving.NewStreamSummary(0.01)
if err != nil {
//...
}
ss.Record(item)
```

## Acknowledge

1. Ahmed Metwally et
   al., [Efficient computation of frequent and top-k elements in data streams](https://doi.org/10.1007/978-3-540-30570-5_27)
   . In Proceedings of the 10th international conference on Database Theory (ICDT'05)
2. Java Implementation of Spacing-Saving algorithm: https://github.com/fzakaria/space-saving