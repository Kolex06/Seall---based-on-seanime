package testmocks

import "seall/internal/api/mediaapi"

type BaseAnimeBuilder struct {
	anime *mediaapi.BaseAnime
}

func NewBaseAnimeBuilder(id int, title string) *BaseAnimeBuilder {
	return &BaseAnimeBuilder{anime: &mediaapi.BaseAnime{
		ID:       id,
		IDMal:    new(501),
		Status:   new(mediaapi.MediaStatusFinished),
		Type:     new(mediaapi.MediaTypeAnime),
		Format:   new(mediaapi.MediaFormatTv),
		Episodes: new(12),
		IsAdult:  new(false),
		Title: &mediaapi.BaseAnime_Title{
			English: new(title),
			Romaji:  new(title),
		},
		Synonyms: []*string{new(title), new("Sample Anime Season 1")},
		StartDate: &mediaapi.BaseAnime_StartDate{
			Year:  new(2024),
			Month: new(1),
			Day:   new(2),
		},
	}}
}

func NewBaseAnime(id int, title string) *mediaapi.BaseAnime {
	return NewBaseAnimeBuilder(id, title).Build()
}

func (b *BaseAnimeBuilder) WithIDMal(idMal int) *BaseAnimeBuilder {
	b.anime.IDMal = new(idMal)
	return b
}

func (b *BaseAnimeBuilder) WithSiteURL(siteURL string) *BaseAnimeBuilder {
	b.anime.SiteURL = new(siteURL)
	return b
}

func (b *BaseAnimeBuilder) WithTitles(english string, romaji string, native string, userPreferred string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.English = new(english)
	b.anime.Title.Romaji = new(romaji)
	b.anime.Title.Native = new(native)
	b.anime.Title.UserPreferred = new(userPreferred)
	return b
}

func (b *BaseAnimeBuilder) WithEnglishTitle(title string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.English = new(title)
	return b
}

func (b *BaseAnimeBuilder) WithRomajiTitle(title string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.Romaji = new(title)
	return b
}

func (b *BaseAnimeBuilder) WithNativeTitle(title string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.Native = new(title)
	return b
}

func (b *BaseAnimeBuilder) WithUserPreferredTitle(title string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.UserPreferred = new(title)
	return b
}

func (b *BaseAnimeBuilder) WithStatus(status mediaapi.MediaStatus) *BaseAnimeBuilder {
	b.anime.Status = new(status)
	return b
}

func (b *BaseAnimeBuilder) WithFormat(format mediaapi.MediaFormat) *BaseAnimeBuilder {
	b.anime.Format = new(format)
	return b
}

func (b *BaseAnimeBuilder) WithEpisodes(episodes int) *BaseAnimeBuilder {
	b.anime.Episodes = new(episodes)
	return b
}

func (b *BaseAnimeBuilder) WithIsAdult(isAdult bool) *BaseAnimeBuilder {
	b.anime.IsAdult = new(isAdult)
	return b
}

func (b *BaseAnimeBuilder) WithSynonyms(synonyms ...string) *BaseAnimeBuilder {
	b.anime.Synonyms = stringPointers(synonyms...)
	return b
}

func (b *BaseAnimeBuilder) WithStartDate(year int, month int, day int) *BaseAnimeBuilder {
	b.anime.StartDate = &mediaapi.BaseAnime_StartDate{
		Year:  new(year),
		Month: new(month),
		Day:   new(day),
	}
	return b
}

func (b *BaseAnimeBuilder) WithEndDate(year int, month int, day int) *BaseAnimeBuilder {
	b.anime.EndDate = &mediaapi.BaseAnime_EndDate{
		Year:  new(year),
		Month: new(month),
		Day:   new(day),
	}
	return b
}

func (b *BaseAnimeBuilder) WithCoverImage(url string) *BaseAnimeBuilder {
	b.anime.CoverImage = &mediaapi.BaseAnime_CoverImage{
		ExtraLarge: new(url),
		Large:      new(url),
		Medium:     new(url),
	}
	return b
}

func (b *BaseAnimeBuilder) WithBannerImage(url string) *BaseAnimeBuilder {
	b.anime.BannerImage = new(url)
	return b
}

func (b *BaseAnimeBuilder) WithNextAiringEpisode(episode int, airingAt int, timeUntilAiring int) *BaseAnimeBuilder {
	b.anime.NextAiringEpisode = &mediaapi.BaseAnime_NextAiringEpisode{
		Episode:         episode,
		AiringAt:        airingAt,
		TimeUntilAiring: timeUntilAiring,
	}
	return b
}

func (b *BaseAnimeBuilder) Build() *mediaapi.BaseAnime {
	return b.anime
}

