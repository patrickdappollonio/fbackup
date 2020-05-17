package main

import "github.com/gdrive-org/gdrive/drive"

type dateComparer struct{}

func (nc *dateComparer) Changed(lf *drive.LocalFile, rf *drive.RemoteFile) bool {
	return lf.Modified().After(rf.Modified())
}
