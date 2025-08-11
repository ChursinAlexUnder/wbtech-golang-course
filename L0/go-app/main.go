package main

import (
	"context"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/go-app/internal"
)

func main() {
	// Если topic-а с необходимым названием для producer не существует,
	// то он создастся по этим параметрам
	var (
		topic             string = "orders"
		partitions        int    = 3
		replicationFactor int    = 1
	)
	ctx := context.Background()
	go internal.Producer(ctx, topic, partitions, replicationFactor)
	internal.Consumer(ctx)
}
