package s3deploy

import (
	"fmt"
	"testing"
)

func TestResolveDeleteObjects(t *testing.T) {
	config := &Config{
		Prefix:           "",
		VersionsFileName: "version.json",
	}

	versions := &Versions{
		Current: &Version{
			Keys: []string{"3.jpg", "4.jpg"},
		},
		Previous: &Version{
			Keys: []string{"2.jpg", "3.jpg"},
		},
		NextDelete: &Version{
			Keys: []string{"1.jpg", "2.jpg"},
		},
	}

	deleteObjects := config.ResolveDeleteObjects(versions)

	if len(deleteObjects) != 1 {
		t.Fatalf(
			fmt.Sprintf("Got delete objects does not contain %d files. got %d",
				1,
				len(deleteObjects)),
		)
	}

	if *deleteObjects[0].Key != "1.jpg" {
		t.Fatalf("Not expected delete object key. expected=%q, got=%q", "1.jpg", *deleteObjects[0].Key)
	}
}
