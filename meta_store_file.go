package fsdup

type fileMetaStore struct {
	// Nothing
}

func NewFileMetaStore() *fileMetaStore {
	debugf("Creating file metadata store\n")
	return &fileMetaStore{}
}

func (s *fileMetaStore) ReadManifest(manifestId string) (*manifest, error) {
	return NewManifestFromFile(manifestId)
}

func (s* fileMetaStore) WriteManifest(manifestId string, manifest *manifest) error {
	return manifest.WriteToFile(manifestId)
}