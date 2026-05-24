package simkl_platform

import (
	"seall/internal/api/mediaapi"
	"sort"
)

type simklStatsBucket struct {
	count          int
	minutesWatched int
	episodes       int
	scoreTotal     float64
	scoreCount     int
	mediaIds       []*int
}

func simklViewerStatsFromCollection(collection *mediaapi.AnimeCollection) *mediaapi.ViewerStats {
	total := &simklStatsBucket{}
	formatBuckets := map[mediaapi.MediaFormat]*simklStatsBucket{}
	statusBuckets := map[mediaapi.MediaListStatus]*simklStatsBucket{}
	scoreBuckets := map[int]*simklStatsBucket{}
	genreBuckets := map[string]*simklStatsBucket{}
	startYearBuckets := map[int]*simklStatsBucket{}
	releaseYearBuckets := map[int]*simklStatsBucket{}

	for _, list := range collection.GetMediaListCollection().GetLists() {
		for _, entry := range list.GetEntries() {
			media := entry.GetMedia()
			if media == nil || media.ID == 0 {
				continue
			}

			id := media.ID
			progress := intValue(entry.Progress)
			minutesWatched := progress * intValue(media.Duration)
			score := scoreValue(entry.Score)
			hasScore := entry.Score != nil && score > 0

			total.add(id, minutesWatched, progress, score, hasScore)

			if media.Format != nil {
				bucketFor(formatBuckets, *media.Format).add(id, minutesWatched, progress, score, hasScore)
			}
			if entry.Status != nil {
				bucketFor(statusBuckets, *entry.Status).add(id, minutesWatched, progress, score, hasScore)
			}
			if hasScore {
				bucketFor(scoreBuckets, int(score)).add(id, minutesWatched, progress, score, true)
			}
			for _, genre := range media.Genres {
				if genre != nil && *genre != "" {
					bucketFor(genreBuckets, *genre).add(id, minutesWatched, progress, score, hasScore)
				}
			}
			if startYear := entryYear(entry); startYear > 0 {
				bucketFor(startYearBuckets, startYear).add(id, minutesWatched, progress, score, hasScore)
			}
			if releaseYear := mediaYear(media); releaseYear > 0 {
				bucketFor(releaseYearBuckets, releaseYear).add(id, minutesWatched, progress, score, hasScore)
			}
		}
	}

	return &mediaapi.ViewerStats{
		Viewer: &mediaapi.ViewerStats_Viewer{
			Statistics: &mediaapi.ViewerStats_Viewer_Statistics{
				Anime: &mediaapi.ViewerStats_Viewer_Statistics_Anime{
					Count:           total.count,
					MinutesWatched:  total.minutesWatched,
					EpisodesWatched: total.episodes,
					MeanScore:       total.meanScore(),
					Formats:         formatStats(formatBuckets),
					Statuses:        statusStats(statusBuckets),
					Scores:          scoreStats(scoreBuckets),
					Genres:          genreStats(genreBuckets),
					StartYears:      startYearStats(startYearBuckets),
					ReleaseYears:    releaseYearStats(releaseYearBuckets),
					Studios:         []*mediaapi.UserStudioStats{},
				},
				Manga: &mediaapi.ViewerStats_Viewer_Statistics_Manga{},
			},
		},
	}
}

func (b *simklStatsBucket) add(mediaId int, minutesWatched int, episodes int, score float64, hasScore bool) {
	b.count++
	b.minutesWatched += minutesWatched
	b.episodes += episodes
	if hasScore {
		b.scoreTotal += score
		b.scoreCount++
	}
	id := mediaId
	b.mediaIds = append(b.mediaIds, &id)
}

func (b *simklStatsBucket) meanScore() float64 {
	if b == nil || b.scoreCount == 0 {
		return 0
	}
	return b.scoreTotal / float64(b.scoreCount)
}

func bucketFor[K comparable](buckets map[K]*simklStatsBucket, key K) *simklStatsBucket {
	bucket := buckets[key]
	if bucket == nil {
		bucket = &simklStatsBucket{}
		buckets[key] = bucket
	}
	return bucket
}

