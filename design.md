# Design

## b+ tree

目前的实现，是有问题的，root 最开始就是叶子节点，可以插入数据，达到 full 状态，分裂后，才会成为内部节点，目前的实现是错误的。

### 内部节点分裂

![splitting](./splitting.png)
