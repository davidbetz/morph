package main

type idxRange struct {
	Low, High int
}

func partition(collectionLen int, partitionSize int) chan idxRange {
	c := make(chan idxRange)
	if partitionSize <= 0 {
		close(c)
		return c
	}

	go func() {
		numFullPartitions := collectionLen / partitionSize
		var i int
		for ; i < numFullPartitions; i++ {
			c <- idxRange{Low: i * partitionSize, High: (i + 1) * partitionSize}
		}

		if collectionLen%partitionSize != 0 { // left over
			c <- idxRange{Low: i * partitionSize, High: collectionLen}
		}

		close(c)
	}()
	return c
}