type BaseMangaBuilder struct {
	manga *mediaapi.BaseManga
}

func NewBaseMangaBuilder(id int, title string) *BaseMangaBuilder {
	return &BaseMangaBuilder{manga: &mediaapi.BaseManga{
		ID:      id,
		Status:  new(mediaapi.MediaStatusFinished),
		Type:    new(mediaapi.MediaTypeManga),
		Format:  new(mediaapi.MediaFormatManga),
		IsAdult: new(false),
		Title: &mediaapi.BaseManga_Title{
			English: new(title),
			Romaji:  new(title),
		},
		Synonyms: []*string{new(title), new(title + " Alternative")},
		StartDate: &mediaapi.BaseManga_StartDate{
			Year: new(2023),
		},
	}}
}

func NewBaseManga(id int, title string) *mediaapi.BaseManga {
	return NewBaseMangaBuilder(id, title).Build()
}

func (b *BaseMangaBuilder) WithIDMal(idMal int) *BaseMangaBuilder {
	b.manga.IDMal = new(idMal)
	return b
}

func (b *BaseMangaBuilder) WithSiteURL(siteURL string) *BaseMangaBuilder {
	b.manga.SiteURL = new(siteURL)
	return b
}

func (b *BaseMangaBuilder) WithTitles(english string, romaji string, native string, userPreferred string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.English = new(english)
	b.manga.Title.Romaji = new(romaji)
	b.manga.Title.Native = new(native)
	b.manga.Title.UserPreferred = new(userPreferred)
	return b
}

func (b *BaseMangaBuilder) WithEnglishTitle(title string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.English = new(title)
	return b
}

func (b *BaseMangaBuilder) WithRomajiTitle(title string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.Romaji = new(title)
	return b
}

func (b *BaseMangaBuilder) WithNativeTitle(title string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.Native = new(title)
	return b
}

func (b *BaseMangaBuilder) WithUserPreferredTitle(title string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.UserPreferred = new(title)
	return b
}

func (b *BaseMangaBuilder) WithStatus(status mediaapi.MediaStatus) *BaseMangaBuilder {
	b.manga.Status = new(status)
	return b
}

func (b *BaseMangaBuilder) WithFormat(format mediaapi.MediaFormat) *BaseMangaBuilder {
	b.manga.Format = new(format)
	return b
}

func (b *BaseMangaBuilder) WithChapters(chapters int) *BaseMangaBuilder {
	b.manga.Chapters = new(chapters)
	return b
}

func (b *BaseMangaBuilder) WithVolumes(volumes int) *BaseMangaBuilder {
	b.manga.Volumes = new(volumes)
	return b
}

func (b *BaseMangaBuilder) WithIsAdult(isAdult bool) *BaseMangaBuilder {
	b.manga.IsAdult = new(isAdult)
	return b
}

func (b *BaseMangaBuilder) WithSynonyms(synonyms ...string) *BaseMangaBuilder {
	b.manga.Synonyms = stringPointers(synonyms...)
	return b
}

func (b *BaseMangaBuilder) WithStartDate(year int, month int, day int) *BaseMangaBuilder {
	b.manga.StartDate = &mediaapi.BaseManga_StartDate{
		Year:  new(year),
		Month: new(month),
		Day:   new(day),
	}
	return b
}

func (b *BaseMangaBuilder) WithEndDate(year int, month int, day int) *BaseMangaBuilder {
	b.manga.EndDate = &mediaapi.BaseManga_EndDate{
		Year:  new(year),
		Month: new(month),
		Day:   new(day),
	}
	return b
}

func (b *BaseMangaBuilder) WithCoverImage(url string) *BaseMangaBuilder {
	b.manga.CoverImage = &mediaapi.BaseManga_CoverImage{
		ExtraLarge: new(url),
		Large:      new(url),
		Medium:     new(url),
	}
	return b
}

func (b *BaseMangaBuilder) WithBannerImage(url string) *BaseMangaBuilder {
	b.manga.BannerImage = new(url)
	return b
}

func (b *BaseMangaBuilder) Build() *mediaapi.BaseManga {
	return b.manga
}

func ensureAnimeTitle(anime *mediaapi.BaseAnime) {
	if anime.Title == nil {
		anime.Title = &mediaapi.BaseAnime_Title{}
	}
}

func ensureMangaTitle(manga *mediaapi.BaseManga) {
	if manga.Title == nil {
		manga.Title = &mediaapi.BaseManga_Title{}
	}
}

func stringPointers(values ...string) []*string {
	ret := make([]*string, 0, len(values))
	for _, value := range values {
		ret = append(ret, new(value))
	}
	return ret
}
