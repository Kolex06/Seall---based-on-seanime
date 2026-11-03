package simkl

import (
	"github.com/goccy/go-json"
	"testing"
)

func TestSIMKLRelationsDecode(t *testing.T) {
	var media StandardMedia
	err := json.Unmarshal([]byte(`{"title":"Original","type":"anime","ids":{"simkl":36433},"relations":[{"title":"Sequel","en_title":"English Sequel","relation_type":"sequel","ids":{"simkl":123456}}]}`), &media)
	if err != nil {
		t.Fatal(err)
	}
	if len(media.Relations) != 1 || media.Relations[0].IDs.PrimarySimklID() != 123456 || media.Relations[0].RelationType != "sequel" {
		t.Fatalf("relations not decoded: %+v", media.Relations)
	}
}
