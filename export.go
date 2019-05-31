package main

import (
	"errors"
	"os"
)

func export(manifestFile, outputFile string) error {
	manifest, err := NewManifestFromFile(manifestFile)
	if err != nil {
		return err
	}

	if err := truncateExportFile(outputFile, manifest.Size()); err != nil {
		return err
	}

	// Open chunk index
	store := NewFileStore("index")

	// Open file
	out, err := os.OpenFile(outputFile, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	buffer := make([]byte, chunkSizeMaxBytes)
	offset := int64(0)

	for _, breakpoint := range manifest.Breakpoints() {
		part := manifest.Get(breakpoint)
		sparse := part.checksum == nil
		length := part.to - part.from

		if sparse {
			Debugf("%013d Skipping sparse section of %d bytes\n", offset, length)
		} else {
			Debugf("%013d Writing chunk %x, offset %d - %d (size %d)\n", offset, part.checksum, part.from, part.to, length)

			read, err := store.ReadAt(part.checksum, buffer[:length], part.from)
			if err != nil {
				return err
			} else if int64(read) != length {
				return errors.New("cannot read all required bytes from chunk")
			}

			written, err := out.WriteAt(buffer[:length], offset)
			if err != nil {
				return err
			} else if int64(written) != length {
				return errors.New("cannot write all bytes to output file")
			}
		}

		offset += length
	}

	err = out.Close()
	if err != nil {
		return err
	}

	return nil
}

// truncateExportFile wipes the output file (truncate to zero, then to target size)
func truncateExportFile(outputFile string, size int64) error {
	if _, err := os.Stat(outputFile); err != nil {
		file, err := os.Create(outputFile)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	} else {
		if err := os.Truncate(outputFile, 0); err != nil {
			return err
		}
	}

	if err := os.Truncate(outputFile, size); err != nil {
		return err
	}

	return nil
}