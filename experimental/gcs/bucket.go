package gcs

import (
	"fmt"
	"k8s.io/contrib/test-utils/utils"
)

func main() {
	fmt.Println("Hello world!")

	bucket := utils.NewBucket(utils.KubekinsBucket)

	items, err := bucket.List("logs/ci-kubernetes-kubemark-500-gce/")

	fmt.Println(items, err)
}
