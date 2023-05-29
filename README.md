# blink-tree-go

blink tree implementation in go

## known issues

- sometimes crashes in tests that perform concurrent deletion and insertion
  - `TestBLTree_deleteManyConcurrently` in `bltree_test.go`
- failed to restart

## reference

- https://github.com/PLW/blink-tree-logic

## related papers

- [Efficient locking for concurrent operations on B-trees](https://dl.acm.org/doi/10.1145/319628.319663)
- [Concurrency control and recovery for balanced Blink trees](https://www.researchgate.net/journal/The-VLDB-Journal-0949-877X)
- [A Blink Tree method and latch protocol for synchronous node deletion in a high
  concurrency environment](https://arxiv.org/ftp/arxiv/papers/1009/1009.2764.pdf)