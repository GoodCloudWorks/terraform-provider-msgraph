package id

import "testing"

func TestID(t *testing.T) {
	t.Run("when new then path is set", func(t *testing.T) {
		id := New("collection", "object-id")

		if id.Path != "collection/object-id" {
			t.Errorf("expected path to be collection/object-id, got %s", id.Path)
		}
	})

	t.Run("when parse then collection and object id are set", func(t *testing.T) {
		id, err := Parse("collection/object-id")

		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}

		if id.Path != "collection/object-id" {
			t.Errorf("expected path to be collection/object-id, got %s", id.Path)
		}

		if id.Collection() != "collection" {
			t.Errorf("expected collection to be collection, got %s", id.Collection())
		}

		if id.ObjectId() != "object-id" {
			t.Errorf("expected object id to be object-id, got %s", id.ObjectId())
		}
	})

	t.Run("when parse with /v1.0/ prefix then api version is set", func(t *testing.T) {
		id, err := Parse("/v1.0/collection/object-id")

		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}

		if id.ApiVersion() != "v1.0" {
			t.Errorf("expected api version to be v1.0, got %s", id.ApiVersion())
		}
	})
}
