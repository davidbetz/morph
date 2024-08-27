package util

type idxRange struct {
	Low, High int
}

func Partition(collectionLen int, PartitionSize int) chan idxRange {
	c := make(chan idxRange)
	if PartitionSize <= 0 {
		close(c)
		return c
	}

	go func() {
		numFullPartitions := collectionLen / PartitionSize
		var i int
		for ; i < numFullPartitions; i++ {
			c <- idxRange{Low: i * PartitionSize, High: (i + 1) * PartitionSize}
		}

		if collectionLen%PartitionSize != 0 { // left over
			c <- idxRange{Low: i * PartitionSize, High: collectionLen}
		}

		close(c)
	}()
	return c
}
