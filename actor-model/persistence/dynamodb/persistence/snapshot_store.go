package persistence

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type SnapshotStore struct {
	// TODO: Add necessary fields, such as database connection or configuration
}

func NewSnapshotStore(/* TODO: Add necessary parameters */) *SnapshotStore {
	return &SnapshotStore{
		// TODO: Initialize fields
	}
}

func (s *SnapshotStore) GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool) {
	// TODO: Implement getting the snapshot from the database
	// 1. Query the database to retrieve the snapshot for the given actorName
	// 2. If a snapshot is found, return it along with its eventIndex and ok = true
	// 3. If no snapshot is found, return nil, 0, and ok = false
	return nil, 0, false
}

func (s *SnapshotStore) PersistSnapshot(actorName string, snapshotIndex int, snapshot protoreflect.ProtoMessage) {
	// TODO: Implement persisting the snapshot to the database
	// 1. Convert the snapshot to a format suitable for storing in the database (e.g., JSON, binary)
	// 2. Store the converted snapshot in the database, associating it with the actorName and snapshotIndex
}

func (s *SnapshotStore) DeleteSnapshots(actorName string, inclusiveToIndex int) {
	// TODO: Implement deleting snapshots from the database
	// 1. Delete all snapshots for the given actorName up to and including the specified inclusiveToIndex
}

