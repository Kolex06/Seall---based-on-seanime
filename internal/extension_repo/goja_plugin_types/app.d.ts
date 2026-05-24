declare namespace $app {

    /**
     * @package mediaapi
     */

    /**
     * @event ListMissedSequelsRequestedEvent
     * @file internal/api/simkl/hook_events.go
     * @description
     * ListMissedSequelsRequestedEvent is triggered when the list missed sequels request is requested.
     * Prevent default to skip the default behavior and return your own data.
     */
    function onListMissedSequelsRequested(cb: (event: ListMissedSequelsRequestedEvent) => void): void;

    interface ListMissedSequelsRequestedEvent {
        next(): void;

        preventDefault(): void;

        animeCollectionWithRelations?: MediaAPI_MediaCollectionWithRelations;
        variables?: Record<string, any>;
        query: string;
        list?: Array<MediaAPI_BaseMedia>;
    }

    /**
     * @event ListMissedSequelsEvent
     * @file internal/api/simkl/hook_events.go
     */
    function onListMissedSequels(cb: (event: ListMissedSequelsEvent) => void): void;

    interface ListMissedSequelsEvent {
        next(): void;

        list?: Array<MediaAPI_BaseMedia>;
    }


    /**
     * @package animap
     */

    /**
     * @event AnimapMediaRequestedEvent
     * @file internal/api/animap/hook_events.go
     * @description
     * AnimapMediaRequestedEvent is triggered when the Animap media is requested.
     * Prevent default to skip the default behavior and return your own data.
     */
    function onAnimapMediaRequested(cb: (event: AnimapMediaRequestedEvent) => void): void;

    interface AnimapMediaRequestedEvent {
        next(): void;

        preventDefault(): void;

        from: string;
        id: number;
        media?: Animap_Anime;
    }

    /**
     * @event AnimapMediaEvent
     * @file internal/api/animap/hook_events.go
     * @description
     * AnimapMediaEvent is triggered after processing AnimapMedia.
     */
    function onAnimapMedia(cb: (event: AnimapMediaEvent) => void): void;

    interface AnimapMediaEvent {
        next(): void;

        media?: Animap_Anime;
    }


    /**
     * @package anime
     */

    /**
     * @event AnimeEntryRequestedEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeEntryRequestedEvent is triggered when an anime entry is requested.
     * Prevent default to skip the default behavior and return the modified entry.
     * This event is triggered before [AnimeEntryEvent].
     * If the modified entry is nil, an error will be returned.
     */
    function onAnimeEntryRequested(cb: (event: AnimeEntryRequestedEvent) => void): void;

    interface AnimeEntryRequestedEvent {
        next(): void;

        preventDefault(): void;

        mediaId: number;
        localFiles?: Array<Media_LocalFile>;
        animeCollection?: MediaAPI_MediaCollection;
        entry?: Media_Entry;
    }

    /**
     * @event AnimeEntryEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeEntryEvent is triggered when the media entry is being returned.
     * This event is triggered after [AnimeEntryRequestedEvent].
     */
    function onAnimeEntry(cb: (event: AnimeEntryEvent) => void): void;

    interface AnimeEntryEvent {
        next(): void;

        entry?: Media_Entry;
    }

    /**
     * @event AnimeEntryFillerHydrationEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeEntryFillerHydrationEvent is triggered when the filler data is being added to the media entry.
     * This event is triggered after [AnimeEntryEvent].
     * Prevent default to skip the filler data.
     */
    function onAnimeEntryFillerHydration(cb: (event: AnimeEntryFillerHydrationEvent) => void): void;

    interface AnimeEntryFillerHydrationEvent {
        next(): void;

        preventDefault(): void;

        entry?: Media_Entry;
    }

    /**
     * @event AnimeEntryLibraryDataRequestedEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeEntryLibraryDataRequestedEvent is triggered when the app requests the library data for a media entry.
     * This is triggered before [AnimeEntryLibraryDataEvent].
     */
    function onAnimeEntryLibraryDataRequested(cb: (event: AnimeEntryLibraryDataRequestedEvent) => void): void;

    interface AnimeEntryLibraryDataRequestedEvent {
        next(): void;

        entryLocalFiles?: Array<Media_LocalFile>;
        mediaId: number;
        currentProgress: number;
    }

    /**
     * @event AnimeEntryLibraryDataEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeEntryLibraryDataEvent is triggered when the library data is being added to the media entry.
     * This is triggered after [AnimeEntryLibraryDataRequestedEvent].
     */
    function onAnimeEntryLibraryData(cb: (event: AnimeEntryLibraryDataEvent) => void): void;

    interface AnimeEntryLibraryDataEvent {
        next(): void;

        entryLibraryData?: Media_EntryLibraryData;
    }

    /**
     * @event MediaEntryManualMatchBeforeSaveEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * MediaEntryManualMatchBeforeSaveEvent is triggered when the user manually matches local files to a media entry.
     * Prevent default to skip saving the local files.
     */
    function onMediaEntryManualMatchBeforeSave(cb: (event: MediaEntryManualMatchBeforeSaveEvent) => void): void;

    interface MediaEntryManualMatchBeforeSaveEvent {
        next(): void;

        preventDefault(): void;

        mediaId: number;
        paths?: Array<string>;
        matchedLocalFiles?: Array<Media_LocalFile>;
    }

    /**
     * @event MissingEpisodesRequestedEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * MissingEpisodesRequestedEvent is triggered when the user requests the missing episodes for the entire library.
     * Prevent default to skip the default process and return the modified missing episodes.
     */
    function onMissingEpisodesRequested(cb: (event: MissingEpisodesRequestedEvent) => void): void;

    interface MissingEpisodesRequestedEvent {
        next(): void;

        preventDefault(): void;

        animeCollection?: MediaAPI_MediaCollection;
        localFiles?: Array<Media_LocalFile>;
        silencedMediaIds?: Array<number>;
        missingEpisodes?: Media_MissingEpisodes;
    }

    /**
     * @event MissingEpisodesEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * MissingEpisodesEvent is triggered when the missing episodes are being returned.
     */
    function onMissingEpisodes(cb: (event: MissingEpisodesEvent) => void): void;

    interface MissingEpisodesEvent {
        next(): void;

        missingEpisodes?: Media_MissingEpisodes;
    }

    /**
     * @event UpcomingEpisodesRequestedEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * UpcomingEpisodesRequestedEvent is triggered when the user requests upcoming episodes.
     * Prevent default to skip the default process and return the modified upcoming episodes.
     */
    function onUpcomingEpisodesRequested(cb: (event: UpcomingEpisodesRequestedEvent) => void): void;

    interface UpcomingEpisodesRequestedEvent {
        next(): void;

        preventDefault(): void;

        animeCollection?: MediaAPI_MediaCollection;
        localFiles?: Array<Media_LocalFile>;
        upcomingEpisodes?: Media_UpcomingEpisodes;
    }

    /**
     * @event UpcomingEpisodesEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * UpcomingEpisodesEvent is triggered when the upcoming episodes are being returned.
     */
    function onUpcomingEpisodes(cb: (event: UpcomingEpisodesEvent) => void): void;

    interface UpcomingEpisodesEvent {
        next(): void;

        upcomingEpisodes?: Media_UpcomingEpisodes;
    }

    /**
     * @event AnimeLibraryCollectionRequestedEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeLibraryCollectionRequestedEvent is triggered when the user requests the library collection.
     * Prevent default to skip the default process and return the modified library collection.
     * If the modified library collection is nil, an error will be returned.
     */
    function onAnimeLibraryCollectionRequested(cb: (event: AnimeLibraryCollectionRequestedEvent) => void): void;

    interface AnimeLibraryCollectionRequestedEvent {
        next(): void;

        preventDefault(): void;

        animeCollection?: MediaAPI_MediaCollection;
        localFiles?: Array<Media_LocalFile>;
        libraryCollection?: Media_LibraryCollection;
    }

    /**
     * @event AnimeLibraryCollectionEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeLibraryCollectionEvent is triggered when the user requests the library collection.
     */
    function onAnimeLibraryCollection(cb: (event: AnimeLibraryCollectionEvent) => void): void;

    interface AnimeLibraryCollectionEvent {
        next(): void;

        libraryCollection?: Media_LibraryCollection;
    }

    /**
     * @event AnimeLibraryStreamCollectionRequestedEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeLibraryStreamCollectionRequestedEvent is triggered when the user requests the library stream collection.
     * This is called when the user enables "Include in library" for either debrid/online/torrent streamings.
     */
    function onAnimeLibraryStreamCollectionRequested(cb: (event: AnimeLibraryStreamCollectionRequestedEvent) => void): void;

    interface AnimeLibraryStreamCollectionRequestedEvent {
        next(): void;

        animeCollection?: MediaAPI_MediaCollection;
        libraryCollection?: Media_LibraryCollection;
    }

    /**
     * @event AnimeLibraryStreamCollectionEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeLibraryStreamCollectionEvent is triggered when the library stream collection is being returned.
     */
    function onAnimeLibraryStreamCollection(cb: (event: AnimeLibraryStreamCollectionEvent) => void): void;

    interface AnimeLibraryStreamCollectionEvent {
        next(): void;

        streamCollection?: Media_StreamCollection;
    }

    /**
     * @event AnimeEntryDownloadInfoRequestedEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeEntryDownloadInfoRequestedEvent is triggered when the app requests the download info for a media entry.
     * This is triggered before [AnimeEntryDownloadInfoEvent].
     */
    function onAnimeEntryDownloadInfoRequested(cb: (event: AnimeEntryDownloadInfoRequestedEvent) => void): void;

    interface AnimeEntryDownloadInfoRequestedEvent {
        next(): void;

        localFiles?: Array<Media_LocalFile>;
        animeMetadata?: Metadata_AnimeMetadata;
        media?: MediaAPI_BaseMedia;
        progress?: number;
        status?: MediaAPI_MediaListStatus;
        entryDownloadInfo?: Media_EntryDownloadInfo;
    }

    /**
     * @event AnimeEntryDownloadInfoEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeEntryDownloadInfoEvent is triggered when the download info is being returned.
     */
    function onAnimeEntryDownloadInfo(cb: (event: AnimeEntryDownloadInfoEvent) => void): void;

    interface AnimeEntryDownloadInfoEvent {
        next(): void;

        entryDownloadInfo?: Media_EntryDownloadInfo;
    }

    /**
     * @event AnimeEpisodeCollectionRequestedEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeEpisodeCollectionRequestedEvent is triggered when the episode collection is being requested.
     * Prevent default to skip the default behavior and return your own data.
     */
    function onAnimeEpisodeCollectionRequested(cb: (event: AnimeEpisodeCollectionRequestedEvent) => void): void;

    interface AnimeEpisodeCollectionRequestedEvent {
        next(): void;

        preventDefault(): void;

        media?: MediaAPI_BaseMedia;
        metadata?: Metadata_AnimeMetadata;
        episodeCollection?: Media_EpisodeCollection;
    }

    /**
     * @event AnimeEpisodeCollectionEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * AnimeEpisodeCollectionEvent is triggered when the episode collection is being returned.
     */
    function onAnimeEpisodeCollection(cb: (event: AnimeEpisodeCollectionEvent) => void): void;

    interface AnimeEpisodeCollectionEvent {
        next(): void;

        episodeCollection?: Media_EpisodeCollection;
    }

    /**
     * @event MediaScheduleItemsEvent
     * @file internal/library/anime/hook_events.go
     * @description
     * MediaScheduleItemsEvent is triggered when the schedule items are being returned.
     */
    function onMediaScheduleItems(cb: (event: MediaScheduleItemsEvent) => void): void;

    interface MediaScheduleItemsEvent {
        next(): void;

        animeCollection?: MediaAPI_MediaCollection;
        items?: Array<Media_ScheduleItem>;
    }


    /**
     * @package anizip
     */

    /**
     * @event AnizipMediaRequestedEvent
     * @file internal/api/anizip/hook_events.go
     * @description
     * AnizipMediaRequestedEvent is triggered when the AniZip media is requested.
     * Prevent default to skip the default behavior and return your own data.
     */
    function onAnizipMediaRequested(cb: (event: AnizipMediaRequestedEvent) => void): void;

    interface AnizipMediaRequestedEvent {
        next(): void;

        preventDefault(): void;

        from: string;
        id: number;
        media?: Anizip_Media;
    }

    /**
     * @event AnizipMediaEvent
     * @file internal/api/anizip/hook_events.go
     * @description
     * AnizipMediaEvent is triggered after processing AnizipMedia.
     */
    function onAnizipMedia(cb: (event: AnizipMediaEvent) => void): void;

    interface AnizipMediaEvent {
        next(): void;

        media?: Anizip_Media;
    }


    /**
     * @package autodownloader
     */

    /**
     * @event AutoDownloaderRunStartedEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderRunStartedEvent is triggered when the autodownloader starts checking for new episodes.
     * Prevent default to abort the run.
     */
    function onAutoDownloaderRunStarted(cb: (event: AutoDownloaderRunStartedEvent) => void): void;

    interface AutoDownloaderRunStartedEvent {
        next(): void;

        preventDefault(): void;

        rules?: Array<Media_AutoDownloaderRule>;
        profiles?: Array<Media_AutoDownloaderProfile>;
        isSimulation: boolean;
    }

    /**
     * @event AutoDownloaderRunCompletedEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderRunCompletedEvent is triggered when the autodownloader finishes a run.
     */
    function onAutoDownloaderRunCompleted(cb: (event: AutoDownloaderRunCompletedEvent) => void): void;

    interface AutoDownloaderRunCompletedEvent {
        next(): void;

        rules?: Array<Media_AutoDownloaderRule>;
        profiles?: Array<Media_AutoDownloaderProfile>;
        isSimulation: boolean;
        downloadedCount: number;
        queuedCount: number;
        delayedCount: number;
    }

    /**
     * @event AutoDownloaderBeforeFetchTorrentsEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderBeforeFetchTorrentsEvent is triggered before the autodownloader fetches torrents from providers.
     * Prevent default to skip native provider retrieval.
     */
    function onAutoDownloaderBeforeFetchTorrents(cb: (event: AutoDownloaderBeforeFetchTorrentsEvent) => void): void;

    interface AutoDownloaderBeforeFetchTorrentsEvent {
        next(): void;

        preventDefault(): void;

        rules?: Array<Media_AutoDownloaderRule>;
        profiles?: Array<Media_AutoDownloaderProfile>;
        providerIds?: Array<string>;
        defaultProvider: string;
        torrents?: Array<AutoDownloader_NormalizedTorrent>;
    }

    /**
     * @event AutoDownloaderTorrentsFetchedEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderTorrentsFetchedEvent is triggered at the beginning of a run, when the autodownloader fetches torrents from the provider.
     */
    function onAutoDownloaderTorrentsFetched(cb: (event: AutoDownloaderTorrentsFetchedEvent) => void): void;

    interface AutoDownloaderTorrentsFetchedEvent {
        next(): void;

        torrents?: Array<AutoDownloader_NormalizedTorrent>;
    }

    /**
     * @event AutoDownloaderMatchVerifiedEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderMatchVerifiedEvent is triggered when a torrent is verified to follow a rule.
     * Changing MatchFound or Episode lets the hook override the verified result.
     * Prevent default to reject the match.
     */
    function onAutoDownloaderMatchVerified(cb: (event: AutoDownloaderMatchVerifiedEvent) => void): void;

    interface AutoDownloaderMatchVerifiedEvent {
        next(): void;

        preventDefault(): void;

        torrent?: AutoDownloader_NormalizedTorrent;
        rule?: Media_AutoDownloaderRule;
        listEntry?: MediaAPI_MediaListEntry;
        localEntry?: Media_LocalFileWrapperEntry;
        episode: number;
        matchFound: boolean;
    }

    /**
     * @event AutoDownloaderBestCandidateSelectedEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderBestCandidateSelectedEvent is triggered when the best candidate for an episode is selected.
     * Prevent default to skip handling the episode.
     */
    function onAutoDownloaderBestCandidateSelected(cb: (event: AutoDownloaderBestCandidateSelectedEvent) => void): void;

    interface AutoDownloaderBestCandidateSelectedEvent {
        next(): void;

        preventDefault(): void;

        rule?: Media_AutoDownloaderRule;
        episode: number;
        candidates?: Array<AutoDownloader_Candidate>;
        candidate?: AutoDownloader_Candidate;
        existingItem?: Models_AutoDownloaderItem;
        isSimulation: boolean;
    }

    /**
     * @event AutoDownloaderSettingsUpdatedEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderSettingsUpdatedEvent is triggered when the autodownloader settings are updated
     */
    function onAutoDownloaderSettingsUpdated(cb: (event: AutoDownloaderSettingsUpdatedEvent) => void): void;

    interface AutoDownloaderSettingsUpdatedEvent {
        next(): void;

        settings?: Models_AutoDownloaderSettings;
    }

    /**
     * @event AutoDownloaderBeforeQueueDelayedTorrentEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderBeforeQueueDelayedTorrentEvent is triggered when the autodownloader is about to queue a torrent with delay.
     * Prevent default to skip the delayed queue behavior.
     */
    function onAutoDownloaderBeforeQueueDelayedTorrent(cb: (event: AutoDownloaderBeforeQueueDelayedTorrentEvent) => void): void;

    interface AutoDownloaderBeforeQueueDelayedTorrentEvent {
        next(): void;

        preventDefault(): void;

        candidate?: AutoDownloader_Candidate;
        rule?: Media_AutoDownloaderRule;
        episode: number;
        delayMinutes: number;
        isSimulation: boolean;
    }

    /**
     * @event AutoDownloaderBeforeDownloadTorrentEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderBeforeDownloadTorrentEvent is triggered when the autodownloader is about to download a torrent.
     * Prevent default to abort the download.
     */
    function onAutoDownloaderBeforeDownloadTorrent(cb: (event: AutoDownloaderBeforeDownloadTorrentEvent) => void): void;

    interface AutoDownloaderBeforeDownloadTorrentEvent {
        next(): void;

        preventDefault(): void;

        torrent?: AutoDownloader_NormalizedTorrent;
        rule?: Media_AutoDownloaderRule;
        episode: number;
        score: number;
        items?: Array<Models_AutoDownloaderItem>;
        existingItem?: Models_AutoDownloaderItem;
        isSimulation: boolean;
    }

    /**
     * @event AutoDownloaderAfterDownloadTorrentEvent
     * @file internal/library/autodownloader/hook_events.go
     * @description
     * AutoDownloaderAfterDownloadTorrentEvent is triggered after the autodownloader queues or downloads a torrent.
     */
    function onAutoDownloaderAfterDownloadTorrent(cb: (event: AutoDownloaderAfterDownloadTorrentEvent) => void): void;

    interface AutoDownloaderAfterDownloadTorrentEvent {
        next(): void;

        torrent?: AutoDownloader_NormalizedTorrent;
        rule?: Media_AutoDownloaderRule;
        episode: number;
        score: number;
        downloaded: boolean;
        item?: Models_AutoDownloaderItem;
        isSimulation: boolean;
    }


    /**
     * @package continuity
     */

    /**
     * @event WatchHistoryItemRequestedEvent
     * @file internal/continuity/hook_events.go
     * @description
     * WatchHistoryItemRequestedEvent is triggered when a watch history item is requested.
     * Prevent default to skip getting the watch history item from the file cache, in this case the event should have a valid WatchHistoryItem object or set it to nil to indicate that the watch history item was not found.
     */
    function onWatchHistoryItemRequested(cb: (event: WatchHistoryItemRequestedEvent) => void): void;

    interface WatchHistoryItemRequestedEvent {
        next(): void;

        preventDefault(): void;

        mediaId: number;
        watchHistoryItem?: Continuity_WatchHistoryItem;
    }

    /**
     * @event WatchHistoryItemUpdatedEvent
     * @file internal/continuity/hook_events.go
     * @description
     * WatchHistoryItemUpdatedEvent is triggered when a watch history item is updated.
     */
    function onWatchHistoryItemUpdated(cb: (event: WatchHistoryItemUpdatedEvent) => void): void;

    interface WatchHistoryItemUpdatedEvent {
        next(): void;

        watchHistoryItem?: Continuity_WatchHistoryItem;
    }

    /**
     * @event WatchHistoryLocalFileEpisodeItemRequestedEvent
     * @file internal/continuity/hook_events.go
     */
    function onWatchHistoryLocalFileEpisodeItemRequested(cb: (event: WatchHistoryLocalFileEpisodeItemRequestedEvent) => void): void;

    interface WatchHistoryLocalFileEpisodeItemRequestedEvent {
        next(): void;

        Path: string;
        LocalFiles?: Array<Media_LocalFile>;
        watchHistoryItem?: Continuity_WatchHistoryItem;
    }

    /**
     * @event WatchHistoryStreamEpisodeItemRequestedEvent
     * @file internal/continuity/hook_events.go
     */
    function onWatchHistoryStreamEpisodeItemRequested(cb: (event: WatchHistoryStreamEpisodeItemRequestedEvent) => void): void;

    interface WatchHistoryStreamEpisodeItemRequestedEvent {
        next(): void;

        Episode: number;
        MediaId: number;
        watchHistoryItem?: Continuity_WatchHistoryItem;
    }


    /**
     * @package debrid_client
     */

    /**
     * @event DebridAutoSelectTorrentsFetchedEvent
     * @file internal/debrid/client/hook_events.go
     * @description
     * DebridAutoSelectTorrentsFetchedEvent is triggered when the torrents are fetched for auto select.
     * The torrents are sorted by seeders from highest to lowest.
     * This event is triggered before the top 3 torrents are analyzed.
     */
    function onDebridAutoSelectTorrentsFetched(cb: (event: DebridAutoSelectTorrentsFetchedEvent) => void): void;

    interface DebridAutoSelectTorrentsFetchedEvent {
        next(): void;

        Torrents?: Array<HibikeTorrent_MediaTorrent>;
    }

    /**
     * @event DebridSkipStreamCheckEvent
     * @file internal/debrid/client/hook_events.go
     * @description
     * DebridSkipStreamCheckEvent is triggered when the debrid client is about to skip the stream check.
     * Prevent default to enable the stream check.
     */
    function onDebridSkipStreamCheck(cb: (event: DebridSkipStreamCheckEvent) => void): void;

    interface DebridSkipStreamCheckEvent {
        next(): void;

        preventDefault(): void;

        streamURL: string;
        retries: number;
    /**
     * in seconds
     */
        retryDelay: number;
    }

    /**
     * @event DebridSendStreamToMediaPlayerEvent
     * @file internal/debrid/client/hook_events.go
     * @description
     * DebridSendStreamToMediaPlayerEvent is triggered when the debrid client is about to send a stream to the media player.
     * Prevent default to skip the playback.
     */
    function onDebridSendStreamToMediaPlayer(cb: (event: DebridSendStreamToMediaPlayerEvent) => void): void;

    interface DebridSendStreamToMediaPlayerEvent {
        next(): void;

        preventDefault(): void;

        windowTitle: string;
        streamURL: string;
        media?: MediaAPI_BaseMedia;
        aniDbEpisode: string;
        playbackType: string;
    }

    /**
     * @event DebridAddTorrentRequestedEvent
     * @file internal/debrid/client/hook_events.go
     * @description
     * DebridAddTorrentRequestedEvent is triggered when Seall is about to add a torrent to the debrid provider.
     * Prevent default to bypass the native add call and provide TorrentItemID yourself.
     */
    function onDebridAddTorrentRequested(cb: (event: DebridAddTorrentRequestedEvent) => void): void;

    interface DebridAddTorrentRequestedEvent {
        next(): void;

        preventDefault(): void;

        options?: Debrid_AddTorrentOptions;
        destination: string;
        mediaId: number;
        torrentItemId: string;
    }

    /**
     * @event DebridAddTorrentEvent
     * @file internal/debrid/client/hook_events.go
     * @description
     * DebridAddTorrentEvent is triggered after Seall adds a torrent to the debrid provider and queues it locally.
     */
    function onDebridAddTorrent(cb: (event: DebridAddTorrentEvent) => void): void;

    interface DebridAddTorrentEvent {
        next(): void;

        options?: Debrid_AddTorrentOptions;
        destination: string;
        mediaId: number;
        torrentItemId: string;
    }

    /**
     * @event DebridLocalDownloadRequestedEvent
     * @file internal/debrid/client/hook_events.go
     * @description
     * DebridLocalDownloadRequestedEvent is triggered when Seall is about to download a debrid torrent locally.
     * Prevent default to skip the default download and override the download.
     */
    function onDebridLocalDownloadRequested(cb: (event: DebridLocalDownloadRequestedEvent) => void): void;

    interface DebridLocalDownloadRequestedEvent {
        next(): void;

        preventDefault(): void;

        torrentName: string;
        destination: string;
        downloadUrl: string;
    }

    /**
     * @event DebridLocalDownloadStartedEvent
     * @file internal/debrid/client/hook_events.go
     * @description
     * DebridLocalDownloadStartedEvent is triggered right after Seall accepts a local debrid download.
     */
    function onDebridLocalDownloadStarted(cb: (event: DebridLocalDownloadStartedEvent) => void): void;

    interface DebridLocalDownloadStartedEvent {
        next(): void;

        torrentItemId: string;
        torrentName: string;
        destination: string;
        downloadUrl: string;
    }

    /**
     * @event DebridLocalDownloadCompletedEvent
     * @file internal/debrid/client/hook_events.go
     * @description
     * DebridLocalDownloadCompletedEvent is triggered when Seall finishes a local debrid download.
     */
    function onDebridLocalDownloadCompleted(cb: (event: DebridLocalDownloadCompletedEvent) => void): void;

    interface DebridLocalDownloadCompletedEvent {
        next(): void;

        torrentItemId: string;
        torrentName: string;
        destination: string;
    }


    /**
     * @package discordrpc_presence
     */

    /**
     * @event DiscordPresenceAnimeActivityRequestedEvent
     * @file internal/discordrpc/presence/hook_events.go
     * @description
     * DiscordPresenceAnimeActivityRequestedEvent is triggered when anime activity is requested, after the [animeActivity] is processed, and right before the activity is sent to queue.
     * There is no guarantee as to when or if the activity will be successfully sent to discord.
     * Note that this event is triggered every 6 seconds or so, avoid heavy processing or perform it only when the activity is changed.
     * Prevent default to stop the activity from being sent to discord.
     */
    function onDiscordPresenceAnimeActivityRequested(cb: (event: DiscordPresenceAnimeActivityRequestedEvent) => void): void;

    interface DiscordPresenceAnimeActivityRequestedEvent {
        next(): void;

        preventDefault(): void;

        animeActivity?: DiscordRPC_AnimeActivity;
        name: string;
        details: string;
        detailsUrl: string;
        state: string;
        startTimestamp?: number;
        endTimestamp?: number;
        largeImage: string;
        largeText: string;
    /**
     * URL to large image, if any
     */
        largeUrl?: string;
        smallImage: string;
        smallText: string;
    /**
     * URL to small image, if any
     */
        smallUrl?: string;
        buttons?: Array<DiscordRPC_Button>;
        instance: boolean;
        type: number;
        statusDisplayType?: number;
    }

    /**
     * @event DiscordPresenceMangaActivityRequestedEvent
     * @file internal/discordrpc/presence/hook_events.go
     * @description
     * DiscordPresenceMangaActivityRequestedEvent is triggered when manga activity is requested, after the [mangaActivity] is processed, and right before the activity is sent to queue.
     * There is no guarantee as to when or if the activity will be successfully sent to discord.
     * Note that this event is triggered every 6 seconds or so, avoid heavy processing or perform it only when the activity is changed.
     * Prevent default to stop the activity from being sent to discord.
     */
    function onDiscordPresenceMangaActivityRequested(cb: (event: DiscordPresenceMangaActivityRequestedEvent) => void): void;

    interface DiscordPresenceMangaActivityRequestedEvent {
        next(): void;

        preventDefault(): void;

        mangaActivity?: DiscordRPC_MangaActivity;
        name: string;
        details: string;
        detailsUrl: string;
        state: string;
        startTimestamp?: number;
        endTimestamp?: number;
        largeImage: string;
        largeText: string;
    /**
     * URL to large image, if any
     */
        largeUrl?: string;
        smallImage: string;
        smallText: string;
    /**
     * URL to small image, if any
     */
        smallUrl?: string;
        buttons?: Array<DiscordRPC_Button>;
        instance: boolean;
        type: number;
        statusDisplayType?: number;
    }

    /**
     * @event DiscordPresenceClientClosedEvent
     * @file internal/discordrpc/presence/hook_events.go
     * @description
     * DiscordPresenceClientClosedEvent is triggered when the discord rpc client is closed.
     */
    function onDiscordPresenceClientClosed(cb: (event: DiscordPresenceClientClosedEvent) => void): void;

    interface DiscordPresenceClientClosedEvent {
        next(): void;

    }


    /**
     * @package fillermanager
     */

    /**
     * @event HydrateFillerDataRequestedEvent
     * @file internal/library/fillermanager/hook_events.go
     * @description
     * HydrateFillerDataRequestedEvent is triggered when the filler manager requests to hydrate the filler data for an entry.
     * This is used by the local file episode list.
     * Prevent default to skip the default behavior and return your own data.
     */
    function onHydrateFillerDataRequested(cb: (event: HydrateFillerDataRequestedEvent) => void): void;

    interface HydrateFillerDataRequestedEvent {
        next(): void;

        preventDefault(): void;

        entry?: Media_Entry;
    }

    /**
     * @event HydrateOnlinestreamFillerDataRequestedEvent
     * @file internal/library/fillermanager/hook_events.go
     * @description
     * HydrateOnlinestreamFillerDataRequestedEvent is triggered when the filler manager requests to hydrate the filler data for online streaming episodes.
     * This is used by the online streaming episode list.
     * Prevent default to skip the default behavior and return your own data.
     */
    function onHydrateOnlinestreamFillerDataRequested(cb: (event: HydrateOnlinestreamFillerDataRequestedEvent) => void): void;

    interface HydrateOnlinestreamFillerDataRequestedEvent {
        next(): void;

        preventDefault(): void;

        episodes?: Array<Onlinestream_Episode>;
    }

    /**
     * @event HydrateEpisodeFillerDataRequestedEvent
     * @file internal/library/fillermanager/hook_events.go
     * @description
     * HydrateEpisodeFillerDataRequestedEvent is triggered when the filler manager requests to hydrate the filler data for specific episodes.
     * This is used by the torrent and debrid streaming episode list.
     * Prevent default to skip the default behavior and return your own data.
     */
    function onHydrateEpisodeFillerDataRequested(cb: (event: HydrateEpisodeFillerDataRequestedEvent) => void): void;

    interface HydrateEpisodeFillerDataRequestedEvent {
        next(): void;

        preventDefault(): void;

        episodes?: Array<Media_Episode>;
    }


    /**
     * @package manga
     */

    /**
     * @event MangaEntryRequestedEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaEntryRequestedEvent is triggered when a manga entry is requested.
     * Prevent default to skip the default behavior and return the modified entry.
     * If the modified entry is nil, an error will be returned.
     */
    function onMangaEntryRequested(cb: (event: MangaEntryRequestedEvent) => void): void;

    interface MangaEntryRequestedEvent {
        next(): void;

        preventDefault(): void;

        mediaId: number;
        mangaCollection?: MediaAPI_MangaCollection;
        entry?: Manga_Entry;
    }

    /**
     * @event MangaEntryEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaEntryEvent is triggered when the manga entry is being returned.
     */
    function onMangaEntry(cb: (event: MangaEntryEvent) => void): void;

    interface MangaEntryEvent {
        next(): void;

        entry?: Manga_Entry;
    }

    /**
     * @event MangaLibraryCollectionRequestedEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaLibraryCollectionRequestedEvent is triggered when the manga library collection is being requested.
     */
    function onMangaLibraryCollectionRequested(cb: (event: MangaLibraryCollectionRequestedEvent) => void): void;

    interface MangaLibraryCollectionRequestedEvent {
        next(): void;

        mangaCollection?: MediaAPI_MangaCollection;
    }

    /**
     * @event MangaLibraryCollectionEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaLibraryCollectionEvent is triggered when the manga library collection is being returned.
     */
    function onMangaLibraryCollection(cb: (event: MangaLibraryCollectionEvent) => void): void;

    interface MangaLibraryCollectionEvent {
        next(): void;

        libraryCollection?: Manga_Collection;
    }

    /**
     * @event MangaDownloadedChapterContainersRequestedEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaDownloadedChapterContainersRequestedEvent is triggered when the manga downloaded chapter containers are being requested.
     * Prevent default to skip the default behavior and return the modified chapter containers.
     * If the modified chapter containers are nil, an error will be returned.
     */
    function onMangaDownloadedChapterContainersRequested(cb: (event: MangaDownloadedChapterContainersRequestedEvent) => void): void;

    interface MangaDownloadedChapterContainersRequestedEvent {
        next(): void;

        preventDefault(): void;

        mangaCollection?: MediaAPI_MangaCollection;
        chapterContainers?: Array<Manga_ChapterContainer>;
    }

    /**
     * @event MangaDownloadedChapterContainersEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaDownloadedChapterContainersEvent is triggered when the manga downloaded chapter containers are being returned.
     */
    function onMangaDownloadedChapterContainers(cb: (event: MangaDownloadedChapterContainersEvent) => void): void;

    interface MangaDownloadedChapterContainersEvent {
        next(): void;

        chapterContainers?: Array<Manga_ChapterContainer>;
    }

    /**
     * @event MangaLatestChapterNumbersMapEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaLatestChapterNumbersMapEvent is triggered when the manga latest chapter numbers map is being returned.
     */
    function onMangaLatestChapterNumbersMap(cb: (event: MangaLatestChapterNumbersMapEvent) => void): void;

    interface MangaLatestChapterNumbersMapEvent {
        next(): void;

        latestChapterNumbersMap?: Record<number, Array<Manga_MangaLatestChapterNumberItem>>;
    }

    /**
     * @event MangaDownloadMapEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaDownloadMapEvent is triggered when the manga download map has been updated.
     * This map is used to tell the client which chapters have been downloaded.
     */
    function onMangaDownloadMap(cb: (event: MangaDownloadMapEvent) => void): void;

    interface MangaDownloadMapEvent {
        next(): void;

        mediaMap?: Manga_MediaMap;
    }

    /**
     * @event MangaChapterContainerRequestedEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaChapterContainerRequestedEvent is triggered when the manga chapter container is being requested.
     * This event happens before the chapter container is fetched from the cache or provider.
     * Prevent default to skip the default behavior and return the modified chapter container.
     * If the modified chapter container is nil, an error will be returned.
     */
    function onMangaChapterContainerRequested(cb: (event: MangaChapterContainerRequestedEvent) => void): void;

    interface MangaChapterContainerRequestedEvent {
        next(): void;

        preventDefault(): void;

        provider: string;
        mediaId: number;
        titles?: Array<string>;
        year: number;
        chapterContainer?: Manga_ChapterContainer;
    }

    /**
     * @event MangaChapterContainerEvent
     * @file internal/manga/hook_events.go
     * @description
     * MangaChapterContainerEvent is triggered when the manga chapter container is being returned.
     * This event happens after the chapter container is fetched from the cache or provider.
     */
    function onMangaChapterContainer(cb: (event: MangaChapterContainerEvent) => void): void;

    interface MangaChapterContainerEvent {
        next(): void;

        chapterContainer?: Manga_ChapterContainer;
    }


    /**
     * @package mediaplayer
     */

    /**
     * @event MediaPlayerLocalFileTrackingRequestedEvent
     * @file internal/mediaplayers/mediaplayer/hook_events.go
     * @description
     * MediaPlayerLocalFileTrackingRequestedEvent is triggered when the playback manager wants to track the progress of a local file.
     * Prevent default to stop tracking.
     */
    function onMediaPlayerLocalFileTrackingRequested(cb: (event: MediaPlayerLocalFileTrackingRequestedEvent) => void): void;

    interface MediaPlayerLocalFileTrackingRequestedEvent {
        next(): void;

        preventDefault(): void;

        startRefreshDelay: number;
        refreshDelay: number;
        maxRetries: number;
    }

    /**
     * @event MediaPlayerStreamTrackingRequestedEvent
     * @file internal/mediaplayers/mediaplayer/hook_events.go
     * @description
     * MediaPlayerStreamTrackingRequestedEvent is triggered when the playback manager wants to track the progress of a stream.
     * Prevent default to stop tracking.
     */
    function onMediaPlayerStreamTrackingRequested(cb: (event: MediaPlayerStreamTrackingRequestedEvent) => void): void;

    interface MediaPlayerStreamTrackingRequestedEvent {
        next(): void;

        preventDefault(): void;

        startRefreshDelay: number;
        refreshDelay: number;
        maxRetries: number;
        maxRetriesAfterStart: number;
    }


    /**
     * @package metadata
     */

    /**
     * @event AnimeMetadataRequestedEvent
     * @file internal/api/metadata/hook_events.go
     * @description
     * AnimeMetadataRequestedEvent is triggered when anime metadata is requested and right before the metadata is processed.
     * This event is followed by [AnimeMetadataEvent] which is triggered when the metadata is available.
     * Prevent default to skip the default behavior and return the modified metadata.
     * If the modified metadata is nil, an error will be returned.
     */
    function onAnimeMetadataRequested(cb: (event: AnimeMetadataRequestedEvent) => void): void;

    interface AnimeMetadataRequestedEvent {
        next(): void;

        preventDefault(): void;

        mediaId: number;
        animeMetadata?: Metadata_AnimeMetadata;
    }

    /**
     * @event AnimeMetadataEvent
     * @file internal/api/metadata/hook_events.go
     * @description
     * AnimeMetadataEvent is triggered when anime metadata is available and is about to be returned.
     * Anime metadata can be requested in many places, ranging from displaying the anime entry to starting a torrent stream.
     * This event is triggered after [AnimeMetadataRequestedEvent].
     * If the modified metadata is nil, an error will be returned.
     */
    function onAnimeMetadata(cb: (event: AnimeMetadataEvent) => void): void;

    interface AnimeMetadataEvent {
        next(): void;

        mediaId: number;
        animeMetadata?: Metadata_AnimeMetadata;
    }

    /**
     * @event AnimeEpisodeMetadataRequestedEvent
     * @file internal/api/metadata/hook_events.go
     * @description
     * AnimeEpisodeMetadataRequestedEvent is triggered when anime episode metadata is requested.
     * Prevent default to skip the default behavior and return the overridden metadata.
     * This event is triggered before [AnimeEpisodeMetadataEvent].
     * If the modified episode metadata is nil, an empty EpisodeMetadata object will be returned.
     */
    function onAnimeEpisodeMetadataRequested(cb: (event: AnimeEpisodeMetadataRequestedEvent) => void): void;

    interface AnimeEpisodeMetadataRequestedEvent {
        next(): void;

        preventDefault(): void;

        animeEpisodeMetadata?: Metadata_EpisodeMetadata;
        episode: string;
        episodeNumber: number;
        mediaId: number;
    }

    /**
     * @event AnimeEpisodeMetadataEvent
     * @file internal/api/metadata/hook_events.go
     * @description
     * AnimeEpisodeMetadataEvent is triggered when anime episode metadata is available and is about to be returned.
     * In the current implementation, episode metadata is requested for display purposes. It is used to get a more complete metadata object since the original AnimeMetadata object is not complete.
     * This event is triggered after [AnimeEpisodeMetadataRequestedEvent].
     * If the modified episode metadata is nil, an empty EpisodeMetadata object will be returned.
     */
    function onAnimeEpisodeMetadata(cb: (event: AnimeEpisodeMetadataEvent) => void): void;

    interface AnimeEpisodeMetadataEvent {
        next(): void;

        animeEpisodeMetadata?: Metadata_EpisodeMetadata;
        episode: string;
        episodeNumber: number;
        mediaId: number;
    }


    /**
     * @package platform
     */

    /**
     * @event GetAnimeEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetAnime(cb: (event: GetAnimeEvent) => void): void;

    interface GetAnimeEvent {
        next(): void;

        anime?: MediaAPI_BaseMedia;
    }

    /**
     * @event GetAnimeDetailsEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetAnimeDetails(cb: (event: GetAnimeDetailsEvent) => void): void;

    interface GetAnimeDetailsEvent {
        next(): void;

        anime?: MediaAPI_MediaDetailsById_Media;
    }

    /**
     * @event GetMangaEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetManga(cb: (event: GetMangaEvent) => void): void;

    interface GetMangaEvent {
        next(): void;

        manga?: MediaAPI_BaseManga;
    }

    /**
     * @event GetMangaDetailsEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetMangaDetails(cb: (event: GetMangaDetailsEvent) => void): void;

    interface GetMangaDetailsEvent {
        next(): void;

        manga?: MediaAPI_MangaDetailsById_Media;
    }

    /**
     * @event GetCachedAnimeCollectionEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetCachedAnimeCollection(cb: (event: GetCachedAnimeCollectionEvent) => void): void;

    interface GetCachedAnimeCollectionEvent {
        next(): void;

        animeCollection?: MediaAPI_MediaCollection;
    }

    /**
     * @event GetCachedMangaCollectionEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetCachedMangaCollection(cb: (event: GetCachedMangaCollectionEvent) => void): void;

    interface GetCachedMangaCollectionEvent {
        next(): void;

        mangaCollection?: MediaAPI_MangaCollection;
    }

    /**
     * @event GetMediaCollectionEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetMediaCollection(cb: (event: GetMediaCollectionEvent) => void): void;

    interface GetMediaCollectionEvent {
        next(): void;

        animeCollection?: MediaAPI_MediaCollection;
    }

    /**
     * @event GetMangaCollectionEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetMangaCollection(cb: (event: GetMangaCollectionEvent) => void): void;

    interface GetMangaCollectionEvent {
        next(): void;

        mangaCollection?: MediaAPI_MangaCollection;
    }

    /**
     * @event GetCachedRawAnimeCollectionEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetCachedRawAnimeCollection(cb: (event: GetCachedRawAnimeCollectionEvent) => void): void;

    interface GetCachedRawAnimeCollectionEvent {
        next(): void;

        animeCollection?: MediaAPI_MediaCollection;
    }

    /**
     * @event GetCachedRawMangaCollectionEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetCachedRawMangaCollection(cb: (event: GetCachedRawMangaCollectionEvent) => void): void;

    interface GetCachedRawMangaCollectionEvent {
        next(): void;

        mangaCollection?: MediaAPI_MangaCollection;
    }

    /**
     * @event GetRawMediaCollectionEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetRawMediaCollection(cb: (event: GetRawMediaCollectionEvent) => void): void;

    interface GetRawMediaCollectionEvent {
        next(): void;

        animeCollection?: MediaAPI_MediaCollection;
    }

    /**
     * @event GetRawMangaCollectionEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetRawMangaCollection(cb: (event: GetRawMangaCollectionEvent) => void): void;

    interface GetRawMangaCollectionEvent {
        next(): void;

        mangaCollection?: MediaAPI_MangaCollection;
    }

    /**
     * @event GetStudioDetailsEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onGetStudioDetails(cb: (event: GetStudioDetailsEvent) => void): void;

    interface GetStudioDetailsEvent {
        next(): void;

        studio?: MediaAPI_StudioDetails;
    }

    /**
     * @event PreUpdateEntryEvent
     * @file internal/platforms/platform/hook_events.go
     * @description
     * PreUpdateEntryEvent is triggered when an entry is about to be updated.
     * Prevent default to skip the default update and override the update.
     */
    function onPreUpdateEntry(cb: (event: PreUpdateEntryEvent) => void): void;

    interface PreUpdateEntryEvent {
        next(): void;

        preventDefault(): void;

        mediaId?: number;
        status?: MediaAPI_MediaListStatus;
        scoreRaw?: number;
        progress?: number;
        startedAt?: MediaAPI_FuzzyDateInput;
        completedAt?: MediaAPI_FuzzyDateInput;
    }

    /**
     * @event PostUpdateEntryEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onPostUpdateEntry(cb: (event: PostUpdateEntryEvent) => void): void;

    interface PostUpdateEntryEvent {
        next(): void;

        mediaId?: number;
    }

    /**
     * @event PreUpdateEntryProgressEvent
     * @file internal/platforms/platform/hook_events.go
     * @description
     * PreUpdateEntryProgressEvent is triggered when an entry's progress is about to be updated.
     * Prevent default to skip the default update and override the update.
     */
    function onPreUpdateEntryProgress(cb: (event: PreUpdateEntryProgressEvent) => void): void;

    interface PreUpdateEntryProgressEvent {
        next(): void;

        preventDefault(): void;

        mediaId?: number;
        progress?: number;
        totalCount?: number;
        status?: MediaAPI_MediaListStatus;
    }

    /**
     * @event PostUpdateEntryProgressEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onPostUpdateEntryProgress(cb: (event: PostUpdateEntryProgressEvent) => void): void;

    interface PostUpdateEntryProgressEvent {
        next(): void;

        mediaId?: number;
    }

    /**
     * @event PreUpdateEntryRepeatEvent
     * @file internal/platforms/platform/hook_events.go
     * @description
     * PreUpdateEntryRepeatEvent is triggered when an entry's repeat is about to be updated.
     * Prevent default to skip the default update and override the update.
     */
    function onPreUpdateEntryRepeat(cb: (event: PreUpdateEntryRepeatEvent) => void): void;

    interface PreUpdateEntryRepeatEvent {
        next(): void;

        preventDefault(): void;

        mediaId?: number;
        repeat?: number;
    }

    /**
     * @event PostUpdateEntryRepeatEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onPostUpdateEntryRepeat(cb: (event: PostUpdateEntryRepeatEvent) => void): void;

    interface PostUpdateEntryRepeatEvent {
        next(): void;

        mediaId?: number;
    }

    /**
     * @event PreDeleteEntryEvent
     * @file internal/platforms/platform/hook_events.go
     * @description
     * PreDeleteEntryEvent is triggered when an entry is about to be deleted.
     * Prevent default to skip the default deletion and override the deletion.
     */
    function onPreDeleteEntry(cb: (event: PreDeleteEntryEvent) => void): void;

    interface PreDeleteEntryEvent {
        next(): void;

        preventDefault(): void;

        mediaId?: number;
        entryId?: number;
    }

    /**
     * @event PostDeleteEntryEvent
     * @file internal/platforms/platform/hook_events.go
     */
    function onPostDeleteEntry(cb: (event: PostDeleteEntryEvent) => void): void;

    interface PostDeleteEntryEvent {
        next(): void;

        mediaId?: number;
        entryId?: number;
    }


    /**
     * @package playbackmanager
     */

    /**
     * @event LocalFilePlaybackRequestedEvent
     * @file internal/library/playbackmanager/hook_events.go
     * @description
     * LocalFilePlaybackRequestedEvent is triggered when a local file is requested to be played.
     * Prevent default to skip the default playback and override the playback.
     */
    function onLocalFilePlaybackRequested(cb: (event: LocalFilePlaybackRequestedEvent) => void): void;

    interface LocalFilePlaybackRequestedEvent {
        next(): void;

        preventDefault(): void;

        path: string;
    }

    /**
     * @event StreamPlaybackRequestedEvent
     * @file internal/library/playbackmanager/hook_events.go
     * @description
     * StreamPlaybackRequestedEvent is triggered when a stream is requested to be played.
     * Prevent default to skip the default playback and override the playback.
     */
    function onStreamPlaybackRequested(cb: (event: StreamPlaybackRequestedEvent) => void): void;

    interface StreamPlaybackRequestedEvent {
        next(): void;

        preventDefault(): void;

        windowTitle: string;
        payload: string;
        media?: MediaAPI_BaseMedia;
        aniDbEpisode: string;
    }

    /**
     * @event PlaybackBeforeTrackingEvent
     * @file internal/library/playbackmanager/hook_events.go
     * @description
     * PlaybackBeforeTrackingEvent is triggered just before the playback tracking starts.
     * Prevent default to skip playback tracking.
     */
    function onPlaybackBeforeTracking(cb: (event: PlaybackBeforeTrackingEvent) => void): void;

    interface PlaybackBeforeTrackingEvent {
        next(): void;

        preventDefault(): void;

        isStream: boolean;
    }

    /**
     * @event PlaybackLocalFileDetailsRequestedEvent
     * @file internal/library/playbackmanager/hook_events.go
     * @description
     * PlaybackLocalFileDetailsRequestedEvent is triggered when the local files details for a specific path are requested.
     * This event is triggered right after the media player loads an episode.
     * The playback manager uses the local files details to track the progress, propose next episodes, etc.
     * In the current implementation, the details are fetched by selecting the local file from the database and making requests to retrieve the media and anime list entry.
     * Prevent default to skip the default fetching and override the details.
     */
    function onPlaybackLocalFileDetailsRequested(cb: (event: PlaybackLocalFileDetailsRequestedEvent) => void): void;

    interface PlaybackLocalFileDetailsRequestedEvent {
        next(): void;

        preventDefault(): void;

        path: string;
        localFiles?: Array<Media_LocalFile>;
        animeListEntry?: MediaAPI_MediaListEntry;
        localFile?: Media_LocalFile;
        localFileWrapperEntry?: Media_LocalFileWrapperEntry;
    }

    /**
     * @event PlaybackStreamDetailsRequestedEvent
     * @file internal/library/playbackmanager/hook_events.go
     * @description
     * PlaybackStreamDetailsRequestedEvent is triggered when the stream details are requested.
     * Prevent default to skip the default fetching and override the details.
     * In the current implementation, the details are fetched by selecting the anime from the anime collection. If nothing is found, the stream is still tracked.
     */
    function onPlaybackStreamDetailsRequested(cb: (event: PlaybackStreamDetailsRequestedEvent) => void): void;

    interface PlaybackStreamDetailsRequestedEvent {
        next(): void;

        preventDefault(): void;

        animeCollection?: MediaAPI_MediaCollection;
        mediaId: number;
        animeListEntry?: MediaAPI_MediaListEntry;
    }


    /**
     * @package scanner
     */

    /**
     * @event ScanStartedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanStartedEvent is triggered when the scanning process begins.
     * Prevent default to skip the rest of the scanning process and return the local files.
     */
    function onScanStarted(cb: (event: ScanStartedEvent) => void): void;

    interface ScanStartedEvent {
        next(): void;

        preventDefault(): void;

        libraryPath: string;
        otherLibraryPaths?: Array<string>;
        enhanced: boolean;
        skipLocked: boolean;
        skipIgnored: boolean;
        localFiles?: Array<Media_LocalFile>;
    }

    /**
     * @event ScanFilePathsRetrievedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanFilePathsRetrievedEvent is triggered when the file paths to scan are retrieved.
     * The event includes file paths from all directories to scan.
     * The event includes file paths of local files that will be skipped.
     */
    function onScanFilePathsRetrieved(cb: (event: ScanFilePathsRetrievedEvent) => void): void;

    interface ScanFilePathsRetrievedEvent {
        next(): void;

        filePaths?: Array<string>;
    }

    /**
     * @event ScanLocalFilesParsedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanLocalFilesParsedEvent is triggered right after the file paths are parsed into local file objects.
     * The event does not include local files that are skipped.
     */
    function onScanLocalFilesParsed(cb: (event: ScanLocalFilesParsedEvent) => void): void;

    interface ScanLocalFilesParsedEvent {
        next(): void;

        localFiles?: Array<Media_LocalFile>;
    }

    /**
     * @event ScanCompletedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanCompletedEvent is triggered when the scanning process finishes.
     * The event includes all the local files (skipped and scanned) to be inserted as a new entry.
     * Right after this event, the local files will be inserted as a new entry.
     */
    function onScanCompleted(cb: (event: ScanCompletedEvent) => void): void;

    interface ScanCompletedEvent {
        next(): void;

        localFiles?: Array<Media_LocalFile>;
    /**
     * in milliseconds
     */
        duration: number;
    }

    /**
     * @event ScanMediaFetcherStartedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanMediaFetcherStartedEvent is triggered right before Seall starts fetching media to be matched against the local files.
     */
    function onScanMediaFetcherStarted(cb: (event: ScanMediaFetcherStartedEvent) => void): void;

    interface ScanMediaFetcherStartedEvent {
        next(): void;

        enhanced: boolean;
        enhanceWithOfflineDatabase: boolean;
        disableAnimeCollection: boolean;
    }

    /**
     * @event ScanMediaFetcherCompletedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanMediaFetcherCompletedEvent is triggered when the media fetcher completes.
     * The event includes all the media fetched from SIMKL.
     * The event includes the media IDs that are not in the user's collection.
     */
    function onScanMediaFetcherCompleted(cb: (event: ScanMediaFetcherCompletedEvent) => void): void;

    interface ScanMediaFetcherCompletedEvent {
        next(): void;

        allMedia?: Array<Media_NormalizedMedia>;
        unknownMediaIds?: Array<number>;
    }

    /**
     * @event ScanMatchingStartedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanMatchingStartedEvent is triggered when the matching process begins.
     * Prevent default to skip the default matching, in which case modified local files will be used.
     */
    function onScanMatchingStarted(cb: (event: ScanMatchingStartedEvent) => void): void;

    interface ScanMatchingStartedEvent {
        next(): void;

        preventDefault(): void;

        localFiles?: Array<Media_LocalFile>;
        normalizedMedia?: Array<Media_NormalizedMedia>;
        algorithm: string;
        threshold: number;
    }

    /**
     * @event ScanLocalFileMatchedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanLocalFileMatchedEvent is triggered when a local file is matched with media and before the match is analyzed.
     * Prevent default to skip the default analysis and override the match.
     */
    function onScanLocalFileMatched(cb: (event: ScanLocalFileMatchedEvent) => void): void;

    interface ScanLocalFileMatchedEvent {
        next(): void;

        preventDefault(): void;

        match?: Media_NormalizedMedia;
        found: boolean;
        localFile?: Media_LocalFile;
        score: number;
    }

    /**
     * @event ScanMatchingCompletedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanMatchingCompletedEvent is triggered when the matching process completes.
     */
    function onScanMatchingCompleted(cb: (event: ScanMatchingCompletedEvent) => void): void;

    interface ScanMatchingCompletedEvent {
        next(): void;

        localFiles?: Array<Media_LocalFile>;
    }

    /**
     * @event ScanHydrationStartedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanHydrationStartedEvent is triggered when the file hydration process begins.
     * Prevent default to skip the rest of the hydration process, in which case the event's local files will be used.
     */
    function onScanHydrationStarted(cb: (event: ScanHydrationStartedEvent) => void): void;

    interface ScanHydrationStartedEvent {
        next(): void;

        preventDefault(): void;

        localFiles?: Array<Media_LocalFile>;
        allMedia?: Array<Media_NormalizedMedia>;
    }

    /**
     * @event ScanLocalFileHydrationStartedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanLocalFileHydrationStartedEvent is triggered when a local file's metadata is about to be hydrated.
     * Prevent default to skip the default hydration and override the hydration.
     */
    function onScanLocalFileHydrationStarted(cb: (event: ScanLocalFileHydrationStartedEvent) => void): void;

    interface ScanLocalFileHydrationStartedEvent {
        next(): void;

        preventDefault(): void;

        localFile?: Media_LocalFile;
        media?: Media_NormalizedMedia;
    }

    /**
     * @event ScanLocalFileHydratedEvent
     * @file internal/library/scanner/hook_events.go
     * @description
     * ScanLocalFileHydratedEvent is triggered when a local file's metadata is hydrated
     */
    function onScanLocalFileHydrated(cb: (event: ScanLocalFileHydratedEvent) => void): void;

    interface ScanLocalFileHydratedEvent {
        next(): void;

        localFile?: Media_LocalFile;
        mediaId: number;
        episode: number;
    }


    /**
     * @package torrent
     */

    /**
     * @event TorrentSearchRequestedEvent
     * @file internal/torrents/torrent/hook_events.go
     * @description
     * TorrentSearchRequestedEvent is triggered before Seall searches anime torrents.
     * Prevent default to skip the native search and return SearchData.
     */
    function onTorrentSearchRequested(cb: (event: TorrentSearchRequestedEvent) => void): void;

    interface TorrentSearchRequestedEvent {
        next(): void;

        preventDefault(): void;

        options: Torrent_AnimeSearchOptions;
        searchData?: Torrent_SearchData;
    }

    /**
     * @event TorrentSearchEvent
     * @file internal/torrents/torrent/hook_events.go
     * @description
     * TorrentSearchEvent is triggered after Seall assembles the torrent search response.
     * Handlers can mutate SearchData before it is cached and returned.
     */
    function onTorrentSearch(cb: (event: TorrentSearchEvent) => void): void;

    interface TorrentSearchEvent {
        next(): void;

        options: Torrent_AnimeSearchOptions;
        searchData?: Torrent_SearchData;
    }


    /**
     * @package torrentstream
     */

    /**
     * @event TorrentStreamAutoSelectTorrentsFetchedEvent
     * @file internal/torrentstream/hook_events.go
     * @description
     * TorrentStreamAutoSelectTorrentsFetchedEvent is triggered when the torrents are fetched for auto select.
     * The torrents are sorted by seeders from highest to lowest.
     * This event is triggered before the top 3 torrents are analyzed.
     */
    function onTorrentStreamAutoSelectTorrentsFetched(cb: (event: TorrentStreamAutoSelectTorrentsFetchedEvent) => void): void;

    interface TorrentStreamAutoSelectTorrentsFetchedEvent {
        next(): void;

        Torrents?: Array<HibikeTorrent_MediaTorrent>;
    }

    /**
     * @event TorrentStreamSendStreamToMediaPlayerEvent
     * @file internal/torrentstream/hook_events.go
     * @description
     * TorrentStreamSendStreamToMediaPlayerEvent is triggered when the torrent stream is about to send a stream to the media player.
     * Prevent default to skip the default playback and override the playback.
     */
    function onTorrentStreamSendStreamToMediaPlayer(cb: (event: TorrentStreamSendStreamToMediaPlayerEvent) => void): void;

    interface TorrentStreamSendStreamToMediaPlayerEvent {
        next(): void;

        preventDefault(): void;

        windowTitle: string;
        streamURL: string;
        media?: MediaAPI_BaseMedia;
        aniDbEpisode: string;
        playbackType: string;
    }

    ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ///////////////////////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollection {
        MediaListCollection?: MediaAPI_MediaCollection_MediaListCollection;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollectionWithRelations {
        MediaListCollection?: MediaAPI_MediaCollectionWithRelations_MediaListCollection;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollectionWithRelations_MediaListCollection {
        lists?: Array<MediaAPI_MediaCollectionWithRelations_MediaListCollection_Lists>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollectionWithRelations_MediaListCollection_Lists {
        entries?: Array<MediaAPI_MediaCollectionWithRelations_MediaListCollection_Lists_Entries>;
        isCustomList?: boolean;
        name?: string;
        status?: MediaAPI_MediaListStatus;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollectionWithRelations_MediaListCollection_Lists_Entries {
        completedAt?: MediaAPI_MediaCollectionWithRelations_MediaListCollection_Lists_Entries_CompletedAt;
        id: number;
        media?: MediaAPI_CompleteMedia;
        notes?: string;
        private?: boolean;
        progress?: number;
        repeat?: number;
        score?: number;
        startedAt?: MediaAPI_MediaCollectionWithRelations_MediaListCollection_Lists_Entries_StartedAt;
        status?: MediaAPI_MediaListStatus;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollectionWithRelations_MediaListCollection_Lists_Entries_CompletedAt {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollectionWithRelations_MediaListCollection_Lists_Entries_StartedAt {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollection_MediaListCollection {
        lists?: Array<MediaAPI_MediaCollection_MediaListCollection_Lists>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollection_MediaListCollection_Lists {
        entries?: Array<MediaAPI_MediaCollection_MediaListCollection_Lists_Entries>;
        isCustomList?: boolean;
        name?: string;
        status?: MediaAPI_MediaListStatus;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollection_MediaListCollection_Lists_Entries {
        completedAt?: MediaAPI_MediaCollection_MediaListCollection_Lists_Entries_CompletedAt;
        id: number;
        media?: MediaAPI_BaseMedia;
        notes?: string;
        private?: boolean;
        progress?: number;
        repeat?: number;
        score?: number;
        startedAt?: MediaAPI_MediaCollection_MediaListCollection_Lists_Entries_StartedAt;
        status?: MediaAPI_MediaListStatus;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollection_MediaListCollection_Lists_Entries_CompletedAt {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaCollection_MediaListCollection_Lists_Entries_StartedAt {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media {
        averageScore?: number;
        characters?: MediaAPI_MediaDetailsById_Media_Characters;
        description?: string;
        duration?: number;
        endDate?: MediaAPI_MediaDetailsById_Media_EndDate;
        genres?: Array<string>;
        id: number;
        meanScore?: number;
        popularity?: number;
        rankings?: Array<MediaAPI_MediaDetailsById_Media_Rankings>;
        recommendations?: MediaAPI_MediaDetailsById_Media_Recommendations;
        relations?: MediaAPI_MediaDetailsById_Media_Relations;
        siteUrl?: string;
        staff?: MediaAPI_MediaDetailsById_Media_Staff;
        startDate?: MediaAPI_MediaDetailsById_Media_StartDate;
        studios?: MediaAPI_MediaDetailsById_Media_Studios;
        trailer?: MediaAPI_MediaDetailsById_Media_Trailer;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Characters {
        edges?: Array<MediaAPI_MediaDetailsById_Media_Characters_Edges>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Characters_Edges {
        id?: number;
        name?: string;
        node?: MediaAPI_BaseCharacter;
        role?: MediaAPI_CharacterRole;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_EndDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Rankings {
        allTime?: boolean;
        context: string;
        format: MediaAPI_MediaFormat;
        rank: number;
        season?: MediaAPI_MediaSeason;
        type: MediaAPI_MediaRankType;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Recommendations {
        edges?: Array<MediaAPI_MediaDetailsById_Media_Recommendations_Edges>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Recommendations_Edges {
        node?: MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node {
        mediaRecommendation?: MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation {
        bannerImage?: string;
        coverImage?: MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_CoverImage;
        description?: string;
        episodes?: number;
        format?: MediaAPI_MediaFormat;
        id: number;
        idMal?: number;
        isAdult?: boolean;
        meanScore?: number;
        season?: MediaAPI_MediaSeason;
        siteUrl?: string;
        startDate?: MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_StartDate;
        status?: MediaAPI_MediaStatus;
        title?: MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_Title;
        trailer?: MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_Trailer;
        type?: MediaAPI_MediaType;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_CoverImage {
        color?: string;
        extraLarge?: string;
        large?: string;
        medium?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_StartDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_Title {
        english?: string;
        native?: string;
        romaji?: string;
        userPreferred?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_Trailer {
        id?: string;
        site?: string;
        thumbnail?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Relations {
        edges?: Array<MediaAPI_MediaDetailsById_Media_Relations_Edges>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Relations_Edges {
        node?: MediaAPI_BaseMedia;
        relationType?: MediaAPI_MediaRelation;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Staff {
        edges?: Array<MediaAPI_MediaDetailsById_Media_Staff_Edges>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Staff_Edges {
        node?: MediaAPI_MediaDetailsById_Media_Staff_Edges_Node;
        role?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Staff_Edges_Node {
        id: number;
        name?: MediaAPI_MediaDetailsById_Media_Staff_Edges_Node_Name;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Staff_Edges_Node_Name {
        full?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_StartDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Studios {
        nodes?: Array<MediaAPI_MediaDetailsById_Media_Studios_Nodes>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Studios_Nodes {
        id: number;
        name: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MediaDetailsById_Media_Trailer {
        id?: string;
        site?: string;
        thumbnail?: string;
    }

    /**
     * - Filepath: internal/api/simkl/collection_helper.go
     */
    export type MediaAPI_MediaListEntry = MediaAPI_MediaCollection_MediaListCollection_Lists_Entries;

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseMedia {
        id: number;
        idMal?: number;
        siteUrl?: string;
        status?: MediaAPI_MediaStatus;
        season?: MediaAPI_MediaSeason;
        type?: MediaAPI_MediaType;
        format?: MediaAPI_MediaFormat;
        seasonYear?: number;
        bannerImage?: string;
        episodes?: number;
        synonyms?: Array<string>;
        isAdult?: boolean;
        countryOfOrigin?: string;
        meanScore?: number;
        description?: string;
        genres?: Array<string>;
        duration?: number;
        trailer?: MediaAPI_BaseMedia_Trailer;
        title?: MediaAPI_BaseMedia_Title;
        coverImage?: MediaAPI_BaseMedia_CoverImage;
        startDate?: MediaAPI_BaseMedia_StartDate;
        endDate?: MediaAPI_BaseMedia_EndDate;
        nextAiringEpisode?: MediaAPI_BaseMedia_NextAiringEpisode;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseMedia_CoverImage {
        color?: string;
        extraLarge?: string;
        large?: string;
        medium?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseMedia_EndDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseMedia_NextAiringEpisode {
        airingAt: number;
        episode: number;
        timeUntilAiring: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseMedia_StartDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseMedia_Title {
        english?: string;
        native?: string;
        romaji?: string;
        userPreferred?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseMedia_Trailer {
        id?: string;
        site?: string;
        thumbnail?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseCharacter {
        id: number;
        isFavourite: boolean;
        gender?: string;
        age?: string;
        dateOfBirth?: MediaAPI_BaseCharacter_DateOfBirth;
        name?: MediaAPI_BaseCharacter_Name;
        image?: MediaAPI_BaseCharacter_Image;
        description?: string;
        siteUrl?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseCharacter_DateOfBirth {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseCharacter_Image {
        large?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseCharacter_Name {
        alternative?: Array<string>;
        full?: string;
        native?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseManga {
        id: number;
        idMal?: number;
        siteUrl?: string;
        status?: MediaAPI_MediaStatus;
        season?: MediaAPI_MediaSeason;
        type?: MediaAPI_MediaType;
        format?: MediaAPI_MediaFormat;
        bannerImage?: string;
        chapters?: number;
        volumes?: number;
        synonyms?: Array<string>;
        isAdult?: boolean;
        countryOfOrigin?: string;
        meanScore?: number;
        description?: string;
        genres?: Array<string>;
        title?: MediaAPI_BaseManga_Title;
        coverImage?: MediaAPI_BaseManga_CoverImage;
        startDate?: MediaAPI_BaseManga_StartDate;
        endDate?: MediaAPI_BaseManga_EndDate;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseManga_CoverImage {
        color?: string;
        extraLarge?: string;
        large?: string;
        medium?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseManga_EndDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseManga_StartDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_BaseManga_Title {
        english?: string;
        native?: string;
        romaji?: string;
        userPreferred?: string;
    }

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     * @description
     *  The role the character plays in the media
     */
    export type MediaAPI_CharacterRole = "MAIN" | "SUPPORTING" | "BACKGROUND";

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_CompleteMedia {
        id: number;
        idMal?: number;
        siteUrl?: string;
        status?: MediaAPI_MediaStatus;
        season?: MediaAPI_MediaSeason;
        seasonYear?: number;
        type?: MediaAPI_MediaType;
        format?: MediaAPI_MediaFormat;
        bannerImage?: string;
        episodes?: number;
        synonyms?: Array<string>;
        isAdult?: boolean;
        countryOfOrigin?: string;
        meanScore?: number;
        description?: string;
        genres?: Array<string>;
        duration?: number;
        trailer?: MediaAPI_CompleteMedia_Trailer;
        title?: MediaAPI_CompleteMedia_Title;
        coverImage?: MediaAPI_CompleteMedia_CoverImage;
        startDate?: MediaAPI_CompleteMedia_StartDate;
        endDate?: MediaAPI_CompleteMedia_EndDate;
        nextAiringEpisode?: MediaAPI_CompleteMedia_NextAiringEpisode;
        relations?: MediaAPI_CompleteMedia_Relations;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_CompleteMedia_CoverImage {
        color?: string;
        extraLarge?: string;
        large?: string;
        medium?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_CompleteMedia_EndDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_CompleteMedia_NextAiringEpisode {
        airingAt: number;
        episode: number;
        timeUntilAiring: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_CompleteMedia_Relations {
        edges?: Array<MediaAPI_CompleteMedia_Relations_Edges>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_CompleteMedia_Relations_Edges {
        node?: MediaAPI_BaseMedia;
        relationType?: MediaAPI_MediaRelation;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_CompleteMedia_StartDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_CompleteMedia_Title {
        english?: string;
        native?: string;
        romaji?: string;
        userPreferred?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_CompleteMedia_Trailer {
        id?: string;
        site?: string;
        thumbnail?: string;
    }

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     * @description
     *  Date object that allows for incomplete date values (fuzzy)
     */
    interface MediaAPI_FuzzyDateInput {
        year?: number;
        month?: number;
        day?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListMedia {
        Page?: MediaAPI_ListMedia_Page;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListMedia_Page {
        media?: Array<MediaAPI_BaseMedia>;
        pageInfo?: MediaAPI_ListMedia_Page_PageInfo;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListMedia_Page_PageInfo {
        currentPage?: number;
        hasNextPage?: boolean;
        lastPage?: number;
        perPage?: number;
        total?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListManga {
        Page?: MediaAPI_ListManga_Page;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListManga_Page {
        media?: Array<MediaAPI_BaseManga>;
        pageInfo?: MediaAPI_ListManga_Page_PageInfo;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListManga_Page_PageInfo {
        currentPage?: number;
        hasNextPage?: boolean;
        lastPage?: number;
        perPage?: number;
        total?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListRecentMedia {
        Page?: MediaAPI_ListRecentMedia_Page;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListRecentMedia_Page {
        airingSchedules?: Array<MediaAPI_ListRecentMedia_Page_AiringSchedules>;
        pageInfo?: MediaAPI_ListRecentMedia_Page_PageInfo;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListRecentMedia_Page_AiringSchedules {
        airingAt: number;
        episode: number;
        id: number;
        media?: MediaAPI_BaseMedia;
        timeUntilAiring: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_ListRecentMedia_Page_PageInfo {
        currentPage?: number;
        hasNextPage?: boolean;
        lastPage?: number;
        perPage?: number;
        total?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaCollection {
        MediaListCollection?: MediaAPI_MangaCollection_MediaListCollection;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaCollection_MediaListCollection {
        lists?: Array<MediaAPI_MangaCollection_MediaListCollection_Lists>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaCollection_MediaListCollection_Lists {
        entries?: Array<MediaAPI_MangaCollection_MediaListCollection_Lists_Entries>;
        isCustomList?: boolean;
        name?: string;
        status?: MediaAPI_MediaListStatus;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaCollection_MediaListCollection_Lists_Entries {
        completedAt?: MediaAPI_MangaCollection_MediaListCollection_Lists_Entries_CompletedAt;
        id: number;
        media?: MediaAPI_BaseManga;
        notes?: string;
        private?: boolean;
        progress?: number;
        repeat?: number;
        score?: number;
        startedAt?: MediaAPI_MangaCollection_MediaListCollection_Lists_Entries_StartedAt;
        status?: MediaAPI_MediaListStatus;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaCollection_MediaListCollection_Lists_Entries_CompletedAt {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaCollection_MediaListCollection_Lists_Entries_StartedAt {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media {
        characters?: MediaAPI_MangaDetailsById_Media_Characters;
        duration?: number;
        genres?: Array<string>;
        id: number;
        rankings?: Array<MediaAPI_MangaDetailsById_Media_Rankings>;
        recommendations?: MediaAPI_MangaDetailsById_Media_Recommendations;
        relations?: MediaAPI_MangaDetailsById_Media_Relations;
        siteUrl?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Characters {
        edges?: Array<MediaAPI_MangaDetailsById_Media_Characters_Edges>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Characters_Edges {
        id?: number;
        name?: string;
        node?: MediaAPI_BaseCharacter;
        role?: MediaAPI_CharacterRole;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Rankings {
        allTime?: boolean;
        context: string;
        format: MediaAPI_MediaFormat;
        rank: number;
        season?: MediaAPI_MediaSeason;
        type: MediaAPI_MediaRankType;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Recommendations {
        edges?: Array<MediaAPI_MangaDetailsById_Media_Recommendations_Edges>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Recommendations_Edges {
        node?: MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node {
        mediaRecommendation?: MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation {
        bannerImage?: string;
        chapters?: number;
        countryOfOrigin?: string;
        coverImage?: MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_CoverImage;
        description?: string;
        endDate?: MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_EndDate;
        format?: MediaAPI_MediaFormat;
        id: number;
        idMal?: number;
        isAdult?: boolean;
        meanScore?: number;
        season?: MediaAPI_MediaSeason;
        siteUrl?: string;
        startDate?: MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_StartDate;
        status?: MediaAPI_MediaStatus;
        synonyms?: Array<string>;
        title?: MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_Title;
        type?: MediaAPI_MediaType;
        volumes?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_CoverImage {
        color?: string;
        extraLarge?: string;
        large?: string;
        medium?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_EndDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_StartDate {
        day?: number;
        month?: number;
        year?: number;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Recommendations_Edges_Node_MediaRecommendation_Title {
        english?: string;
        native?: string;
        romaji?: string;
        userPreferred?: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Relations {
        edges?: Array<MediaAPI_MangaDetailsById_Media_Relations_Edges>;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_MangaDetailsById_Media_Relations_Edges {
        node?: MediaAPI_BaseManga;
        relationType?: MediaAPI_MediaRelation;
    }

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     * @description
     *  The format the media was released in
     */
    export type MediaAPI_MediaFormat = "TV" |
    "TV_SHORT" |
    "MOVIE" |
    "SPECIAL" |
    "OVA" |
    "ONA" |
    "MUSIC" |
    "MANGA" |
    "NOVEL" |
    "ONE_SHOT";

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     * @description
     *  Media list watching/reading status enum.
     */
    export type MediaAPI_MediaListStatus = "CURRENT" |
    "PLANNING" |
    "COMPLETED" |
    "DROPPED" |
    "PAUSED" |
    "REPEATING";

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     * @description
     *  The type of ranking
     */
    export type MediaAPI_MediaRankType = "RATED" | "POPULAR";

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     * @description
     *  Type of relation media has to its parent.
     */
    export type MediaAPI_MediaRelation = "ADAPTATION" |
    "PREQUEL" |
    "SEQUEL" |
    "PARENT" |
    "SIDE_STORY" |
    "CHARACTER" |
    "SUMMARY" |
    "ALTERNATIVE" |
    "SPIN_OFF" |
    "OTHER" |
    "SOURCE" |
    "COMPILATION" |
    "CONTAINS";

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     */
    export type MediaAPI_MediaSeason = "WINTER" | "SPRING" | "SUMMER" | "FALL";

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     * @description
     *  Media sort enums
     */
    export type MediaAPI_MediaSort = "ID" |
    "ID_DESC" |
    "TITLE_ROMAJI" |
    "TITLE_ROMAJI_DESC" |
    "TITLE_ENGLISH" |
    "TITLE_ENGLISH_DESC" |
    "TITLE_NATIVE" |
    "TITLE_NATIVE_DESC" |
    "TYPE" |
    "TYPE_DESC" |
    "FORMAT" |
    "FORMAT_DESC" |
    "START_DATE" |
    "START_DATE_DESC" |
    "END_DATE" |
    "END_DATE_DESC" |
    "SCORE" |
    "SCORE_DESC" |
    "POPULARITY" |
    "POPULARITY_DESC" |
    "TRENDING" |
    "TRENDING_DESC" |
    "EPISODES" |
    "EPISODES_DESC" |
    "DURATION" |
    "DURATION_DESC" |
    "STATUS" |
    "STATUS_DESC" |
    "CHAPTERS" |
    "CHAPTERS_DESC" |
    "VOLUMES" |
    "VOLUMES_DESC" |
    "UPDATED_AT" |
    "UPDATED_AT_DESC" |
    "SEARCH_MATCH" |
    "FAVOURITES" |
    "FAVOURITES_DESC";

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     * @description
     *  The current releasing status of the media
     */
    export type MediaAPI_MediaStatus = "FINISHED" | "RELEASING" | "NOT_YET_RELEASED" | "CANCELLED" | "HIATUS";

    /**
     * - Filepath: internal/api/simkl/models_gen.go
     * @description
     *  Media type enum, anime or manga.
     */
    export type MediaAPI_MediaType = "ANIME" | "MANGA";

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_StudioDetails {
        Studio?: MediaAPI_StudioDetails_Studio;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_StudioDetails_Studio {
        id: number;
        isAnimationStudio: boolean;
        media?: MediaAPI_StudioDetails_Studio_Media;
        name: string;
    }

    /**
     * - Filepath: internal/api/simkl/client_gen.go
     */
    interface MediaAPI_StudioDetails_Studio_Media {
        nodes?: Array<MediaAPI_BaseMedia>;
    }

    /**
     * - Filepath: internal/api/animap/animap.go
     */
    interface Animap_Anime {
        title: string;
        titles?: Record<string, string>;
        /**
         * YYYY-MM-DD
         */
        startDate?: string;
        /**
         * YYYY-MM-DD
         */
        endDate?: string;
        /**
         * Finished, Airing, Upcoming, etc.
         */
        status: string;
        /**
         * TV, OVA, Movie, etc.
         */
        type: string;
        /**
         * Indexed by AniDB episode number, "1", "S1", etc.
         */
        episodes?: Record<string, Animap_Episode>;
        mappings?: Animap_AnimeMapping;
    }

    /**
     * - Filepath: internal/api/animap/animap.go
     */
    interface Animap_AnimeMapping {
        anidb_id?: number;
        simkl_id?: number;
        kitsu_id?: number;
        thetvdb_id?: number;
        /**
         * Can be int or string, forced to string
         */
        themoviedb_id?: string;
        mal_id?: number;
        livechart_id?: number;
        /**
         * Can be int or string, forced to string
         */
        animeplanet_id?: string;
        anisearch_id?: number;
        simkl_id?: number;
        notifymoe_id?: string;
        animecountdown_id?: number;
        type?: string;
    }

    /**
     * - Filepath: internal/api/animap/animap.go
     */
    interface Animap_Episode {
        anidbEpisode: string;
        anidbEid: number;
        tvdbEid?: number;
        tvdbShowId?: number;
        /**
         * YYYY-MM-DD
         */
        airDate?: string;
        /**
         * Title of the episode from AniDB
         */
        anidbTitle?: string;
        /**
         * Title of the episode from TVDB
         */
        tvdbTitle?: string;
        overview?: string;
        image?: string;
        /**
         * minutes
         */
        runtime?: number;
        /**
         * Xm
         */
        length?: string;
        seasonNumber?: number;
        seasonName?: string;
        number: number;
        absoluteNumber?: number;
    }

    /**
     * - Filepath: internal/library/anime/autodownloader_types.go
     */
    interface Media_AutoDownloaderCondition {
        id: string;
        term: string;
        isRegex: boolean;
        action: Media_AutoDownloaderProfileRuleFormatAction;
        /**
         * Only used if Action == "score"
         */
        score: number;
    }

    /**
     * - Filepath: internal/library/anime/autodownloader_types.go
     */
    interface Media_AutoDownloaderProfile {
        dbId: number;
        name: string;
        global: boolean;
        releaseGroups?: Array<string>;
        resolutions?: Array<string>;
        conditions?: Array<Media_AutoDownloaderCondition>;
        minimumScore: number;
        minSeeders?: number;
        minSize?: string;
        maxSize?: string;
        delayMinutes: number;
        skipDelayScore: number;
        providers?: Array<string>;
    }

    /**
     * - Filepath: internal/library/anime/autodownloader_types.go
     */
    export type Media_AutoDownloaderProfileRuleFormatAction = "score" | "block" | "require";

    /**
     * - Filepath: internal/library/anime/autodownloader_types.go
     */
    interface Media_AutoDownloaderRule {
        dbId: number;
        enabled: boolean;
        mediaId: number;
        destination: string;
        profileId?: number;
        releaseGroups?: Array<string>;
        resolutions?: Array<string>;
        episodeNumbers?: Array<number>;
        episodeType: Media_AutoDownloaderRuleEpisodeType;
        comparisonTitle: string;
        titleComparisonType: Media_AutoDownloaderRuleTitleComparisonType;
        additionalTerms?: Array<string>;
        excludeTerms?: Array<string>;
        minSeeders: number;
        minSize: string;
        maxSize: string;
        customEpisodeNumberAbsoluteOffset?: number;
        providers?: Array<string>;
    }

    /**
     * - Filepath: internal/library/anime/autodownloader_types.go
     */
    export type Media_AutoDownloaderRuleEpisodeType = "recent" | "selected";

    /**
     * - Filepath: internal/library/anime/autodownloader_types.go
     */
    export type Media_AutoDownloaderRuleTitleComparisonType = "contains" | "likely";

    /**
     * - Filepath: internal/library/anime/entry.go
     */
    interface Media_Entry {
        mediaId: number;
        media?: MediaAPI_BaseMedia;
        listData?: Media_EntryListData;
        libraryData?: Media_EntryLibraryData;
        downloadInfo?: Media_EntryDownloadInfo;
        episodes?: Array<Media_Episode>;
        nextEpisode?: Media_Episode;
        localFiles?: Array<Media_LocalFile>;
        anidbId: number;
        currentEpisodeCount: number;
        _isNakamaEntry: boolean;
        nakamaLibraryData?: Media_NakamaEntryLibraryData;
    }

    /**
     * - Filepath: internal/library/anime/entry_download_info.go
     */
    interface Media_EntryDownloadEpisode {
        episodeNumber: number;
        aniDBEpisode: string;
        episode?: Media_Episode;
    }

    /**
     * - Filepath: internal/library/anime/entry_download_info.go
     */
    interface Media_EntryDownloadInfo {
        episodesToDownload?: Array<Media_EntryDownloadEpisode>;
        canBatch: boolean;
        batchAll: boolean;
        hasInaccurateSchedule: boolean;
        rewatch: boolean;
        absoluteOffset: number;
    }

    /**
     * - Filepath: internal/library/anime/entry_library_data.go
     */
    interface Media_EntryLibraryData {
        allFilesLocked: boolean;
        sharedPath: string;
        unwatchedCount: number;
        mainFileCount: number;
    }

    /**
     * - Filepath: internal/library/anime/entry.go
     */
    interface Media_EntryListData {
        progress?: number;
        score?: number;
        status?: MediaAPI_MediaListStatus;
        repeat?: number;
        startedAt?: string;
        completedAt?: string;
    }

    /**
     * - Filepath: internal/library/anime/episode.go
     */
    interface Media_Episode {
        type: Media_LocalFileType;
        /**
         * e.g, Show: "Episode 1", Movie: "Violet Evergarden The Movie"
         */
        displayTitle: string;
        /**
         * e.g, "Shibuya Incident - Gate, Open"
         */
        episodeTitle: string;
        episodeNumber: number;
        /**
         * AniDB episode number
         */
        aniDBEpisode?: string;
        absoluteEpisodeNumber: number;
        /**
         * Usually the same as EpisodeNumber, unless there is a discrepancy between SIMKL and AniDB
         */
        progressNumber: number;
        localFile?: Media_LocalFile;
        /**
         * Is in the local files
         */
        isDownloaded: boolean;
        /**
         * (image, airDate, length, summary, overview)
         */
        episodeMetadata?: Media_EpisodeMetadata;
        /**
         * (episode, aniDBEpisode, type...)
         */
        fileMetadata?: Media_LocalFileMetadata;
        /**
         * No AniDB data
         */
        isInvalid: boolean;
        /**
         * Alerts the user that there is a discrepancy between SIMKL and AniDB
         */
        metadataIssue?: string;
        baseAnime?: MediaAPI_BaseMedia;
        _isNakamaEpisode: boolean;
    }

    /**
     * - Filepath: internal/library/anime/episode_collection.go
     */
    interface Media_EpisodeCollection {
        hasMappingError: boolean;
        episodes?: Array<Media_Episode>;
        metadata?: Metadata_AnimeMetadata;
    }

    /**
     * - Filepath: internal/library/anime/episode.go
     */
    interface Media_EpisodeMetadata {
        anidbId?: number;
        image?: string;
        airDate?: string;
        length?: number;
        summary?: string;
        overview?: string;
        isFiller?: boolean;
        /**
         * Indicates if the episode has a real image
         */
        hasImage?: boolean;
        title?: string;
    }

    /**
     * - Filepath: internal/library/anime/collection.go
     */
    interface Media_LibraryCollection {
        continueWatchingList?: Array<Media_Episode>;
        lists?: Array<Media_LibraryCollectionList>;
        unmatchedLocalFiles?: Array<Media_LocalFile>;
        unmatchedGroups?: Array<Media_UnmatchedGroup>;
        ignoredLocalFiles?: Array<Media_LocalFile>;
        unknownGroups?: Array<Media_UnknownGroup>;
        stats?: Media_LibraryCollectionStats;
        /**
         * Hydrated by the route handler
         */
        stream?: Media_StreamCollection;
    }

    /**
     * - Filepath: internal/library/anime/collection.go
     */
    interface Media_LibraryCollectionEntry {
        media?: MediaAPI_BaseMedia;
        mediaId: number;
        /**
         * Library data
         */
        libraryData?: Media_EntryLibraryData;
        /**
         * Library data from Nakama
         */
        nakamaLibraryData?: Media_NakamaEntryLibraryData;
        /**
         * SIMKL list data
         */
        listData?: Media_EntryListData;
    }

    /**
     * - Filepath: internal/library/anime/collection.go
     */
    interface Media_LibraryCollectionList {
        type?: MediaAPI_MediaListStatus;
        status?: MediaAPI_MediaListStatus;
        entries?: Array<Media_LibraryCollectionEntry>;
    }

    /**
     * - Filepath: internal/library/anime/collection.go
     */
    interface Media_LibraryCollectionStats {
        totalEntries: number;
        totalFiles: number;
        totalShows: number;
        totalMovies: number;
        totalSpecials: number;
        totalSize: string;
    }

    /**
     * - Filepath: internal/library/anime/localfile.go
     */
    interface Media_LocalFile {
        path: string;
        name: string;
        parsedInfo?: Media_LocalFileParsedData;
        parsedFolderInfo?: Array<Media_LocalFileParsedData>;
        metadata?: Media_LocalFileMetadata;
        locked: boolean;
        /**
         * Unused for now
         */
        ignored: boolean;
        mediaId: number;
    }

    /**
     * - Filepath: internal/library/anime/localfile.go
     */
    interface Media_LocalFileMetadata {
        episode: number;
        aniDBEpisode: string;
        type: Media_LocalFileType;
    }

    /**
     * - Filepath: internal/library/anime/localfile.go
     */
    interface Media_LocalFileParsedData {
        original: string;
        title?: string;
        releaseGroup?: string;
        season?: string;
        seasonRange?: Array<string>;
        part?: string;
        partRange?: Array<string>;
        episode?: string;
        episodeRange?: Array<string>;
        episodeTitle?: string;
        year?: string;
    }

    /**
     * - Filepath: internal/library/anime/localfile.go
     */
    export type Media_LocalFileType = "main" | "special" | "nc";

    /**
     * - Filepath: internal/library/anime/localfile_wrapper.go
     */
    interface Media_LocalFileWrapperEntry {
        mediaId: number;
        localFiles?: Array<Media_LocalFile>;
    }

    /**
     * - Filepath: internal/library/anime/missing_episodes.go
     */
    interface Media_MissingEpisodes {
        episodes?: Array<Media_Episode>;
        silencedEpisodes?: Array<Media_Episode>;
    }

    /**
     * - Filepath: internal/library/anime/entry_library_data.go
     */
    interface Media_NakamaEntryLibraryData {
        unwatchedCount: number;
        mainFileCount: number;
    }

    /**
     * - Filepath: internal/library/anime/normalized_media.go
     */
    interface Media_NormalizedMedia {
        ID: number;
        IdMal?: number;
        Title?: Media_NormalizedMediaTitle;
        Synonyms?: Array<string>;
        Format?: MediaAPI_MediaFormat;
        Status?: MediaAPI_MediaStatus;
        Season?: MediaAPI_MediaSeason;
        Year?: number;
        StartDate?: Media_NormalizedMediaDate;
        Episodes?: number;
        BannerImage?: string;
        CoverImage?: Media_NormalizedMediaCoverImage;
        NextAiringEpisode?: Media_NormalizedMediaNextAiringEpisode;
        fetched: boolean;
    }

    /**
     * - Filepath: internal/library/anime/normalized_media.go
     */
    interface Media_NormalizedMediaCoverImage {
        ExtraLarge?: string;
        Large?: string;
        Medium?: string;
        Color?: string;
    }

    /**
     * - Filepath: internal/library/anime/normalized_media.go
     */
    interface Media_NormalizedMediaDate {
        Year?: number;
        Month?: number;
        Day?: number;
    }

    /**
     * - Filepath: internal/library/anime/normalized_media.go
     */
    interface Media_NormalizedMediaNextAiringEpisode {
        AiringAt: number;
        TimeUntilAiring: number;
        Episode: number;
    }

    /**
     * - Filepath: internal/library/anime/normalized_media.go
     */
    interface Media_NormalizedMediaTitle {
        Romaji?: string;
        English?: string;
        Native?: string;
        UserPreferred?: string;
    }

    /**
     * - Filepath: internal/library/anime/schedule.go
     */
    interface Media_ScheduleItem {
        mediaId: number;
        title: string;
        time: string;
        dateTime?: string;
        image: string;
        episodeNumber: number;
        isMovie: boolean;
        isSeasonFinale: boolean;
    }

    /**
     * - Filepath: internal/library/anime/collection.go
     */
    interface Media_StreamCollection {
        continueWatchingList?: Array<Media_Episode>;
        anime?: Array<MediaAPI_BaseMedia>;
        listData?: Record<number, Media_EntryListData>;
    }

    /**
     * - Filepath: internal/library/anime/collection.go
     */
    interface Media_UnknownGroup {
        mediaId: number;
        localFiles?: Array<Media_LocalFile>;
    }

    /**
     * - Filepath: internal/library/anime/collection.go
     */
    interface Media_UnmatchedGroup {
        dir: string;
        localFiles?: Array<Media_LocalFile>;
        suggestions?: Array<MediaAPI_BaseMedia>;
    }

    /**
     * - Filepath: internal/library/anime/upcoming_episodes.go
     */
    interface Media_UpcomingEpisode {
        mediaId: number;
        episodeNumber: number;
        airingAt: number;
        timeUntilAiring: number;
        baseAnime?: MediaAPI_BaseMedia;
        episodeMetadata?: Media_EpisodeMetadata;
    }

    /**
     * - Filepath: internal/library/anime/upcoming_episodes.go
     */
    interface Media_UpcomingEpisodes {
        episodes?: Array<Media_UpcomingEpisode>;
    }

    /**
     * - Filepath: internal/api/anizip/anizip.go
     */
    interface Anizip_Episode {
        tvdbEid?: number;
        airdate?: string;
        seasonNumber?: number;
        episodeNumber?: number;
        absoluteEpisodeNumber?: number;
        title?: Record<string, string>;
        image?: string;
        summary?: string;
        overview?: string;
        runtime?: number;
        length?: number;
        episode?: string;
        anidbEid?: number;
        rating?: string;
    }

    /**
     * - Filepath: internal/api/anizip/anizip.go
     */
    interface Anizip_Mappings {
        animeplanet_id?: string;
        kitsu_id?: number;
        mal_id?: number;
        type?: string;
        simkl_id?: number;
        anisearch_id?: number;
        anidb_id?: number;
        notifymoe_id?: string;
        livechart_id?: number;
        thetvdb_id?: number;
        imdb_id?: string;
        themoviedb_id?: string;
    }

    /**
     * - Filepath: internal/api/anizip/anizip.go
     */
    interface Anizip_Media {
        titles?: Record<string, string>;
        episodes?: Record<string, Anizip_Episode>;
        episodeCount: number;
        specialCount: number;
        mappings?: Anizip_Mappings;
    }

    /**
     * - Filepath: internal/library/autodownloader/autodownloader.go
     * @description
     *  Candidate represents a potential torrent to download with its score
     */
    interface AutoDownloader_Candidate {
        Torrent?: AutoDownloader_NormalizedTorrent;
        Score: number;
    }

    /**
     * - Filepath: internal/library/autodownloader/autodownloader_torrent.go
     */
    interface AutoDownloader_NormalizedTorrent {
        parsedData?: $habari.Metadata;
        /**
         * Access using GetMagnet()
         */
        magnet: string;
        ExtensionID: string;
        provider?: string;
        name: string;
        date: string;
        size: number;
        formattedSize: string;
        seeders: number;
        leechers: number;
        downloadCount: number;
        link: string;
        downloadUrl: string;
        magnetLink?: string;
        infoHash?: string;
        resolution?: string;
        isBatch?: boolean;
        episodeNumber?: number;
        releaseGroup?: string;
        isBestRelease: boolean;
        confirmed: boolean;
    }

    /**
     * - Filepath: internal/continuity/manager.go
     */
    export type Continuity_Kind = "onlinestream" | "mediastream" | "external_player";

    /**
     * - Filepath: internal/continuity/history.go
     */
    interface Continuity_UpdateWatchHistoryItemOptions {
        currentTime: number;
        duration: number;
        mediaId: number;
        episodeNumber: number;
        filepath?: string;
        kind: Continuity_Kind;
    }

    /**
     * - Filepath: internal/continuity/history.go
     */
    export type Continuity_WatchHistory = Record<number, Continuity_WatchHistoryItem>;

    /**
     * - Filepath: internal/continuity/history.go
     */
    interface Continuity_WatchHistoryItem {
        kind: Continuity_Kind;
        filepath: string;
        mediaId: number;
        episodeNumber: number;
        currentTime: number;
        duration: number;
        timeAdded?: string;
        timeUpdated?: string;
    }

    /**
     * - Filepath: internal/continuity/history.go
     */
    interface Continuity_WatchHistoryItemResponse {
        item?: Continuity_WatchHistoryItem;
        found: boolean;
    }

    /**
     * - Filepath: internal/debrid/debrid/debrid.go
     */
    interface Debrid_AddTorrentOptions {
        magnetLink: string;
        infoHash: string;
        /**
         * Real-Debrid only, ID, IDs, or "all"
         */
        selectFileId: string;
    }

    /**
     * - Filepath: internal/debrid/debrid/debrid.go
     */
    interface Debrid_CachedFile {
        size: number;
        name: string;
    }

    /**
     * - Filepath: internal/debrid/debrid/debrid.go
     */
    interface Debrid_TorrentItemInstantAvailability {
        /**
         * Key is the file ID (or index)
         */
        cachedFiles?: Record<string, Debrid_CachedFile>;
    }

    /**
     * - Filepath: internal/discordrpc/presence/presence.go
     */
    interface DiscordRPC_AnimeActivity {
        id: number;
        title: string;
        image: string;
        isMovie: boolean;
        episodeNumber: number;
        paused: boolean;
        progress: number;
        duration: number;
        totalEpisodes?: number;
        currentEpisodeCount?: number;
        episodeTitle?: string;
    }

    /**
     * - Filepath: internal/discordrpc/client/activity.go
     */
    interface DiscordRPC_Button {
        label?: string;
        url?: string;
    }

    /**
     * - Filepath: internal/discordrpc/presence/presence.go
     */
    interface DiscordRPC_LegacyAnimeActivity {
        id: number;
        title: string;
        image: string;
        isMovie: boolean;
        episodeNumber: number;
    }

    /**
     * - Filepath: internal/discordrpc/presence/presence.go
     */
    interface DiscordRPC_MangaActivity {
        id: number;
        title: string;
        image: string;
        chapter: string;
    }

    /**
     * - Filepath: internal/extension/hibike/manga/types.go
     */
    interface HibikeManga_ChapterDetails {
        provider: string;
        id: string;
        url: string;
        title: string;
        chapter: string;
        index: number;
        scanlator?: string;
        language?: string;
        rating?: number;
        updatedAt?: string;
        localIsPDF?: boolean;
    }

    /**
     * - Filepath: internal/extension/hibike/torrent/types.go
     */
    interface HibikeTorrent_MediaTorrent {
        provider?: string;
        name: string;
        date: string;
        size: number;
        formattedSize: string;
        seeders: number;
        leechers: number;
        downloadCount: number;
        link: string;
        downloadUrl: string;
        magnetLink?: string;
        infoHash?: string;
        resolution?: string;
        isBatch?: boolean;
        episodeNumber?: number;
        releaseGroup?: string;
        isBestRelease: boolean;
        confirmed: boolean;
    }

    /**
     * - Filepath: internal/manga/chapter_container.go
     */
    interface Manga_ChapterContainer {
        mediaId: number;
        provider: string;
        chapters?: Array<HibikeManga_ChapterDetails>;
    }

    /**
     * - Filepath: internal/manga/collection.go
     */
    interface Manga_Collection {
        lists?: Array<Manga_CollectionList>;
    }

    /**
     * - Filepath: internal/manga/collection.go
     */
    interface Manga_CollectionEntry {
        media?: MediaAPI_BaseManga;
        mediaId: number;
        /**
         * SIMKL list data
         */
        listData?: Manga_EntryListData;
    }

    /**
     * - Filepath: internal/manga/collection.go
     */
    interface Manga_CollectionList {
        type?: MediaAPI_MediaListStatus;
        status?: MediaAPI_MediaListStatus;
        entries?: Array<Manga_CollectionEntry>;
    }

    /**
     * - Filepath: internal/manga/manga_entry.go
     */
    interface Manga_Entry {
        mediaId: number;
        media?: MediaAPI_BaseManga;
        listData?: Manga_EntryListData;
    }

    /**
     * - Filepath: internal/manga/manga_entry.go
     */
    interface Manga_EntryListData {
        progress?: number;
        score?: number;
        status?: MediaAPI_MediaListStatus;
        repeat?: number;
        startedAt?: string;
        completedAt?: string;
    }

    /**
     * - Filepath: internal/manga/chapter_container.go
     */
    interface Manga_MangaLatestChapterNumberItem {
        provider: string;
        scanlator: string;
        language: string;
        number: number;
    }

    /**
     * - Filepath: internal/manga/download.go
     */
    export type Manga_MediaMap = Record<number, Manga_ProviderDownloadMap>;

    /**
     * - Filepath: internal/manga/download.go
     */
    export type Manga_ProviderDownloadMap = Record<string, Array<Manga_ProviderDownloadMapChapterInfo>>;

    /**
     * - Filepath: internal/manga/download.go
     */
    interface Manga_ProviderDownloadMapChapterInfo {
        chapterId: string;
        chapterNumber: string;
    }

    /**
     * - Filepath: internal/api/metadata/types.go
     */
    interface Metadata_AnimeMappings {
        animeplanetId?: string;
        kitsuId?: number;
        malId?: number;
        type?: string;
        simklId?: number;
        anisearchId?: number;
        anidbId?: number;
        notifymoeId?: string;
        livechartId?: number;
        thetvdbId?: number;
        imdbId?: string;
        themoviedbId?: string;
    }

    /**
     * - Filepath: internal/api/metadata/types.go
     */
    interface Metadata_AnimeMetadata {
        titles?: Record<string, string>;
        episodes?: Record<string, Metadata_EpisodeMetadata>;
        episodeCount: number;
        specialCount: number;
        mappings?: Metadata_AnimeMappings;
    }

    /**
     * - Filepath: internal/api/metadata/types.go
     */
    interface Metadata_EpisodeMetadata {
        anidbId: number;
        tvdbId: number;
        title: string;
        image: string;
        airDate: string;
        length: number;
        summary: string;
        overview: string;
        episodeNumber: number;
        episode: string;
        seasonNumber: number;
        absoluteEpisodeNumber: number;
        anidbEid: number;
        /**
         * Indicates if the episode has a real image
         */
        hasImage: boolean;
    }

    /**
     * - Filepath: internal/database/models/models.go
     */
    interface Models_AutoDownloaderItem {
        ruleId: number;
        mediaId: number;
        episode: number;
        link: string;
        hash: string;
        magnet: string;
        torrentName: string;
        downloaded: boolean;
        isDelayed: boolean;
        delayUntil?: string;
        score: number;
        id: number;
        createdAt?: string;
        updatedAt?: string;
    }

    /**
     * - Filepath: internal/database/models/models.go
     */
    interface Models_AutoDownloaderSettings {
        provider: string;
        interval: number;
        enabled: boolean;
        downloadAutomatically: boolean;
        enableEnhancedQueries: boolean;
        enableSeasonCheck: boolean;
        useDebrid: boolean;
    }

    /**
     * - Filepath: internal/onlinestream/repository.go
     */
    interface Onlinestream_Episode {
        number: number;
        title?: string;
        image?: string;
        description?: string;
        isFiller?: boolean;
        metadata?: Media_Episode;
    }

    /**
     * - Filepath: internal/torrent_clients/torrent_client/torrent.go
     */
    interface TorrentClient_Torrent {
        name: string;
        hash: string;
        seeds: number;
        upSpeed: string;
        downSpeed: string;
        progress: number;
        size: string;
        eta: string;
        status: TorrentClient_TorrentStatus;
        contentPath: string;
    }

    /**
     * - Filepath: internal/torrent_clients/torrent_client/torrent.go
     */
    export type TorrentClient_TorrentStatus = "downloading" | "seeding" | "paused" | "other" | "stopped";

    /**
     * - Filepath: internal/torrents/torrent/search.go
     */
    interface Torrent_AnimeSearchOptions {
        provider: string;
        type?: Torrent_AnimeSearchType;
        media?: MediaAPI_BaseMedia;
        query?: string;
        batch?: boolean;
        episodeNumber?: number;
        bestReleases?: boolean;
        resolution?: string;
        includeSpecialProviders?: boolean;
        skipPreviews?: boolean;
    }

    /**
     * - Filepath: internal/torrents/torrent/search.go
     */
    export type Torrent_AnimeSearchType = "smart" | "simple";

    /**
     * - Filepath: internal/torrents/torrent/search.go
     */
    interface Torrent_Preview {
        /**
         * nil if batch
         */
        episode?: Media_Episode;
        torrent?: HibikeTorrent_MediaTorrent;
    }

    /**
     * - Filepath: internal/torrents/torrent/search.go
     */
    interface Torrent_SearchData {
        /**
         * Torrents found
         */
        torrents?: Array<HibikeTorrent_MediaTorrent>;
        /**
         * TorrentPreview for each torrent
         */
        previews?: Array<Torrent_Preview>;
        /**
         * Torrent metadata
         */
        torrentMetadata?: Record<string, Torrent_TorrentMetadata>;
        /**
         * Debrid instant availability
         */
        debridInstantAvailability?: Record<string, Debrid_TorrentItemInstantAvailability>;
        /**
         * Animap media
         */
        animeMetadata?: Metadata_AnimeMetadata;
        includedSpecialProviders?: Array<string>;
    }

    /**
     * - Filepath: internal/torrents/torrent/search.go
     */
    interface Torrent_TorrentMetadata {
        distance: number;
        metadata?: $habari.Metadata;
    }

}
