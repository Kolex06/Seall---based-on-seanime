package db_bridge

import (
	"seall/internal/database/db"
	"seall/internal/database/models"
	hibiketorrent "seall/internal/extension/hibike/torrent"

	"github.com/goccy/go-json"
)

func GetTorrentstreamHistory(db *db.Database, mId int) (*hibiketorrent.MediaTorrent, *hibiketorrent.BatchEpisodeFiles, error) {
	var history models.TorrentstreamHistory
	if err := db.Gorm().Where("media_id = ?", mId).First(&history).Error; err != nil {
		return nil, nil, err
	}

	var torrent hibiketorrent.MediaTorrent
	if err := json.Unmarshal(history.Torrent, &torrent); err != nil {
		return nil, nil, err
	}

	var files *hibiketorrent.BatchEpisodeFiles
	if len(history.BatchEpisodeFiles) > 0 {
		_ = json.Unmarshal(history.BatchEpisodeFiles, &files)
	}

	return &torrent, files, nil
}

func InsertTorrentstreamHistory(db *db.Database, mId int, torrent *hibiketorrent.MediaTorrent, files *hibiketorrent.BatchEpisodeFiles) error {
	if torrent == nil {
		return nil
	}

	// Marshal the data
	bytes, err := json.Marshal(torrent)
	if err != nil {
		return err
	}

	var filesBytes []byte
	if files != nil {
		filesBytes, err = json.Marshal(files)
		if err != nil {
			return err
		}
	}

	// Get current history
	var history models.TorrentstreamHistory
	if err := db.Gorm().Where("media_id = ?", mId).First(&history).Error; err == nil {
		// Update the history
		history.Torrent = bytes
		history.BatchEpisodeFiles = filesBytes
		return db.Gorm().Save(&history).Error
	}

	return db.Gorm().Create(&models.TorrentstreamHistory{
		MediaId:           mId,
		Torrent:           bytes,
		BatchEpisodeFiles: filesBytes,
	}).Error
}

func DeleteTorrentstreamHistory(db *db.Database, mId int) error {
	if mId == 0 {
		return nil
	}

	return db.Gorm().Where("media_id = ?", mId).Delete(&models.TorrentstreamHistory{}).Error
}