func formatStats(buckets map[mediaapi.MediaFormat]*simklStatsBucket) []*mediaapi.UserFormatStats {
	keys := sortedKeys(buckets)
	ret := make([]*mediaapi.UserFormatStats, 0, len(keys))
	for _, key := range keys {
		bucket := buckets[key]
		format := key
		ret = append(ret, &mediaapi.UserFormatStats{
			Format:         &format,
			Count:          bucket.count,
			MinutesWatched: bucket.minutesWatched,
			MeanScore:      bucket.meanScore(),
			MediaIds:       bucket.mediaIds,
		})
	}
	return ret
}

func statusStats(buckets map[mediaapi.MediaListStatus]*simklStatsBucket) []*mediaapi.UserStatusStats {
	keys := sortedKeys(buckets)
	ret := make([]*mediaapi.UserStatusStats, 0, len(keys))
	for _, key := range keys {
		bucket := buckets[key]
		status := key
		ret = append(ret, &mediaapi.UserStatusStats{
			Status:         &status,
			Count:          bucket.count,
			MinutesWatched: bucket.minutesWatched,
			MeanScore:      bucket.meanScore(),
			MediaIds:       bucket.mediaIds,
		})
	}
	return ret
}

func scoreStats(buckets map[int]*simklStatsBucket) []*mediaapi.UserScoreStats {
	keys := sortedKeys(buckets)
	ret := make([]*mediaapi.UserScoreStats, 0, len(keys))
	for _, key := range keys {
		bucket := buckets[key]
		score := key
		ret = append(ret, &mediaapi.UserScoreStats{
			Score:          &score,
			Count:          bucket.count,
			MinutesWatched: bucket.minutesWatched,
			MeanScore:      bucket.meanScore(),
			MediaIds:       bucket.mediaIds,
		})
	}
	return ret
}

func genreStats(buckets map[string]*simklStatsBucket) []*mediaapi.UserGenreStats {
	keys := sortedKeys(buckets)
	ret := make([]*mediaapi.UserGenreStats, 0, len(keys))
	for _, key := range keys {
		bucket := buckets[key]
		genre := key
		ret = append(ret, &mediaapi.UserGenreStats{
			Genre:          &genre,
			Count:          bucket.count,
			MinutesWatched: bucket.minutesWatched,
			MeanScore:      bucket.meanScore(),
			MediaIds:       bucket.mediaIds,
		})
	}
	return ret
}

func startYearStats(buckets map[int]*simklStatsBucket) []*mediaapi.UserStartYearStats {
	keys := sortedKeys(buckets)
	ret := make([]*mediaapi.UserStartYearStats, 0, len(keys))
	for _, key := range keys {
		bucket := buckets[key]
		year := key
		ret = append(ret, &mediaapi.UserStartYearStats{
			StartYear:      &year,
			Count:          bucket.count,
			MinutesWatched: bucket.minutesWatched,
			MeanScore:      bucket.meanScore(),
			MediaIds:       bucket.mediaIds,
		})
	}
	return ret
}

func releaseYearStats(buckets map[int]*simklStatsBucket) []*mediaapi.UserReleaseYearStats {
	keys := sortedKeys(buckets)
	ret := make([]*mediaapi.UserReleaseYearStats, 0, len(keys))
	for _, key := range keys {
		bucket := buckets[key]
		year := key
		ret = append(ret, &mediaapi.UserReleaseYearStats{
			ReleaseYear:    &year,
			Count:          bucket.count,
			MinutesWatched: bucket.minutesWatched,
			MeanScore:      bucket.meanScore(),
			MediaIds:       bucket.mediaIds,
		})
	}
	return ret
}

func sortedKeys[K ~int | ~string, V any](input map[K]V) []K {
	keys := make([]K, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i int, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func intValue(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func scoreValue(value *float64) float64 {
	if value == nil {
		return 0
	}
	if *value > 100 {
		return 100
	}
	return *value
}

func mediaYear(media *mediaapi.BaseAnime) int {
	if media == nil {
		return 0
	}
	if media.SeasonYear != nil && *media.SeasonYear > 0 {
		return *media.SeasonYear
	}
	return intValue(media.GetStartDate().Year)
}

func entryYear(entry *mediaapi.AnimeCollection_MediaListCollection_Lists_Entries) int {
	if entry == nil {
		return 0
	}
	if entry.CompletedAt != nil && entry.CompletedAt.Year != nil {
		return *entry.CompletedAt.Year
	}
	if entry.StartedAt != nil && entry.StartedAt.Year != nil {
		return *entry.StartedAt.Year
	}
	return 0
}
