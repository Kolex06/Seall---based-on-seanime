/// <reference path="./plugin.d.ts" />
/// <reference path="./system.d.ts" />
/// <reference path="./app.d.ts" />
/// <reference path="./core.d.ts" />

interface AmaSettings {
    betterMarketplace: boolean
    carousels: boolean
    carouselsSearch: boolean
    carouselsExtensions: boolean
    carouselsLists: boolean
    carouselsManga: boolean
    carouselsOther: boolean
    subDubIcons: boolean
    hideFileNames: boolean
    useBas1874Marketplace: boolean
}

function init() {
    $ui.register(function(ctx) {

        const SETTINGS_KEY = "ama-ui-tweaks.settings"

        const DEFAULT_SETTINGS: AmaSettings = {
            betterMarketplace: true,
            carousels: true,
            carouselsSearch: false,
            carouselsExtensions: true,
            carouselsLists: true,
            carouselsManga: true,
            carouselsOther: true,
            subDubIcons: true,
            hideFileNames: false,
            useBas1874Marketplace: false,
        }

        function getStorageApi(): any | null {
            try {
                // @ts-ignore
                if (typeof $storage !== "undefined" && $storage) {
                    // @ts-ignore
                    return $storage
                }
            } catch (_) {}

            return null
        }

        function normalizeSettings(value: Partial<AmaSettings> | null | undefined): AmaSettings {
            const saved = value || {}
            const betterMarketplace = saved.betterMarketplace !== false

            return {
                betterMarketplace,
                carousels: saved.carousels !== false,
                carouselsSearch: saved.carouselsSearch === true,
                carouselsExtensions: saved.carouselsExtensions !== false,
                carouselsLists: saved.carouselsLists !== false,
                carouselsManga: saved.carouselsManga !== false,
                carouselsOther: saved.carouselsOther !== false,
                subDubIcons: saved.subDubIcons !== false,
                hideFileNames: saved.hideFileNames === true,
                useBas1874Marketplace: saved.useBas1874Marketplace === true,
            }
        }

        function loadSettings(): AmaSettings {
            try {
                const storage = getStorageApi()

                if (storage && typeof storage.get === "function") {
                    return normalizeSettings(storage.get<Partial<AmaSettings>>(SETTINGS_KEY))
                }
            } catch (_) {}

            return normalizeSettings(null)
        }

        function saveSettings(settings: AmaSettings) {
            try {
                const storage = getStorageApi()

                if (storage && typeof storage.set === "function") {
                    storage.set(SETTINGS_KEY, settings)
                }
            } catch (_) {}
        }

        const settingsState = ctx.state<AmaSettings>(loadSettings())

        const betterMarketplaceRef = ctx.fieldRef<boolean>(settingsState.get().betterMarketplace)
        const carouselsRef = ctx.fieldRef<boolean>(settingsState.get().carousels)
        const carouselsSearchRef = ctx.fieldRef<boolean>(settingsState.get().carouselsSearch)
        const carouselsExtensionsRef = ctx.fieldRef<boolean>(settingsState.get().carouselsExtensions)
        const carouselsListsRef = ctx.fieldRef<boolean>(settingsState.get().carouselsLists)
        const carouselsMangaRef = ctx.fieldRef<boolean>(settingsState.get().carouselsManga)
        const carouselsOtherRef = ctx.fieldRef<boolean>(settingsState.get().carouselsOther)
        const subDubIconsRef = ctx.fieldRef<boolean>(settingsState.get().subDubIcons)
        const hideFileNamesRef = ctx.fieldRef<boolean>(settingsState.get().hideFileNames)
        const useBas1874MarketplaceRef = ctx.fieldRef<boolean>(settingsState.get().useBas1874Marketplace)

        const trayIconUrl = "https://raw.githubusercontent.com/Kolex06/Seanime-Stuff/main/icons/SeaUtils-Kolex06-Version.png"

        const tray = ctx.newTray({
            iconUrl: trayIconUrl,
            withContent: true,
            width: "320px",
            minHeight: "190px",
        })

        async function pushSettingsToClient(settings: AmaSettings) {
            const body = await ctx.dom.queryOne("body")
            if (!body) return

            const bridge = await ctx.dom.createElement("script")
            bridge.setText(
                "window.__AMA_SAVE_SETTINGS__ && window.__AMA_SAVE_SETTINGS__(" +
                JSON.stringify(settings) +
                ");"
            )

            body.append(bridge)
        }

        function saveAndApplySettings(next: AmaSettings) {
            saveSettings(next)
            settingsState.set(next)
            pushSettingsToClient(next)
        }

        function setFieldRefValue<T>(fieldRef: any, value: T) {
            try {
                if (fieldRef && typeof fieldRef.set === "function") {
                    fieldRef.set(value)
                }
            } catch (_) {}
        }

        function updateSetting<K extends keyof AmaSettings>(key: K, value: AmaSettings[K]) {
            const current = settingsState.get()
            const next: AmaSettings = {
                ...current,
                [key]: value,
            }

            saveAndApplySettings(next)
        }

        betterMarketplaceRef.onValueChange((value) => {
            updateSetting("betterMarketplace", !!value)
        })

        carouselsRef.onValueChange((value) => {
            updateSetting("carousels", !!value)
        })

        carouselsSearchRef.onValueChange((value) => {
            updateSetting("carouselsSearch", !!value)
        })

        carouselsExtensionsRef.onValueChange((value) => {
            updateSetting("carouselsExtensions", !!value)
        })

        carouselsListsRef.onValueChange((value) => {
            updateSetting("carouselsLists", !!value)
        })

        carouselsMangaRef.onValueChange((value) => {
            updateSetting("carouselsManga", !!value)
        })

        carouselsOtherRef.onValueChange((value) => {
            updateSetting("carouselsOther", !!value)
        })

        subDubIconsRef.onValueChange((value) => {
            updateSetting("subDubIcons", !!value)
        })

        hideFileNamesRef.onValueChange((value) => {
            updateSetting("hideFileNames", !!value)
        })

        useBas1874MarketplaceRef.onValueChange((value) => {
            updateSetting("useBas1874Marketplace", !!value)
        })

        tray.render(() => {
            const settings = settingsState.get()
            const disabledValueText = (label: string, value: boolean, reason: string) => {
                return tray.text(label + ": " + (value ? "On" : "Off") + " (disabled - " + reason + ")")
            }

            const carouselSearchControl = settings.carousels
                ? tray.switch("Carousels: Search", {
                    fieldRef: carouselsSearchRef,
                })
                : disabledValueText("Carousels: Search", settings.carouselsSearch, "turn on Carousels")

            const carouselExtensionsControl = settings.carousels && settings.betterMarketplace
                ? tray.switch("Carousels: Extensions", {
                    fieldRef: carouselsExtensionsRef,
                })
                : disabledValueText(
                    "Carousels: Extensions",
                    settings.carouselsExtensions,
                    settings.carousels ? "turn on Better Marketplace" : "turn on Carousels"
                )

            const carouselListsControl = settings.carousels
                ? tray.switch("Carousels: My Lists", {
                    fieldRef: carouselsListsRef,
                })
                : disabledValueText("Carousels: My Lists", settings.carouselsLists, "turn on Carousels")

            const carouselMangaControl = settings.carousels
                ? tray.switch("Carousels: Manga", {
                    fieldRef: carouselsMangaRef,
                })
                : disabledValueText("Carousels: Manga", settings.carouselsManga, "turn on Carousels")

            const carouselOtherControl = settings.carousels
                ? tray.switch("Carousels: Other Pages", {
                    fieldRef: carouselsOtherRef,
                })
                : disabledValueText("Carousels: Other Pages", settings.carouselsOther, "turn on Carousels")

            return tray.stack([
                tray.text("SeaUtils Kolex06-Version Settings"),
                tray.switch("Better Marketplace", {
                    fieldRef: betterMarketplaceRef,
                }),
                tray.switch("Carousels", {
                    fieldRef: carouselsRef,
                }),
                carouselSearchControl,
                carouselExtensionsControl,
                carouselListsControl,
                carouselMangaControl,
                carouselOtherControl,
                tray.switch("Sub/Dub Icons", {
                    fieldRef: subDubIconsRef,
                }),
                tray.switch("Hide File Names", {
                    fieldRef: hideFileNamesRef,
                }),
                tray.switch("Bas1874 Marketplace", {
                    fieldRef: useBas1874MarketplaceRef,
                }),
            ])
        })

        const initialFeatureSettings = settingsState.get()

        const SCHEDULE_EVENT_QUERY = '[data-schedule-calendar-event-item-link], [data-schedule-calendar-event-item-content], [data-schedule-calendar-event-item-name], [data-schedule-calendar-event-item-episode], [data-schedule-calendar-event-item-finale-icon], [data-schedule-calendar-event-item], [data-schedule-event-item], [data-schedule-media-id]'
        const SERVER_SCHEDULE_EVENT_QUERY = '[data-schedule-calendar-event-item-link], [data-schedule-calendar-event-item-content], [data-schedule-calendar-event-item-name], [data-schedule-calendar-event-item-episode], [data-schedule-calendar-event-item-finale-icon], [data-schedule-calendar-event-item], [data-schedule-calendar-event-item-root], [data-schedule-event-item], [data-schedule-media-id]'
        const SCHEDULE_SETTINGS_BUTTON_QUERY = '[data-schedule-calendar-header-button-settings="true"]'
        const SCHEDULE_TOKEN_STORAGE_KEY = 'ama-ui-tweaks.animeScheduleApiToken'
        const ANIME_SCHEDULE_API_BASE_URL = 'https://animeschedule.net/api/v3'

        function normalizeServerAnimeScheduleToken(value: any): string {
            return String(value || "")
                .trim()
                .replace(/^["']|["']$/g, "")
                .replace(/^authorization\s*:\s*/i, "")
                .replace(/^bearer\s+/i, "")
                .trim()
        }

        function getServerStorageString(key: string): string {
            try {
                const storage = getStorageApi()
                if (storage && typeof storage.get === "function") {
                    const value = storage.get<string>(key)
                    return normalizeServerAnimeScheduleToken(value)
                }
            } catch (_) {}

            return ""
        }

        function setServerStorageString(key: string, value: string) {
            try {
                const storage = getStorageApi()
                if (storage && typeof storage.set === "function") {
                    storage.set(key, normalizeServerAnimeScheduleToken(value))
                }
            } catch (_) {}
        }

        let serverAnimeScheduleApiToken = getServerStorageString(SCHEDULE_TOKEN_STORAGE_KEY)
        const serverDubTimetableCache: Record<string, any[] | null> = {}
        const serverDubTimetablePromises: Record<string, Promise<any[] | null>> = {}
        let serverScheduleRefreshRunning = false

        async function setServerScheduleApiStatus(status: string, detail: string) {
            try {
                const body = await ctx.dom.queryOne("body")
                if (!body) return
                body.setAttribute("data-ama-server-schedule-api-status", status)
                body.setAttribute("data-ama-server-schedule-api-detail", detail || "")
            } catch (_) {}
        }

        function normalizeServerId(value: any): string {
            if (value === null || value === undefined) return ""
            const text = String(value).trim()
            if (!text) return ""
            const numeric = text.match(/\d+/)
            return numeric ? numeric[0] : text
        }

        function serverGetAnilistIdFromHref(href: string): string {
            if (!href) return ""

            try {
                const parsed = new URL(href, "http://localhost")
                const names = ["anilistId", "anilist_id", "aniListId", "mediaId", "media_id", "animeId", "anime_id", "id"]
                for (const name of names) {
                    const id = normalizeServerId(parsed.searchParams.get(name))
                    if (id) return id
                }

                const path = parsed.pathname || ""
                const match = path.match(/(?:\/anime\/|\/entry\/|\/media\/|\/details\/)(?:anime\/)?(\d+)(?:\/|$)/i)
                return match && match[1] ? match[1] : ""
            } catch (_) {
                const match = String(href).match(/[?&#](?:anilistId|anilist_id|aniListId|mediaId|media_id|animeId|anime_id|id)=(\d+)/i)
                return match && match[1] ? match[1] : ""
            }
        }

        function normalizeServerEpisodeNumber(value: any): string {
            const text = String(value === null || value === undefined ? "" : value).trim()
            if (!text) return ""
            const direct = text.match(/^\d+(?:\.\d+)?$/)
            if (direct) return String(Number(direct[0])).replace(/\.0$/, "")
            const match = text.match(/(?:episode|ep\.?|e)\s*#?\s*(\d+(?:\.\d+)?)/i)
            return match && match[1] ? String(Number(match[1])).replace(/\.0$/, "") : ""
        }

        function serverParseDate(value: any): Date | null {
            const text = String(value || "").trim()
            if (!text || text.indexOf("0001-01-01") === 0) return null
            const date = new Date(text)
            return Number.isNaN(date.getTime()) ? null : date
        }

        function serverDatesSameLocalDay(a: Date | null, b: Date | null): boolean {
            if (!a || !b) return false
            return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate()
        }

        function canUseServerScheduleFallback(scheduleDate: Date | null): boolean {
            if (!scheduleDate || Number.isNaN(scheduleDate.getTime())) return false

            const currentWeekStart = new Date()
            currentWeekStart.setHours(0, 0, 0, 0)
            const dayOffset = (currentWeekStart.getDay() + 6) % 7
            currentWeekStart.setDate(currentWeekStart.getDate() - dayOffset)

            const scheduleDay = new Date(scheduleDate.getTime())
            scheduleDay.setHours(0, 0, 0, 0)

            return scheduleDay.getTime() < currentWeekStart.getTime()
        }

        function serverIsoWeekInfo(date: Date | null): { year: number, week: number } {
            const source = date && !Number.isNaN(date.getTime()) ? date : new Date()
            const utc = new Date(Date.UTC(source.getUTCFullYear(), source.getUTCMonth(), source.getUTCDate()))
            const day = utc.getUTCDay() || 7
            utc.setUTCDate(utc.getUTCDate() + 4 - day)
            const yearStart = new Date(Date.UTC(utc.getUTCFullYear(), 0, 1))
            const week = Math.ceil((((utc.getTime() - yearStart.getTime()) / 86400000) + 1) / 7)
            return { year: utc.getUTCFullYear(), week }
        }

        function serverFlattenObjects(value: any, result: any[], depth: number) {
            if (!value || depth > 8) return
            if (Array.isArray(value)) {
                value.forEach(item => serverFlattenObjects(item, result, depth + 1))
                return
            }
            if (typeof value !== "object") return
            result.push(value)
            Object.keys(value).forEach(key => {
                const child = value[key]
                if (child && typeof child === "object") serverFlattenObjects(child, result, depth + 1)
            })
        }

        function serverEntryAnilistId(entry: any): string {
            if (!entry || typeof entry !== "object") return ""
            const numericId = normalizeServerId(entry.id)
            if (numericId && (entry.idMal || entry.episode || entry.format || entry.duration)) return numericId
            const direct = normalizeServerId(entry.anilistId || entry.anilistID || entry.aniListId || entry.anilist_id || entry.mediaId || entry.media_id || entry.idAniList)
            if (direct) return direct
            const websites = entry.websites || {}
            return serverGetAnilistIdFromHref(websites.aniList || websites.anilist || entry.aniList || entry.anilist || "")
        }

        function serverEntryEpisodeNumber(entry: any): string {
            if (!entry || typeof entry !== "object") return ""
            const values = [
                entry.episodeNumber,
                entry.episode,
                entry.episodeNum,
                entry.number,
                entry.dubEpisode,
                entry.dubEpisodeNumber,
                entry.latestDubEpisode,
                entry.latestDubEpisodeNumber,
                entry.dubbedEpisode,
                entry.currentEpisode,
                entry.airingEpisode,
                entry.episode && entry.episode.aired,
                entry.episode && entry.episode.number,
                entry.episode && entry.episode.episodeNumber,
            ]

            for (const value of values) {
                const episode = normalizeServerEpisodeNumber(value)
                if (episode) return episode
            }

            return ""
        }

        function serverEntryDate(entry: any): Date | null {
            if (!entry || typeof entry !== "object") return null
            return serverParseDate(entry.episodeDate || entry.date || entry.datetime || entry.airingAt || entry.airedAt || (entry.episode && (entry.episode.airedAt || entry.episode.date || entry.episode.datetime)))
        }

        function serverNormalizeTimetableEntries(data: any): any[] {
            const objects: any[] = []
            serverFlattenObjects(data, objects, 0)
            return objects.filter(entry => serverEntryEpisodeNumber(entry) && serverEntryAnilistId(entry))
        }

        async function loadServerDubTimetableEntries(scheduleDate: Date | null): Promise<any[] | null> {
            const token = normalizeServerAnimeScheduleToken(serverAnimeScheduleApiToken || await syncServerAnimeScheduleTokenFromDom())
            if (!token) {
                await setServerScheduleApiStatus("missing-token", "")
                return null
            }

            const timezone = (() => {
                try {
                    return Intl.DateTimeFormat().resolvedOptions().timeZone || "Europe/Oslo"
                } catch (_) {
                    return "Europe/Oslo"
                }
            })()
            const weekInfo = serverIsoWeekInfo(scheduleDate)
            const cacheKey = weekInfo.year + "-w" + weekInfo.week + "-" + timezone + "-" + (token ? "token" : "public")
            if (Object.prototype.hasOwnProperty.call(serverDubTimetableCache, cacheKey)) return serverDubTimetableCache[cacheKey]
            if (Object.prototype.hasOwnProperty.call(serverDubTimetablePromises, cacheKey)) return serverDubTimetablePromises[cacheKey]

            serverDubTimetablePromises[cacheKey] = (async () => {
                const headers: Record<string, string> = {
                    Accept: "application/json",
                }
                if (token) headers.Authorization = "Bearer " + token

                try {
                    const response = await ctx.fetch(ANIME_SCHEDULE_API_BASE_URL + "/timetables/dub?year=" + encodeURIComponent(String(weekInfo.year)) + "&week=" + encodeURIComponent(String(weekInfo.week)) + "&tz=" + encodeURIComponent(timezone), {
                        headers,
                        noCloudflareBypass: true,
                        timeout: 12,
                    })

                    if (!response.ok) {
                        await setServerScheduleApiStatus("http-" + String(response.status), response.statusText || "")
                        serverDubTimetableCache[cacheKey] = null
                        return null
                    }

                    const entries = serverNormalizeTimetableEntries(response.json())
                    await setServerScheduleApiStatus("loaded", String(entries.length))
                    serverDubTimetableCache[cacheKey] = entries
                    return entries
                } catch (_) {
                    await setServerScheduleApiStatus("failed", "")
                    serverDubTimetableCache[cacheKey] = null
                    return null
                } finally {
                    delete serverDubTimetablePromises[cacheKey]
                }
            })()

            return serverDubTimetablePromises[cacheKey]
        }

        async function syncServerAnimeScheduleTokenFromDom(): Promise<string> {
            try {
                const body = await ctx.dom.queryOne("body")
                const bodyToken = body ? normalizeServerAnimeScheduleToken(await body.getAttribute("data-ama-anime-schedule-api-token")) : ""
                if (bodyToken) {
                    serverAnimeScheduleApiToken = bodyToken
                    setServerStorageString(SCHEDULE_TOKEN_STORAGE_KEY, bodyToken)
                    return bodyToken
                }
            } catch (_) {}

            try {
                const input = await ctx.dom.queryOne('[data-ama-schedule-token-input="true"], [data-ama-schedule-token-field="true"] input')
                const inputToken = input ? normalizeServerAnimeScheduleToken(await input.getProperty("value")) : ""
                if (inputToken) {
                    serverAnimeScheduleApiToken = inputToken
                    setServerStorageString(SCHEDULE_TOKEN_STORAGE_KEY, inputToken)
                    return inputToken
                }
            } catch (_) {}

            return serverAnimeScheduleApiToken
        }

        async function getServerScheduleEventRoot(event: $ui.DOMElement): Promise<$ui.DOMElement> {
            let current: $ui.DOMElement | null = event

            for (let i = 0; current && i < 7; i++) {
                const href = await current.getAttribute("href")
                const isLink = !!href && (href.includes("/entry") || href.includes("id="))
                const hasScheduleLink = await current.hasDataAttribute("schedule-calendar-event-item-link")
                const hasScheduleItem = await current.hasDataAttribute("schedule-calendar-event-item")
                const time = await current.queryOne("time[datetime], [datetime]")
                const name = await current.queryOne('[data-schedule-calendar-event-item-name], [data-schedule-calendar-event-item-title], [data-schedule-event-title], [data-media-title]')

                if (isLink || hasScheduleLink || hasScheduleItem || (time && name)) return current

                current = await current.getParent({ identifyChildren: true })
            }

            return event
        }

        async function getServerScheduleEventInfo(event: $ui.DOMElement): Promise<{ root: $ui.DOMElement, anilistId: string, title: string, scheduleDate: Date | null }> {
            const root = await getServerScheduleEventRoot(event)
            let anilistId = normalizeServerId(await root.getDataAttribute("amaDubAnilistId"))
            const link = await root.queryOne("a[href]")
            const eventHref = await root.getAttribute("href")
            const linkHref = link ? await link.getAttribute("href") : ""
            anilistId = anilistId || serverGetAnilistIdFromHref(eventHref || "") || serverGetAnilistIdFromHref(linkHref || "")

            const name = await root.queryOne('[data-schedule-calendar-event-item-name], [data-schedule-calendar-event-item-title], [data-schedule-event-title], [data-media-title]')
            const title = name ? String(await name.getText() || "").trim() : String(await root.getText() || "").trim()

            const time = await root.queryOne("time[datetime], [datetime]")
            const scheduleDate = serverParseDate(time ? await time.getAttribute("datetime") : "")

            return { root, anilistId, title, scheduleDate }
        }

        async function addServerScheduleDubBadge(event: $ui.DOMElement, details: { episodeNumber: string, episodeDate: string }) {
            const existing = await event.queryOne(".ama-schedule-dub-badge-server")
            if (existing) {
                existing.setText("DUB Ep. " + details.episodeNumber)
                existing.setAttribute("title", "Dub Ep. " + details.episodeNumber + (details.episodeDate ? " - " + details.episodeDate : ""))
                return
            }

            const anchor = await event.queryOne('[data-schedule-calendar-event-item-name], [data-schedule-calendar-event-item-text], [data-schedule-calendar-event-item-title], [data-schedule-event-title], [data-media-title], a[href]')
            if (!anchor) return

            const badge = await ctx.dom.createElement("span")
            badge.setText("DUB Ep. " + details.episodeNumber)
            badge.setAttribute("class", "ama-media-badge dub ama-schedule-dub-badge-server")
            badge.setAttribute("title", "Dub Ep. " + details.episodeNumber + (details.episodeDate ? " - " + details.episodeDate : ""))
            badge.setAttribute("aria-label", "Dub Ep. " + details.episodeNumber)
            badge.setCssText("display:inline-flex;align-items:center;justify-content:center;height:20px;min-width:32px;padding:0 6px;margin-left:6px;border-radius:999px;border:1px solid rgba(255,183,197,.7);color:#ffb7c5;background:rgba(255,183,197,.13);font-size:11px;font-weight:800;line-height:1;vertical-align:middle;")
            anchor.after(badge)
        }

        async function removeServerScheduleDubBadge(event: $ui.DOMElement) {
            const existing = await event.queryOne(".ama-schedule-dub-badge-server")
            if (existing) existing.remove()
        }

        async function enhanceServerScheduleEvent(event: $ui.DOMElement): Promise<string> {
            const info = await getServerScheduleEventInfo(event)
            if (!info.anilistId || !info.scheduleDate) {
                return "missing-info:id=" + (info.anilistId ? "yes" : "no") + ";date=" + (info.scheduleDate ? "yes" : "no")
            }

            if (canUseServerScheduleFallback(info.scheduleDate)) {
                return "past-fallback"
            }

            const entries = await loadServerDubTimetableEntries(info.scheduleDate)
            if (!Array.isArray(entries)) return "api-unavailable"

            const match = entries.find(entry => {
                if (String(serverEntryAnilistId(entry)) !== String(info.anilistId)) return false
                const entryDate = serverEntryDate(entry)
                return serverDatesSameLocalDay(entryDate, info.scheduleDate)
            })

            if (!match) {
                await removeServerScheduleDubBadge(info.root)
                return "no-match"
            }

            await addServerScheduleDubBadge(info.root, {
                episodeNumber: serverEntryEpisodeNumber(match),
                episodeDate: String((serverEntryDate(match) || {}).toISOString ? serverEntryDate(match)!.toISOString() : ""),
            })
            return "matched"
        }

        async function attachServerScheduleTokenInput(dialog: $ui.DOMElement) {
            let input = await dialog.queryOne('[data-ama-schedule-token-input="true"]')
            if (!input) {
                input = await dialog.queryOne('[data-ama-schedule-token-field="true"] input')
                if (input) input.setAttribute("data-ama-schedule-token-input", "true")
            }

            if (!input) {
                return
            } else {
                const existingValue = normalizeServerAnimeScheduleToken(await input.getProperty("value"))
                if (!serverAnimeScheduleApiToken && existingValue) {
                    serverAnimeScheduleApiToken = existingValue
                    setServerStorageString(SCHEDULE_TOKEN_STORAGE_KEY, existingValue)
                }
                input.setProperty("value", serverAnimeScheduleApiToken)
            }

            input.addEventListener("input", async () => {
                const value = normalizeServerAnimeScheduleToken(await input!.getProperty("value"))
                serverAnimeScheduleApiToken = value
                setServerStorageString(SCHEDULE_TOKEN_STORAGE_KEY, value)
                Object.keys(serverDubTimetableCache).forEach(key => delete serverDubTimetableCache[key])
                Object.keys(serverDubTimetablePromises).forEach(key => delete serverDubTimetablePromises[key])
            })
            input.addEventListener("change", async () => {
                const value = normalizeServerAnimeScheduleToken(await input!.getProperty("value"))
                serverAnimeScheduleApiToken = value
                setServerStorageString(SCHEDULE_TOKEN_STORAGE_KEY, value)
                Object.keys(serverDubTimetableCache).forEach(key => delete serverDubTimetableCache[key])
                Object.keys(serverDubTimetablePromises).forEach(key => delete serverDubTimetablePromises[key])
            })
        }

        async function refreshServerScheduleDubBadges() {
            if (serverScheduleRefreshRunning) return
            serverScheduleRefreshRunning = true

            try {
                await setServerScheduleApiStatus("scanning", "")
                await syncServerAnimeScheduleTokenFromDom()
                const events = await ctx.dom.query(SERVER_SCHEDULE_EVENT_QUERY, { identifyChildren: true })
                if (!events.length) {
                    await setServerScheduleApiStatus("no-events", "")
                    return
                }

                const candidates = await Promise.all(events.slice(0, 240).map(async event => {
                    try {
                        const info = await getServerScheduleEventInfo(event)
                        return { event: info.root || event, scheduleDate: info.scheduleDate }
                    } catch (_) {
                        return { event, scheduleDate: null }
                    }
                }))

                const seenCandidateIds: Record<string, boolean> = {}
                const uniqueCandidates = candidates.filter(candidate => {
                    const key = candidate.event && candidate.event.id ? candidate.event.id : ""
                    if (!key) return true
                    if (seenCandidateIds[key]) return false
                    seenCandidateIds[key] = true
                    return true
                })

                uniqueCandidates.sort((a, b) => {
                    const aPast = canUseServerScheduleFallback(a.scheduleDate) ? 1 : 0
                    const bPast = canUseServerScheduleFallback(b.scheduleDate) ? 1 : 0
                    return aPast - bPast
                })

                const results = await Promise.all(uniqueCandidates.slice(0, 100).map(candidate => {
                    return enhanceServerScheduleEvent(candidate.event).catch(() => "event-failed")
                }))

                const counts: Record<string, number> = {}
                results.forEach(result => {
                    const key = String(result || "unknown").split(":")[0]
                    counts[key] = (counts[key] || 0) + 1
                })

                if (counts.matched) {
                    await setServerScheduleApiStatus("matched", JSON.stringify(counts))
                } else if (counts["no-match"]) {
                    await setServerScheduleApiStatus("loaded-no-match", JSON.stringify(counts))
                } else if (counts["api-unavailable"]) {
                    await setServerScheduleApiStatus("api-unavailable", JSON.stringify(counts))
                } else if (counts["past-fallback"]) {
                    await setServerScheduleApiStatus("past-fallback", JSON.stringify(counts))
                } else if (counts["missing-info"]) {
                    await setServerScheduleApiStatus("missing-info", JSON.stringify(counts))
                } else {
                    await setServerScheduleApiStatus("scanned", JSON.stringify(counts))
                }
            } catch (_) {
                await setServerScheduleApiStatus("scan-failed", "")
            } finally {
                serverScheduleRefreshRunning = false
            }
        }

        // Disabled while the Schedule API integration is reworked; the browser-side
        // Schedule fallback remains active and should not throw during Seall startup.

        const carouselCSS = `
            body[data-ama-hide-file-names="true"] [data-episode-grid-item-filename="true"] {
                display: none !important;
            }

            body[data-ama-better-marketplace="true"][data-ama-carousels-active="true"] .group\\/extension-card .absolute.top-3.right-3.flex.flex-col.gap-1 {
                flex-direction: row !important;
                align-items: center !important;
                justify-content: center !important;
                gap: 4px !important;
                padding: 4px !important;
                min-width: 0 !important;
                width: auto !important;
                max-width: max-content !important;
            }

            body[data-ama-better-marketplace="true"][data-ama-carousels-active="true"] .group\\/extension-card .absolute.top-3.right-3.flex.flex-col.gap-1 > button:not(.UI-IconButton_root) {
                width: 32px !important;
                min-width: 32px !important;
                max-width: 32px !important;
                height: 32px !important;
                padding: 0 !important;
                flex: 0 0 32px !important;
            }

            body[data-ama-better-marketplace="true"][data-ama-carousels-active="true"] .group\\/extension-card .absolute.top-3.right-3.flex.flex-col.gap-1 > button:not(.UI-IconButton_root) > .UI-Button__icon {
                margin-inline-end: 0 !important;
            }

            body[data-ama-better-marketplace="true"][data-ama-carousels-active="true"] .group\\/extension-card .absolute.top-3.right-3.flex.flex-col.gap-1 > button:not(.UI-IconButton_root) > span:not(.UI-Button__icon) {
                display: none !important;
                opacity: 0 !important;
                visibility: hidden !important;
                width: 0 !important;
                max-width: 0 !important;
                overflow: hidden !important;
            }

            .ama-carousel-nav-btn,
            button.ama-carousel-nav-btn,
            .ama-manga-carousel-parent > .ama-carousel-nav-btn,
            svg[data-ama-random-search-icon="true"] {
                display: none !important;
                opacity: 0 !important;
                visibility: hidden !important;
                pointer-events: none !important;
            }

            *::-webkit-scrollbar-button {
                display: none !important;
                width: 0 !important;
                height: 0 !important;
            }

            body[data-ama-better-marketplace="true"][data-ama-better-marketplace-page="true"] input[placeholder*="Search"]:not(.ama-search-input),
            body[data-ama-better-marketplace="true"][data-ama-better-marketplace-page="true"] input[placeholder*="search"]:not(.ama-search-input) {
                background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='rgba(255,255,255,0.45)' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='11' cy='11' r='8'/%3E%3Cline x1='21' y1='21' x2='16.65' y2='16.65'/%3E%3C/svg%3E") !important;
                background-repeat: no-repeat !important;
                background-position: 10px center !important;
                background-size: 16px 16px !important;
                padding-left: 34px !important;
            }

            [data-anilist-collection-lists="true"] [data-page-wrapper-container="true"],
            [data-anilist-collection-lists="true"] .py-6 {
                display: block !important;
                width: 100% !important;
                max-width: 100% !important;
            }

            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"],
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"],
            body[data-ama-carousels-active="true"] [data-manga-page-container="true"] [data-media-card-grid="true"],
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-grid="true"],
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-lazy-grid="true"],
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists-tabs] ~ div .grid {
                display: flex !important;
                flex-wrap: nowrap !important;
                overflow-x: auto !important;
                overflow-y: hidden !important;
                gap: 16px !important;
                padding: 16px 8px 22px 8px !important;
                margin: 0 !important;
                width: 100% !important;
                scroll-snap-type: none !important;
                scroll-behavior: auto !important;
                overscroll-behavior-x: contain !important;
                -webkit-overflow-scrolling: touch !important;
                touch-action: pan-x pan-y !important;
                contain: layout paint style !important;
            }

            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"],
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"],
            body[data-ama-carousels-active="true"] [data-manga-page-container="true"] [data-media-card-grid="true"],
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-grid="true"],
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-lazy-grid="true"],
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists-tabs] ~ div .grid,
            body[data-ama-better-marketplace="true"] .ama-carousel-row {
                cursor: grab !important;
                scrollbar-width: none !important;
                -ms-overflow-style: none !important;
            }

            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"]::-webkit-scrollbar,
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"]::-webkit-scrollbar,
            body[data-ama-carousels-active="true"] [data-manga-page-container="true"] [data-media-card-grid="true"]::-webkit-scrollbar,
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-grid="true"]::-webkit-scrollbar,
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-lazy-grid="true"]::-webkit-scrollbar,
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists-tabs] ~ div .grid::-webkit-scrollbar,
            body[data-ama-better-marketplace="true"] .ama-carousel-row::-webkit-scrollbar {
                display: none !important;
                width: 0 !important;
                height: 0 !important;
            }

            .ama-drag-pending,
            .ama-dragging {
                cursor: grabbing !important;
                user-select: none !important;
            }

            .ama-drag-pending *,
            .ama-dragging * {
                cursor: grabbing !important;
                user-select: none !important;
            }

            .ama-dragging a,
            .ama-dragging img,
            .ama-dragging [role="button"] {
                pointer-events: none !important;
            }

            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"] a,
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"] a,
            body[data-ama-carousels-active="true"] [data-manga-page-container="true"] [data-media-card-grid="true"] a,
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-grid="true"] a,
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-lazy-grid="true"] a,
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists-tabs] ~ div .grid a,
            body[data-ama-better-marketplace="true"] .ama-carousel-row a,
            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"] img,
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"] img,
            body[data-ama-better-marketplace="true"] .ama-carousel-row img {
                -webkit-user-drag: none !important;
                user-drag: none !important;
            }

            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"] > div,
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"] > div,
            body[data-ama-carousels-active="true"] [data-manga-page-container="true"] [data-media-card-grid="true"] > div {
                flex: 0 0 160px !important;
                display: block !important;
                width: 160px !important;
                max-width: 160px !important;
                contain: layout paint style !important;
            }

            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-grid="true"] > div,
            body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-lazy-grid="true"] > div {
                flex: 0 0 155px !important;
                display: block !important;
                width: 155px !important;
                max-width: 155px !important;
                contain: layout paint style !important;
            }

            @media (min-width: 768px) {
                body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"] > div,
                body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"] > div,
                body[data-ama-carousels-active="true"] [data-manga-page-container="true"] [data-media-card-grid="true"] > div {
                    flex: 0 0 185px !important;
                    width: 185px !important;
                    max-width: 185px !important;
                }

                body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-grid="true"] > div,
                body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-lazy-grid="true"] > div {
                    flex: 0 0 175px !important;
                    width: 175px !important;
                    max-width: 175px !important;
                }
            }

            @media (min-width: 1280px) {
                body[data-ama-carousels-active="true"] [data-manga-page-container="true"] [data-media-card-grid="true"] > div {
                    flex: 0 0 210px !important;
                    width: 210px !important;
                    max-width: 210px !important;
                }

                body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-grid="true"] > div,
                body[data-ama-carousels-active="true"] [data-anilist-collection-lists="true"] [data-media-card-lazy-grid="true"] > div {
                    flex: 0 0 195px !important;
                    width: 195px !important;
                    max-width: 195px !important;
                }
            }

            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"] [data-media-entry-card-hover-popup="true"],
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"] [data-media-entry-card-hover-popup="true"] {
                display: none !important;
                opacity: 0 !important;
                visibility: hidden !important;
                pointer-events: none !important;
            }

            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"] [data-media-entry-card-body-image="true"],
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"] [data-media-entry-card-body-image="true"] {
                transition: none !important;
                transform: none !important;
                will-change: auto !important;
            }

            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"] .group\\/media-entry-card:hover [data-media-entry-card-body-image="true"],
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"] .group\\/media-entry-card:hover [data-media-entry-card-body-image="true"] {
                transform: none !important;
            }

            body[data-ama-carousels-active="true"] .grid[data-media-card-grid="true"] *,
            body[data-ama-carousels-active="true"] .grid[data-media-card-lazy-grid="true"] * {
                scroll-behavior: auto !important;
            }

            .ama-optimized-carousel {
                display: flex !important;
                flex-wrap: nowrap !important;
                overflow-x: auto !important;
                overflow-y: hidden !important;
                gap: 16px !important;
                padding: 16px 8px 22px 8px !important;
                margin: 0 !important;
                width: 100% !important;
                scroll-snap-type: none !important;
                scroll-behavior: auto !important;
                overscroll-behavior-x: contain !important;
                -webkit-overflow-scrolling: touch !important;
                touch-action: pan-x pan-y !important;
                contain: layout paint style !important;
                cursor: grab !important;
                scrollbar-width: none !important;
                -ms-overflow-style: none !important;
            }

            .ama-optimized-carousel::-webkit-scrollbar {
                display: none !important;
                width: 0 !important;
                height: 0 !important;
            }

            .ama-optimized-carousel a,
            .ama-optimized-carousel img {
                -webkit-user-drag: none !important;
                user-drag: none !important;
            }

            .ama-optimized-carousel > div {
                flex: 0 0 160px !important;
                display: block !important;
                width: 160px !important;
                max-width: 160px !important;
                contain: layout paint style !important;
            }

            @media (min-width: 768px) {
                .ama-optimized-carousel > div {
                    flex: 0 0 185px !important;
                    width: 185px !important;
                    max-width: 185px !important;
                }
            }

            @media (min-width: 1280px) {
                .ama-optimized-carousel > div {
                    flex: 0 0 210px !important;
                    width: 210px !important;
                    max-width: 210px !important;
                }
            }

            .ama-optimized-carousel [data-media-entry-card-hover-popup="true"] {
                display: none !important;
                opacity: 0 !important;
                visibility: hidden !important;
                pointer-events: none !important;
            }

            .ama-optimized-carousel [data-media-entry-card-body-image="true"] {
                transition: none !important;
                transform: none !important;
                will-change: auto !important;
            }

            .ama-optimized-carousel .group\/media-entry-card:hover [data-media-entry-card-body-image="true"] {
                transform: none !important;
            }

            .ama-optimized-carousel * {
                scroll-behavior: auto !important;
            }

            body[data-ama-subdub-icons="true"] .ama-media-meta-row {
                display: flex !important;
                align-items: center !important;
                justify-content: space-between !important;
                gap: 8px !important;
                width: 100% !important;
                min-width: 0 !important;
            }

            body[data-ama-subdub-icons="true"] .ama-media-meta-row > [data-media-entry-card-title-section-year-season="true"] {
                min-width: 0 !important;
                overflow: hidden !important;
                text-overflow: ellipsis !important;
                white-space: nowrap !important;
            }

            body[data-ama-subdub-icons="true"] .ama-media-badges {
                margin-left: auto !important;
                display: inline-flex !important;
                align-items: center !important;
                justify-content: flex-end !important;
                gap: 6px !important;
                flex: 0 0 auto !important;
            }

            body[data-ama-subdub-icons="true"] .ama-media-badge {
                display: inline-flex !important;
                align-items: center !important;
                justify-content: center !important;
                height: 22px !important;
                min-width: 28px !important;
                padding: 0 6px !important;
                border-radius: 999px !important;
                border: 1px solid currentColor !important;
                background: rgba(255,255,255,0.06) !important;
                line-height: 1 !important;
                flex: 0 0 auto !important;
            }

            body[data-ama-subdub-icons="true"] .ama-media-badge svg {
                width: 17px !important;
                height: 17px !important;
                display: block !important;
            }

            body[data-ama-subdub-icons="true"] .ama-media-badge.cc {
                color: #a7e7ff !important;
                background: rgba(167, 231, 255, 0.13) !important;
                border-color: rgba(167, 231, 255, 0.65) !important;
                box-shadow: 0 0 10px rgba(167, 231, 255, 0.12) !important;
            }

            body[data-ama-subdub-icons="true"] .ama-media-badge.dub {
                color: #ffb7c5 !important;
                background: rgba(255, 183, 197, 0.13) !important;
                border-color: rgba(255, 183, 197, 0.7) !important;
                box-shadow: 0 0 10px rgba(255, 183, 197, 0.12) !important;
            }

            body[data-ama-subdub-icons="true"] .ama-media-badge.dub[hidden] {
                display: none !important;
            }

            body[data-ama-subdub-icons="true"] .ama-schedule-dub-badge {
                margin-left: 6px !important;
                vertical-align: middle !important;
                flex: 0 0 auto !important;
                min-width: 28px !important;
            }

            body[data-ama-subdub-icons="false"] .ama-media-badges {
                display: none !important;
            }

            body[data-ama-better-marketplace="true"] .group\\/extension-card {
                flex: 0 0 300px !important;
                background: rgba(255, 255, 255, 0.05) !important;
                border: 1px solid rgba(255, 255, 255, 0.08) !important;
                border-radius: 24px !important;
                padding: 20px !important;
                transition: background 0.12s ease !important;
                contain: layout paint style !important;
            }

            body[data-ama-better-marketplace="true"] .group\\/extension-card:hover {
                background: rgba(255, 255, 255, 0.08) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-marketplace-extra,
            body[data-ama-better-marketplace="true"] .ama-installed-marketplace-extra {
                display: flex !important;
                flex-wrap: wrap !important;
                align-items: center !important;
                gap: 6px !important;
                margin-top: 10px !important;
                padding-top: 10px !important;
                border-top: 1px solid rgba(255,255,255,0.08) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-marketplace-extra-badge {
                display: inline-flex !important;
                align-items: center !important;
                min-height: 20px !important;
                padding: 3px 7px !important;
                border-radius: 999px !important;
                border: 1px solid rgba(255,255,255,0.12) !important;
                background: rgba(255,255,255,0.06) !important;
                color: rgba(255,255,255,0.78) !important;
                font-size: 11px !important;
                font-weight: 700 !important;
                line-height: 1 !important;
                white-space: nowrap !important;
            }

            body[data-ama-better-marketplace="true"] .ama-marketplace-extra-badge.working {
                color: #86efac !important;
                border-color: rgba(74,222,128,0.45) !important;
                background: rgba(34,197,94,0.12) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-marketplace-extra-badge.broken {
                color: #fca5a5 !important;
                border-color: rgba(248,113,113,0.5) !important;
                background: rgba(239,68,68,0.14) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-marketplace-extra-badge.deprecated {
                color: #fde68a !important;
                border-color: rgba(250,204,21,0.45) !important;
                background: rgba(234,179,8,0.12) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-marketplace-extra-badge.scan {
                color: #bae6fd !important;
                border-color: rgba(56,189,248,0.45) !important;
                background: rgba(14,165,233,0.12) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-marketplace-extra-badge.official {
                color: #c4b5fd !important;
                border-color: rgba(167,139,250,0.45) !important;
                background: rgba(124,58,237,0.12) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-marketplace-extra-badge.link {
                text-decoration: none !important;
                cursor: pointer !important;
            }

            body[data-ama-better-marketplace="true"] .ama-marketplace-extra-note {
                width: 100% !important;
                color: rgba(255,255,255,0.48) !important;
                font-size: 11px !important;
                line-height: 1.35 !important;
            }

            body[data-ama-better-marketplace="true"] .group\\/extension-card[data-ama-update-available="true"],
            body[data-ama-better-marketplace="true"] .ama-catalog-card-wrap[data-ama-update-available="true"] > .group\\/extension-card,
            body[data-ama-better-marketplace="true"] .UI-Card__root[data-ama-update-available="true"] {
                border-color: rgba(56, 189, 248, 0.95) !important;
                outline: 1px solid rgba(56, 189, 248, 0.85) !important;
                outline-offset: 0 !important;
                background: rgba(14, 165, 233, 0.12) !important;
                box-shadow: 0 0 0 1px rgba(56, 189, 248, 0.75), 0 0 18px rgba(56, 189, 248, 0.24) !important;
            }

            body[data-ama-better-marketplace="true"] .group\\/extension-card[data-ama-update-available="true"]:hover,
            body[data-ama-better-marketplace="true"] .ama-catalog-card-wrap[data-ama-update-available="true"] > .group\\/extension-card:hover,
            body[data-ama-better-marketplace="true"] .UI-Card__root[data-ama-update-available="true"]:hover {
                background: rgba(14, 165, 233, 0.18) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-extension-carousel {
                display: flex !important;
                flex-wrap: nowrap !important;
                overflow-x: auto !important;
                overflow-y: hidden !important;
                gap: 16px !important;
                padding: 10px 40px 22px 40px !important;
                margin: 0 !important;
                width: 100% !important;
                scroll-behavior: auto !important;
                overscroll-behavior-x: contain !important;
                -webkit-overflow-scrolling: touch !important;
                touch-action: pan-x pan-y !important;
                contain: layout paint style !important;
                cursor: grab !important;
                scrollbar-width: none !important;
                -ms-overflow-style: none !important;
            }

            body[data-ama-better-marketplace="true"] .ama-extension-carousel::-webkit-scrollbar {
                display: none !important;
                width: 0 !important;
                height: 0 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-extension-carousel a,
            body[data-ama-better-marketplace="true"] .ama-extension-carousel img {
                -webkit-user-drag: none !important;
                user-drag: none !important;
            }

            body[data-ama-better-marketplace="true"] .ama-status-section {
                margin-top: 10px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-status-normal-grid {
                display: grid !important;
                grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)) !important;
                gap: 16px !important;
                padding: 10px 40px 22px 40px !important;
                margin: 0 !important;
                width: 100% !important;
                overflow: visible !important;
                cursor: default !important;
            }

            body[data-ama-better-marketplace="true"] .ama-status-normal-grid .ama-catalog-card-wrap {
                flex: none !important;
                width: 100% !important;
                max-width: none !important;
            }

            body[data-ama-better-marketplace="true"] .ama-status-source-hidden {
                display: none !important;
            }

            body[data-ama-better-marketplace="true"] .ama-header-container {
                display: flex !important;
                align-items: center !important;
                justify-content: space-between !important;
                gap: 20px !important;
                padding: 14px 40px 0px 40px !important;
                margin-bottom: 8px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-header-left {
                display: flex !important;
                align-items: center !important;
                gap: 12px !important;
                min-width: 0 !important;
                flex-shrink: 1 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-header-left h3 {
                margin: 0 !important;
                white-space: nowrap !important;
            }

            body[data-ama-better-marketplace="true"] .ama-header-right {
                display: flex !important;
                align-items: center !important;
                gap: 12px !important;
                margin-left: auto !important;
            }

            body[data-ama-better-marketplace="true"] .ama-search-wrapper {
                position: relative !important;
                display: flex !important;
                align-items: center !important;
            }

            body[data-ama-better-marketplace="true"] .ama-search-icon {
                position: absolute !important;
                left: 10px !important;
                width: 15px !important;
                height: 15px !important;
                pointer-events: none !important;
                opacity: 0.45 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-search-input {
                background: rgba(255,255,255,0.06) !important;
                border: 1px solid rgba(255,255,255,0.08) !important;
                border-radius: 14px !important;
                padding: 8px 14px 8px 32px !important;
                color: white !important;
                width: 220px !important;
                outline: none !important;
            }

            body[data-ama-better-marketplace="true"] .ama-view-btn {
                background: rgba(255,255,255,0.08) !important;
                border: 1px solid rgba(255,255,255,0.08) !important;
                padding: 8px 16px !important;
                border-radius: 14px !important;
                color: #fff !important;
                cursor: pointer !important;
                font-weight: 600 !important;
                white-space: nowrap !important;
            }

            body[data-ama-better-marketplace="true"] .ama-global-catalog-bar {
                display: flex !important;
                align-items: center !important;
                justify-content: flex-end !important;
                gap: 12px !important;
                padding: 0 40px 12px 40px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-global-catalog-btn {
                background: rgba(14,165,233,0.14) !important;
                border: 1px solid rgba(56,189,248,0.45) !important;
                color: #e0f7ff !important;
                padding: 9px 16px !important;
                border-radius: 14px !important;
                cursor: pointer !important;
                font-weight: 700 !important;
                white-space: nowrap !important;
            }

            body[data-ama-better-marketplace="true"] .ama-global-catalog-btn:hover {
                background: rgba(14,165,233,0.22) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-modal {
                position: fixed !important;
                top: 5% !important;
                left: 5% !important;
                width: 90% !important;
                height: 90% !important;
                background: rgba(18,18,22,0.96) !important;
                border: 1px solid rgba(255,255,255,0.08) !important;
                border-radius: 28px !important;
                padding: 30px !important;
                z-index: 999999 !important;
                overflow-y: auto !important;
                color: white !important;
                contain: layout paint style !important;
            }

            body[data-ama-better-marketplace="true"] .ama-modal-close {
                float: right !important;
                background: rgba(255,255,255,0.08) !important;
                border: none !important;
                color: white !important;
                padding: 8px 14px !important;
                border-radius: 12px !important;
                cursor: pointer !important;
            }

            body[data-ama-better-marketplace="true"] .modal-row {
                margin-bottom: 28px !important;
            }

            body[data-ama-better-marketplace="true"] .modal-row-title {
                font-weight: 700 !important;
                font-size: 18px !important;
                margin-bottom: 12px !important;
                padding-left: 40px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-carousel-row {
                display: flex !important;
                overflow-x: auto !important;
                overflow-y: hidden !important;
                gap: 20px !important;
                padding: 10px 40px 22px 40px !important;
                scroll-behavior: auto !important;
                overscroll-behavior-x: contain !important;
                -webkit-overflow-scrolling: touch !important;
                contain: layout paint style !important;
            }

            body[data-ama-better-marketplace="true"] .ama-catalog-card-wrap {
                position: relative !important;
                flex: 0 0 300px !important;
                width: 300px !important;
                max-width: 300px !important;
                display: block !important;
            }

            body[data-ama-better-marketplace="true"] .ama-catalog-card-wrap > .group\\/extension-card {
                width: 100% !important;
                max-width: 100% !important;
                height: 100% !important;
                box-sizing: border-box !important;
            }

            body[data-ama-better-marketplace="true"] .ama-catalog-card-wrap[hidden] {
                display: none !important;
            }

            body[data-ama-better-marketplace="true"] .ama-clone-actions {
                position: absolute !important;
                top: 12px !important;
                right: 12px !important;
                display: grid !important;
                grid-template-columns: 32px 32px !important;
                grid-auto-rows: 32px !important;
                gap: 4px !important;
                z-index: 9999999 !important;
                pointer-events: auto !important;
            }

            body[data-ama-better-marketplace="true"] .ama-clone-action[data-ama-action="preferences"] {
                grid-column: 1 !important;
                grid-row: 1 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-clone-action[data-ama-action="more"] {
                grid-column: 2 !important;
                grid-row: 1 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-clone-action[data-ama-action="documentation"] {
                grid-column: 1 !important;
                grid-row: 2 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-clone-action[data-ama-action="code"] {
                grid-column: 2 !important;
                grid-row: 2 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-clone-action {
                width: 32px !important;
                height: 32px !important;
                border-radius: 10px !important;
                border: 1px solid rgba(255,255,255,0.12) !important;
                background: rgba(255,255,255,0.08) !important;
                color: white !important;
                display: inline-flex !important;
                align-items: center !important;
                justify-content: center !important;
                cursor: pointer !important;
                pointer-events: auto !important;
                position: relative !important;
                z-index: 9999999 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-clone-action:hover {
                background: rgba(255,255,255,0.14) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-code-preview {
                margin: 16px 0 0 0 !important;
                padding: 16px !important;
                background: rgba(0,0,0,0.35) !important;
                border: 1px solid rgba(255,255,255,0.08) !important;
                border-radius: 14px !important;
                max-height: 70vh !important;
                overflow: auto !important;
                white-space: pre-wrap !important;
                word-break: break-word !important;
                font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace !important;
                font-size: 12px !important;
                line-height: 1.5 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-detail-grid {
                display: grid !important;
                grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)) !important;
                gap: 12px !important;
                margin-top: 16px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-detail-item {
                padding: 12px !important;
                border-radius: 12px !important;
                background: rgba(255,255,255,0.06) !important;
                border: 1px solid rgba(255,255,255,0.08) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-detail-label {
                color: rgba(255,255,255,0.5) !important;
                font-size: 12px !important;
                margin-bottom: 4px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-action-panel {
                display: flex !important;
                flex-wrap: wrap !important;
                gap: 10px !important;
                margin-top: 18px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-action-button {
                border: 1px solid rgba(255,255,255,0.12) !important;
                background: rgba(255,255,255,0.08) !important;
                color: white !important;
                padding: 10px 14px !important;
                border-radius: 12px !important;
                cursor: pointer !important;
                font-weight: 700 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-action-button:hover {
                background: rgba(255,255,255,0.14) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-action-button[disabled] {
                cursor: not-allowed !important;
                opacity: .55 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-action-button.danger {
                background: rgba(239,68,68,0.14) !important;
                border-color: rgba(239,68,68,0.35) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-action-status {
                width: 100% !important;
                min-height: 20px !important;
                color: rgba(255,255,255,.68) !important;
                font-size: 13px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-config-form {
                display: flex !important;
                flex-direction: column !important;
                gap: 14px !important;
                margin-top: 18px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-config-field {
                display: flex !important;
                flex-direction: column !important;
                gap: 6px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-config-field label,
            body[data-ama-better-marketplace="true"] .ama-config-switch-label {
                font-weight: 700 !important;
                color: rgba(255,255,255,0.9) !important;
            }

            body[data-ama-better-marketplace="true"] .ama-config-input,
            body[data-ama-better-marketplace="true"] .ama-config-select {
                background: rgba(255,255,255,0.06) !important;
                border: 1px solid rgba(255,255,255,0.1) !important;
                border-radius: 12px !important;
                color: #fff !important;
                padding: 10px 12px !important;
                outline: none !important;
                width: 100% !important;
            }

            body[data-ama-better-marketplace="true"] .ama-config-switch-row {
                display: flex !important;
                align-items: center !important;
                justify-content: space-between !important;
                gap: 14px !important;
                padding: 10px 12px !important;
                background: rgba(255,255,255,0.04) !important;
                border: 1px solid rgba(255,255,255,0.08) !important;
                border-radius: 12px !important;
            }

            .ama-config-switch-row[data-ama-disabled="true"] {
                opacity: .42 !important;
                filter: grayscale(1) !important;
                cursor: not-allowed !important;
                background: rgba(255,255,255,0.025) !important;
                border-color: rgba(255,255,255,0.04) !important;
            }

            .ama-config-switch-row[data-ama-disabled="true"] * {
                cursor: not-allowed !important;
            }

            .ama-config-switch-row[data-ama-disabled="true"] input {
                opacity: .5 !important;
                accent-color: #64748b !important;
                filter: grayscale(1) saturate(0) !important;
            }

            .ama-config-switch-row[data-ama-disabled="true"] input[type="checkbox"] {
                appearance: none !important;
                -webkit-appearance: none !important;
                width: 38px !important;
                height: 22px !important;
                min-width: 38px !important;
                border-radius: 999px !important;
                border: 1px solid rgba(148,163,184,.45) !important;
                background: #1f2937 !important;
                position: relative !important;
                opacity: .9 !important;
                filter: none !important;
            }

            .ama-config-switch-row[data-ama-disabled="true"] input[type="checkbox"]:checked {
                background: #475569 !important;
                border-color: rgba(148,163,184,.7) !important;
            }

            .ama-config-switch-row[data-ama-disabled="true"] input[type="checkbox"]::before {
                content: "" !important;
                position: absolute !important;
                top: 2px !important;
                left: 2px !important;
                width: 16px !important;
                height: 16px !important;
                border-radius: 999px !important;
                background: #94a3b8 !important;
                transition: none !important;
            }

            .ama-config-switch-row[data-ama-disabled="true"] input[type="checkbox"]:checked::before {
                transform: translateX(16px) !important;
                background: #cbd5e1 !important;
            }

            .ama-disabled-switch-visual {
                min-width: 44px !important;
                height: 24px !important;
                border-radius: 999px !important;
                border: 1px solid rgba(180,180,180,.45) !important;
                background: #262626 !important;
                color: #d4d4d4 !important;
                display: inline-flex !important;
                align-items: center !important;
                justify-content: center !important;
                font-size: 11px !important;
                font-weight: 700 !important;
                line-height: 1 !important;
                pointer-events: none !important;
                box-shadow: none !important;
            }

            .ama-disabled-switch-visual::before {
                content: none !important;
            }

            .ama-disabled-switch-visual[data-ama-checked="true"] {
                background: #3f3f46 !important;
                border-color: rgba(212,212,216,.65) !important;
                color: #f4f4f5 !important;
            }

            .ama-disabled-switch-visual[data-ama-checked="true"]::before {
                content: none !important;
            }

            .ama-config-toggle {
                min-width: 52px !important;
                height: 26px !important;
                border-radius: 999px !important;
                border: 1px solid rgba(148,163,184,.35) !important;
                background: rgba(15,23,42,.9) !important;
                color: #cbd5e1 !important;
                display: inline-flex !important;
                align-items: center !important;
                justify-content: center !important;
                font-size: 11px !important;
                font-weight: 800 !important;
                cursor: pointer !important;
                box-shadow: none !important;
            }

            .ama-config-toggle[data-ama-value="true"] {
                background: rgba(14,165,233,.22) !important;
                border-color: rgba(56,189,248,.7) !important;
                color: #e0f7ff !important;
            }

            .ama-config-switch-row[data-ama-disabled="true"] .ama-config-toggle {
                background: #262626 !important;
                border-color: rgba(180,180,180,.45) !important;
                color: #d4d4d4 !important;
                cursor: not-allowed !important;
            }

            .ama-config-switch-row[data-ama-disabled="true"] .ama-config-toggle[data-ama-value="true"] {
                background: #3f3f46 !important;
                border-color: rgba(212,212,216,.65) !important;
                color: #f4f4f5 !important;
            }

            body[data-ama-better-marketplace="true"] .ama-config-help {
                color: rgba(255,255,255,0.5) !important;
                font-size: 12px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-config-actions {
                display: flex !important;
                align-items: center !important;
                gap: 12px !important;
                margin-top: 4px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-config-save {
                background: #fff !important;
                border: 0 !important;
                color: #111 !important;
                border-radius: 12px !important;
                cursor: pointer !important;
                font-weight: 700 !important;
                padding: 10px 16px !important;
            }

            body[data-ama-better-marketplace="true"] .ama-config-status {
                color: rgba(255,255,255,0.65) !important;
                font-size: 13px !important;
            }

            body[data-ama-better-marketplace="true"] .flex.gap-2.flex-wrap.pt-4 .UI-Badge__root:nth-child(1),
            body[data-ama-better-marketplace="true"] .group\\/extension-card .UI-Badge__root:nth-child(1) {
                color: aqua !important;
                border-color: aqua !important;
            }

            body[data-ama-better-marketplace="true"] .flex.gap-2.flex-wrap.pt-4 .UI-Badge__root:nth-child(2),
            body[data-ama-better-marketplace="true"] .group\\/extension-card .UI-Badge__root:nth-child(2) {
                color: #4ade80 !important;
                border-color: #4ade80 !important;
            }

            body[data-ama-better-marketplace="true"] .flex.gap-2.flex-wrap.pt-4 .UI-Badge__root:nth-child(3),
            body[data-ama-better-marketplace="true"] .group\\/extension-card .UI-Badge__root:nth-child(3) {
                color: #f87171 !important;
                border-color: #f87171 !important;
                padding-left: 8px !important;
                padding-right: 8px !important;
            }

            @media (max-width: 900px) {
                body[data-ama-better-marketplace="true"] .ama-header-container {
                    flex-direction: column !important;
                    align-items: stretch !important;
                }

                body[data-ama-better-marketplace="true"] .ama-header-right {
                    width: 100% !important;
                }

                body[data-ama-better-marketplace="true"] .ama-search-input {
                    flex: 1 !important;
                    width: 100% !important;
                }
            }
        `;

        ctx.dom.onReady(async () => {
            const body = await ctx.dom.queryOne("body");
            if (!body) return;

            const style = await ctx.dom.createElement("style");
            style.setText(carouselCSS);
            body.append(style);

            const script = await ctx.dom.createElement("script");
            script.setText(`
                (() => {
                    const SETTINGS_KEY = 'ama-ui-tweaks.settings';

                    const targetGridsQuery = '.grid[data-media-card-grid="true"], .grid[data-media-card-lazy-grid="true"]';
                    const mediaEntryCardQuery = '[data-media-entry-card-container="true"]';
                    const scheduleEventQuery = '[data-schedule-calendar-event-item-link], [data-schedule-calendar-event-item-content], [data-schedule-calendar-event-item-name], [data-schedule-calendar-event-item-episode], [data-schedule-calendar-event-item-finale-icon], [data-schedule-calendar-event-item], [data-schedule-calendar-mobile-list-day-item-event-link], [data-schedule-calendar-mobile-list-day-item-event-content], [data-schedule-calendar-mobile-list-day-item-event-text], [data-schedule-calendar-mobile-list-day-item-event-episode], [data-schedule-calendar-mobile-list-day-item-event-icons], [data-schedule-event-item], [data-schedule-media-id]';
                    const scheduleEntryLinkQuery = 'a[href^="/entry?id="], a[href*="/entry?id="]';
                    const cardQuery = '.UI-Card__root';
                    const arrowQuery = '.ama-carousel-nav-btn, .ama-manga-carousel-parent';
                    const randomSearchIconPath = 'M10 18a7.952 7.952 0 0 0 4.897-1.688l4.396 4.396 1.414-1.414-4.396-4.396A7.952 7.952 0 0 0 18 10c0-4.411-3.589-8-8-8s-8 3.589-8 8 3.589 8 8 8zm0-14c3.309 0 6 2.691 6 6s-2.691 6-6 6-6-2.691-6-6 2.691-6 6-6z';
                    const dubAniScheduleUrl = 'https://raw.githubusercontent.com/Bas1874/AniSchedule/refs/heads/master/raw/dub-episode-feed.json';
                    const dubMappingUrl = 'https://raw.githubusercontent.com/Joelis57/MyDubList/refs/heads/main/dubs/mappings/mappings_anilist.jsonl';
                    const dubEnglishSourceUrls = [
                        'https://raw.githubusercontent.com/Joelis57/MyDubList/refs/heads/main/dubs/confidence/very-high/dubbed_english.json',
                        'https://raw.githubusercontent.com/Joelis57/MyDubList/refs/heads/main/dubs/confidence/high/dubbed_english.json'
                    ];
                    const dubFeedCacheKey = 'ama-hybrid-dub-anilist-ids-v1';
                    const dubFeedCacheTTL = 1000 * 60 * 60 * 12;
                    const animeScheduleApiBaseUrl = 'https://animeschedule.net/api/v3';
                    const asunaTracksScheduleApiBaseUrl = 'https://asunatracks.space/public/api';
                    const animeScheduleApiTokenKey = 'ama-ui-tweaks.animeScheduleApiToken';

                    const defaultSettings = {
                        betterMarketplace: true,
                        carousels: true,
                        carouselsSearch: false,
                        carouselsExtensions: true,
                        carouselsLists: true,
                        carouselsManga: true,
                        carouselsOther: true,
                        subDubIcons: true,
                        hideFileNames: false,
                        useBas1874Marketplace: false,
                    };

                    const BAS1874_MARKETPLACE_URL = 'https://raw.githubusercontent.com/Bas1874/Seanime-Marketplace/refs/heads/main/Marketplace/Main.json';

                    function readBrowserSettings() {
                        try {
                            const raw = window.localStorage.getItem(SETTINGS_KEY);
                            if (!raw) return {};
                            const parsed = JSON.parse(raw);
                            if (!parsed || typeof parsed !== 'object') return {};
                            return parsed;
                        } catch (_) {
                            return {};
                        }
                    }

                    function writeBrowserSettings(settings) {
                        try {
                            window.localStorage.setItem(SETTINGS_KEY, JSON.stringify(settings));
                        } catch (_) {}
                    }

                    function clearAnimeScheduleApiTokenSetting() {
                        try {
                            window.localStorage.removeItem(animeScheduleApiTokenKey);
                        } catch (_) {}
                    }

                    function normalizeFeatureSettings(settings) {
                        const next = Object.assign({}, defaultSettings, settings || {});
                        delete next.animeScheduleApiToken;
                        return next;
                    }

                    let featureSettings = normalizeFeatureSettings(Object.assign(
                        {},
                        defaultSettings,
                        ${JSON.stringify(initialFeatureSettings)},
                        readBrowserSettings()
                    ));
                    clearAnimeScheduleApiTokenSetting();

                    const CC_ICON =
                        '<svg stroke="currentColor" fill="currentColor" stroke-width="0" viewBox="0 0 24 24" aria-hidden="true" xmlns="http://www.w3.org/2000/svg">' +
                            '<path d="M19 4H5c-1.103 0-2 .897-2 2v12c0 1.103.897 2 2 2h14c1.103 0 2-.897 2-2V6c0-1.103-.897-2-2-2zM5 18V6h14l.002 12H5z"></path>' +
                            '<path d="M8 15h2v-2H8v-2h2V9H8c-1.103 0-2 .897-2 2v2c0 1.103.897 2 2 2zm6 0h2v-2h-2v-2h2V9h-2c-1.103 0-2 .897-2 2v2c0 1.103.897 2 2 2z"></path>' +
                        '</svg>';

                    const MIC_ICON =
                        '<svg stroke="currentColor" fill="currentColor" stroke-width="0" viewBox="0 0 24 24" aria-hidden="true" xmlns="http://www.w3.org/2000/svg">' +
                            '<path d="M12 14c1.654 0 3-1.346 3-3V5c0-1.654-1.346-3-3-3S9 3.346 9 5v6c0 1.654 1.346 3 3 3z"></path>' +
                            '<path d="M17 11c0 2.757-2.243 5-5 5s-5-2.243-5-5H5c0 3.519 2.613 6.432 6 6.92V21h2v-3.08c3.387-.488 6-3.401 6-6.92h-2z"></path>' +
                        '</svg>';

                    const DOWNLOAD_ICON = '<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><path d="M7 10l5 5 5-5"></path><path d="M12 15V3"></path></svg>';
                    const MORE_ICON = '<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="1"></circle><circle cx="12" cy="5" r="1"></circle><circle cx="12" cy="19" r="1"></circle></svg>';
                    const CODE_ICON = '<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m16 18 6-6-6-6"></path><path d="m8 6-6 6 6 6"></path></svg>';
                    const DOC_ICON = '<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path><path d="M4 4.5A2.5 2.5 0 0 1 6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5z"></path></svg>';
                    const SETTINGS_ICON = '<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="21" x2="14" y1="4" y2="4"></line><line x1="10" x2="3" y1="4" y2="4"></line><line x1="21" x2="12" y1="12" y2="12"></line><line x1="8" x2="3" y1="12" y2="12"></line><line x1="21" x2="16" y1="20" y2="20"></line><line x1="12" x2="3" y1="20" y2="20"></line><line x1="14" x2="14" y1="2" y2="6"></line><line x1="8" x2="8" y1="10" y2="14"></line><line x1="16" x2="16" y1="18" y2="22"></line></svg>';

                    let dubIdSetPromise = null;
                    let malToAnilistMapPromise = null;
                    let animeScheduleDubApiEntriesPromises = new Map();
                    let animeScheduleDubFeedEntriesPromise = null;
                    const animeScheduleAnimeDetailsPromises = new Map();
                    const dragScrollEnhancementVersion = 'v8';
                    const marketplaceEnhancementVersion = 'v11';
                    const catalogActionSources = new Map();
                    let allExtensionsPromise = null;
                    let marketplaceExtensionsPromise = null;
                    let bas1874MarketplaceMetadataPromise = null;
                    let bas1874MarketplaceMetadataCache = [];
                    let marketplaceStatusSectionsScheduled = false;
                    let marketplaceStatusSectionsRendering = false;
                    let catalogActionSourceCounter = 0;
                    let catalogActionHandlerBound = false;
                    let marketplaceSearchVisibilityHandlerBound = false;

                    function isElement(node) {
                        return node && node.nodeType === 1;
                    }

                    function hasSearchPageSignal() {
                        const path = String(window.location.pathname || '').toLowerCase();
                        const search = String(window.location.search || '').toLowerCase();
                        const hash = String(window.location.hash || '').toLowerCase();
                        const title = String(document.title || '').toLowerCase();
                        const combined = path + ' ' + search + ' ' + hash + ' ' + title;

                        if (/(^|[\/#?&=-])search($|[\/#?&=-])/.test(combined)) return true;
                        if (document.querySelector('[data-search-page-container], [data-search-page-title], [data-search-page-list], [data-search-page="true"], [data-media-search-page="true"]')) return true;
                        if (document.querySelector('input[placeholder*="Search anime"], input[placeholder*="Search Anime"], input[placeholder*="Search manga"], input[placeholder*="Search Manga"], input[placeholder*="Search media"], input[placeholder*="Search Media"], input[placeholder*="Search for anime"], input[placeholder*="Search for Anime"], input[placeholder*="Search for manga"], input[placeholder*="Search for Manga"]')) return true;

                        return false;
                    }

                    function getCarouselPageKey() {
                        const path = String(window.location.pathname || '').toLowerCase();
                        const hash = String(window.location.hash || '').toLowerCase();
                        const combined = path + ' ' + hash;

                        if (hasSearchPageSignal()) {
                            return 'search';
                        }

                        if (isBetterMarketplacePage()) {
                            return 'extensions';
                        }

                        if (
                            combined.includes('list') ||
                            combined.includes('collection') ||
                            combined.includes('library') ||
                            document.querySelector('[data-anilist-collection-lists="true"], [data-anilist-collection-lists-tabs]')
                        ) {
                            return 'lists';
                        }

                        if (
                            combined.includes('manga') ||
                            combined.includes('novel') ||
                            combined.includes('book') ||
                            document.querySelector('[data-manga-page-container="true"], [data-manga-entry-screen="true"], [data-manga-entry-page="true"]')
                        ) {
                            return 'manga';
                        }

                        return 'other';
                    }

                    function areCarouselsEnabledForCurrentPage() {
                        if (!featureSettings.carousels) return false;

                        const pageKey = getCarouselPageKey();

                        if (pageKey === 'search') return featureSettings.carouselsSearch === true;
                        if (pageKey === 'extensions') return areExtensionCarouselsEnabled();
                        if (pageKey === 'lists') return featureSettings.carouselsLists !== false;
                        if (pageKey === 'manga') return featureSettings.carouselsManga !== false;

                        return featureSettings.carouselsOther !== false;
                    }

                    function areExtensionCarouselsEnabled() {
                        return !!featureSettings.betterMarketplace && !!featureSettings.carousels && featureSettings.carouselsExtensions !== false;
                    }

                    function setBodyFlags() {
                        if (!document.body) return;

                        document.body.setAttribute('data-ama-better-marketplace', String(!!featureSettings.betterMarketplace));
                        document.body.setAttribute('data-ama-better-marketplace-page', String(isBetterMarketplacePage()));
                        document.body.setAttribute('data-ama-carousels', String(!!featureSettings.carousels));
                        document.body.setAttribute('data-ama-carousel-page', getCarouselPageKey());
                        document.body.setAttribute('data-ama-carousels-active', String(areCarouselsEnabledForCurrentPage()));
                        document.body.setAttribute('data-ama-subdub-icons', String(!!featureSettings.subDubIcons));
                        document.body.setAttribute('data-ama-hide-file-names', String(!!featureSettings.hideFileNames));
                        document.body.setAttribute('data-ama-bas1874-marketplace', String(!!featureSettings.useBas1874Marketplace));
                    }

                    function escapeHtml(value) {
                        return String(value || '')
                            .replace(/&/g, '&amp;')
                            .replace(/</g, '&lt;')
                            .replace(/>/g, '&gt;')
                            .replace(/"/g, '&quot;')
                            .replace(/'/g, '&#039;');
                    }

                    function normalizeId(value) {
                        if (value === null || value === undefined) return '';

                        const id = String(value).trim();

                        if (!id) return '';

                        const numeric = id.match(/\\d+/);

                        return numeric ? numeric[0] : id;
                    }

                    function addId(set, value) {
                        const id = normalizeId(value);
                        if (!id) return;

                        set.add(id);
                    }

                    function isIdKey(key) {
                        const normalized = String(key || '').replace(/[-_\\s]/g, '').toLowerCase();

                        return (
                            normalized === 'id' ||
                            normalized === 'animeid' ||
                            normalized === 'mediaid' ||
                            normalized === 'anilistid' ||
                            normalized === 'aniid' ||
                            normalized === 'malid' ||
                            normalized === 'showid'
                        );
                    }

                    function extractDubIdsFromData(data) {
                        const ids = new Set();

                        if (!data) return ids;

                        if (Array.isArray(data)) {
                            data.forEach(item => {
                                if (typeof item === 'string' || typeof item === 'number') {
                                    addId(ids, item);
                                    return;
                                }

                                if (!item || typeof item !== 'object') return;

                                Object.keys(item).forEach(key => {
                                    if (isIdKey(key)) addId(ids, item[key]);
                                });
                            });

                            return ids;
                        }

                        if (typeof data === 'object') {
                            Object.keys(data).forEach(key => {
                                if (/^\\d+$/.test(key)) {
                                    addId(ids, key);
                                }

                                const value = data[key];

                                if (value && typeof value === 'object' && !Array.isArray(value)) {
                                    Object.keys(value).forEach(innerKey => {
                                        if (isIdKey(innerKey)) addId(ids, value[innerKey]);
                                    });
                                }

                                if (Array.isArray(value)) {
                                    value.forEach(item => {
                                        if (!item || typeof item !== 'object') return;

                                        Object.keys(item).forEach(innerKey => {
                                            if (isIdKey(innerKey)) addId(ids, item[innerKey]);
                                        });
                                    });
                                }
                            });
                        }

                        return ids;
                    }

                    function splitTextLines(text) {
                        return String(text || '').split(String.fromCharCode(13)).join('').split(String.fromCharCode(10));
                    }

                    function collectIdsFromJsonValue(value, ids, depth) {
                        if (value === null || value === undefined || depth > 8) return;

                        if (typeof value === 'string' || typeof value === 'number') {
                            addId(ids, value);
                            return;
                        }

                        if (Array.isArray(value)) {
                            value.forEach(item => collectIdsFromJsonValue(item, ids, depth + 1));
                            return;
                        }

                        if (typeof value !== 'object') return;

                        Object.keys(value).forEach(key => {
                            const normalizedKey = String(key || '').replace(/[-_\\s]/g, '').toLowerCase();
                            const child = value[key];

                            if (/^\\d+$/.test(String(key)) && child) {
                                addId(ids, key);
                                return;
                            }

                            if (
                                normalizedKey === 'id' ||
                                normalizedKey === 'mal' ||
                                normalizedKey === 'malid' ||
                                normalizedKey === 'myanimelistid' ||
                                normalizedKey === 'animeid'
                            ) {
                                addId(ids, child);
                                return;
                            }

                            collectIdsFromJsonValue(child, ids, depth + 1);
                        });
                    }

                    function extractDubMalIdsFromJsonText(text) {
                        const ids = new Set();
                        if (!text) return ids;

                        try {
                            collectIdsFromJsonValue(JSON.parse(String(text)), ids, 0);
                            return ids;
                        } catch (_) {}

                        splitTextLines(text).forEach(line => {
                            let cleaned = String(line || '').trim();
                            ['[', ']', ',', '"', "'"].forEach(ch => {
                                cleaned = cleaned.split(ch).join('');
                            });

                            if (/^\\d+$/.test(cleaned)) addId(ids, cleaned);
                        });

                        return ids;
                    }

                    function extractMalToAnilistMapFromJsonl(text) {
                        const map = new Map();
                        if (!text) return map;

                        splitTextLines(text).forEach(line => {
                            const trimmed = String(line || '').trim();
                            if (!trimmed || trimmed[0] !== '{') return;

                            try {
                                const item = JSON.parse(trimmed);
                                if (!item || typeof item !== 'object') return;

                                const malId = normalizeId(item.mal_id || item.malId || item.malID || item.myAnimeListId || item.idMal);
                                const anilistId = normalizeId(item.anilist_id || item.anilistId || item.anilistID || item.aniListId || item.idAniList);

                                if (malId && anilistId) map.set(malId, anilistId);
                            } catch (_) {}
                        });

                        return map;
                    }

                    function loadMalToAnilistMap() {
                        if (malToAnilistMapPromise) return malToAnilistMapPromise;

                        malToAnilistMapPromise = fetchTextWithTimeout(dubMappingUrl, 25000)
                            .then(text => extractMalToAnilistMapFromJsonl(text))
                            .catch(() => new Map());

                        return malToAnilistMapPromise;
                    }

                    function getCachedDubIds() {
                        try {
                            const raw = window.localStorage.getItem(dubFeedCacheKey);
                            if (!raw) return null;

                            const parsed = JSON.parse(raw);

                            if (!parsed || typeof parsed !== 'object') return null;
                            if (!Array.isArray(parsed.ids)) return null;
                            if (!parsed.time || Date.now() - parsed.time > dubFeedCacheTTL) return null;

                            return new Set(parsed.ids.map(normalizeId).filter(Boolean));
                        } catch (_) {
                            return null;
                        }
                    }

                    function setCachedDubIds(ids) {
                        try {
                            if (!ids || ids.size === 0) {
                                window.localStorage.removeItem(dubFeedCacheKey);
                                return;
                            }

                            window.localStorage.setItem(dubFeedCacheKey, JSON.stringify({
                                time: Date.now(),
                                ids: Array.from(ids)
                            }));
                        } catch (_) {}
                    }

                    async function fetchJsonWithTimeout(url, timeoutMs) {
                        const controller = new AbortController();
                        const timeout = setTimeout(() => controller.abort(), timeoutMs);

                        try {
                            const response = await fetch(url, {
                                signal: controller.signal,
                                cache: 'force-cache'
                            });

                            if (!response.ok) return null;

                            return await response.json();
                        } catch (_) {
                            return null;
                        } finally {
                            clearTimeout(timeout);
                        }
                    }

                    async function fetchTextWithTimeout(url, timeoutMs) {
                        const controller = new AbortController();
                        const timeout = setTimeout(() => controller.abort(), timeoutMs);

                        try {
                            const response = await fetch(url, {
                                signal: controller.signal,
                                cache: 'no-store'
                            });

                            if (!response.ok) return '';

                            return await response.text();
                        } catch (_) {
                            return '';
                        } finally {
                            clearTimeout(timeout);
                        }
                    }

                    function getAsunaTracksScheduleUrlFromEndpoint(endpoint) {
                        try {
                            if (!String(endpoint || '').startsWith('/timetables/dub')) return '';

                            const params = new URLSearchParams(String(endpoint).split('?')[1] || '');
                            const next = new URLSearchParams();
                            next.set('type', 'dub');

                            const year = params.get('year');
                            const week = params.get('week');
                            const timezone = params.get('tz') || params.get('timezone');

                            if (year) next.set('year', year);
                            if (week) next.set('week', week);
                            if (timezone) next.set('tz', timezone);

                            return asunaTracksScheduleApiBaseUrl + '/anime-schedule?' + next.toString();
                        } catch (_) {
                            return '';
                        }
                    }

                    async function fetchJsonRequestCandidates(candidates, endpoint, timeoutMs) {
                        let lastStatus = '';
                        let lastStatusText = '';
                        let lastMessage = '';

                        for (const candidate of candidates) {
                            const controller = new AbortController();
                            const timeout = setTimeout(() => controller.abort(), timeoutMs);

                            try {
                                const response = await fetch(candidate.url, {
                                    signal: controller.signal,
                                    cache: 'no-store',
                                    credentials: candidate.credentials || 'include',
                                    headers: candidate.headers || { Accept: 'application/json' }
                                });
                                syncSeanimeIdentity(response);

                                if (!response.ok) {
                                    lastStatus = 'http-' + response.status;
                                    lastStatusText = response.statusText || '';
                                    continue;
                                }

                                const json = await response.json();
                                window.__AMA_ANIME_SCHEDULE_API_STATUS__ = {
                                    endpoint,
                                    source: candidate.source || 'api',
                                    status: 'loaded',
                                    proxied: !!candidate.proxied,
                                    rows: Array.isArray(json) ? json.length : (json && Array.isArray(json.items) ? json.items.length : (json && typeof json === 'object' ? Object.keys(json).length : 0)),
                                };
                                return json;
                            } catch (error) {
                                lastStatus = 'failed';
                                lastMessage = error && error.message ? error.message : String(error || '');
                            } finally {
                                clearTimeout(timeout);
                            }
                        }

                        window.__AMA_ANIME_SCHEDULE_API_STATUS__ = {
                            endpoint,
                            source: candidates[0] && candidates[0].source ? candidates[0].source : 'api',
                            status: lastStatus || 'failed',
                            statusText: lastStatusText,
                            message: lastMessage,
                            proxied: candidates.some(candidate => candidate.proxied),
                        };

                        return null;
                    }

                    async function fetchAnimeScheduleApiJson(endpoint, token, timeoutMs) {
                        const cleanToken = normalizeAnimeScheduleApiToken(token);
                        const asunaTracksScheduleUrl = getAsunaTracksScheduleUrlFromEndpoint(endpoint);
                        window.__AMA_ANIME_SCHEDULE_API_STATUS__ = {
                            endpoint,
                            status: 'starting',
                            source: asunaTracksScheduleUrl ? 'asunatracks' : 'animeschedule',
                            proxied: !!cleanToken || !!asunaTracksScheduleUrl,
                        };

                        if (asunaTracksScheduleUrl) {
                            const headers = {
                                Accept: 'application/json'
                            };
                            const proxiedUrl = '/api/v1/proxy?' + new URLSearchParams({
                                url: asunaTracksScheduleUrl,
                                headers: JSON.stringify(headers)
                            }).toString();

                            return fetchJsonRequestCandidates([
                                {
                                    url: asunaTracksScheduleUrl,
                                    headers,
                                    credentials: 'omit',
                                    proxied: false,
                                    source: 'asunatracks'
                                },
                                {
                                    url: proxiedUrl,
                                    headers: getSeanimeHeaders({ Accept: 'application/json' }),
                                    credentials: 'include',
                                    proxied: true,
                                    source: 'asunatracks'
                                }
                            ], endpoint, timeoutMs);
                        }

                        const controller = new AbortController();
                        const timeout = setTimeout(() => controller.abort(), timeoutMs);

                        try {
                            const headers = {
                                Accept: 'application/json'
                            };

                            let requestUrl = animeScheduleApiBaseUrl + endpoint;

                            if (cleanToken) {
                                headers.Authorization = 'Bearer ' + cleanToken;
                                requestUrl = '/api/v1/proxy?' + new URLSearchParams({
                                    url: requestUrl,
                                    headers: JSON.stringify(headers)
                                }).toString();
                            }

                            const response = await fetch(requestUrl, {
                                signal: controller.signal,
                                cache: 'no-store',
                                credentials: 'include',
                                headers: cleanToken ? getSeanimeHeaders({ Accept: 'application/json' }) : headers
                            });
                            syncSeanimeIdentity(response);

                            if (!response.ok) {
                                window.__AMA_ANIME_SCHEDULE_API_STATUS__ = {
                                    endpoint,
                                    status: 'http-' + response.status,
                                    statusText: response.statusText || '',
                                    proxied: !!cleanToken,
                                };
                                return null;
                            }

                            const json = await response.json();
                            window.__AMA_ANIME_SCHEDULE_API_STATUS__ = {
                                endpoint,
                                status: 'loaded',
                                proxied: !!cleanToken,
                                rows: Array.isArray(json) ? json.length : (json && typeof json === 'object' ? Object.keys(json).length : 0),
                            };
                            return json;
                        } catch (error) {
                            window.__AMA_ANIME_SCHEDULE_API_STATUS__ = {
                                endpoint,
                                status: 'failed',
                                message: error && error.message ? error.message : String(error || ''),
                                proxied: !!cleanToken,
                            };
                            return null;
                        } finally {
                            clearTimeout(timeout);
                        }
                    }

                    function normalizeAnimeScheduleApiToken(value) {
                        return String(value || '')
                            .trim()
                            .replace(/^["']|["']$/g, '')
                            .replace(/^authorization\s*:\s*/i, '')
                            .replace(/^bearer\s+/i, '')
                            .trim();
                    }

                    function getAnimeScheduleApiToken() {
                        return '';
                    }

                    function syncAnimeScheduleApiTokenToDom() {
                        try {
                            if (document.body) {
                                document.body.removeAttribute('data-ama-anime-schedule-api-token');
                            }
                        } catch (_) {}
                    }

                    function normalizeComparableTitle(value) {
                        return String(value || '')
                            .toLowerCase()
                            .replace(/&amp;/g, '&')
                            .replace(/[^a-z0-9]+/g, ' ')
                            .trim()
                            .replace(/\\s+/g, ' ');
                    }

                    function normalizeEpisodeNumber(value) {
                        const text = String(value === null || value === undefined ? '' : value).trim();
                        if (!text) return '';

                        const numeric = Number(text);
                        if (!Number.isNaN(numeric) && Number.isFinite(numeric)) {
                            return String(numeric).replace(/\\.0$/, '');
                        }

                        return text;
                    }

                    function getIsoWeekInfo(date) {
                        const source = date instanceof Date && !Number.isNaN(date.getTime()) ? date : new Date();
                        const utc = new Date(Date.UTC(source.getFullYear(), source.getMonth(), source.getDate()));
                        const day = utc.getUTCDay() || 7;
                        utc.setUTCDate(utc.getUTCDate() + 4 - day);
                        const yearStart = new Date(Date.UTC(utc.getUTCFullYear(), 0, 1));
                        const week = Math.ceil((((utc.getTime() - yearStart.getTime()) / 86400000) + 1) / 7);

                        return {
                            year: utc.getUTCFullYear(),
                            week
                        };
                    }

                    function parseAnimeScheduleDate(value) {
                        const text = String(value || '').trim();
                        if (!text || text.indexOf('0001-01-01') === 0) return null;

                        const date = new Date(text);
                        return date && !Number.isNaN(date.getTime()) ? date : null;
                    }

                    function getWholeWeeksBetween(startDate, endDate) {
                        if (!(startDate instanceof Date) || !(endDate instanceof Date)) return -1;
                        if (Number.isNaN(startDate.getTime()) || Number.isNaN(endDate.getTime())) return -1;

                        const diff = endDate.getTime() - startDate.getTime();
                        if (diff < 0) return -1;

                        return Math.floor(diff / (7 * 24 * 60 * 60 * 1000));
                    }

                    function flattenObjects(value, result, depth) {
                        if (!value || depth > 8) return;

                        if (Array.isArray(value)) {
                            value.forEach(item => flattenObjects(item, result, depth + 1));
                            return;
                        }

                        if (typeof value !== 'object') return;

                        result.push(value);

                        Object.keys(value).forEach(key => {
                            const child = value[key];
                            if (child && typeof child === 'object') flattenObjects(child, result, depth + 1);
                        });
                    }

                    function getAnimeScheduleEntryTitleCandidates(entry) {
                        if (!entry || typeof entry !== 'object') return [];

                        const names = entry.names || entry.name || {};
                        const candidates = [
                            entry.title,
                            entry.name,
                            entry.animeTitle,
                            entry.anime_title,
                            entry.mediaTitle,
                            entry.media_title,
                            entry.english,
                            entry.romaji,
                            entry.native,
                            entry.title_english,
                            entry.title_romaji,
                            entry.title_native,
                            entry.media && entry.media.title_english,
                            entry.media && entry.media.title,
                            entry.media && entry.media.name,
                            names.title,
                            names.english,
                            names.romaji,
                            names.native,
                            Array.isArray(names.synonyms) ? names.synonyms[0] : '',
                            entry.route
                        ];

                        return candidates.map(value => {
                            if (value && typeof value === 'object') {
                                return String(value.english || value.romaji || value.native || value.title || value.name || '').trim();
                            }

                            return String(value || '').trim();
                        }).filter(Boolean);
                    }

                    function getAnimeScheduleEntryTitle(entry) {
                        const candidates = getAnimeScheduleEntryTitleCandidates(entry);

                        return candidates[0] || '';
                    }

                    function animeScheduleTitlesMatch(entry, title) {
                        const scheduleTitle = normalizeComparableTitle(title);
                        if (!scheduleTitle) return false;

                        const entryTitles = getAnimeScheduleEntryTitleCandidates(entry)
                            .map(normalizeComparableTitle)
                            .filter(Boolean);

                        if (entryTitles.some(entryTitle => entryTitle === scheduleTitle)) return true;

                        const scheduleWords = new Set(scheduleTitle.split(' ').filter(word => word.length > 2));
                        if (scheduleWords.size >= 3) {
                            const hasStrongWordMatch = entryTitles.some(entryTitle => {
                                const entryWords = entryTitle.split(' ').filter(word => word.length > 2);
                                if (entryWords.length < 3) return false;

                                let shared = 0;
                                entryWords.forEach(word => {
                                    if (scheduleWords.has(word)) shared += 1;
                                });

                                return shared >= Math.min(4, Math.ceil(Math.min(scheduleWords.size, entryWords.length) * 0.75));
                            });

                            if (hasStrongWordMatch) return true;
                        }

                        return entryTitles.some(entryTitle => {
                            if (entryTitle.length < 18 || scheduleTitle.length < 18) return false;

                            const longer = entryTitle.length >= scheduleTitle.length ? entryTitle : scheduleTitle;
                            const shorter = entryTitle.length < scheduleTitle.length ? entryTitle : scheduleTitle;
                            if (shorter.length / longer.length < 0.85) return false;

                            return longer.includes(shorter);
                        });
                    }

                    function getAnimeScheduleEntryAnilistId(entry) {
                        if (!entry || typeof entry !== 'object') return '';

                        const mapped = normalizeId(entry.__amaAnilistId || entry.amaAnilistId);
                        if (mapped) return mapped;

                        const numericId = normalizeId(entry.id);
                        if (numericId && (entry.idMal || entry.episode || entry.format || entry.duration)) return numericId;

                        const direct = normalizeId(
                            entry.anilistId ||
                            entry.anilistID ||
                            entry.aniListId ||
                            entry.anilist_id ||
                            entry.ani_list_id ||
                            (entry.external_ids && (entry.external_ids.anilist || entry.external_ids.aniList || entry.external_ids.anilist_id)) ||
                            (entry.externalIds && (entry.externalIds.anilist || entry.externalIds.aniList || entry.externalIds.anilist_id)) ||
                            (entry.media && entry.media.external_ids && (entry.media.external_ids.anilist || entry.media.external_ids.aniList || entry.media.external_ids.anilist_id)) ||
                            (entry.media && entry.media.externalIds && (entry.media.externalIds.anilist || entry.media.externalIds.aniList || entry.media.externalIds.anilist_id)) ||
                            entry.mediaId ||
                            entry.media_id ||
                            entry.idAniList ||
                            entry.id_anilist
                        );

                        if (direct) return direct;

                        const websites = entry.websites || entry.website || {};
                        const websiteValue = websites.aniList || websites.anilist || websites.ani_list || entry.aniList || entry.anilist || '';
                        const websiteId = normalizeId(websiteValue) || getAnilistIdFromHref(websiteValue);

                        if (websiteId) return websiteId;

                        return readAnilistIdFromObject(entry, 0, new WeakSet(), false);
                    }

                    function getAnimeScheduleEntryMalId(entry) {
                        if (!entry || typeof entry !== 'object') return '';

                        return normalizeId(
                            entry.mal_id ||
                            entry.malId ||
                            entry.malID ||
                            entry.myAnimeListId ||
                            entry.idMal ||
                            (entry.media && (entry.media.mal_id || entry.media.malId || entry.media.malID || entry.media.myAnimeListId || entry.media.idMal))
                        );
                    }

                    function getAsunaTracksScheduleEntryAnilistId(entry) {
                        if (!entry || typeof entry !== 'object') return '';

                        return normalizeId(
                            entry.__amaAnilistId ||
                            entry.amaAnilistId ||
                            (entry.external_ids && (entry.external_ids.anilist || entry.external_ids.aniList || entry.external_ids.anilist_id)) ||
                            (entry.externalIds && (entry.externalIds.anilist || entry.externalIds.aniList || entry.externalIds.anilist_id)) ||
                            (entry.media && entry.media.external_ids && (entry.media.external_ids.anilist || entry.media.external_ids.aniList || entry.media.external_ids.anilist_id)) ||
                            (entry.media && entry.media.externalIds && (entry.media.externalIds.anilist || entry.media.externalIds.aniList || entry.media.externalIds.anilist_id))
                        );
                    }

                    function getEpisodeNumberFromValue(value) {
                        const text = String(value === null || value === undefined ? '' : value).trim();
                        if (!text) return '';

                        const direct = text.match(/^\\d+(?:\\.\\d+)?$/);
                        if (direct) return normalizeEpisodeNumber(direct[0]);

                        const patterns = [
                            /(?:episode|ep\\.?|e)\\s*#?\\s*(\\d+(?:\\.\\d+)?)/i,
                            /#\\s*(\\d+(?:\\.\\d+)?)/,
                            /\\b(\\d+(?:\\.\\d+)?)\\s*(?:sub|dub|raw)?\\s*(?:aired|airing|episode|ep)\\b/i
                        ];

                        for (const pattern of patterns) {
                            const match = text.match(pattern);
                            if (match && match[1]) return normalizeEpisodeNumber(match[1]);
                        }

                        return '';
                    }

                    function getAnimeScheduleEntryEpisodeNumber(entry) {
                        if (!entry || typeof entry !== 'object') return '';

                        const values = [
                            entry.episodeNumber,
                            entry.episode_number,
                            entry.episode,
                            entry.episodeNum,
                            entry.episode_num,
                            entry.episode_no,
                            entry.number,
                            entry.dubEpisode,
                            entry.dub_episode,
                            entry.dubEpisodeNumber,
                            entry.dub_episode_number,
                            entry.latestDubEpisode,
                            entry.latest_dub_episode,
                            entry.latestDubEpisodeNumber,
                            entry.latest_dub_episode_number,
                            entry.lastDubEpisode,
                            entry.last_dub_episode,
                            entry.lastDubEpisodeNumber,
                            entry.last_dub_episode_number,
                            entry.dub,
                            entry.dubbedEpisode,
                            entry.dubbed_episode,
                            entry.dubbedEpisodeNumber,
                            entry.dubbed_episode_number,
                            entry.currentEpisode,
                            entry.current_episode,
                            entry.airingEpisode,
                            entry.airing_episode,
                            entry.episode && entry.episode.aired,
                            entry.episode && entry.episode.number,
                            entry.episode && entry.episode.episodeNumber,
                            entry.episode && entry.episode.episode_number
                        ];

                        for (const value of values) {
                            const episode = getEpisodeNumberFromValue(value);
                            if (episode) return episode;
                        }

                        return '';
                    }

                    function isAnimeScheduleDubEntry(entry) {
                        if (!entry || typeof entry !== 'object') return true;

                        const airType = String(entry.airType || entry.air_type || entry.type || entry.releaseType || entry.release_type || entry.subOrDub || entry.sub_or_dub || entry.audio || '').toLowerCase();
                        if (airType) return airType.includes('dub');

                        const text = String(entry.title || entry.name || entry.label || '').toLowerCase();
                        return !text || text.includes('dub');
                    }

                    function getScheduleEpisodeNumber(event) {
                        if (!event) return '';

                        const attrNames = [
                            'data-episode-number',
                            'data-schedule-episode-number',
                            'data-schedule-calendar-event-item-episode-number',
                            'data-episode',
                            'data-ep'
                        ];

                        const attrEpisode = getEpisodeNumberFromValue(getAttributeAnyCase(event, attrNames));
                        if (attrEpisode) return attrEpisode;

                        const data = event.dataset || {};
                        const datasetNames = ['episodeNumber', 'scheduleEpisodeNumber', 'scheduleCalendarEventItemEpisodeNumber', 'episode', 'ep'];

                        for (const name of datasetNames) {
                            const episode = getEpisodeNumberFromValue(data[name]);
                            if (episode) return episode;
                        }

                        const textTargets = [event];
                        if (event.querySelectorAll) {
                            event.querySelectorAll('[data-schedule-calendar-event-item-episode], [data-schedule-calendar-event-item-episode-time], [data-schedule-calendar-mobile-list-day-item-event-episode], [data-schedule-calendar-mobile-list-day-item-event-episode-time], [data-schedule-event-episode], [data-episode-number], [data-episode], span, p, div').forEach(el => {
                                if (textTargets.length < 80) textTargets.push(el);
                            });
                        }

                        for (const target of textTargets) {
                            const episode = getEpisodeNumberFromValue(target.textContent || '');
                            if (episode) return episode;
                        }

                        return '';
                    }

                    function getScheduleDate(event) {
                        if (!event) return null;

                        const time = event.querySelector && event.querySelector('time[datetime], [datetime]');
                        const datetime = time && time.getAttribute ? time.getAttribute('datetime') : '';
                        const date = datetime ? new Date(datetime) : null;

                        if (date && !Number.isNaN(date.getTime())) return date;

                        const datedParent = event.closest && event.closest('[datetime], [data-date], [data-day], [data-schedule-date], [data-schedule-day], [data-calendar-date]');
                        if (datedParent && datedParent.getAttribute) {
                            const attrDate = getAttributeAnyCase(datedParent, ['datetime', 'data-date', 'data-day', 'data-schedule-date', 'data-schedule-day', 'data-calendar-date']);
                            const parentDate = parseAnimeScheduleDate(attrDate);
                            if (parentDate) return parentDate;
                        }

                        return null;
                    }

                    function getScheduleTitle(event) {
                        if (!event) return '';

                        const titleTarget = event.querySelector && event.querySelector('[data-schedule-calendar-event-item-name], [data-schedule-calendar-event-item-title], [data-schedule-calendar-mobile-list-day-item-event-text], [data-schedule-calendar-mobile-list-day-item-event-image], [data-schedule-event-title], [data-media-title], [title], a[href]');
                        const candidates = [
                            titleTarget && titleTarget.getAttribute && titleTarget.getAttribute('title'),
                            titleTarget && titleTarget.getAttribute && titleTarget.getAttribute('alt'),
                            titleTarget && titleTarget.textContent,
                            event.getAttribute && event.getAttribute('title'),
                            event.textContent
                        ];

                        return candidates.map(value => String(value || '').trim()).find(Boolean) || '';
                    }

                    function datesAreSameLocalDay(a, b) {
                        if (!(a instanceof Date) || !(b instanceof Date)) return false;
                        if (Number.isNaN(a.getTime()) || Number.isNaN(b.getTime())) return false;

                        return a.getFullYear() === b.getFullYear() &&
                            a.getMonth() === b.getMonth() &&
                            a.getDate() === b.getDate();
                    }

                    function canUseScheduleFallback(scheduleDate) {
                        if (!(scheduleDate instanceof Date) || Number.isNaN(scheduleDate.getTime())) return false;

                        const currentWeekStart = new Date();
                        currentWeekStart.setHours(0, 0, 0, 0);
                        const dayOffset = (currentWeekStart.getDay() + 6) % 7;
                        currentWeekStart.setDate(currentWeekStart.getDate() - dayOffset);

                        const scheduleDay = new Date(scheduleDate.getTime());
                        scheduleDay.setHours(0, 0, 0, 0);

                        return scheduleDay.getTime() < currentWeekStart.getTime();
                    }

                    function getAnimeScheduleEntryDate(entry) {
                        if (!entry || typeof entry !== 'object') return null;

                        return parseAnimeScheduleDate(
                            entry.episodeDate ||
                            entry.episode_date ||
                            entry.date ||
                            entry.datetime ||
                            entry.airingAt ||
                            entry.airing_at ||
                            entry.airedAt ||
                            entry.aired_at ||
                            entry.releaseTime ||
                            entry.release_time ||
                            entry.startsAt ||
                            entry.starts_at ||
                            entry.startTime ||
                            entry.start_time ||
                            (entry.episode && (entry.episode.airedAt || entry.episode.aired_at || entry.episode.date || entry.episode.datetime || entry.episode.episode_date))
                        );
                    }

                    function animeScheduleEntryMatchesSchedule(entry, anilistId, title, episodeNumber, scheduleDate) {
                        if (!entry) return false;
                        if (!isAnimeScheduleDubEntry(entry)) return false;

                        const entryEpisode = getAnimeScheduleEntryEpisodeNumber(entry);
                        if (!entryEpisode) return false;

                        const entryDate = getAnimeScheduleEntryDate(entry);
                        if (entryDate && scheduleDate && !datesAreSameLocalDay(entryDate, scheduleDate)) return false;

                        const entryAnilistId = getAnimeScheduleEntryAnilistId(entry);
                        if (anilistId && entryAnilistId && String(entryAnilistId) === String(anilistId)) return true;

                        return animeScheduleTitlesMatch(entry, title);
                    }

                    function asunaTracksScheduleEntryMatchesSchedule(entry, anilistId, title, scheduleDate) {
                        if (!entry || typeof entry !== 'object') return false;
                        if (!isAnimeScheduleDubEntry(entry)) return false;
                        if (!getAnimeScheduleEntryEpisodeNumber(entry)) return false;

                        const entryAnilistId = getAsunaTracksScheduleEntryAnilistId(entry);
                        if (anilistId && entryAnilistId && String(entryAnilistId) === String(anilistId)) return true;

                        const entryDate = getAnimeScheduleEntryDate(entry);
                        if (entryDate && scheduleDate && !datesAreSameLocalDay(entryDate, scheduleDate)) return false;

                        return animeScheduleTitlesMatch(entry, title);
                    }

                    function animeScheduleEntryMatchesAnime(entry, anilistId, title) {
                        if (!entry) return false;
                        if (!isAnimeScheduleDubEntry(entry)) return false;

                        const entryAnilistId = getAnimeScheduleEntryAnilistId(entry);
                        if (anilistId && entryAnilistId && String(entryAnilistId) === String(anilistId)) return true;

                        return animeScheduleTitlesMatch(entry, title);
                    }

                    function getBestDubEntryForAnime(entries, anilistId, title) {
                        if (!Array.isArray(entries)) return null;

                        const matches = entries.filter(entry => animeScheduleEntryMatchesAnime(entry, anilistId, title));
                        if (!matches.length) return null;

                        matches.sort((a, b) => {
                            const aEpisode = Number(getAnimeScheduleEntryEpisodeNumber(a) || 0);
                            const bEpisode = Number(getAnimeScheduleEntryEpisodeNumber(b) || 0);

                            return bEpisode - aEpisode;
                        });

                        return matches[0] || null;
                    }

                    function getScheduleDubMatchDetails(entry, source) {
                        if (!entry) return null;

                        const episodeNumber = getAnimeScheduleEntryEpisodeNumber(entry);
                        const title = getAnimeScheduleEntryTitle(entry);
                        const episodeDate = entry.episodeDate || entry.episode_date || entry.date || entry.datetime || entry.airingAt || entry.airing_at || entry.airedAt || entry.aired_at || (entry.episode && (entry.episode.airedAt || entry.episode.aired_at || entry.episode.date || entry.episode.datetime || entry.episode.episode_date)) || '';

                        return {
                            source,
                            episodeNumber,
                            title,
                            episodeDate,
                            exact: source === 'api' || source === 'fallback'
                        };
                    }

                    function getAnimeScheduleAnimeItems(data) {
                        if (!data) return [];

                        if (Array.isArray(data)) return data;
                        if (Array.isArray(data.anime)) return data.anime;
                        if (Array.isArray(data.items)) return data.items;
                        if (Array.isArray(data.results)) return data.results;
                        if (Array.isArray(data.data)) return data.data;

                        return [data];
                    }

                    function loadAnimeScheduleAnimeByAnilistId(anilistId) {
                        const token = getAnimeScheduleApiToken();
                        const id = normalizeId(anilistId);
                        if (!id) return Promise.resolve(null);

                        if (animeScheduleAnimeDetailsPromises.has(id)) return animeScheduleAnimeDetailsPromises.get(id);

                        const promise = fetchAnimeScheduleApiJson('/anime?anilist-ids=' + encodeURIComponent(id), token, 9000).then(data => {
                            const items = getAnimeScheduleAnimeItems(data);
                            return items.find(item => {
                                return item && String(getAnimeScheduleEntryAnilistId(item)) === String(id);
                            }) || items[0] || null;
                        }).catch(() => null);

                        animeScheduleAnimeDetailsPromises.set(id, promise);

                        return promise;
                    }

                    function getDubDetailsFromAnimeScheduleAnime(anime, scheduleDate, title) {
                        if (!anime || typeof anime !== 'object') return null;

                        const dubPremier = parseAnimeScheduleDate(anime.dubPremier || anime.dubPremierDate || anime.dubStartDate);
                        if (!dubPremier) return null;

                        const targetDate = scheduleDate instanceof Date && !Number.isNaN(scheduleDate.getTime()) ? scheduleDate : new Date();
                        const weeks = getWholeWeeksBetween(dubPremier, targetDate);
                        if (weeks < 0) return null;

                        const episodeNumber = String(weeks + 1);
                        const totalEpisodes = Number(anime.episodes || 0);
                        if (totalEpisodes && Number(episodeNumber) > totalEpisodes) return null;

                        return {
                            source: 'anime-api',
                            episodeNumber,
                            title: getAnimeScheduleEntryTitle(anime) || title,
                            episodeDate: dubPremier.toISOString(),
                            exact: false
                        };
                    }

                    function normalizeAnimeScheduleEntries(data) {
                        const objects = [];
                        flattenObjects(data, objects, 0);

                        return objects.filter(entry => {
                            return getAnimeScheduleEntryEpisodeNumber(entry) && (getAnimeScheduleEntryAnilistId(entry) || getAnimeScheduleEntryTitle(entry));
                        });
                    }

                    function annotateAnimeScheduleEntriesWithAnilistIds(entries, malToAnilist) {
                        if (!Array.isArray(entries) || !malToAnilist || !malToAnilist.size) return entries;

                        entries.forEach(entry => {
                            if (!entry || typeof entry !== 'object' || getAsunaTracksScheduleEntryAnilistId(entry)) return;

                            const malId = getAnimeScheduleEntryMalId(entry);
                            const anilistId = malId ? malToAnilist.get(malId) : '';
                            if (anilistId) entry.__amaAnilistId = anilistId;
                        });

                        return entries;
                    }

                    async function loadAnimeScheduleApiDubEntries(scheduleDate) {
                        const token = getAnimeScheduleApiToken();

                        const timezone = (() => {
                            try {
                                return Intl.DateTimeFormat().resolvedOptions().timeZone || 'Europe/Oslo';
                            } catch (_) {
                                return 'Europe/Oslo';
                            }
                        })();
                        const weekInfo = getIsoWeekInfo(scheduleDate);
                        const cacheKey = weekInfo.year + '-w' + weekInfo.week + '-' + timezone + '-' + (token ? 'token' : 'public');

                        if (animeScheduleDubApiEntriesPromises.has(cacheKey)) return animeScheduleDubApiEntriesPromises.get(cacheKey);

                        const promise = (async () => {
                            const data = await fetchAnimeScheduleApiJson('/timetables/dub?year=' + encodeURIComponent(weekInfo.year) + '&week=' + encodeURIComponent(weekInfo.week) + '&tz=' + encodeURIComponent(timezone), token, 25000);
                            if (!data) return null;

                            const entries = normalizeAnimeScheduleEntries(data);
                            const malToAnilist = await loadMalToAnilistMap();

                            return annotateAnimeScheduleEntriesWithAnilistIds(entries, malToAnilist);
                        })();

                        animeScheduleDubApiEntriesPromises.set(cacheKey, promise);

                        return promise;
                    }

                    async function loadAnimeScheduleApiDubEntriesNear(scheduleDate) {
                        if (scheduleDate instanceof Date && !Number.isNaN(scheduleDate.getTime())) {
                            return loadAnimeScheduleApiDubEntries(scheduleDate);
                        }

                        const base = new Date();
                        const dates = [-7, 0, 7, 14].map(offsetDays => {
                            const date = new Date(base.getTime());
                            date.setDate(date.getDate() + offsetDays);
                            return date;
                        });

                        const groups = await Promise.all(dates.map(date => loadAnimeScheduleApiDubEntries(date).catch(() => null)));
                        const merged = [];
                        const seen = new Set();

                        groups.forEach(entries => {
                            if (!Array.isArray(entries)) return;

                            entries.forEach(entry => {
                                const key = [
                                    getAnimeScheduleEntryTitle(entry),
                                    getAnimeScheduleEntryEpisodeNumber(entry),
                                    getAnimeScheduleEntryDate(entry) ? getAnimeScheduleEntryDate(entry).toISOString() : '',
                                    getAnimeScheduleEntryMalId(entry),
                                    getAnimeScheduleEntryAnilistId(entry)
                                ].join('|');

                                if (seen.has(key)) return;
                                seen.add(key);
                                merged.push(entry);
                            });
                        });

                        return merged.length ? merged : null;
                    }

                    function loadAnimeScheduleFeedDubEntries() {
                        if (animeScheduleDubFeedEntriesPromise) return animeScheduleDubFeedEntriesPromise;

                        animeScheduleDubFeedEntriesPromise = (async () => {
                            const data = await fetchJsonWithTimeout(dubAniScheduleUrl, 7000);
                            if (!data) return [];

                            return normalizeAnimeScheduleEntries(data);
                        })();

                        return animeScheduleDubFeedEntriesPromise;
                    }

                    function loadDubIds() {
                        if (dubIdSetPromise) return dubIdSetPromise;

                        dubIdSetPromise = (async () => {
                            const cached = getCachedDubIds();
                            if (cached) return cached;

                            const ids = new Set();

                            const malToAnilist = await loadMalToAnilistMap();

                            if (malToAnilist && malToAnilist.size) {
                                const dubTexts = await Promise.all(
                                    dubEnglishSourceUrls.map(url => fetchTextWithTimeout(url, 9000))
                                );

                                dubTexts.forEach(text => {
                                    const malIds = extractDubMalIdsFromJsonText(text);
                                    malIds.forEach(malId => {
                                        const anilistId = malToAnilist.get(malId);
                                        if (anilistId) ids.add(anilistId);
                                    });
                                });
                            }

                            const aniScheduleData = await fetchJsonWithTimeout(dubAniScheduleUrl, 7000);
                            extractDubIdsFromData(aniScheduleData).forEach(id => ids.add(id));

                            setCachedDubIds(ids);

                            return ids;
                        })();

                        return dubIdSetPromise;
                    }

                    window.__AMA_DUB_CLEAR_CACHE__ = function() {
                        try {
                            window.localStorage.removeItem(dubFeedCacheKey);
                        } catch (_) {}

                        dubIdSetPromise = null;
                        animeScheduleDubApiEntriesPromises = new Map();
                        animeScheduleDubFeedEntriesPromise = null;
                        animeScheduleAnimeDetailsPromises.clear();
                        malToAnilistMapPromise = null;
                    };

                    window.__AMA_DUB_RESCAN_SCHEDULE__ = function() {
                        try {
                            document.querySelectorAll('[data-ama-schedule-dub-enhanced]').forEach(event => {
                                delete event.dataset.amaScheduleDubEnhanced;
                            });
                            scheduleRoot(document.body || document.documentElement);
                        } catch (_) {}
                    };

                    function getAttributeAnyCase(element, names) {
                        if (!element || !element.getAttribute) return '';

                        for (const name of names) {
                            const value = element.getAttribute(name);
                            if (value) return value;
                        }

                        const attributes = element.attributes ? Array.from(element.attributes) : [];
                        const lowerNames = names.map(name => String(name).toLowerCase());

                        for (const attr of attributes) {
                            if (lowerNames.includes(String(attr.name || '').toLowerCase())) return attr.value;
                        }

                        return '';
                    }

                    function getStrictIdFromElement(element, allowGenericId) {
                        if (!element) return '';

                        const attributeNames = [
                            'data-media-id',
                            'data-anilist-id',
                            'data-anilistid',
                            'data-ani-list-id',
                            'data-anilist-media-id',
                            'data-media-anilist-id',
                            'data-entry-anilist-id',
                            'data-entry-media-id',
                            'data-anime-id',
                            'data-media-entry-id',
                            'data-mediaentryid',
                            'data-entry-id',
                            'data-schedule-media-id'
                        ];

                        if (allowGenericId) attributeNames.push('data-id');

                        const attrId = normalizeId(getAttributeAnyCase(element, attributeNames));
                        if (attrId) return attrId;

                        const data = element.dataset || {};
                        const datasetNames = [
                            'amaDubAnilistId',
                            'mediaId',
                            'anilistId',
                            'anilistid',
                            'aniListId',
                            'anilistMediaId',
                            'mediaAnilistId',
                            'entryAnilistId',
                            'entryMediaId',
                            'animeId',
                            'mediaEntryId',
                            'entryId',
                            'scheduleMediaId'
                        ];

                        if (allowGenericId) datasetNames.push('id');

                        for (const name of datasetNames) {
                            const id = normalizeId(data[name]);
                            if (id) return id;
                        }

                        return '';
                    }

                    function getAnilistIdFromHref(href) {
                        if (!href) return '';

                        let pathAndQuery = String(href || '');

                        try {
                            const parsed = new URL(pathAndQuery, window.location.origin);
                            pathAndQuery = (parsed.pathname || '') + (parsed.search || '') + (parsed.hash || '');

                            const paramNames = ['anilistId', 'anilist_id', 'aniListId', 'mediaId', 'media_id', 'animeId', 'anime_id', 'id'];
                            for (const name of paramNames) {
                                const id = normalizeId(parsed.searchParams.get(name));
                                if (id) return id;
                            }
                        } catch (_) {}

                        const patterns = [
                            /(?:^|\\/)anime\\/([0-9]+)(?:[\\/\\?#]|$)/i,
                            /(?:^|\\/)anime\\/details\\/([0-9]+)(?:[\\/\\?#]|$)/i,
                            /(?:^|\\/)media\\/(?:anime\\/)?([0-9]+)(?:[\\/\\?#]|$)/i,
                            /(?:^|\\/)entry\\/(?:anime\\/)?([0-9]+)(?:[\\/\\?#]|$)/i,
                            /(?:^|\\/)details\\/(?:anime\\/)?([0-9]+)(?:[\\/\\?#]|$)/i,
                            /[?&#](?:anilistId|anilist_id|aniListId|mediaId|media_id|animeId|anime_id)=([0-9]+)(?:&|#|$)/i
                        ];

                        for (const pattern of patterns) {
                            const match = pathAndQuery.match(pattern);
                            if (match && match[1]) return match[1];
                        }

                        return '';
                    }

                    function objectLooksLikeAniListMediaObject(value) {
                        if (!value || typeof value !== 'object') return false;

                        const mediaType = String(value.type || value.mediaType || value.kind || value.format || value.__typename || '').toLowerCase();
                        if (mediaType.includes('anime')) return true;
                        if (mediaType === 'media' || mediaType === 'baseanime') return true;

                        return !!(value.title || value.coverImage || value.bannerImage || value.season || value.seasonYear || value.episodes || value.nextAiringEpisode);
                    }

                    function readAnilistIdFromObject(value, depth, seen, contextual) {
                        if (!value || depth > 5) return '';
                        if (typeof value === 'string' || typeof value === 'number') return contextual ? normalizeId(value) : '';
                        if (typeof value !== 'object') return '';
                        if (seen.has(value)) return '';

                        seen.add(value);

                        const exactKeys = ['anilistId', 'anilistID', 'aniListId', 'anilist_id', 'ani_list_id', 'mediaId', 'media_id', 'animeId', 'anime_id'];
                        for (const key of exactKeys) {
                            if (!Object.prototype.hasOwnProperty.call(value, key)) continue;
                            const id = normalizeId(value[key]);
                            if (id) return id;
                        }

                        if (Object.prototype.hasOwnProperty.call(value, 'id') && (contextual || objectLooksLikeAniListMediaObject(value))) {
                            const id = normalizeId(value.id);
                            if (id) return id;
                        }

                        for (const key of Object.keys(value)) {
                            const nextContextual = contextual || /anilist|ani[_-]?list|anime|media|entry|card|details|schedule/.test(String(key).toLowerCase());
                            if (!nextContextual && depth >= 3) continue;

                            const id = readAnilistIdFromObject(value[key], depth + 1, seen, nextContextual);
                            if (id) return id;
                        }

                        return '';
                    }

                    function getAnilistIdFromReactData(root) {
                        if (!root) return '';

                        const elements = [root];
                        if (root.querySelectorAll) {
                            root.querySelectorAll('a[href], [data-media-entry-card-body="true"], [data-media-entry-card-title-section-title="true"], [data-media-entry-card-body-image="true"], [data-schedule-calendar-event-item-finale-icon]').forEach(el => elements.push(el));
                        }

                        for (const element of elements) {
                            const keys = Object.keys(element).filter(key => {
                                return key.indexOf('__reactProps$') === 0 ||
                                    key.indexOf('__reactFiber$') === 0 ||
                                    key.indexOf('__reactInternalInstance$') === 0;
                            });

                            for (const key of keys) {
                                const id = readAnilistIdFromObject(element[key], 0, new WeakSet(), false);
                                if (id) return id;
                            }
                        }

                        return '';
                    }

                    function getAnilistIdFromElement(root) {
                        if (!root) return '';

                        const directId = getStrictIdFromElement(root, true);
                        if (directId) return directId;

                        const child = root.querySelector && root.querySelector('[data-media-id], [data-anilist-id], [data-anilistid], [data-ani-list-id], [data-anilist-media-id], [data-media-anilist-id], [data-entry-anilist-id], [data-entry-media-id], [data-anime-id], [data-media-entry-id], [data-mediaentryid], [data-entry-id], [data-id], [data-schedule-media-id]');
                        const childId = getStrictIdFromElement(child, true);
                        if (childId) return childId;

                        const links = root.querySelectorAll ? Array.from(root.querySelectorAll('a[href]')) : [];
                        for (const link of links) {
                            const hrefId = getAnilistIdFromHref(link.getAttribute('href'));
                            if (hrefId) return hrefId;
                        }

                        const closestLink = root.closest ? root.closest('a[href]') : null;
                        const closestHrefId = closestLink ? getAnilistIdFromHref(closestLink.getAttribute('href')) : '';
                        if (closestHrefId) return closestHrefId;

                        return getAnilistIdFromReactData(root);
                    }

                    async function hasDubForAnimeCard(card) {
                        if (!card) return false;

                        const ids = await loadDubIds();

                        const anilistId = getAnilistIdFromElement(card) || normalizeId(card.getAttribute('data-media-id') || card.dataset.amaDubAnilistId || '');
                        const malId = normalizeId(card.getAttribute('data-media-mal-id') || card.dataset.amaDubMalId || '');

                        if (anilistId) card.dataset.amaDubAnilistId = anilistId;

                        return ids.has(anilistId) || ids.has(malId);
                    }

                    async function hasDubForScheduleEvent(event) {
                        if (!event) return false;

                        const details = await getDubForScheduleEvent(event);
                        return !!details;
                    }

                    async function getDubForScheduleEvent(event) {
                        if (!event) return null;

                        const episodeNumber = getScheduleEpisodeNumber(event);
                        if (!episodeNumber) return null;

                        const anilistId = getAnilistIdFromElement(event) || normalizeId(event.getAttribute('data-media-id') || event.dataset.amaDubAnilistId || '');
                        const title = getScheduleTitle(event);
                        const scheduleDate = getScheduleDate(event);

                        if (anilistId) event.dataset.amaDubAnilistId = anilistId;

                        const apiEntries = await loadAnimeScheduleApiDubEntriesNear(scheduleDate);
                        if (Array.isArray(apiEntries)) {
                            const apiMatch = apiEntries.find(entry => asunaTracksScheduleEntryMatchesSchedule(entry, anilistId, title, scheduleDate));
                            if (apiMatch) {
                                event.dataset.amaScheduleDubMatch = 'episode-api';
                                return getScheduleDubMatchDetails(apiMatch, 'api');
                            }
                        }

                        if (canUseScheduleFallback(scheduleDate)) {
                            const fallbackEntries = await loadAnimeScheduleFeedDubEntries();
                            const fallbackMatch = fallbackEntries.find(entry => animeScheduleEntryMatchesSchedule(entry, anilistId, title, episodeNumber, scheduleDate));
                            if (fallbackMatch) {
                                event.dataset.amaScheduleDubMatch = 'episode-fallback';
                                return getScheduleDubMatchDetails(fallbackMatch, 'fallback');
                            }
                        }

                        event.dataset.amaScheduleDubMatch = 'none';
                        return null;
                    }

                    function setAnimeScheduleApiToken(value) {
                        clearAnimeScheduleApiTokenSetting();
                        syncAnimeScheduleApiTokenToDom();
                        animeScheduleDubApiEntriesPromises = new Map();
                        animeScheduleAnimeDetailsPromises.clear();
                        scheduleRoot(document.body || document.documentElement);
                    }

                    function getOpenScheduleSettingsDialog(dialogId) {
                        return null;
                    }

                    function enhanceScheduleSettingsDialog(dialogId) {
                        clearAnimeScheduleApiTokenSetting();
                    }

                    window.__AMA_DUB_DEBUG_SCHEDULE__ = async function() {
                        const seen = new Set();
                        let events = Array.from(document.querySelectorAll(scheduleEventQuery))
                            .map(event => getScheduleEventElement(event))
                            .filter(event => {
                                if (!event || seen.has(event)) return false;
                                seen.add(event);
                                return isVisibleScheduleEvent(event);
                            })
                            .slice(0, 40);

                        if (!events.length) {
                            seen.clear();
                            events = Array.from(document.querySelectorAll(scheduleEventQuery))
                                .map(event => getScheduleEventElement(event))
                                .filter(event => {
                                    if (!event || seen.has(event)) return false;
                                    seen.add(event);
                                    return true;
                                })
                                .slice(0, 40);
                        }

                        return Promise.all(events.map(async event => {
                            const item = getScheduleEventElement(event);
                            const episodeNumber = getScheduleEpisodeNumber(item);
                            const anilistId = getAnilistIdFromElement(item);
                            const title = getScheduleTitle(item);
                            const scheduleDate = getScheduleDate(item);
                            const weekInfo = getIsoWeekInfo(scheduleDate);
                            const apiTokenConfigured = !!getAnimeScheduleApiToken();
                            const apiEntries = await loadAnimeScheduleApiDubEntriesNear(scheduleDate);
                            const fallbackAllowed = canUseScheduleFallback(scheduleDate);
                            const fallbackEntries = fallbackAllowed && !Array.isArray(apiEntries) ? await loadAnimeScheduleFeedDubEntries() : [];
                            const sourceEntries = Array.isArray(apiEntries) ? apiEntries : fallbackEntries;
                            const matchedEntry = Array.isArray(apiEntries)
                                ? sourceEntries.find(entry => asunaTracksScheduleEntryMatchesSchedule(entry, anilistId, title, scheduleDate))
                                : sourceEntries.find(entry => animeScheduleEntryMatchesSchedule(entry, anilistId, title, episodeNumber, scheduleDate));
                            const exactMatched = !!matchedEntry;
                            const latestEntry = getBestDubEntryForAnime(apiEntries, anilistId, title) || getBestDubEntryForAnime(fallbackEntries, anilistId, title);
                            const matchDetails = matchedEntry
                                ? getScheduleDubMatchDetails(matchedEntry, Array.isArray(apiEntries) ? 'api' : 'fallback')
                                : null;

                            return {
                                title,
                                anilistId,
                                episodeNumber,
                                week: weekInfo.week,
                                year: weekInfo.year,
                                apiTokenConfigured,
                                animeScheduleApiStatus: window.__AMA_ANIME_SCHEDULE_API_STATUS__ || null,
                                serverApiStatus: document.body ? document.body.getAttribute('data-ama-server-schedule-api-status') || '' : '',
                                serverApiDetail: document.body ? document.body.getAttribute('data-ama-server-schedule-api-detail') || '' : '',
                                apiStatus: Array.isArray(apiEntries) ? (apiTokenConfigured ? 'loaded-with-token' : 'loaded-public') : (apiTokenConfigured ? 'failed-or-empty' : 'public-failed'),
                                apiEntryCount: Array.isArray(apiEntries) ? apiEntries.length : null,
                                source: Array.isArray(apiEntries) ? 'api' : (fallbackAllowed ? 'fallback' : 'api'),
                                fallbackAllowed,
                                exactMatched,
                                latestDubEpisode: latestEntry ? getAnimeScheduleEntryEpisodeNumber(latestEntry) : '',
                                latestDubDate: latestEntry ? String((getAnimeScheduleEntryDate(latestEntry) || {}).toISOString ? getAnimeScheduleEntryDate(latestEntry).toISOString() : '') : '',
                                matchDetails,
                                matched: !!matchDetails,
                                mode: exactMatched ? 'same-day-dub-for-anime' : (latestEntry ? 'no-same-day-dub-match' : 'none')
                            };
                        }));
                    };

                    function removeArrowArtifacts(root) {
                        if (!root || !isElement(root)) return;

                        if (root.matches('.ama-carousel-nav-btn')) {
                            root.remove();
                            return;
                        }

                        if (root.matches('.ama-manga-carousel-parent')) {
                            root.classList.remove('ama-manga-carousel-parent');
                        }

                        root.querySelectorAll('.ama-carousel-nav-btn').forEach(btn => btn.remove());

                        root.querySelectorAll('.ama-manga-carousel-parent').forEach(parent => {
                            parent.classList.remove('ama-manga-carousel-parent');
                        });
                    }

                    function cleanupCarousels(root) {
                        if (!root || !root.querySelectorAll) return;

                        if (root.matches && root.matches(targetGridsQuery)) {
                            root.classList.remove('ama-optimized-carousel');
                            root.dataset.amaCarouselEnhanced = "false";
                        }

                        root.querySelectorAll(targetGridsQuery).forEach(grid => {
                            grid.classList.remove('ama-optimized-carousel');
                            grid.dataset.amaCarouselEnhanced = "false";
                        });
                    }

                    function isRandomSearchSvg(svg) {
                        if (!svg || svg.tagName.toLowerCase() !== 'svg') return false;

                        const path = svg.querySelector('path');
                        if (!path) return false;

                        return path.getAttribute('d') === randomSearchIconPath;
                    }

                    function shouldRemoveRandomSearchSvg(svg) {
                        if (!svg || svg.classList.contains('ama-search-icon')) return false;
                        if (svg.closest('.ama-search-wrapper')) return false;
                        if (svg.closest('button, a, [role="button"], [role="tab"], [role="menuitem"]')) return false;

                        const parent = svg.parentElement;
                        if (!parent) return false;

                        const parentRect = parent.getBoundingClientRect();
                        const svgRect = svg.getBoundingClientRect();

                        if (!parentRect || !svgRect || parentRect.width < 240) return false;

                        const svgCenterX = svgRect.left + (svgRect.width / 2) - parentRect.left;
                        const isCentered = svgCenterX > parentRect.width * 0.35 && svgCenterX < parentRect.width * 0.65;

                        return isCentered;
                    }

                    function removeRandomSearchIcons(root) {
                        if (!featureSettings.betterMarketplace) return;
                        if (!root) return;

                        const svgs = [];

                        if (root.matches && root.matches('svg')) {
                            svgs.push(root);
                        }

                        if (root.querySelectorAll) {
                            root.querySelectorAll('svg').forEach(svg => svgs.push(svg));
                        }

                        svgs.forEach(svg => {
                            if (!isRandomSearchSvg(svg)) return;
                            if (!shouldRemoveRandomSearchSvg(svg)) return;

                            svg.setAttribute('data-ama-random-search-icon', 'true');
                            svg.remove();
                        });
                    }

                    function optimizeImages(root) {
                        if (!root || !root.querySelectorAll) return;

                        root.querySelectorAll('img:not([data-ama-image-optimized="true"])').forEach(img => {
                            img.dataset.amaImageOptimized = "true";
                            img.loading = 'lazy';
                            img.decoding = 'async';

                            try {
                                img.fetchPriority = 'low';
                            } catch (_) {}

                            try {
                                img.draggable = false;
                            } catch (_) {}
                        });
                    }

                    function makeDraggableScroller(el, featureKey) {
                        if (!el) return;
                        if (el.dataset.amaDragScrollEnhanced === dragScrollEnhancementVersion) return;

                        el.dataset.amaDragScrollEnhanced = dragScrollEnhancementVersion;
                        el.dataset.amaDragFeature = featureKey || (
                            el.classList.contains('ama-carousel-row') || el.classList.contains('ama-extension-carousel')
                                ? 'betterMarketplace'
                                : 'carousels'
                        );

                        if (el.querySelectorAll) {
                            el.querySelectorAll('a, img').forEach(node => {
                                try {
                                    node.draggable = false;
                                } catch (_) {}
                            });
                        }

                        let isDown = false;
                        let didDrag = false;
                        let suppressClick = false;
                        let startX = 0;
                        let startY = 0;
                        let startScrollLeft = 0;

                        function isRealControl(target) {
                            return !!(
                                target &&
                                target.closest &&
                                target.closest('.ama-clone-action, .ama-clone-actions, button, input, textarea, select, summary, [contenteditable="true"]')
                            );
                        }

                        function isEnabled() {
                            if (el.dataset.amaDragFeature === 'betterMarketplace') {
                                return !!featureSettings.betterMarketplace;
                            }

                            if (el.dataset.amaDragFeature === 'extensionsCarousel') {
                                return areExtensionCarouselsEnabled();
                            }

                            return areCarouselsEnabledForCurrentPage();
                        }

                        function stopDrag(pointerId) {
                            isDown = false;
                            el.classList.remove('ama-drag-pending');
                            el.classList.remove('ama-dragging');

                            setTimeout(() => {
                                suppressClick = false;
                                didDrag = false;
                            }, 350);

                            try {
                                if (pointerId !== undefined) el.releasePointerCapture(pointerId);
                            } catch (_) {}
                        }

                        el.addEventListener('pointerdown', event => {
                            if (!isEnabled()) return;
                            if (isRealControl(event.target)) return;
                            if (event.button !== 0) return;
                            if (el.scrollWidth <= el.clientWidth + 2) return;

                            isDown = true;
                            didDrag = false;
                            suppressClick = false;
                            startX = event.clientX;
                            startY = event.clientY;
                            startScrollLeft = el.scrollLeft;
                            el.classList.add('ama-drag-pending');
                        }, true);

                        el.addEventListener('pointermove', event => {
                            if (!isDown) return;

                            const dx = event.clientX - startX;
                            const dy = event.clientY - startY;
                            const isClearHorizontalDrag = Math.abs(dx) > 24 && Math.abs(dx) > Math.abs(dy) + 8;

                            if (isClearHorizontalDrag) {
                                if (!didDrag) {
                                    try {
                                        el.setPointerCapture(event.pointerId);
                                    } catch (_) {}
                                }

                                didDrag = true;
                                suppressClick = true;
                                el.classList.remove('ama-drag-pending');
                                el.classList.add('ama-dragging');
                                event.preventDefault();
                                event.stopPropagation();

                                try {
                                    window.getSelection && window.getSelection().removeAllRanges();
                                } catch (_) {}
                            }

                            if (didDrag) {
                                el.scrollLeft = startScrollLeft - dx;
                                event.preventDefault();
                                event.stopPropagation();
                            }
                        }, true);

                        el.addEventListener('pointerup', event => {
                            stopDrag(event.pointerId);
                        });

                        el.addEventListener('pointercancel', event => {
                            stopDrag(event.pointerId);
                        });

                        el.addEventListener('click', event => {
                            suppressClick = false;
                            didDrag = false;
                        }, true);

                        el.addEventListener('dragstart', event => {
                            event.preventDefault();
                            event.stopPropagation();
                        }, true);
                    }

                    function enhanceCarouselGrid(grid) {
                        if (!areCarouselsEnabledForCurrentPage()) return;
                        if (!grid || grid.dataset.amaCarouselEnhanced === "true") return;

                        grid.dataset.amaCarouselEnhanced = "true";
                        grid.classList.add('ama-optimized-carousel');

                        makeDraggableScroller(grid, 'carousels');
                        optimizeImages(grid);
                    }

                    const dubObserver = new IntersectionObserver(entries => {
                        entries.forEach(entry => {
                            if (!entry.isIntersecting) return;

                            const card = entry.target;
                            dubObserver.unobserve(card);

                            if (!featureSettings.subDubIcons) return;

                            const badge = card.querySelector('.ama-media-badge.dub');
                            if (!badge) return;

                            hasDubForAnimeCard(card).then(hasDub => {
                                if (!badge.isConnected) return;
                                if (!featureSettings.subDubIcons) return;

                                if (hasDub) {
                                    badge.hidden = false;
                                    badge.setAttribute('aria-hidden', 'false');
                                } else {
                                    badge.hidden = true;
                                    badge.setAttribute('aria-hidden', 'true');
                                }
                            });
                        });
                    }, {
                        root: null,
                        rootMargin: '500px 900px',
                        threshold: 0.01
                    });

                    function isAnimeCard(card) {
                        if (!card) return false;

                        const explicitType = String(
                            card.getAttribute('data-media-type') ||
                            card.getAttribute('data-entry-type') ||
                            card.getAttribute('data-media-kind') ||
                            ''
                        ).toLowerCase();

                        if (explicitType.includes('manga') || explicitType.includes('novel') || explicitType.includes('book')) {
                            return false;
                        }

                        if (explicitType.includes('anime')) {
                            return true;
                        }

                        if (card.closest('[data-manga-page-container="true"], [data-manga-entry-screen="true"], [data-manga-entry-page="true"]')) {
                            return false;
                        }

                        const lowerPath = String(window.location.pathname || '').toLowerCase();

                        if (lowerPath.includes('manga') || lowerPath.includes('novel') || lowerPath.includes('book')) {
                            return false;
                        }

                        const link = card.querySelector('a[href]');
                        const href = link ? String(link.getAttribute('href') || '').toLowerCase() : '';

                        if (href.includes('manga') || href.includes('novel') || href.includes('book')) {
                            return false;
                        }

                        return true;
                    }

                    function removeMediaBadges(card) {
                        if (!card) return;

                        const yearEl = card.querySelector('[data-media-entry-card-title-section-year-season="true"]');
                        const row = yearEl ? yearEl.parentElement : null;

                        if (row) {
                            row.classList.remove('ama-media-meta-row');
                            row.querySelectorAll('.ama-media-badges').forEach(el => el.remove());
                        }

                        card.dataset.amaMediaBadgeEnhanced = "false";
                    }

                    function cleanupAllMediaBadges(root) {
                        if (!root) return;

                        if (root.matches && root.matches(mediaEntryCardQuery)) {
                            removeMediaBadges(root);
                        }

                        if (root.querySelectorAll) {
                            root.querySelectorAll(mediaEntryCardQuery).forEach(removeMediaBadges);
                        }
                    }

                    function createDubMediaBadge(extraClass) {
                        const dubBadge = document.createElement('span');
                        dubBadge.className = 'ama-media-badge dub' + (extraClass ? ' ' + extraClass : '');
                        dubBadge.title = 'Dub available';
                        dubBadge.setAttribute('aria-label', 'Dub available');
                        dubBadge.setAttribute('aria-hidden', 'true');
                        dubBadge.hidden = true;
                        dubBadge.innerHTML = MIC_ICON;

                        return dubBadge;
                    }

                    function updateScheduleDubBadge(badge, details) {
                        const episodeLabel = details && details.episodeNumber ? ' Ep. ' + details.episodeNumber : '';
                        const sourceLabel = details && details.source ? ' (' + details.source + ')' : '';
                        const dateLabel = details && details.episodeDate ? ' - ' + details.episodeDate : '';
                        badge.textContent = 'DUB' + episodeLabel;
                        badge.title = 'Dub' + episodeLabel + sourceLabel + dateLabel;
                        badge.setAttribute('aria-label', badge.title);
                    }

                    function isInsideScheduleArea(node) {
                        if (!node || !node.closest) return false;

                        if (node.closest('[data-schedule-calendar-event-item-link], [data-schedule-calendar-mobile-list-day-item-event-link], [data-schedule-calendar-event-item], [data-schedule-calendar-mobile-list-day-item-event-content], [data-schedule-event-item], [data-schedule-media-id]')) return true;
                        if (node.closest('[data-schedule-calendar], [data-schedule-calendar-week], [data-schedule-calendar-mobile-list], [data-schedule-calendar-day], [data-schedule-page], [data-route="schedule"]')) return true;

                        const path = String(window.location.pathname || '').toLowerCase();
                        const hash = String(window.location.hash || '').toLowerCase();
                        if (!path.includes('schedule') && !hash.includes('schedule')) return false;

                        return !!node.closest('main, [role="main"], body');
                    }

                    function getScheduleEventElement(root) {
                        if (!root || !root.closest) return root;

                        const event = root.closest('[data-schedule-calendar-event-item-link], [data-schedule-calendar-mobile-list-day-item-event-link], [data-schedule-calendar-event-item], [data-schedule-calendar-event-item-root], [data-schedule-calendar-mobile-list-day-item-event-content], [data-schedule-event-item], [data-schedule-media-id]');
                        if (event) return event;

                        const link = root.closest(scheduleEntryLinkQuery);
                        if (link && isInsideScheduleArea(link)) return link;

                        return root;
                    }

                    function isVisibleScheduleEvent(root) {
                        const event = getScheduleEventElement(root);
                        if (!event || !event.isConnected) return false;

                        const style = window.getComputedStyle ? window.getComputedStyle(event) : null;
                        if (style && (style.display === 'none' || style.visibility === 'hidden' || style.opacity === '0')) return false;

                        const rects = event.getClientRects ? event.getClientRects() : [];
                        if (!rects || !rects.length) return false;

                        const rect = event.getBoundingClientRect ? event.getBoundingClientRect() : null;
                        if (!rect || rect.width <= 0 || rect.height <= 0) return false;

                        return true;
                    }

                    function getScheduleBadgeAnchor(event) {
                        if (!event || !event.querySelector) return null;

                        const finaleIcon = event.querySelector('[data-schedule-calendar-event-item-finale-icon]');
                        if (finaleIcon) return finaleIcon.parentElement || finaleIcon;

                        const mobileIcons = event.querySelector('[data-schedule-calendar-mobile-list-day-item-event-icons]');
                        if (mobileIcons) return mobileIcons;

                        return event.querySelector('[data-schedule-calendar-event-item-name], [data-schedule-calendar-event-item-text], [data-schedule-calendar-event-item-title], [data-schedule-calendar-mobile-list-day-item-event-text], [data-schedule-event-title], [data-media-title], a[href]') || event;
                    }

                    function insertScheduleDubBadge(anchor, badge) {
                        if (!anchor || !badge || badge.isConnected) return;

                        if (anchor.matches && anchor.matches('[data-schedule-calendar-mobile-list-day-item-event-icons]')) {
                            anchor.appendChild(badge);
                            return;
                        }

                        anchor.insertAdjacentElement('afterend', badge);
                    }

                    function removeScheduleDubBadges(root) {
                        if (!root) return;

                        if (root.matches && root.matches('.ama-schedule-dub-badge, .ama-schedule-dub-badge-server')) {
                            root.remove();
                        }

                        if (root.querySelectorAll) {
                            root.querySelectorAll('.ama-schedule-dub-badge, .ama-schedule-dub-badge-server').forEach(badge => badge.remove());
                            root.querySelectorAll('[data-ama-schedule-dub-enhanced]').forEach(event => {
                                delete event.dataset.amaScheduleDubEnhanced;
                            });
                        }
                    }

                    function cleanupStrayScheduleDubBadges(root) {
                        if (!root || !root.querySelectorAll) return;

                        root.querySelectorAll('.ama-schedule-dub-badge, .ama-schedule-dub-badge-server').forEach(badge => {
                            if (!isInsideScheduleArea(badge)) badge.remove();
                        });
                    }

                    function enhanceScheduleEvent(root) {
                        if (!featureSettings.subDubIcons) {
                            removeScheduleDubBadges(root);
                            return;
                        }

                        const event = getScheduleEventElement(root);
                        if (!event || !event.querySelector) return;
                        if (!isInsideScheduleArea(event)) return;
                        if (event.dataset.amaScheduleDubEnhanced === "true") return;

                        const anchor = getScheduleBadgeAnchor(event);
                        if (!anchor) return;

                        const badge = event.querySelector('.ama-schedule-dub-badge:not(.ama-schedule-dub-badge-server)') || createDubMediaBadge('ama-schedule-dub-badge');
                        badge.hidden = true;
                        badge.setAttribute('aria-hidden', 'true');

                        insertScheduleDubBadge(anchor, badge);

                        event.dataset.amaScheduleDubEnhanced = "true";

                        getDubForScheduleEvent(event).then(details => {
                            if (!badge.isConnected) return;
                            if (!featureSettings.subDubIcons) return;

                            if (details) {
                                updateScheduleDubBadge(badge, details);
                                badge.hidden = false;
                                badge.setAttribute('aria-hidden', 'false');
                            } else {
                                badge.hidden = true;
                                badge.setAttribute('aria-hidden', 'true');
                            }
                        });
                    }

                    function enhanceMediaEntryCard(card) {
                        if (!featureSettings.subDubIcons) {
                            removeMediaBadges(card);
                            return;
                        }

                        if (!card) return;

                        if (!isAnimeCard(card)) {
                            removeMediaBadges(card);
                            return;
                        }

                        if (card.dataset.amaMediaBadgeEnhanced === "true") return;

                        const yearEl = card.querySelector('[data-media-entry-card-title-section-year-season="true"]');
                        const titleEl = card.querySelector('[data-media-entry-card-title-section-title="true"]');

                        if (!yearEl || !titleEl) return;

                        const row = yearEl.parentElement;
                        if (!row) return;

                        card.dataset.amaMediaBadgeEnhanced = "true";

                        row.classList.add('ama-media-meta-row');

                        row.querySelectorAll('.ama-media-badges').forEach(el => el.remove());

                        const badges = document.createElement('span');
                        badges.className = 'ama-media-badges';

                        const dubBadge = createDubMediaBadge('');

                        const ccBadge = document.createElement('span');
                        ccBadge.className = 'ama-media-badge cc';
                        ccBadge.title = 'Captions';
                        ccBadge.setAttribute('aria-label', 'Captions');
                        ccBadge.innerHTML = CC_ICON;

                        badges.appendChild(dubBadge);
                        badges.appendChild(ccBadge);

                        row.appendChild(badges);

                        card.dataset.amaDubTitle = titleEl.textContent.trim();
                        card.dataset.amaDubMalId = card.getAttribute('data-media-mal-id') || '';
                        card.dataset.amaDubAnilistId = card.getAttribute('data-media-id') || '';

                        dubObserver.observe(card);
                    }

                    function filterCards(container, term) {
                        const normalized = String(term || '').toLowerCase();
                        const isModalSearch = !!(
                            container &&
                            (
                                (container.classList && container.classList.contains('ama-modal')) ||
                                (container.closest && container.closest('.ama-modal'))
                            )
                        );

                        container.querySelectorAll('.group\\\\/extension-card').forEach(card => {
                            const isVisible = card.innerText.toLowerCase().includes(normalized);
                            const wrapper = card.closest('.ama-catalog-card-wrap');

                            if (wrapper) {
                                wrapper.hidden = !isVisible;
                            } else {
                                card.style.display = isVisible ? 'flex' : 'none';
                            }
                        });

                        if (!isModalSearch && container && container.querySelectorAll) {
                            getMarketplaceSections(false).forEach(marketplaceSection => {
                                const section = marketplaceSection.card;
                                const grid = section.querySelector('.ama-extension-carousel, .grid');
                                if (!grid) return;

                                if (!normalized) {
                                    section.hidden = false;
                                    return;
                                }

                                const hasVisible = Array.from(grid.querySelectorAll('.group\\\\/extension-card')).some(card => {
                                    if (card.classList && card.classList.contains('ama-status-source-hidden')) return false;
                                    const wrapper = card.closest('.ama-catalog-card-wrap');
                                    return wrapper ? !wrapper.hidden : card.style.display !== 'none';
                                });

                                section.hidden = !hasVisible;
                            });

                            container.querySelectorAll('.ama-status-section').forEach(section => {
                                const grid = section.querySelector('.ama-extension-carousel, .grid');
                                if (!grid) {
                                    section.hidden = true;
                                    return;
                                }

                                if (!normalized) {
                                    section.hidden = !grid.querySelector('.group\\\\/extension-card');
                                    return;
                                }

                                const hasVisible = Array.from(grid.querySelectorAll('.ama-catalog-card-wrap, .group\\\\/extension-card')).some(item => {
                                    return !item.hidden && item.style.display !== 'none';
                                });

                                section.hidden = !hasVisible;
                            });

                            syncMarketplaceSectionSearchVisibility();
                        }
                    }

                    function cleanupBetterMarketplace(root) {
                        if (!root) return;

                        const cards = [];

                        if (root.matches && root.matches(cardQuery)) {
                            cards.push(root);
                        }

                        if (root.querySelectorAll) {
                            root.querySelectorAll(cardQuery).forEach(card => cards.push(card));
                        }

                        cards.forEach(card => {
                            const header = card.querySelector(':scope > .ama-header-container');

                            if (header) {
                                const title = header.querySelector('h3');

                                if (title) {
                                    card.insertBefore(title, header);
                                }

                                header.remove();
                            }

                            card.dataset.amaEnhanced = "false";

                            card.querySelectorAll('.ama-extension-carousel').forEach(grid => {
                                grid.classList.remove('ama-extension-carousel');
                                grid.dataset.amaDragScrollEnhanced = "false";
                            });
                        });

                        if (root.querySelectorAll) {
                            root.querySelectorAll('.ama-global-catalog-bar').forEach(bar => bar.remove());
                            root.querySelectorAll('.ama-live-card-actions').forEach(actions => actions.remove());
                            root.querySelectorAll('.ama-extension-carousel').forEach(grid => {
                                grid.classList.remove('ama-extension-carousel');
                                grid.dataset.amaDragScrollEnhanced = "false";
                            });
                            root.querySelectorAll('.ama-modal').forEach(modal => modal.remove());
                        }
                    }

                    function cleanupExtensionCarousels(root) {
                        if (!root) return;

                        const grids = [];

                        if (root.matches && root.matches('.ama-extension-carousel')) {
                            grids.push(root);
                        }

                        if (root.querySelectorAll) {
                            root.querySelectorAll('.ama-extension-carousel').forEach(grid => grids.push(grid));
                        }

                        grids.forEach(grid => {
                            grid.classList.remove('ama-extension-carousel');
                            grid.dataset.amaDragScrollEnhanced = "false";
                            grid.classList.remove('ama-drag-pending');
                            grid.classList.remove('ama-dragging');
                        });
                    }

                    function hasMarketplaceExtensionCards(root) {
                        return !!(root && root.querySelector && root.querySelector('.group\\\\/extension-card'));
                    }

                    function isMarketplaceSectionCard(card) {
                        if (!card || !card.querySelector) return false;

                        const grid = card.querySelector('.ama-extension-carousel, .grid');
                        return hasMarketplaceExtensionCards(grid || card);
                    }

                    function isBetterMarketplacePage() {
                        return hasMarketplaceExtensionCards(document);
                    }

                    function getExtensionCardData(card) {
                        const paragraphs = Array.from(card.querySelectorAll('p')).map(p => p.textContent.trim()).filter(Boolean);
                        const badges = Array.from(card.querySelectorAll('.UI-Badge__root')).map(badge => badge.textContent.trim()).filter(Boolean);
                        const getBadgeValue = (label) => {
                            const pattern = new RegExp('^' + label + '\\s*:\\s*(.+)$', 'i');
                            const match = badges.map(value => String(value || '').match(pattern)).find(Boolean);

                            return match ? match[1].trim() : '';
                        };
                        const idEl = card.querySelector('.text-xs .opacity-30, [data-extension-id], [data-extension-card-id]');
                        const titleEl = card.querySelector('.font-semibold');
                        const manifestLinkPattern = new RegExp('plugins/[^/?#]+[.]json', 'i');
                        const manifestLink = Array.from(card.querySelectorAll('a[href]')).map(link => link.getAttribute('href') || '').find(href => manifestLinkPattern.test(href)) || '';
                        const manifestId = getExtensionIdFromManifestUri(manifestLink);
                        const firstPlainBadge = badges.find(value => !/^(id|author|language|manifest url)\s*:/i.test(value)) || '';

                        return {
                            id: getBadgeValue('ID') || manifestId || (idEl && idEl.textContent.trim()) || getKnownExtensionIdFromText(card.innerText || card.textContent || ''),
                            name: (titleEl && titleEl.textContent.trim()) || paragraphs[0] || 'Extension',
                            description: paragraphs[2] || paragraphs[1] || '',
                            version: getBadgeValue('Version') || firstPlainBadge,
                            author: getBadgeValue('Author') || badges[1] || '',
                            language: getBadgeValue('Language') || badges[2] || '',
                        };
                    }

                    function hasExtensionUpdateText(value) {
                        const text = String(value || '').trim();
                        if (!text) return false;

                        const normalized = text.replace(/\s+/g, ' ');
                        const arrowChars = [8594, 8658, 10140, 10230].map(code => String.fromCharCode(code));

                        if (normalized.includes('->')) return true;
                        if (arrowChars.some(char => normalized.includes(char))) return true;

                        return new RegExp('\\bv?\\d+(?:\\.\\d+){1,3}\\s+(?:to|updated? to)\\s+v?\\d+(?:\\.\\d+){1,3}\\b', 'i').test(normalized);
                    }

                    function getExtensionUpdateText(card) {
                        if (!card || !card.querySelectorAll) return '';

                        const parts = [];

                        card.querySelectorAll('.UI-Badge__root, [class*="Badge"], [class*="badge"]').forEach(node => {
                            const value = node.textContent && node.textContent.trim();
                            if (value) parts.push(value);
                        });

                        card.querySelectorAll('button, [aria-label], [title]').forEach(node => {
                            const text = node.textContent && node.textContent.trim();
                            const aria = node.getAttribute && node.getAttribute('aria-label');
                            const title = node.getAttribute && node.getAttribute('title');

                            if (text) parts.push(text);
                            if (aria) parts.push(aria);
                            if (title) parts.push(title);
                        });

                        const fullText = card.innerText || card.textContent || '';
                        if (fullText) parts.push(fullText);

                        return parts.join(' | ');
                    }

                    function normalizeVersionText(value) {
                        const match = String(value || '').match(/\d+(?:\.\d+){1,3}/);
                        return match ? match[0] : '';
                    }

                    function getLastVersionText(value) {
                        const matches = String(value || '').match(/\d+(?:\.\d+){1,3}/g);
                        return matches && matches.length ? matches[matches.length - 1] : '';
                    }

                    function compareVersionText(a, b) {
                        const left = normalizeVersionText(a);
                        const right = normalizeVersionText(b);

                        if (!left || !right) return 0;

                        const leftParts = left.split('.').map(part => parseInt(part, 10) || 0);
                        const rightParts = right.split('.').map(part => parseInt(part, 10) || 0);
                        const length = Math.max(leftParts.length, rightParts.length);

                        for (let i = 0; i < length; i += 1) {
                            const diff = (leftParts[i] || 0) - (rightParts[i] || 0);
                            if (diff !== 0) return diff > 0 ? 1 : -1;
                        }

                        return 0;
                    }

                    function getExtensionIdFromManifestUri(value) {
                        const text = String(value || '');
                        const match = text.match(new RegExp('plugins/([^/?#]+)[.]json', 'i'));
                        if (!match) return '';

                        const file = match[1].toLowerCase();
                        if (file === 'seautils-kolex06-version') return 'SeaUtils-Kolex06-Version';
                        if (file === 'asunatracks-sync') return 'asunatracks-sync';

                        return match[1];
                    }

                    function getKnownExtensionIdFromText(value) {
                        const text = normalizeExtensionText(value);

                        if (text.includes('asunatracks sync') || text.includes('asunatracks-sync')) {
                            return 'asunatracks-sync';
                        }

                        if (text.includes('seautils kolex06-version') || text.includes('seautils-kolex06-version')) {
                            return 'SeaUtils-Kolex06-Version';
                        }

                        return '';
                    }

                    function applyExtensionUpdateStyle(card, hasUpdate) {
                        if (!card) return;

                        const wrapper = card.closest && card.closest('.ama-catalog-card-wrap');
                        const targets = wrapper ? [card, wrapper] : [card];

                        if (hasUpdate) {
                            card.dataset.amaUpdateAvailable = 'true';
                            if (wrapper) wrapper.dataset.amaUpdateAvailable = 'true';

                            targets.forEach(target => {
                                if (!target || !target.style) return;

                                target.style.setProperty('border-color', 'rgba(56, 189, 248, 0.95)', 'important');
                                target.style.setProperty('outline', '1px solid rgba(56, 189, 248, 0.85)', 'important');
                                target.style.setProperty('outline-offset', '0', 'important');
                                target.style.setProperty('box-shadow', '0 0 0 1px rgba(56, 189, 248, 0.75), 0 0 18px rgba(56, 189, 248, 0.24)', 'important');
                                target.style.setProperty('background', 'rgba(14, 165, 233, 0.12)', 'important');
                            });
                        } else {
                            delete card.dataset.amaUpdateAvailable;
                            if (wrapper) delete wrapper.dataset.amaUpdateAvailable;

                            targets.forEach(target => {
                                if (!target || !target.style) return;

                                target.style.removeProperty('border-color');
                                target.style.removeProperty('outline');
                                target.style.removeProperty('outline-offset');
                                target.style.removeProperty('box-shadow');
                                target.style.removeProperty('background');
                            });
                        }
                    }

                    const updateGlowSuppressions = new Map();

                    function getUpdateSuppressionKeys(data) {
                        const id = normalizeMarketplaceLookupText(data && data.id);
                        const name = normalizeMarketplaceLookupText(data && data.name);

                        return [id, name].filter((value, index, list) => value && list.indexOf(value) === index);
                    }

                    function isUpdateGlowSuppressed(data) {
                        const keys = getUpdateSuppressionKeys(data);
                        if (!keys.length) return false;

                        return keys.some(key => {
                            const expiresAt = updateGlowSuppressions.get(key) || 0;
                            if (!expiresAt) return false;

                            if (Date.now() > expiresAt) {
                                updateGlowSuppressions.delete(key);
                                return false;
                            }

                            return true;
                        });
                    }

                    function suppressUpdateGlowForData(data) {
                        getUpdateSuppressionKeys(data).forEach(key => {
                            updateGlowSuppressions.set(key, Date.now() + 30000);
                        });
                    }

                    function markExtensionUpdateState(card) {
                        if (!card || !card.querySelectorAll) return false;

                        const data = getExtensionCardData(card);
                        const hasUpdate = hasExtensionUpdateText(getExtensionUpdateText(card));
                        const shouldShowUpdate = hasUpdate && !isUpdateGlowSuppressed(data);
                        const value = shouldShowUpdate ? 'true' : 'false';

                        applyExtensionUpdateStyle(card, shouldShowUpdate);

                        const wrapper = card.closest('.ama-catalog-card-wrap');
                        if (wrapper) {
                            if (shouldShowUpdate) {
                                wrapper.dataset.amaUpdateAvailable = 'true';
                            } else {
                                delete wrapper.dataset.amaUpdateAvailable;
                            }
                        }

                        return shouldShowUpdate;
                    }

                    function clearExtensionUpdateStateForData(data) {
                        const targetId = normalizeMarketplaceLookupText(data && data.id);
                        const targetName = normalizeMarketplaceLookupText(data && data.name);

                        document.querySelectorAll('.group\\\\/extension-card').forEach(card => {
                            const cardData = getExtensionCardData(card);
                            const cardId = normalizeMarketplaceLookupText(cardData.id);
                            const cardName = normalizeMarketplaceLookupText(cardData.name);
                            const matches = (
                                (targetId && cardId === targetId) ||
                                (targetName && cardName === targetName)
                            );

                            if (!matches) return;

                            applyExtensionUpdateStyle(card, false);

                            const wrapper = card.closest && card.closest('.ama-catalog-card-wrap');
                            if (wrapper) {
                                delete wrapper.dataset.amaUpdateAvailable;
                            }
                        });
                    }

                    function refreshExtensionUpdateStateForData(data) {
                        suppressUpdateGlowForData(data);
                        clearExtensionUpdateStateForData(data);

                        const recheck = () => {
                            clearExtensionCache();

                            document.querySelectorAll('.group\\\\/extension-card').forEach(card => {
                                const cardData = getExtensionCardData(card);
                                const matches = (
                                    normalizeMarketplaceLookupText(data && data.id) && normalizeMarketplaceLookupText(cardData.id) === normalizeMarketplaceLookupText(data && data.id)
                                ) || (
                                    normalizeMarketplaceLookupText(data && data.name) && normalizeMarketplaceLookupText(cardData.name) === normalizeMarketplaceLookupText(data && data.name)
                                );

                                if (matches) markExtensionUpdateState(card);
                            });
                        };

                        setTimeout(recheck, 350);
                        setTimeout(recheck, 1200);
                        setTimeout(recheck, 2500);
                    }

                    async function refreshExtensionUpdateStateIfInstalled(card, data) {
                        if (!card || !data) return;

                        const updateTargetVersion = getLastVersionText(getExtensionUpdateText(card));

                        if (!updateTargetVersion) {
                            refreshExtensionUpdateStateForData(data);
                            return;
                        }

                        clearExtensionCache();

                        try {
                            const extension = await findExtensionForCard(card);
                            const installedVersion = extension && extension.version;

                            if (installedVersion && compareVersionText(installedVersion, updateTargetVersion) >= 0) {
                                refreshExtensionUpdateStateForData(data);
                            }
                        } catch (_) {}
                    }

                    function isNativeExtensionUpdateButton(button) {
                        if (!button) return false;

                        const text = [
                            button.textContent || '',
                            button.getAttribute && button.getAttribute('aria-label') || '',
                            button.getAttribute && button.getAttribute('title') || ''
                        ].join(' ').replace(/\s+/g, ' ').trim().toLowerCase();

                        if (!text) return false;
                        if (text.includes('check for update')) return false;

                        return (
                            text.includes('install update') ||
                            text.includes('install extension update') ||
                            text === 'update' ||
                            text.includes(' update ')
                        );
                    }

                    function handleNativeExtensionUpdateClick(event) {
                        const target = event && event.target;
                        const button = target && target.closest ? target.closest('button, [role="button"]') : null;
                        if (!button) return;
                        if (button.closest && button.closest('.ama-clone-action, .ama-clone-actions')) return;

                        const card = button.closest && button.closest('.group\\\\/extension-card');
                        if (!card) return;
                        if (card.closest && card.closest('.ama-modal')) return;
                        if (!card.dataset.amaUpdateAvailable && !(card.closest && card.closest('[data-ama-update-available="true"]'))) return;

                        const isExplicitUpdateButton = isNativeExtensionUpdateButton(button);
                        const hasUpdateText = hasExtensionUpdateText(getExtensionUpdateText(card));
                        if (!isExplicitUpdateButton && !hasUpdateText) return;

                        const data = getExtensionCardData(card);
                        if (!getUpdateSuppressionKeys(data).length) return;

                        [600, 1400, 2800, 5000].forEach(delay => {
                            setTimeout(() => refreshExtensionUpdateStateIfInstalled(card, data), delay);
                        });
                    }

                    function markMarketplaceExtensionCards(root) {
                        if (!root || !root.querySelectorAll) return;

                        const cards = [];

                        if (root.matches && root.matches('.group\\\\/extension-card')) {
                            cards.push(root);
                        }

                        root.querySelectorAll('.group\\\\/extension-card').forEach(card => cards.push(card));

                        cards.forEach(card => {
                            markExtensionUpdateState(card);
                            enhanceMarketplaceExtraInfo(card);
                            enhanceInstalledMarketplaceExtraInfo(card);
                            enhanceInstalledMarketplaceCardActions(card);
                        });
                    }

                    function openAmaModal(title, contentHtml) {
                        const modal = document.createElement('div');
                        modal.className = 'ama-modal';
                        modal.innerHTML =
                            '<button type="button" class="ama-modal-close">Close</button>' +
                            '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">' + escapeHtml(title) + '</h2>' +
                            contentHtml;

                        modal.querySelector('.ama-modal-close').onclick = () => modal.remove();

                        modal.onclick = event => {
                            if (event.target === modal) modal.remove();
                        };

                        document.body.appendChild(modal);

                        return modal;
                    }

                    function readStoredString(storage, key) {
                        try {
                            const raw = storage.getItem(key);
                            if (!raw) return '';

                            try {
                                const parsed = JSON.parse(raw);
                                return typeof parsed === 'string' ? parsed : String(raw || '');
                            } catch (_) {
                                return String(raw || '');
                            }
                        } catch (_) {
                            return '';
                        }
                    }

                    function getOrCreateSeanimeClientId() {
                        const existing = readStoredString(window.localStorage, 'seall-client-id').trim();
                        if (existing) return existing;

                        let next = '';

                        try {
                            if (window.crypto && typeof window.crypto.randomUUID === 'function') {
                                next = window.crypto.randomUUID();
                            }
                        } catch (_) {}

                        if (!next) {
                            next = 'seall-' + Math.random().toString(16).slice(2) + '-' + Date.now().toString(16);
                        }

                        try {
                            window.localStorage.setItem('seall-client-id', next);
                        } catch (_) {}

                        return next;
                    }

                    function getSeanimeHeaders(extraHeaders) {
                        const headers = Object.assign({}, extraHeaders || {});
                        const token = readStoredString(window.localStorage, 'sea-server-auth-token').trim();
                        const clientId = getOrCreateSeanimeClientId();
                        const clientProof = readStoredString(window.sessionStorage, 'seall-client-id-proof').trim();

                        if (token) headers['X-Seall-Token'] = token;
                        if (clientId) headers['X-Seall-Client-Id'] = clientId;
                        if (clientProof) headers['X-Seall-Client-Id-Proof'] = clientProof;
                        if (window.electron) headers['X-Seall-Client-Platform'] = 'denshi';

                        return headers;
                    }

                    function syncSeanimeIdentity(response) {
                        if (!response || !response.headers) return;

                        try {
                            const clientId = response.headers.get('X-Seall-Client-Id') || response.headers.get('x-seall-client-id') || '';
                            const clientProof = response.headers.get('X-Seall-Client-Id-Proof') || response.headers.get('x-seall-client-id-proof') || '';

                            if (clientId.trim()) {
                                window.localStorage.setItem('seall-client-id', clientId.trim());
                            }

                            if (clientProof.trim()) {
                                window.sessionStorage.setItem('seall-client-id-proof', clientProof.trim());
                            }
                        } catch (_) {}
                    }

                    function getSeanimeApiCandidates(endpoint) {
                        const raw = String(endpoint || '');

                        if (/^https?:\\/\\//i.test(raw)) return [raw];

                        const path = raw.charAt(0) === '/' ? raw : '/' + raw;
                        const candidates = [path];

                        try {
                            const protocol = window.location && window.location.protocol && window.location.protocol !== 'file:'
                                ? window.location.protocol
                                : 'http:';
                            const hostname = window.location && window.location.hostname ? window.location.hostname : '';

                            if (hostname) {
                                candidates.push(protocol + '//' + hostname + ':43211' + path);
                            }
                        } catch (_) {}

                        candidates.push('http://127.0.0.1:43211' + path);
                        candidates.push('http://localhost:43211' + path);

                        return candidates.filter((candidate, index, list) => candidate && list.indexOf(candidate) === index);
                    }

                    async function parseSeanimeResponse(response) {
                        const text = await response.text();
                        let parsed = text;

                        if (text) {
                            try {
                                parsed = JSON.parse(text);
                            } catch (_) {}
                        }

                        if (!response.ok) {
                            const message = parsed && typeof parsed === 'object' && parsed.error
                                ? parsed.error
                                : (typeof parsed === 'string' && parsed.trim() ? parsed : 'Request failed (HTTP ' + response.status + ').');
                            throw new Error(message);
                        }

                        if (typeof parsed === 'string' && /^\\s*<!doctype|^\\s*<html/i.test(parsed)) {
                            throw new Error('This address did not return Seall API data.');
                        }

                        if (parsed && typeof parsed === 'object') {
                            if (parsed.error) throw new Error(parsed.error);
                            if (Object.prototype.hasOwnProperty.call(parsed, 'data')) return parsed.data;
                        }

                        return parsed;
                    }

                    async function fetchSeanime(endpoint, options) {
                        const candidates = getSeanimeApiCandidates(endpoint);
                        let lastError = null;

                        for (const candidate of candidates) {
                            const nextOptions = Object.assign({}, options || {});
                            nextOptions.credentials = 'include';
                            nextOptions.cache = nextOptions.cache || 'no-store';
                            nextOptions.headers = getSeanimeHeaders(nextOptions.headers);

                            try {
                                const response = await fetch(candidate, nextOptions);
                                syncSeanimeIdentity(response);
                                return await parseSeanimeResponse(response);
                            } catch (error) {
                                lastError = error;
                            }
                        }

                        throw lastError || new Error('Request failed.');
                    }

                    function flattenExtensions(data) {
                        const result = [];

                        if (!data || typeof data !== 'object') return result;

                        if (Array.isArray(data.extensions)) {
                            data.extensions.forEach(extension => {
                                if (extension) result.push(Object.assign({}, extension, { __amaDisabled: false }));
                            });
                        }

                        if (Array.isArray(data.disabledExtensions)) {
                            data.disabledExtensions.forEach(extension => {
                                if (extension) result.push(Object.assign({}, extension, { __amaDisabled: true }));
                            });
                        }

                        ['invalidExtensions', 'invalidUserConfigExtensions'].forEach(key => {
                            if (Array.isArray(data[key])) {
                                data[key].forEach(item => {
                                    if (item && item.extension) result.push(item.extension);
                                });
                            }
                        });

                        return result;
                    }

                    function loadAllExtensions() {
                        if (allExtensionsPromise) return allExtensionsPromise;

                        allExtensionsPromise = fetchSeanime('/api/v1/extensions/all', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json'
                            },
                            body: JSON.stringify({
                                withUpdates: false
                            })
                        }).then(flattenExtensions).catch(() => {
                            allExtensionsPromise = null;
                            return [];
                        });

                        return allExtensionsPromise;
                    }

                    function normalizeExtensionText(value) {
                        return String(value || '').trim().toLowerCase();
                    }

                    async function findExtensionForCard(card) {
                        const data = getExtensionCardData(card);
                        const extensions = await loadAllExtensions();

                        if (data.id) {
                            const byId = extensions.find(extension => extension && extension.id === data.id);
                            if (byId) return byId;
                        }

                        const targetName = normalizeExtensionText(data.name);
                        const targetAuthor = normalizeExtensionText(data.author);
                        const targetVersion = normalizeExtensionText(data.version).replace(/^v/, '');
                        const targetDescription = normalizeExtensionText(data.description);

                        let best = null;
                        let bestScore = 0;

                        extensions.forEach(extension => {
                            if (!extension) return;

                            let score = 0;

                            if (targetName && normalizeExtensionText(extension.name) === targetName) score += 6;
                            if (targetAuthor && normalizeExtensionText(extension.author) === targetAuthor) score += 3;
                            if (targetVersion && normalizeExtensionText(extension.version).replace(/^v/, '') === targetVersion) score += 2;
                            if (targetDescription && normalizeExtensionText(extension.description).includes(targetDescription)) score += 1;

                            if (score > bestScore) {
                                best = extension;
                                bestScore = score;
                            }
                        });

                        return bestScore >= 6 ? best : null;
                    }

                    function showInstalledDetails(card) {
                        const data = getExtensionCardData(card);

                        openAmaModal(data.name, [
                            '<p style="color:rgba(255,255,255,.65); margin:0 0 16px 0;">' + escapeHtml(data.description) + '</p>',
                            '<div class="ama-detail-grid">',
                                '<div class="ama-detail-item"><div class="ama-detail-label">ID</div><div>' + escapeHtml(data.id) + '</div></div>',
                                '<div class="ama-detail-item"><div class="ama-detail-label">Version</div><div>' + escapeHtml(data.version || 'Unknown') + '</div></div>',
                                '<div class="ama-detail-item"><div class="ama-detail-label">Author</div><div>' + escapeHtml(data.author || 'Unknown') + '</div></div>',
                                '<div class="ama-detail-item"><div class="ama-detail-label">Language</div><div>' + escapeHtml(data.language || 'Unknown') + '</div></div>',
                            '</div>'
                        ].join(''));
                    }

                    function bindAmaModalClose(modal) {
                        const close = modal && modal.querySelector ? modal.querySelector('.ama-modal-close') : null;
                        if (close) close.onclick = () => modal.remove();
                    }

                    function getExtensionActionId(data, extension) {
                        return (extension && extension.id) || (data && data.id) || '';
                    }

                    function getExtensionManifestUri(extension) {
                        return extension && (extension.manifestURI || extension.manifestUri) ? (extension.manifestURI || extension.manifestUri) : '';
                    }

                    function getExtensionDocumentationUrl(extension) {
                        if (!extension) return '';

                        return String(extension.readme || extension.documentation || extension.docs || extension.website || '').trim();
                    }

                    function isExtensionDisabled(card, extension) {
                        if (extension && extension.__amaDisabled) return true;
                        return !!(card && /\\bDisabled\\b/i.test(card.innerText || ''));
                    }

                    function isKolex06VersionExtension(data, extension) {
                        const id = normalizeExtensionText(getExtensionActionId(data, extension));
                        const name = normalizeExtensionText((extension && extension.name) || (data && data.name));
                        return id === 'seautils' || id === 'seautils-kolex06-version' || id === 'seautils kolex06-version' || name === 'seautils' || name === 'seautils kolex06-version' || name === 'seautils kolex-version';
                    }

                    function postExtensionAction(endpoint, body) {
                        return fetchSeanime(endpoint, {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json'
                            },
                            body: JSON.stringify(body || {})
                        });
                    }

                    function clearExtensionCache() {
                        allExtensionsPromise = null;
                    }

                    function renderKolex06VersionPreferencesModal(modal) {
                        const current = Object.assign({}, defaultSettings, featureSettings || {});

                        modal.innerHTML =
                            '<button type="button" class="ama-modal-close">Close</button>' +
                            '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">Preferences</h2>' +
                            '<p style="color:rgba(255,255,255,.65); margin:0 0 16px 0;">SeaUtils Kolex06-Version</p>' +
                            '<div class="ama-config-form">' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Better Marketplace</label><input data-ama-pref="betterMarketplace" type="checkbox"></div>' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Carousels</label><input data-ama-pref="carousels" type="checkbox"></div>' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Carousels: Search</label><button type="button" class="ama-config-toggle" data-ama-pref="carouselsSearch"></button></div>' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Carousels: Extensions</label><button type="button" class="ama-config-toggle" data-ama-pref="carouselsExtensions"></button></div>' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Carousels: My Lists</label><button type="button" class="ama-config-toggle" data-ama-pref="carouselsLists"></button></div>' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Carousels: Manga</label><button type="button" class="ama-config-toggle" data-ama-pref="carouselsManga"></button></div>' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Carousels: Other Pages</label><button type="button" class="ama-config-toggle" data-ama-pref="carouselsOther"></button></div>' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Sub/Dub Icons</label><input data-ama-pref="subDubIcons" type="checkbox"></div>' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Hide File Names</label><input data-ama-pref="hideFileNames" type="checkbox"></div>' +
                                '<div class="ama-config-switch-row"><label class="ama-config-switch-label">Bas1874 Marketplace</label><input data-ama-pref="useBas1874Marketplace" type="checkbox"></div>' +
                                '<div class="ama-config-help">Uses the extended marketplace feed with status and VirusTotal scan metadata. Refresh the Extensions page after switching.</div>' +
                                '<div class="ama-config-actions"><button type="button" class="ama-config-save">Save</button><span class="ama-config-status"></span></div>' +
                            '</div>';

                        bindAmaModalClose(modal);

                        ['betterMarketplace', 'carousels', 'carouselsSearch', 'carouselsExtensions', 'carouselsLists', 'carouselsManga', 'carouselsOther', 'subDubIcons', 'hideFileNames', 'useBas1874Marketplace'].forEach(key => {
                            const control = modal.querySelector('[data-ama-pref="' + key + '"]');
                            if (!control) return;

                            const value = current[key] !== false;

                            if (control.matches('button')) {
                                control.dataset.amaValue = value ? 'true' : 'false';
                                control.textContent = value ? 'On' : 'Off';
                                control.setAttribute('aria-pressed', value ? 'true' : 'false');
                                return;
                            }

                            control.checked = value;
                        });

                        modal.querySelectorAll('.ama-config-toggle').forEach(button => {
                            button.addEventListener('click', () => {
                                if (button.disabled || button.dataset.amaDisabled === 'true') return;

                                const nextValue = button.dataset.amaValue !== 'true';
                                button.dataset.amaValue = nextValue ? 'true' : 'false';
                                button.textContent = nextValue ? 'On' : 'Off';
                                button.setAttribute('aria-pressed', nextValue ? 'true' : 'false');
                            });
                        });

                        const betterMarketplaceInput = modal.querySelector('[data-ama-pref="betterMarketplace"]');
                        const carouselsInput = modal.querySelector('[data-ama-pref="carousels"]');
                        const carouselSettingKeys = ['carouselsSearch', 'carouselsExtensions', 'carouselsLists', 'carouselsManga', 'carouselsOther'];

                        const getPreferenceValue = (control) => {
                            if (!control) return false;
                            if (control.matches && control.matches('button')) return control.dataset.amaValue === 'true';
                            if (control.dataset && control.dataset.amaSavedChecked) return control.dataset.amaSavedChecked === 'true';
                            return !!control.checked;
                        };

                        const setPreferenceDisabled = (input, disabled, title) => {
                            if (!input) return;

                            const row = input.closest('.ama-config-switch-row');
                            let visual = row ? row.querySelector('.ama-disabled-switch-visual') : null;
                            const savedChecked = getPreferenceValue(input);

                            if (input.matches && input.matches('button')) {
                                input.disabled = !!disabled;
                                input.dataset.amaDisabled = disabled ? 'true' : 'false';

                                if (row) {
                                    row.dataset.amaDisabled = disabled ? 'true' : 'false';
                                    row.title = disabled ? (title || '') : '';
                                    row.style.opacity = disabled ? '0.42' : '';
                                    row.style.filter = disabled ? 'grayscale(1)' : '';
                                    row.style.background = disabled ? 'rgba(255,255,255,0.025)' : '';
                                    row.style.borderColor = disabled ? 'rgba(255,255,255,0.04)' : '';
                                }

                                return;
                            }

                            input.disabled = !!disabled;
                            input.style.visibility = disabled ? 'hidden' : '';
                            input.style.position = disabled ? 'absolute' : '';
                            input.style.right = disabled ? '12px' : '';
                            input.style.pointerEvents = disabled ? 'none' : '';
                            input.style.width = disabled ? '0' : '';
                            input.style.height = disabled ? '0' : '';
                            input.style.minWidth = disabled ? '0' : '';
                            input.style.opacity = disabled ? '0' : '';
                            input.style.accentColor = disabled ? '#64748b' : '';
                            input.style.filter = disabled ? 'grayscale(1) saturate(0)' : '';
                            input.style.borderRadius = disabled ? '999px' : '';
                            input.style.border = disabled ? '1px solid rgba(148,163,184,.55)' : '';
                            input.style.background = disabled ? '#262626' : '';
                            input.style.boxShadow = disabled ? 'none' : '';

                            if (!row) return;

                            if (disabled) {
                                input.dataset.amaSavedChecked = savedChecked ? 'true' : 'false';
                                input.checked = false;

                                if (!visual) {
                                    visual = document.createElement('span');
                                    visual.className = 'ama-disabled-switch-visual';
                                    row.appendChild(visual);
                                }

                                visual.dataset.amaChecked = savedChecked ? 'true' : 'false';
                                visual.textContent = savedChecked ? 'On' : 'Off';
                            } else if (visual) {
                                visual.remove();
                            }

                            if (!disabled && input.dataset.amaSavedChecked) {
                                input.checked = input.dataset.amaSavedChecked === 'true';
                                delete input.dataset.amaSavedChecked;
                            }

                            row.dataset.amaDisabled = disabled ? 'true' : 'false';
                            row.title = disabled ? (title || '') : '';
                            row.style.opacity = disabled ? '0.42' : '';
                            row.style.filter = disabled ? 'grayscale(1)' : '';
                            row.style.background = disabled ? 'rgba(255,255,255,0.025)' : '';
                            row.style.borderColor = disabled ? 'rgba(255,255,255,0.04)' : '';
                        };

                        const syncCarouselControls = () => {
                            const carouselsEnabled = !carouselsInput || carouselsInput.checked;
                            const betterMarketplaceEnabled = !betterMarketplaceInput || betterMarketplaceInput.checked;

                            carouselSettingKeys.forEach(key => {
                                const input = modal.querySelector('[data-ama-pref="' + key + '"]');
                                let disabled = !carouselsEnabled;
                                let title = 'Turn on Carousels to use this setting.';

                                if (key === 'carouselsExtensions' && !betterMarketplaceEnabled) {
                                    disabled = true;
                                    title = 'Turn on Better Marketplace to use this setting.';
                                }

                                setPreferenceDisabled(input, disabled, title);
                            });
                        };

                        if (betterMarketplaceInput) {
                            betterMarketplaceInput.addEventListener('change', syncCarouselControls);
                        }

                        if (carouselsInput) {
                            carouselsInput.addEventListener('change', syncCarouselControls);
                        }

                        syncCarouselControls();

                        const save = modal.querySelector('.ama-config-save');
                        const status = modal.querySelector('.ama-config-status');

                        if (save) {
                            save.onclick = () => {
                                const readPreferenceChecked = (key) => {
                                    return getPreferenceValue(modal.querySelector('[data-ama-pref="' + key + '"]'));
                                };

                                const next = {
                                    betterMarketplace: readPreferenceChecked('betterMarketplace'),
                                    carousels: readPreferenceChecked('carousels'),
                                    carouselsSearch: readPreferenceChecked('carouselsSearch'),
                                    carouselsExtensions: readPreferenceChecked('carouselsExtensions'),
                                    carouselsLists: readPreferenceChecked('carouselsLists'),
                                    carouselsManga: readPreferenceChecked('carouselsManga'),
                                    carouselsOther: readPreferenceChecked('carouselsOther'),
                                    subDubIcons: readPreferenceChecked('subDubIcons'),
                                    hideFileNames: readPreferenceChecked('hideFileNames'),
                                    useBas1874Marketplace: readPreferenceChecked('useBas1874Marketplace'),
                                };

                                if (typeof window.__AMA_SAVE_SETTINGS__ === 'function') {
                                    window.__AMA_SAVE_SETTINGS__(next);
                                } else {
                                    featureSettings = normalizeFeatureSettings(next);
                                    writeBrowserSettings(featureSettings);
                                    applyBas1874MarketplacePreference();
                                    setBodyFlags();
                                }

                                if (status) status.textContent = 'Preferences saved.';
                            };
                        }
                    }

                    async function showKolex06VersionPreferences() {
                        const modal = openAmaModal('Preferences', '<p>Loading SeaUtils Kolex06-Version preferences...</p>');
                        renderKolex06VersionPreferencesModal(modal);
                    }

                    async function showInstalledMore(card) {
                        const data = getExtensionCardData(card);
                        const modal = openAmaModal(data.name, '<p>Loading extension actions...</p>');

                        try {
                            const extension = await findExtensionForCard(card);
                            const actionId = getExtensionActionId(data, extension);
                            const manifestUri = getExtensionManifestUri(extension);
                            const disabled = isExtensionDisabled(card, extension);
                            const name = (extension && extension.name) || data.name;
                            const version = (extension && extension.version) || data.version || 'Unknown';
                            const author = (extension && extension.author) || data.author || 'Unknown';
                            const language = (extension && extension.language) || data.language || 'Unknown';
                            const description = (extension && extension.description) || data.description || '';
                            const canManage = !!actionId && manifestUri !== 'builtin';
                            const canShowPreferences = !!(extension && extension.userConfig);

                            modal.innerHTML =
                                '<button type="button" class="ama-modal-close">Close</button>' +
                                '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">' + escapeHtml(name) + '</h2>' +
                                '<p style="color:rgba(255,255,255,.65); margin:0 0 16px 0;">' + escapeHtml(description) + '</p>' +
                                '<div class="ama-detail-grid">' +
                                    '<div class="ama-detail-item"><div class="ama-detail-label">ID</div><div>' + escapeHtml(actionId || 'Unknown') + '</div></div>' +
                                    '<div class="ama-detail-item"><div class="ama-detail-label">Version</div><div>' + escapeHtml(version) + '</div></div>' +
                                    '<div class="ama-detail-item"><div class="ama-detail-label">Author</div><div>' + escapeHtml(author) + '</div></div>' +
                                    '<div class="ama-detail-item"><div class="ama-detail-label">Language</div><div>' + escapeHtml(language) + '</div></div>' +
                                '</div>' +
                                '<div class="ama-action-panel">' +
                                    (canShowPreferences ? '<button type="button" class="ama-action-button" data-ama-more-action="preferences">Preferences</button>' : '') +
                                    (manifestUri && manifestUri !== 'builtin' ? '<button type="button" class="ama-action-button" data-ama-more-action="check-update">Check for updates</button>' : '') +
                                    (canManage ? '<button type="button" class="ama-action-button" data-ama-more-action="toggle-disabled">' + (disabled ? 'Enable' : 'Disable') + '</button>' : '') +
                                    (canManage ? '<button type="button" class="ama-action-button danger" data-ama-more-action="uninstall">Uninstall</button>' : '') +
                                    '<div class="ama-action-status"></div>' +
                                '</div>';

                            bindAmaModalClose(modal);

                            const status = modal.querySelector('.ama-action-status');
                            const setStatus = (message) => {
                                if (status) status.textContent = message || '';
                            };

                            const preferences = modal.querySelector('[data-ama-more-action="preferences"]');
                            if (preferences) {
                                preferences.onclick = () => {
                                    modal.remove();
                                    showInstalledPreferences(card);
                                };
                            }

                            const checkUpdate = modal.querySelector('[data-ama-more-action="check-update"]');
                            if (checkUpdate) {
                                checkUpdate.onclick = async () => {
                                    checkUpdate.disabled = true;
                                    setStatus('Checking for updates...');

                                    try {
                                        const fetched = await postExtensionAction('/api/v1/extensions/external/fetch', {
                                            manifestUri
                                        });

                                        if (fetched && fetched.version && fetched.version !== version) {
                                            setStatus('Update available: ' + fetched.version);
                                        } else {
                                            setStatus('The extension is up to date.');
                                        }
                                    } catch (error) {
                                        setStatus(error && error.message ? error.message : 'Could not check for updates.');
                                    } finally {
                                        checkUpdate.disabled = false;
                                    }
                                };
                            }

                            const toggleDisabled = modal.querySelector('[data-ama-more-action="toggle-disabled"]');
                            if (toggleDisabled) {
                                toggleDisabled.onclick = async () => {
                                    toggleDisabled.disabled = true;
                                    setStatus((disabled ? 'Enabling' : 'Disabling') + ' extension...');

                                    try {
                                        await postExtensionAction('/api/v1/extensions/external/disabled', {
                                            id: actionId,
                                            disabled: !disabled
                                        });
                                        clearExtensionCache();
                                        setStatus(disabled ? 'Extension enabled.' : 'Extension disabled.');
                                    } catch (error) {
                                        setStatus(error && error.message ? error.message : 'Could not update this extension.');
                                    } finally {
                                        toggleDisabled.disabled = false;
                                    }
                                };
                            }

                            const uninstall = modal.querySelector('[data-ama-more-action="uninstall"]');
                            if (uninstall) {
                                uninstall.onclick = async () => {
                                    if (!window.confirm('Uninstall ' + name + '?')) return;

                                    uninstall.disabled = true;
                                    setStatus('Uninstalling extension...');

                                    try {
                                        await postExtensionAction('/api/v1/extensions/external/uninstall', {
                                            id: actionId
                                        });
                                        clearExtensionCache();
                                        setStatus('Extension uninstalled.');
                                    } catch (error) {
                                        setStatus(error && error.message ? error.message : 'Could not uninstall this extension.');
                                    } finally {
                                        uninstall.disabled = false;
                                    }
                                };
                            }
                        } catch (error) {
                            if (modal) {
                                modal.innerHTML =
                                    '<button type="button" class="ama-modal-close">Close</button>' +
                                    '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">' + escapeHtml(data.name) + '</h2>' +
                                    '<p>' + escapeHtml(error && error.message ? error.message : 'Could not load extension actions.') + '</p>';
                                bindAmaModalClose(modal);
                            }
                        }
                    }

                    async function hasPreferencesForCard(card) {
                        const extension = await findExtensionForCard(card);
                        return !!(extension && extension.userConfig);
                    }

                    function getKnownDocumentationUrl(data) {
                        const id = normalizeExtensionText(data && data.id);
                        const name = normalizeExtensionText(data && data.name);

                        if (id === 'seautils-kolex06-version' || name === 'seautils kolex06-version') {
                            return 'https://raw.githubusercontent.com/Kolex06/Seanime-Stuff/refs/heads/main/plugins/SeaUtilsKolex06Version/README.md';
                        }

                        if (id === 'asunatracks-sync' || name === 'asunatracks sync') {
                            return 'https://raw.githubusercontent.com/Kolex06/Seanime-Stuff/refs/heads/main/plugins/AsunaTracksSync/README.md';
                        }

                        return '';
                    }

                    async function getDocumentationUrlForCard(card) {
                        const data = getExtensionCardData(card);
                        let url = getKnownDocumentationUrl(data);
                        if (url) return url;

                        const extension = await findExtensionForCard(card);
                        url = getExtensionDocumentationUrl(extension);
                        if (url) return url;

                        if (data.id) {
                            const marketplaceExtension = await fetchMarketplaceExtensionById(data.id);
                            url = getExtensionDocumentationUrl(marketplaceExtension);
                        }

                        return url;
                    }

                    async function hasDocumentationForCard(card) {
                        return !!(await getDocumentationUrlForCard(card));
                    }

                    function sourceCardHasNativePreferences(card) {
                        if (!card || !card.querySelectorAll) return false;

                        const topRight = card.querySelector('[class*="top-3"][class*="right-3"]');
                        if (!topRight) return false;

                        const rows = Array.from(topRight.children).filter(row => {
                            return row && row.querySelectorAll && row.querySelectorAll('button').length > 0;
                        });

                        const actionRow = rows[0];
                        if (!actionRow) return false;

                        return actionRow.querySelectorAll('button').length > 1;
                    }

                    function getSavedConfigValue(extUserConfig, field) {
                        const savedValues = extUserConfig && extUserConfig.savedUserConfig && extUserConfig.savedUserConfig.values
                            ? extUserConfig.savedUserConfig.values
                            : {};

                        if (Object.prototype.hasOwnProperty.call(savedValues, field.name)) {
                            return savedValues[field.name];
                        }

                        return field.default || '';
                    }

                    function appendConfigHelp(parent, field) {
                        if (!field.default) return;

                        const help = document.createElement('div');
                        help.className = 'ama-config-help';
                        help.textContent = 'Default: ' + String(field.default);
                        parent.appendChild(help);
                    }

                    function renderPreferencesModal(modal, extension, extUserConfig) {
                        const fields = extUserConfig && extUserConfig.userConfig && Array.isArray(extUserConfig.userConfig.fields)
                            ? extUserConfig.userConfig.fields
                            : [];

                        modal.innerHTML =
                            '<button type="button" class="ama-modal-close">Close</button>' +
                            '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">Preferences</h2>' +
                            '<p style="color:rgba(255,255,255,.65); margin:0 0 16px 0;">' + escapeHtml(extension.name) + '</p>' +
                            '<div class="ama-config-form"></div>';

                        modal.querySelector('.ama-modal-close').onclick = () => modal.remove();

                        const form = modal.querySelector('.ama-config-form');
                        const controls = {};

                        fields.forEach(field => {
                            const savedValue = getSavedConfigValue(extUserConfig, field);

                            if (field.type === 'switch') {
                                const row = document.createElement('div');
                                row.className = 'ama-config-switch-row';

                                const label = document.createElement('label');
                                label.className = 'ama-config-switch-label';
                                label.textContent = field.label || field.name;

                                const input = document.createElement('input');
                                input.type = 'checkbox';
                                input.checked = String(savedValue || field.default || '') === 'true';

                                row.appendChild(label);
                                row.appendChild(input);
                                form.appendChild(row);
                                appendConfigHelp(form, field);

                                controls[field.name] = () => input.checked ? 'true' : 'false';
                                return;
                            }

                            const wrapper = document.createElement('div');
                            wrapper.className = 'ama-config-field';

                            const label = document.createElement('label');
                            label.textContent = field.label || field.name;
                            wrapper.appendChild(label);

                            if (field.type === 'select' && Array.isArray(field.options)) {
                                const select = document.createElement('select');
                                select.className = 'ama-config-select';

                                field.options.forEach(option => {
                                    const item = document.createElement('option');
                                    item.value = String(option.value || '');
                                    item.textContent = String(option.label || option.value || '');
                                    select.appendChild(item);
                                });

                                select.value = String(savedValue || field.default || '');
                                wrapper.appendChild(select);
                                controls[field.name] = () => select.value;
                            } else {
                                const input = document.createElement('input');
                                input.className = 'ama-config-input';
                                input.type = 'text';
                                input.value = String(savedValue || field.default || '');
                                wrapper.appendChild(input);
                                controls[field.name] = () => input.value;
                            }

                            appendConfigHelp(wrapper, field);
                            form.appendChild(wrapper);
                        });

                        const actions = document.createElement('div');
                        actions.className = 'ama-config-actions';

                        const save = document.createElement('button');
                        save.type = 'button';
                        save.className = 'ama-config-save';
                        save.textContent = 'Save';

                        const status = document.createElement('span');
                        status.className = 'ama-config-status';

                        save.onclick = async () => {
                            const values = {};

                            fields.forEach(field => {
                                values[field.name] = controls[field.name] ? controls[field.name]() : String(field.default || '');
                            });

                            save.disabled = true;
                            status.textContent = 'Saving...';

                            try {
                                await fetchSeanime('/api/v1/extensions/user-config', {
                                    method: 'POST',
                                    headers: {
                                        'Content-Type': 'application/json'
                                    },
                                    body: JSON.stringify({
                                        id: extension.id,
                                        version: extUserConfig && extUserConfig.userConfig ? extUserConfig.userConfig.version || 0 : 0,
                                        values
                                    })
                                });

                                status.textContent = 'Preferences saved.';
                            } catch (error) {
                                status.textContent = error && error.message ? error.message : 'Could not save preferences.';
                            } finally {
                                save.disabled = false;
                            }
                        };

                        actions.appendChild(save);
                        actions.appendChild(status);
                        form.appendChild(actions);
                    }

                    async function showInstalledPreferences(card) {
                        const data = getExtensionCardData(card);
                        const modal = openAmaModal('Preferences', '<p>Loading preferences for ' + escapeHtml(data.name) + '...</p>');

                        try {
                            const extension = await findExtensionForCard(card);

                            if (!extension || !extension.id) {
                                throw new Error('Could not find this extension.');
                            }

                            if (!extension.userConfig) {
                                throw new Error('This extension has no preferences.');
                            }

                            const extUserConfig = await fetchSeanime('/api/v1/extensions/user-config/' + encodeURIComponent(extension.id), {
                                method: 'GET'
                            });

                            renderPreferencesModal(modal, extension, extUserConfig);
                        } catch (error) {
                            if (modal) {
                                modal.innerHTML =
                                    '<button type="button" class="ama-modal-close">Close</button>' +
                                    '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">Preferences</h2>' +
                                    '<p>' + escapeHtml(error && error.message ? error.message : 'Could not load preferences for this extension.') + '</p>';
                                modal.querySelector('.ama-modal-close').onclick = () => modal.remove();
                            }
                        }
                    }


                    function normalizeCodePayload(value) {
                        if (!value) return '';

                        if (typeof value === 'string') return value;

                        if (value && typeof value === 'object') {
                            if (typeof value.payload === 'string') return value.payload;
                            if (typeof value.code === 'string') return value.code;
                            if (typeof value.source === 'string') return value.source;
                            if (typeof value.data === 'string') return value.data;
                        }

                        return '';
                    }

                    async function fetchManifestPayload(manifestUri) {
                        const payloadUri = String(manifestUri || '').trim();
                        if (!payloadUri || payloadUri === 'builtin') return '';

                        const text = await fetchTextWithTimeout(payloadUri, 10000);
                        if (!text) return '';

                        try {
                            const manifest = JSON.parse(text);
                            return normalizeCodePayload(manifest.payload);
                        } catch (_) {
                            return '';
                        }
                    }

                    async function resolveCodePayloadForCard(card, extension, data) {
                        const extensionId = (extension && extension.id) || data.id;
                        const manifestUri = getExtensionManifestUri(extension);
                        const directPayload = normalizeCodePayload(extension && extension.payload);

                        if (directPayload) return directPayload;

                        if (extensionId) {
                            try {
                                const apiPayload = await fetchSeanime('/api/v1/extensions/payload/' + encodeURIComponent(extensionId), {
                                    method: 'GET'
                                });

                                const normalizedApiPayload = normalizeCodePayload(apiPayload);
                                if (normalizedApiPayload) return normalizedApiPayload;
                            } catch (_) {}
                        }

                        const manifestPayload = await fetchManifestPayload(manifestUri);
                        if (manifestPayload) return manifestPayload;

                        if (isKolex06VersionExtension(data, extension)) {
                            const fallbackPayload = await fetchManifestPayload('https://raw.githubusercontent.com/Kolex06/Seanime-Stuff/refs/heads/main/plugins/SeaUtils-Kolex06-Version.json');
                            if (fallbackPayload) return fallbackPayload;
                        }

                        if (!extensionId) {
                            throw new Error('Could not find this extension ID.');
                        }

                        return '';
                    }

                    async function showInstalledDocumentation(card) {
                        const data = getExtensionCardData(card);

                        try {
                            const url = await getDocumentationUrlForCard(card);

                            if (!url) {
                                openAmaModal('Documentation', '<p>No documentation found for ' + escapeHtml(data.name) + '.</p>');
                                return;
                            }

                            window.open(url, '_blank', 'noopener,noreferrer');
                        } catch (error) {
                            openAmaModal('Documentation', '<p>' + escapeHtml(error && error.message ? error.message : 'Could not open documentation.') + '</p>');
                        }
                    }

                    async function showInstalledCode(card) {
                        const data = getExtensionCardData(card);

                        const modal = openAmaModal('Code', '<p>Loading code for ' + escapeHtml(data.name) + '...</p>');

                        try {
                            const extension = await findExtensionForCard(card);
                            const payload = await resolveCodePayloadForCard(card, extension, data);

                            if (modal) {
                                modal.innerHTML =
                                    '<button type="button" class="ama-modal-close">Close</button>' +
                                    '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">Code</h2>' +
                                    '<p style="color:rgba(255,255,255,.65); margin:0;">' + escapeHtml(data.name) + '</p>' +
                                    '<pre class="ama-code-preview">' + escapeHtml(payload || 'No code found.') + '</pre>';
                                modal.querySelector('.ama-modal-close').onclick = () => modal.remove();
                            }
                        } catch (error) {
                            if (modal) {
                                modal.innerHTML =
                                    '<button type="button" class="ama-modal-close">Close</button>' +
                                    '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">Code</h2>' +
                                    '<p>' + escapeHtml(error && error.message ? error.message : 'Could not load code for this extension.') + '</p>';
                                modal.querySelector('.ama-modal-close').onclick = () => modal.remove();
                            }
                        }
                    }


                    function readMarketplaceUrl() {
                        try {
                            const raw = window.localStorage.getItem('marketplace-url');
                            if (!raw) return '';

                            try {
                                return JSON.parse(raw) || '';
                            } catch (_) {
                                return raw || '';
                            }
                        } catch (_) {
                            return '';
                        }
                    }

                    function writeMarketplaceUrl(value) {
                        try {
                            if (value) {
                                window.localStorage.setItem('marketplace-url', JSON.stringify(value));
                            } else {
                                window.localStorage.removeItem('marketplace-url');
                            }
                        } catch (_) {}
                    }

                    function applyBas1874MarketplacePreference() {
                        const current = readMarketplaceUrl();

                        if (featureSettings.useBas1874Marketplace) {
                            if (current !== BAS1874_MARKETPLACE_URL) {
                                writeMarketplaceUrl(BAS1874_MARKETPLACE_URL);
                                marketplaceExtensionsPromise = null;
                                bas1874MarketplaceMetadataPromise = null;
                            }
                            return;
                        }

                        if (current === BAS1874_MARKETPLACE_URL) {
                            writeMarketplaceUrl('');
                            marketplaceExtensionsPromise = null;
                            bas1874MarketplaceMetadataPromise = null;
                        }
                    }

                    function getMarketplaceEndpoint() {
                        const marketplaceUrl = readMarketplaceUrl();

                        return marketplaceUrl
                            ? '/api/v1/extensions/marketplace?marketplace=' + encodeURIComponent(marketplaceUrl)
                            : '/api/v1/extensions/marketplace';
                    }

                    function loadMarketplaceExtensions() {
                        if (marketplaceExtensionsPromise) return marketplaceExtensionsPromise;

                        marketplaceExtensionsPromise = fetchSeanime(getMarketplaceEndpoint(), {
                            method: 'GET'
                        }).then(extensions => {
                            return Array.isArray(extensions) ? extensions : [];
                        }).catch(() => {
                            marketplaceExtensionsPromise = null;
                            return [];
                        });

                        return marketplaceExtensionsPromise;
                    }

                    async function fetchMarketplaceExtensionById(extensionId) {
                        const extensions = await loadMarketplaceExtensions();

                        return extensions.find(extension => extension && extension.id === extensionId) || null;
                    }

                    function normalizeMarketplaceLookupText(value) {
                        return String(value || '')
                            .trim()
                            .replace(/^(id|author|language|version|manifest url)\s*:\s*/i, '')
                            .trim()
                            .toLowerCase();
                    }

                    function normalizeMarketplaceSlug(value) {
                        return normalizeMarketplaceLookupText(value)
                            .replace(/&/g, ' and ')
                            .replace(/[^a-z0-9]+/g, '-')
                            .replace(/^-+|-+$/g, '');
                    }

                    function slugContainsSlug(containerSlug, itemSlug) {
                        if (!containerSlug || !itemSlug || itemSlug.length < 3) return false;

                        return ('-' + containerSlug + '-').includes('-' + itemSlug + '-');
                    }

                    function loadBas1874MarketplaceMetadata() {
                        if (bas1874MarketplaceMetadataPromise) return bas1874MarketplaceMetadataPromise;

                        bas1874MarketplaceMetadataPromise = fetch(BAS1874_MARKETPLACE_URL, {
                            cache: 'no-store'
                        }).then(response => {
                            if (!response.ok) throw new Error('Could not load Bas1874 marketplace metadata.');
                            return response.json();
                        }).then(data => {
                            if (Array.isArray(data)) {
                                bas1874MarketplaceMetadataCache = data;
                                return data;
                            }

                            if (data && Array.isArray(data.extensions)) {
                                bas1874MarketplaceMetadataCache = data.extensions;
                                return data.extensions;
                            }

                            if (data && Array.isArray(data.items)) {
                                bas1874MarketplaceMetadataCache = data.items;
                                return data.items;
                            }

                            bas1874MarketplaceMetadataCache = [];
                            return [];
                        }).catch(() => {
                            bas1874MarketplaceMetadataPromise = null;
                            bas1874MarketplaceMetadataCache = [];
                            return [];
                        });

                        return bas1874MarketplaceMetadataPromise;
                    }

                    function isInstalledExtensionsView() {
                        return !!document.querySelector('input[placeholder="Search installed extensions..."]');
                    }

                    function isMarketplaceExtraCard(card) {
                        if (!card || !card.isConnected) return false;
                        if (card.closest && card.closest('.ama-modal')) return true;
                        if (isInstalledExtensionsView()) return false;

                        return true;
                    }

                    async function fetchBas1874MetadataForMarketplaceCard(card) {
                        const data = getExtensionCardData(card);
                        const cardId = normalizeMarketplaceLookupText(data.id);

                        if (!cardId) return null;

                        const extensions = await loadBas1874MarketplaceMetadata();

                        return extensions.find(extension => normalizeMarketplaceLookupText(extension && extension.id) === cardId) || null;
                    }

                    function getBas1874MetadataForCardData(data) {
                        if (!featureSettings.useBas1874Marketplace || !Array.isArray(bas1874MarketplaceMetadataCache)) return null;

                        const cardId = normalizeMarketplaceLookupText(data && data.id);
                        const cardName = normalizeMarketplaceLookupText(data && data.name);
                        const cardAuthor = normalizeMarketplaceLookupText(data && data.author);
                        const cardIdSlug = normalizeMarketplaceSlug(data && data.id);
                        const cardNameSlug = normalizeMarketplaceSlug(data && data.name);

                        if (cardId) {
                            const byId = bas1874MarketplaceMetadataCache.find(extension => {
                                const extensionId = normalizeMarketplaceLookupText(extension && extension.id);
                                const extensionIdSlug = normalizeMarketplaceSlug(extension && extension.id);
                                const extensionNameSlug = normalizeMarketplaceSlug(extension && extension.name);

                                return extensionId === cardId || extensionIdSlug === cardIdSlug || extensionNameSlug === cardIdSlug;
                            });

                            if (byId) return byId;
                        }

                        if (cardName) {
                            return bas1874MarketplaceMetadataCache.find(extension => {
                                const nameMatches = normalizeMarketplaceLookupText(extension && extension.name) === cardName;
                                const slugMatches = normalizeMarketplaceSlug(extension && extension.name) === cardNameSlug || normalizeMarketplaceSlug(extension && extension.id) === cardNameSlug;
                                const authorMatches = !cardAuthor || normalizeMarketplaceLookupText(extension && extension.author) === cardAuthor;

                                return (nameMatches && authorMatches) || slugMatches;
                            }) || null;
                        }

                        return null;
                    }

                    function getBas1874StatusGroupForCard(item, data) {
                        if (item && item.querySelector) {
                            if (item.querySelector('.ama-marketplace-extra-badge.broken')) return 'Broken Extensions';
                            if (item.querySelector('.ama-marketplace-extra-badge.deprecated')) return 'Deprecated Extensions';
                        }

                        const extension = getBas1874MetadataForCardData(data);
                        if (!extension) return '';

                        if (extension.brokenTag) return 'Broken Extensions';
                        if (extension.deprecatedTag) return 'Deprecated Extensions';

                        return '';
                    }

                    function createMarketplaceExtraBadge(text, className, href) {
                        const el = href ? document.createElement('a') : document.createElement('span');
                        el.className = 'ama-marketplace-extra-badge' + (className ? ' ' + className : '');
                        el.textContent = text;

                        if (href) {
                            el.href = href;
                            el.target = '_blank';
                            el.rel = 'noreferrer';
                            el.classList.add('link');
                        }

                        return el;
                    }

                    function appendBas1874ExtraBadges(extra, extension, includeNote) {
                        if (!extra || !extension) return false;

                        if (extension.official) extra.appendChild(createMarketplaceExtraBadge('Official', 'official'));
                        if (extension.workingTag) extra.appendChild(createMarketplaceExtraBadge('Working', 'working'));
                        if (extension.brokenTag) extra.appendChild(createMarketplaceExtraBadge('Broken', 'broken'));
                        if (extension.deprecatedTag) extra.appendChild(createMarketplaceExtraBadge('Deprecated', 'deprecated'));

                        if (extension.flags) {
                            extra.appendChild(createMarketplaceExtraBadge('VT ' + String(extension.flags), 'scan', extension.permalink || ''));
                        } else if (extension.permalink) {
                            extra.appendChild(createMarketplaceExtraBadge('VirusTotal', 'scan', extension.permalink));
                        }

                        if (extension.scannedOnVersion) {
                            extra.appendChild(createMarketplaceExtraBadge('Scanned ' + String(extension.scannedOnVersion), 'scan'));
                        }

                        if (extension.lastWorkingVersion) {
                            extra.appendChild(createMarketplaceExtraBadge('Last working ' + String(extension.lastWorkingVersion), 'working'));
                        }

                        if (!extra.children.length) return false;

                        extra.appendChild(createMarketplaceExtraBadge('Bas1874 metadata', 'scan'));

                        if (includeNote) {
                            const note = document.createElement('div');
                            note.className = 'ama-marketplace-extra-note';
                            note.textContent = 'Status and scan data are hints. Review plugins yourself before installing.';
                            extra.appendChild(note);
                        }

                        return true;
                    }

                    function renderMarketplaceExtraInfo(card, extension) {
                        if (!card || !extension || !isMarketplaceExtraCard(card)) return;

                        card.querySelectorAll('.ama-marketplace-extra').forEach(el => el.remove());

                        if (!featureSettings.useBas1874Marketplace) return;

                        const extra = document.createElement('div');
                        extra.className = 'ama-marketplace-extra';

                        if (!appendBas1874ExtraBadges(extra, extension, true)) return;

                        card.appendChild(extra);
                    }

                    function getInstalledBadgeRow(card) {
                        if (!card || !card.querySelectorAll) return null;

                        const rows = Array.from(card.querySelectorAll('.flex')).filter(row => {
                            return row && row.querySelectorAll && row.querySelectorAll('.UI-Badge__root').length > 0;
                        });

                        if (!rows.length) return null;

                        const idRow = rows.find(row => /(^|\s)ID\s*:/i.test(row.textContent || ''));
                        return idRow || rows[rows.length - 1];
                    }

                    async function fetchBas1874MetadataForInstalledCard(card) {
                        const data = getExtensionCardData(card);

                        await loadBas1874MarketplaceMetadata();

                        const directMatch = getBas1874MetadataForCardData(data);
                        if (directMatch) return directMatch;

                        const cardSlug = normalizeMarketplaceSlug(card.innerText || card.textContent || '');
                        if (!cardSlug) return null;

                        return bas1874MarketplaceMetadataCache.find(extension => {
                            const idSlug = normalizeMarketplaceSlug(extension && extension.id);
                            const nameSlug = normalizeMarketplaceSlug(extension && extension.name);

                            return slugContainsSlug(cardSlug, idSlug) || slugContainsSlug(cardSlug, nameSlug);
                        }) || null;
                    }

                    function renderInstalledMarketplaceExtraInfo(card, extension) {
                        if (!card || !extension || !isInstalledExtensionsView()) return;

                        card.querySelectorAll('.ama-installed-marketplace-extra').forEach(el => el.remove());

                        if (!featureSettings.useBas1874Marketplace) return;

                        const extra = document.createElement('div');
                        extra.className = 'ama-installed-marketplace-extra';

                        if (!appendBas1874ExtraBadges(extra, extension, false)) return;

                        const badgeRow = getInstalledBadgeRow(card);

                        if (badgeRow && badgeRow.parentElement) {
                            badgeRow.parentElement.insertBefore(extra, badgeRow.nextSibling);
                        } else {
                            card.appendChild(extra);
                        }
                    }

                    function enhanceInstalledMarketplaceExtraInfo(card) {
                        if (!card || !isInstalledExtensionsView()) return;

                        if (!featureSettings.useBas1874Marketplace) {
                            card.querySelectorAll('.ama-installed-marketplace-extra').forEach(el => el.remove());
                            delete card.dataset.amaInstalledMarketplaceExtraEnhanced;
                            return;
                        }

                        if (card.dataset.amaInstalledMarketplaceExtraEnhanced === marketplaceEnhancementVersion) return;
                        card.dataset.amaInstalledMarketplaceExtraEnhanced = marketplaceEnhancementVersion;

                        fetchBas1874MetadataForInstalledCard(card).then(extension => {
                            if (!card.isConnected || !extension) return;
                            renderInstalledMarketplaceExtraInfo(card, extension);
                            hideStatusSourceCard(card, extension);
                            scheduleMarketplaceStatusSections();
                        }).catch(() => {});
                    }

                    function hideStatusSourceCard(card, extension) {
                        if (!card || !extension || card.closest('.ama-status-section') || card.closest('.ama-modal')) return;

                        let group = '';

                        if (extension.brokenTag) {
                            group = 'Broken Extensions';
                        } else if (extension.deprecatedTag) {
                            group = 'Deprecated Extensions';
                        }

                        if (!group) return;

                        card.classList.add('ama-status-source-hidden');
                        card.dataset.amaStatusGroup = group;
                    }

                    function enhanceMarketplaceExtraInfo(card) {
                        if (!card || !featureSettings.useBas1874Marketplace || !isMarketplaceExtraCard(card)) {
                            if (card && card.querySelectorAll) {
                                card.querySelectorAll('.ama-marketplace-extra').forEach(el => el.remove());
                            }
                            return;
                        }

                        if (card.dataset.amaMarketplaceExtraEnhanced === marketplaceEnhancementVersion) return;
                        card.dataset.amaMarketplaceExtraEnhanced = marketplaceEnhancementVersion;

                        fetchBas1874MetadataForMarketplaceCard(card).then(extension => {
                            if (!card.isConnected || !extension) return;
                            renderMarketplaceExtraInfo(card, extension);
                            hideStatusSourceCard(card, extension);
                            scheduleMarketplaceStatusSections();
                        }).catch(() => {});
                    }

                    async function installMarketplaceExtension(card) {
                        const data = getExtensionCardData(card);

                        if (!data.id) {
                            openAmaModal('Download', '<p>Could not find this extension ID.</p>');
                            return;
                        }

                        const modal = openAmaModal('Download', '<p>Preparing download for ' + escapeHtml(data.name) + '...</p>');

                        try {
                            const extension = await fetchMarketplaceExtensionById(data.id);
                            const manifestUri = extension && (extension.manifestURI || extension.manifestUri);

                            if (!manifestUri) {
                                throw new Error('Could not find manifest URL.');
                            }

                            const json = await fetchSeanime('/api/v1/extensions/external/install', {
                                method: 'POST',
                                headers: {
                                    'Content-Type': 'application/json'
                                },
                                body: JSON.stringify({
                                    manifestUri
                                })
                            });

                            clearExtensionCache();
                            refreshExtensionUpdateStateForData(data);

                            if (modal) {
                                const message = json && json.message ? json.message : 'Extension installed.';
                                modal.innerHTML =
                                    '<button type="button" class="ama-modal-close">Close</button>' +
                                    '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">Downloaded</h2>' +
                                    '<p>' + escapeHtml(message) + '</p>';
                                modal.querySelector('.ama-modal-close').onclick = () => modal.remove();
                            }
                        } catch (error) {
                            if (modal) {
                                modal.innerHTML =
                                    '<button type="button" class="ama-modal-close">Close</button>' +
                                    '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">Download failed</h2>' +
                                    '<p>' + escapeHtml(error && error.message ? error.message : 'Could not install this extension.') + '</p>';
                                modal.querySelector('.ama-modal-close').onclick = () => modal.remove();
                            }
                        }
                    }

                    function getCatalogActionSourceId(sourceCard) {
                        if (!sourceCard.dataset.amaCatalogSourceId) {
                            catalogActionSourceCounter += 1;
                            sourceCard.dataset.amaCatalogSourceId = String(catalogActionSourceCounter);
                        }

                        catalogActionSources.set(sourceCard.dataset.amaCatalogSourceId, sourceCard);

                        return sourceCard.dataset.amaCatalogSourceId;
                    }

                    function getFallbackActionCard(button) {
                        const wrapper = button && button.closest ? button.closest('.ama-catalog-card-wrap') : null;
                        if (!wrapper || !wrapper.querySelector) return null;

                        return wrapper.querySelector('.group\\\\/extension-card');
                    }

                    function getFallbackActionId(button) {
                        const wrapper = button && button.closest ? button.closest('.ama-catalog-card-wrap') : null;
                        return (wrapper && wrapper.dataset && wrapper.dataset.amaActionExtensionId) || '';
                    }

                    function findCurrentInstalledCardById(extensionId) {
                        const id = normalizeMarketplaceLookupText(extensionId);
                        if (!id) return null;

                        const cards = Array.from(document.querySelectorAll('.group\\\\/extension-card'));
                        return cards.find(card => {
                            if (card.closest && card.closest('.ama-modal')) return false;
                            if (card.closest && card.closest('.ama-status-section')) return false;
                            const data = getExtensionCardData(card);
                            return normalizeMarketplaceLookupText(data.id) === id;
                        }) || null;
                    }

                    function runCatalogAction(action, sourceId, fallbackCard, extensionId) {
                        const storedSourceCard = catalogActionSources.get(sourceId);
                        let sourceCard = (
                            storedSourceCard && storedSourceCard.isConnected
                                ? storedSourceCard
                                : (findCurrentInstalledCardById(extensionId) || fallbackCard)
                        );

                        if (isInstalledExtensionsView() && action !== 'download') {
                            sourceCard = findCurrentInstalledCardById(extensionId) || sourceCard;
                        }

                        if (!sourceCard || !sourceCard.isConnected) {
                            openAmaModal('Action unavailable', '<p>This card is no longer available. Close Full Catalog, open it again, and try once more.</p>');
                            return;
                        }

                        if (action === 'details' || action === 'more') {
                            showInstalledMore(sourceCard);
                            return;
                        }

                        if (action === 'code') {
                            showInstalledCode(sourceCard);
                            return;
                        }

                        if (action === 'documentation') {
                            showInstalledDocumentation(sourceCard);
                            return;
                        }

                        if (action === 'preferences') {
                            showInstalledPreferences(sourceCard);
                            return;
                        }

                        if (action === 'download') {
                            installMarketplaceExtension(sourceCard);
                        }
                    }

                    function ensureCatalogActionHandler() {
                        if (catalogActionHandlerBound) return;

                        catalogActionHandlerBound = true;

                        document.addEventListener('click', event => {
                            const target = event.target;
                            const button = target && target.closest ? target.closest('.ama-clone-action') : null;

                            if (!button) return;
                            if (button.dataset.amaDirectAction === 'true') return;

                            event.preventDefault();
                            event.stopPropagation();
                            if (event.stopImmediatePropagation) event.stopImmediatePropagation();

                            runCatalogAction(button.dataset.amaAction, button.dataset.amaSourceId, getFallbackActionCard(button), button.dataset.amaExtensionId || getFallbackActionId(button));
                        }, true);
                    }

                    function addCloneAction(actions, title, icon, action, sourceId, extensionId) {
                        if (actions.querySelector('[data-ama-action="' + action + '"]')) return;

                        const button = document.createElement('button');
                        button.type = 'button';
                        button.className = 'ama-clone-action';
                        button.title = title;
                        button.setAttribute('aria-label', title);
                        button.dataset.amaAction = action;
                        button.dataset.amaSourceId = sourceId;
                        if (extensionId) button.dataset.amaExtensionId = extensionId;
                        button.innerHTML = icon;
                        button.onpointerdown = event => {
                            event.stopPropagation();
                            if (event.stopImmediatePropagation) event.stopImmediatePropagation();
                        };
                        button.onmousedown = event => {
                            event.stopPropagation();
                            if (event.stopImmediatePropagation) event.stopImmediatePropagation();
                        };
                        button.onmouseup = event => {
                            event.stopPropagation();
                            if (event.stopImmediatePropagation) event.stopImmediatePropagation();
                        };
                        button.addEventListener('click', event => {
                            event.preventDefault();
                            event.stopPropagation();
                            if (event.stopImmediatePropagation) event.stopImmediatePropagation();
                            if (button.dataset.amaDirectAction === 'true') return;
                            runCatalogAction(action, sourceId, getFallbackActionCard(button), button.dataset.amaExtensionId || getFallbackActionId(button));
                        }, true);

                        actions.appendChild(button);
                    }

                    function bindDirectStatusActions(actions) {
                        if (!actions || !actions.querySelectorAll) return;

                        function runDirectStatusAction(button, event) {
                            event.preventDefault();
                            event.stopPropagation();
                            if (event.stopImmediatePropagation) event.stopImmediatePropagation();

                            const action = button.dataset.amaAction;
                            const extensionId = button.dataset.amaExtensionId || getFallbackActionId(button);
                            const card = (
                                isInstalledExtensionsView() && action !== 'download'
                                    ? findCurrentInstalledCardById(extensionId)
                                    : null
                            ) || getFallbackActionCard(button);
                            if (!card) return;

                            if (action === 'details' || action === 'more') {
                                showInstalledMore(card);
                                return;
                            }

                            if (action === 'code') {
                                showInstalledCode(card);
                                return;
                            }

                            if (action === 'documentation') {
                                showInstalledDocumentation(card);
                                return;
                            }

                            if (action === 'preferences') {
                                showInstalledPreferences(card);
                            }
                        }

                        actions.querySelectorAll('.ama-clone-action').forEach(button => {
                            button.dataset.amaDirectAction = 'true';
                            button.onpointerdown = event => runDirectStatusAction(button, event);
                            button.addEventListener('click', event => {
                                runDirectStatusAction(button, event);
                            }, true);
                        });
                    }

                    function createMarketplaceCloneActions(sourceCard, isInstalledCatalog, forceMarketplaceActions) {
                        const actions = document.createElement('div');
                        actions.className = 'ama-clone-actions';

                        const liveButtons = Array.from(sourceCard.querySelectorAll('button'));
                        const isInstalled = sourceCard.dataset.amaInstalledMarketplaceCard === 'true' || !!sourceCard.querySelector('button[disabled]');
                        const data = getExtensionCardData(sourceCard);
                        const sourceId = getCatalogActionSourceId(sourceCard);
                        const extensionId = data.id || '';

                        ensureCatalogActionHandler();

                        if (isInstalledCatalog || isInstalled) {
                            if (sourceCardHasNativePreferences(sourceCard)) {
                                addCloneAction(actions, 'Preferences', SETTINGS_ICON, 'preferences', sourceId, extensionId);
                            } else {
                                hasPreferencesForCard(sourceCard).then(hasPreferences => {
                                    if (!hasPreferences) return;

                                    addCloneAction(actions, 'Preferences', SETTINGS_ICON, 'preferences', sourceId, extensionId);
                                });
                            }

                            addCloneAction(actions, 'More', MORE_ICON, 'more', sourceId, extensionId);

                            if (getKnownDocumentationUrl(data)) {
                                addCloneAction(actions, 'Documentation', DOC_ICON, 'documentation', sourceId, extensionId);
                            } else {
                                hasDocumentationForCard(sourceCard).then(hasDocumentation => {
                                    if (!hasDocumentation) return;

                                    addCloneAction(actions, 'Documentation', DOC_ICON, 'documentation', sourceId, extensionId);
                                });
                            }

                            addCloneAction(actions, 'Code', CODE_ICON, 'code', sourceId, extensionId);
                        } else if (!isInstalled && (liveButtons.length || forceMarketplaceActions)) {
                            addCloneAction(actions, 'Download', DOWNLOAD_ICON, 'download', sourceId, extensionId);
                        }

                        return actions.children.length ? actions : null;
                    }

                    function enhanceInstalledMarketplaceCardActions(card) {
                        if (!card || !featureSettings.betterMarketplace) return;
                        if (isInstalledExtensionsView()) return;
                        if (card.closest && (card.closest('.ama-modal') || card.closest('.ama-catalog-card-wrap') || card.closest('.ama-status-section'))) return;
                        if (!isMarketplaceExtraCard(card)) return;
                        if (!card.querySelector('button[disabled]')) return;

                        const existing = card.querySelector(':scope > .ama-live-card-actions');
                        if (existing) return;

                        const actions = createMarketplaceCloneActions(card, true);
                        if (!actions) return;

                        actions.classList.add('ama-live-card-actions');
                        card.appendChild(actions);
                    }

                    function getMarketplaceSections(includeStatusSections) {
                        return Array.from(document.querySelectorAll(cardQuery)).map(card => {
                            if (!isMarketplaceSectionCard(card)) return null;
                            if (!includeStatusSections && card.dataset && card.dataset.amaStatusSection) return null;

                            const titleEl = card.querySelector('.ama-header-left h3, h3');
                            const grid = card.querySelector('.ama-extension-carousel, .grid');

                            if (!titleEl || !grid) return null;

                            const items = Array.from(grid.querySelectorAll('.group\\\\/extension-card'));
                            if (!items.length) return null;

                            return {
                                title: titleEl.textContent.trim() || 'Extensions',
                                card,
                                grid,
                                items,
                            };
                        }).filter(Boolean);
                    }

                    function groupMarketplaceSectionsByAuthor(sections, useStatusGroups) {
                        const grouped = new Map();
                        const hasStatusSections = (sections || []).some(section => section && section.card && section.card.dataset && section.card.dataset.amaStatusSection);
                        const shouldUseStatusGroups = useStatusGroups !== false;

                        (sections || []).forEach(section => {
                            section.items.forEach(item => {
                                if (hasStatusSections && item.classList && item.classList.contains('ama-status-source-hidden')) return;

                                const data = getExtensionCardData(item);
                                const author = (shouldUseStatusGroups ? getBas1874StatusGroupForCard(item, data) : '') || data.author || 'Unknown';

                                if (!grouped.has(author)) {
                                    grouped.set(author, {
                                        title: author,
                                        card: section.card,
                                        grid: section.grid,
                                        items: [],
                                    });
                                }

                                grouped.get(author).items.push(item);
                            });
                        });

                        const priority = (title) => {
                            if (title === 'Broken Extensions') return 2;
                            if (title === 'Deprecated Extensions') return 3;
                            return 0;
                        };

                        return Array.from(grouped.values()).sort((a, b) => {
                            const leftPriority = priority(a.title);
                            const rightPriority = priority(b.title);

                            if (leftPriority !== rightPriority) return leftPriority - rightPriority;

                            return a.title.localeCompare(b.title);
                        });
                    }

                    function syncMarketplaceStatusSectionSearchVisibility() {
                        const hasSearchTerm = Array.from(document.querySelectorAll('input')).some(input => {
                            if (!input || input.closest('.ama-modal')) return false;
                            if (!/search/i.test(String(input.placeholder || ''))) return false;

                            return String(input.value || '').trim().length > 0;
                        });

                        document.querySelectorAll('.ama-status-section').forEach(section => {
                            const grid = section.querySelector('.ama-extension-carousel, .grid');
                            if (!grid) {
                                section.hidden = true;
                                return;
                            }

                            if (!hasSearchTerm) {
                                section.hidden = !grid.querySelector('.group\\\\/extension-card');
                                return;
                            }

                            const hasVisible = Array.from(grid.querySelectorAll('.ama-catalog-card-wrap, .group\\\\/extension-card')).some(item => {
                                return !item.hidden && item.style.display !== 'none';
                            });

                            section.hidden = !hasVisible;
                        });
                    }

                    function isVisibleMarketplaceCard(card) {
                        if (!card || (card.classList && card.classList.contains('ama-status-source-hidden'))) return false;

                        const wrapper = card.closest && card.closest('.ama-catalog-card-wrap');
                        if (wrapper && wrapper.hidden) return false;

                        if (card.hidden || card.style.display === 'none') return false;

                        try {
                            const style = window.getComputedStyle(card);
                            if (style && style.display === 'none') return false;
                        } catch (_) {}

                        return true;
                    }

                    function syncMarketplaceSectionSearchVisibility() {
                        const hasSearchTerm = Array.from(document.querySelectorAll('input')).some(input => {
                            if (!input || input.closest('.ama-modal')) return false;
                            if (!/search/i.test(String(input.placeholder || ''))) return false;

                            return String(input.value || '').trim().length > 0;
                        });

                        getMarketplaceSections(false).forEach(section => {
                            if (!section || !section.card || (section.card.dataset && section.card.dataset.amaStatusSection)) return;

                            if (!hasSearchTerm) {
                                section.card.hidden = false;
                                return;
                            }

                            const hasVisible = section.items.some(isVisibleMarketplaceCard);
                            section.card.hidden = !hasVisible;
                        });
                    }

                    function ensureMarketplaceSearchVisibilityHandler() {
                        if (marketplaceSearchVisibilityHandlerBound) return;

                        marketplaceSearchVisibilityHandlerBound = true;

                        document.addEventListener('input', event => {
                            const input = event.target;
                            if (!input || !input.matches || !input.matches('input')) return;
                            if (input.closest('.ama-modal')) return;
                            if (!/search/i.test(String(input.placeholder || ''))) return;

                            setTimeout(syncMarketplaceStatusSectionSearchVisibility, 0);
                            setTimeout(syncMarketplaceSectionSearchVisibility, 80);
                            setTimeout(syncMarketplaceSectionSearchVisibility, 220);
                        }, true);
                    }

                    function getAuthorFullCatalogSections() {
                        return groupMarketplaceSectionsByAuthor(getMarketplaceSections(true));
                    }

                    function bindMarketplaceStatusViewAll(section, title) {
                        if (!section || !section.querySelector) return;

                        const viewBtn = section.querySelector('.ama-view-btn');
                        if (!viewBtn) return;

                        const openStatusCatalog = event => {
                            if (event) {
                                event.preventDefault();
                                event.stopPropagation();
                                if (event.stopImmediatePropagation) event.stopImmediatePropagation();
                            }

                            const grid = section.querySelector('.ama-extension-carousel, .grid');
                            const statusSection = {
                                title,
                                card: section,
                                grid,
                                items: grid ? Array.from(grid.querySelectorAll('.group\\\\/extension-card')) : [],
                            };

                            openFullCatalogModal(groupMarketplaceSectionsByAuthor([statusSection], false), 'Full Catalog');
                        };

                        viewBtn.onclick = openStatusCatalog;
                        viewBtn.onpointerdown = openStatusCatalog;
                    }

                    function syncMarketplaceStatusSectionLayout(section) {
                        if (!section || !section.querySelector) return;

                        const grid = section.querySelector('.grid');
                        if (!grid) return;

                        if (areExtensionCarouselsEnabled()) {
                            grid.classList.add('ama-extension-carousel');
                            grid.classList.remove('ama-status-normal-grid');
                            makeDraggableScroller(grid, 'extensionsCarousel');
                            return;
                        }

                        grid.classList.remove('ama-extension-carousel');
                        grid.classList.remove('ama-drag-pending');
                        grid.classList.remove('ama-dragging');
                        grid.classList.add('ama-status-normal-grid');
                        grid.dataset.amaDragFeature = 'extensionsCarousel';
                    }

                    function getOrCreateMarketplaceStatusSection(title) {
                        let section = document.querySelector('[data-ama-status-section="' + title + '"]');
                        if (section) {
                            syncMarketplaceStatusSectionLayout(section);
                            bindMarketplaceStatusViewAll(section, title);
                            return section;
                        }

                        const sections = getMarketplaceSections(false);
                        const lastCard = sections.length ? sections[sections.length - 1].card : null;
                        const parent = lastCard && lastCard.parentElement;
                        if (!parent) return null;

                        const card = document.createElement('div');
                        card.className = 'UI-Card__root ama-status-section';
                        card.dataset.amaStatusSection = title;

                        const header = document.createElement('div');
                        header.className = 'ama-header-container';

                        const left = document.createElement('div');
                        left.className = 'ama-header-left';

                        const heading = document.createElement('h3');
                        heading.textContent = title;

                        left.appendChild(heading);

                        const viewBtn = document.createElement('button');
                        viewBtn.type = 'button';
                        viewBtn.className = 'ama-view-btn';
                        viewBtn.innerText = 'View All';

                        left.appendChild(viewBtn);
                        header.appendChild(left);

                        const grid = document.createElement('div');
                        grid.className = 'grid';

                        card.appendChild(header);
                        card.appendChild(grid);
                        parent.appendChild(card);
                        syncMarketplaceStatusSectionLayout(card);
                        bindMarketplaceStatusViewAll(card, title);

                        return card;
                    }

                    function regroupMarketplaceStatusSections() {
                        if (marketplaceStatusSectionsRendering) return;
                        if (!featureSettings.useBas1874Marketplace || !bas1874MarketplaceMetadataCache.length) return;

                        marketplaceStatusSectionsRendering = true;

                        const brokenSection = getOrCreateMarketplaceStatusSection('Broken Extensions');
                        const deprecatedSection = getOrCreateMarketplaceStatusSection('Deprecated Extensions');
                        const brokenGrid = brokenSection && brokenSection.querySelector('.ama-extension-carousel, .grid');
                        const deprecatedGrid = deprecatedSection && deprecatedSection.querySelector('.ama-extension-carousel, .grid');

                        if (!brokenGrid || !deprecatedGrid) {
                            marketplaceStatusSectionsRendering = false;
                            return;
                        }

                        brokenGrid.innerHTML = '';
                        deprecatedGrid.innerHTML = '';

                        document.querySelectorAll('.ama-status-source-hidden').forEach(item => {
                            item.classList.remove('ama-status-source-hidden');
                            delete item.dataset.amaStatusGroup;
                        });

                        getMarketplaceSections(false).forEach(section => {
                            if (section.card === brokenSection || section.card === deprecatedSection) return;

                            Array.from(section.grid.querySelectorAll('.group\\\\/extension-card')).forEach(item => {
                                const data = getExtensionCardData(item);
                                const group = getBas1874StatusGroupForCard(item, data);
                                let targetGrid = null;

                                if (group === 'Broken Extensions') {
                                    targetGrid = brokenGrid;
                                } else if (group === 'Deprecated Extensions') {
                                    targetGrid = deprecatedGrid;
                                }

                                if (targetGrid) {
                                    item.classList.add('ama-status-source-hidden');
                                    item.dataset.amaStatusGroup = group;

                                    const clone = item.cloneNode(true);
                                    const wrapper = document.createElement('div');
                                    wrapper.className = 'ama-catalog-card-wrap';
                                    wrapper.dataset.amaActionExtensionId = data.id || '';
                                    clone.classList.remove('ama-status-source-hidden');
                                    delete clone.dataset.amaStatusGroup;
                                    if (item.querySelector('button[disabled]')) {
                                        clone.dataset.amaInstalledMarketplaceCard = 'true';
                                    } else {
                                        delete clone.dataset.amaInstalledMarketplaceCard;
                                    }
                                    clone.querySelectorAll('button').forEach(button => button.remove());
                                    markExtensionUpdateState(clone);
                                    const extension = getBas1874MetadataForCardData(data);
                                    renderMarketplaceExtraInfo(clone, extension);
                                    optimizeImages(clone);

                                    wrapper.appendChild(clone);

                                    const actionSource = isInstalledExtensionsView() ? clone : item;
                                    const actions = createMarketplaceCloneActions(actionSource, isInstalledExtensionsView(), true);
                                    if (actions) {
                                        bindDirectStatusActions(actions);
                                        wrapper.appendChild(actions);
                                    }

                                    targetGrid.appendChild(wrapper);
                                }
                            });
                        });

                        [brokenSection, deprecatedSection].forEach(section => {
                            const grid = section.querySelector('.ama-extension-carousel, .grid');
                            if (!grid || !grid.querySelector('.group\\\\/extension-card')) {
                                section.hidden = true;
                                return;
                            }

                            section.hidden = false;
                            syncMarketplaceStatusSectionLayout(section);
                        });

                        const sections = getMarketplaceSections(false);
                        const lastCard = sections.length ? sections[sections.length - 1].card : null;
                        const parent = lastCard && lastCard.parentElement;

                        if (parent) {
                            parent.appendChild(brokenSection);
                            parent.appendChild(deprecatedSection);
                        }

                        syncMarketplaceStatusSectionSearchVisibility();
                        marketplaceStatusSectionsRendering = false;
                    }

                    function cleanupMarketplaceStatusSources(root) {
                        const scope = root || document;
                        if (!scope || !scope.querySelectorAll) return;

                        scope.querySelectorAll('.ama-status-source-hidden').forEach(item => {
                            item.classList.remove('ama-status-source-hidden');
                            delete item.dataset.amaStatusGroup;
                        });
                    }

                    function cleanupMarketplaceStatusSections(root) {
                        const scope = root || document;
                        if (!scope || !scope.querySelectorAll) return;

                        scope.querySelectorAll('.ama-status-section').forEach(section => section.remove());
                    }

                    function scheduleMarketplaceStatusSections() {
                        if (marketplaceStatusSectionsScheduled) return;

                        marketplaceStatusSectionsScheduled = true;

                        setTimeout(() => {
                            marketplaceStatusSectionsScheduled = false;
                            regroupMarketplaceStatusSections();
                        }, 30);
                    }

                    function openAuthorFullCatalog(title) {
                        const open = () => openFullCatalogModal(getAuthorFullCatalogSections(), title || 'Full Catalog');

                        if (featureSettings.useBas1874Marketplace) {
                            loadBas1874MarketplaceMetadata().then(open).catch(open);
                            return;
                        }

                        open();
                    }

                    function appendCatalogCards(rowGrid, items, isInstalledCatalog) {
                        items.forEach(item => {
                            const clone = item.cloneNode(true);
                            const wrapper = document.createElement('div');
                            wrapper.className = 'ama-catalog-card-wrap';
                            const isStatusItem = !!(item.closest && item.closest('.ama-status-section'));
                            const data = getExtensionCardData(item);
                            wrapper.dataset.amaActionExtensionId = data.id || '';
                            clone.querySelectorAll('button').forEach(button => button.remove());
                            markExtensionUpdateState(clone);
                            enhanceMarketplaceExtraInfo(clone);
                            optimizeImages(clone);

                            wrapper.appendChild(clone);

                            const actions = createMarketplaceCloneActions(item, isInstalledCatalog, isStatusItem);
                            if (actions) {
                                wrapper.appendChild(actions);
                            }

                            rowGrid.appendChild(wrapper);
                        });
                    }

                    function bindFullCatalogSearch(modal) {
                        const modalSearch = modal.querySelector('#ama-modal-search');
                        if (!modalSearch) return;

                        let modalSearchTimer = null;

                        modalSearch.oninput = (e) => {
                            clearTimeout(modalSearchTimer);

                            modalSearchTimer = setTimeout(() => {
                                filterCards(modal, e.target.value);

                                modal.querySelectorAll('.modal-row').forEach(row => {
                                    const visible = Array.from(row.querySelectorAll('.ama-catalog-card-wrap, .group\\/extension-card')).some(card => !card.hidden && card.style.display !== 'none');
                                    row.style.display = visible ? '' : 'none';
                                });
                            }, 140);
                        };
                    }

                    function openFullCatalogModal(sections, title) {
                        const modal = document.createElement('div');
                        modal.className = 'ama-modal';
                        const isInstalledCatalog = !!document.querySelector('input[placeholder="Search installed extensions..."]');

                        modal.innerHTML =
                            '<button type="button" class="ama-modal-close">Close</button>' +
                            '<h2 style="margin-top:0; margin-bottom:16px; font-size:32px;">' + escapeHtml(title || 'Full Catalog') + '</h2>' +
                            '<div class="ama-search-wrapper" style="margin-bottom:24px;">' +
                                '<svg class="ama-search-icon" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">' +
                                    '<circle cx="11" cy="11" r="8"/>' +
                                    '<line x1="21" y1="21" x2="16.65" y2="16.65"/>' +
                                '</svg>' +
                                '<input id="ama-modal-search" class="ama-search-input" placeholder="Search..." style="width:300px;"/>' +
                            '</div>' +
                            '<div id="ama-modal-content"></div>';

                        const content = modal.querySelector('#ama-modal-content');
                        const fragment = document.createDocumentFragment();

                        sections.forEach(section => {
                            const row = document.createElement('div');
                            row.className = 'modal-row';

                            const rowTitle = document.createElement('div');
                            rowTitle.className = 'modal-row-title';
                            rowTitle.textContent = section.title;

                            const rowGrid = document.createElement('div');
                            rowGrid.className = 'ama-carousel-row';

                            appendCatalogCards(rowGrid, section.items, isInstalledCatalog);
                            makeDraggableScroller(rowGrid, 'betterMarketplace');

                            row.appendChild(rowTitle);
                            row.appendChild(rowGrid);
                            fragment.appendChild(row);
                        });

                        content.appendChild(fragment);

                        bindFullCatalogSearch(modal);

                        modal.querySelector('.ama-modal-close').onclick = () => modal.remove();

                        modal.onclick = (e) => {
                            if (e.target === modal) modal.remove();
                        };

                        document.body.appendChild(modal);
                    }

                    function ensureGlobalFullCatalogButton() {
                        if (!featureSettings.betterMarketplace) return;
                        const existing = document.querySelector('.ama-global-catalog-bar');
                        const sections = getMarketplaceSections(false);
                        if (existing) {
                            if (!sections.length) {
                                existing.remove();
                                return;
                            }

                            existing.querySelector('.ama-global-catalog-btn').onclick = () => openAuthorFullCatalog('Full Catalog');
                            return;
                        }
                        if (!sections.length) return;

                        const firstCard = sections[0].card;
                        if (!firstCard || !firstCard.parentElement) return;

                        const bar = document.createElement('div');
                        bar.className = 'ama-global-catalog-bar';

                        const button = document.createElement('button');
                        button.type = 'button';
                        button.className = 'ama-global-catalog-btn';
                        button.textContent = 'Full Catalog';
                        button.onclick = () => openAuthorFullCatalog('Full Catalog');

                        bar.appendChild(button);
                        firstCard.parentElement.insertBefore(bar, firstCard);
                    }

                    function enhanceExtensionCard(card) {
                        if (!featureSettings.betterMarketplace) return;
                        if (!card) return;
                        if (card.dataset && card.dataset.amaStatusSection) return;

                        if (!isMarketplaceSectionCard(card)) {
                            if (card.dataset.amaEnhanced || card.querySelector(':scope > .ama-header-container')) {
                                cleanupBetterMarketplace(card);
                            }

                            return;
                        }

                        if (card.dataset.amaEnhanced && card.dataset.amaEnhanced !== marketplaceEnhancementVersion) {
                            cleanupBetterMarketplace(card);
                        }

                        markMarketplaceExtensionCards(card);

                        if (card.dataset.amaEnhanced === marketplaceEnhancementVersion) {
                            const enhancedGrid = card.querySelector('.grid');
                            if (enhancedGrid && areExtensionCarouselsEnabled()) {
                                enhancedGrid.classList.add('ama-extension-carousel');
                                makeDraggableScroller(enhancedGrid, 'betterMarketplace');
                            } else {
                                cleanupExtensionCarousels(card);
                            }

                            return;
                        }

                        const titleEl = card.querySelector('h3');
                        const grid = card.querySelector('.grid');

                        if (!titleEl || !grid) return;

                        card.dataset.amaEnhanced = marketplaceEnhancementVersion;

                        const header = document.createElement('div');
                        header.className = 'ama-header-container';

                        const left = document.createElement('div');
                        left.className = 'ama-header-left';
                        left.appendChild(titleEl);

                        const viewBtn = document.createElement('button');
                        viewBtn.type = 'button';
                        viewBtn.className = 'ama-view-btn';
                        viewBtn.innerText = 'View All';
                        left.appendChild(viewBtn);

                        const right = document.createElement('div');
                        right.className = 'ama-header-right';

                        const searchWrapper = document.createElement('div');
                        searchWrapper.className = 'ama-search-wrapper';

                        const searchIcon = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
                        searchIcon.setAttribute('class', 'ama-search-icon');
                        searchIcon.setAttribute('viewBox', '0 0 24 24');
                        searchIcon.setAttribute('fill', 'none');
                        searchIcon.setAttribute('stroke', 'white');
                        searchIcon.setAttribute('stroke-width', '2');
                        searchIcon.setAttribute('stroke-linecap', 'round');
                        searchIcon.setAttribute('stroke-linejoin', 'round');
                        searchIcon.innerHTML = '<circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>';

                        const search = document.createElement('input');
                        search.className = 'ama-search-input';
                        search.placeholder = 'Search...';

                        searchWrapper.appendChild(searchIcon);
                        searchWrapper.appendChild(search);
                        right.appendChild(searchWrapper);

                        header.appendChild(left);
                        header.appendChild(right);
                        card.prepend(header);
                        if (areExtensionCarouselsEnabled()) {
                            grid.classList.add('ama-extension-carousel');
                            makeDraggableScroller(grid, 'betterMarketplace');
                        } else {
                            cleanupExtensionCarousels(card);
                        }

                        let searchTimer = null;

                        search.oninput = (e) => {
                            clearTimeout(searchTimer);
                            searchTimer = setTimeout(() => {
                                filterCards(grid, e.target.value);
                            }, 140);
                        };

                        viewBtn.onclick = () => {
                            const section = {
                                title: titleEl.textContent.trim() || 'Full Catalog',
                                card,
                                grid,
                                items: Array.from(grid.querySelectorAll('.group\\\\/extension-card')),
                            };

                            const open = () => openFullCatalogModal(groupMarketplaceSectionsByAuthor([section]), 'Full Catalog');

                            if (featureSettings.useBas1874Marketplace) {
                                loadBas1874MarketplaceMetadata().then(open).catch(open);
                                return;
                            }

                            open();
                        };

                    }

                    function processRoot(root) {
                        if (!root) return;

                        setBodyFlags();

                        if (!featureSettings.betterMarketplace) {
                            cleanupBetterMarketplace(root);
                        }

                        if (!areCarouselsEnabledForCurrentPage()) {
                            cleanupCarousels(root);
                        }

                        if (!areExtensionCarouselsEnabled()) {
                            cleanupExtensionCarousels(root);
                        }

                        if (!featureSettings.subDubIcons) {
                            cleanupAllMediaBadges(root);
                            removeScheduleDubBadges(root);
                        }

                        if (!featureSettings.useBas1874Marketplace) {
                            cleanupMarketplaceStatusSources(root);
                            cleanupMarketplaceStatusSections(root);
                        }

                        if (root === document) {
                            document.querySelectorAll('.ama-carousel-nav-btn').forEach(btn => btn.remove());

                            document.querySelectorAll('.ama-manga-carousel-parent').forEach(parent => {
                                parent.classList.remove('ama-manga-carousel-parent');
                            });

                            if (areCarouselsEnabledForCurrentPage()) {
                                document.querySelectorAll(targetGridsQuery).forEach(grid => {
                                    enhanceCarouselGrid(grid);
                                });
                            }

                            if (featureSettings.subDubIcons) {
                                document.querySelectorAll(mediaEntryCardQuery).forEach(card => {
                                    enhanceMediaEntryCard(card);
                                });

                                document.querySelectorAll(scheduleEventQuery).forEach(event => {
                                    enhanceScheduleEvent(event);
                                });

                                document.querySelectorAll(scheduleEntryLinkQuery).forEach(link => {
                                    if (isInsideScheduleArea(link)) enhanceScheduleEvent(link);
                                });
                            }

                            if (featureSettings.betterMarketplace) {
                                document.querySelectorAll(cardQuery).forEach(card => {
                                    enhanceExtensionCard(card);
                                });
                                markMarketplaceExtensionCards(document);
                                scheduleMarketplaceStatusSections();
                                ensureMarketplaceSearchVisibilityHandler();
                                ensureGlobalFullCatalogButton();
                            } else {
                                cleanupBetterMarketplace(document);
                            }

                            optimizeImages(document);
                            removeRandomSearchIcons(document);
                            cleanupStrayScheduleDubBadges(document);
                            return;
                        }

                        removeArrowArtifacts(root);
                        removeRandomSearchIcons(root);

                        if (root.matches && root.matches(targetGridsQuery)) {
                            enhanceCarouselGrid(root);
                            return;
                        }

                        if (root.matches && root.matches(mediaEntryCardQuery)) {
                            enhanceMediaEntryCard(root);
                            return;
                        }

                        if (root.matches && root.matches(scheduleEventQuery)) {
                            enhanceScheduleEvent(root);
                            return;
                        }

                        if (root.matches && root.matches(scheduleEntryLinkQuery)) {
                            if (isInsideScheduleArea(root)) enhanceScheduleEvent(root);
                            return;
                        }

                        if (root.matches && root.matches(cardQuery)) {
                            enhanceExtensionCard(root);
                            return;
                        }

                        if (root.querySelectorAll) {
                            if (areCarouselsEnabledForCurrentPage()) {
                                root.querySelectorAll(targetGridsQuery).forEach(grid => {
                                    enhanceCarouselGrid(grid);
                                });
                            }

                            if (featureSettings.subDubIcons) {
                                root.querySelectorAll(mediaEntryCardQuery).forEach(card => {
                                    enhanceMediaEntryCard(card);
                                });

                                root.querySelectorAll(scheduleEventQuery).forEach(event => {
                                    enhanceScheduleEvent(event);
                                });

                                root.querySelectorAll(scheduleEntryLinkQuery).forEach(link => {
                                    if (isInsideScheduleArea(link)) enhanceScheduleEvent(link);
                                });
                            }

                            cleanupStrayScheduleDubBadges(root);

                            if (featureSettings.betterMarketplace) {
                                root.querySelectorAll(cardQuery).forEach(card => {
                                    enhanceExtensionCard(card);
                                });
                                markMarketplaceExtensionCards(root);
                                scheduleMarketplaceStatusSections();
                                ensureMarketplaceSearchVisibilityHandler();
                                ensureGlobalFullCatalogButton();
                            } else {
                                cleanupBetterMarketplace(root);
                            }

                            optimizeImages(root);
                        }
                    }

                    const queuedRoots = new Set();
                    let scheduled = false;

                    function flushQueue() {
                        scheduled = false;

                        const roots = Array.from(queuedRoots);
                        queuedRoots.clear();

                        for (const root of roots) {
                            try {
                                processRoot(root);
                            } catch (_) {}
                        }
                    }

                    function scheduleRoot(root) {
                        if (!root || !isElement(root)) return;

                        queuedRoots.add(root);

                        if (scheduled) return;
                        scheduled = true;

                        if ('requestIdleCallback' in window) {
                            window.requestIdleCallback(flushQueue, { timeout: 700 });
                        } else {
                            setTimeout(flushQueue, 220);
                        }
                    }

                    window.__AMA_SAVE_SETTINGS__ = function(nextSettings) {
                        const savedSettings = readBrowserSettings();
                        featureSettings = normalizeFeatureSettings(Object.assign({}, savedSettings, featureSettings, nextSettings || {}));
                        writeBrowserSettings(featureSettings);
                        clearAnimeScheduleApiTokenSetting();
                        syncAnimeScheduleApiTokenToDom();
                        animeScheduleDubApiEntriesPromises = new Map();
                        animeScheduleDubFeedEntriesPromise = null;
                        animeScheduleAnimeDetailsPromises.clear();

                        applyBas1874MarketplacePreference();
                        setBodyFlags();

                        if (!featureSettings.betterMarketplace) {
                            cleanupBetterMarketplace(document);
                        }

                        if (!areCarouselsEnabledForCurrentPage()) {
                            cleanupCarousels(document);
                        }

                        if (!areExtensionCarouselsEnabled()) {
                            cleanupExtensionCarousels(document);
                        }

                        document.querySelectorAll('.ama-status-section').forEach(section => {
                            syncMarketplaceStatusSectionLayout(section);
                        });

                        if (!featureSettings.subDubIcons) {
                            cleanupAllMediaBadges(document);
                            removeScheduleDubBadges(document);
                        }

                        if (!featureSettings.useBas1874Marketplace) {
                            cleanupMarketplaceStatusSources(document);
                            cleanupMarketplaceStatusSections(document);
                            document.querySelectorAll('.ama-marketplace-extra').forEach(el => el.remove());
                            document.querySelectorAll('.ama-installed-marketplace-extra').forEach(el => el.remove());
                            document.querySelectorAll('[data-ama-marketplace-extra-enhanced]').forEach(el => {
                                delete el.dataset.amaMarketplaceExtraEnhanced;
                            });
                            document.querySelectorAll('[data-ama-installed-marketplace-extra-enhanced]').forEach(el => {
                                delete el.dataset.amaInstalledMarketplaceExtraEnhanced;
                            });
                        }

                        scheduleRoot(document.body || document.documentElement);
                    };

                    window.__AMA_APPLY_SETTINGS__ = window.__AMA_SAVE_SETTINGS__;

                    let lastAmaRouteKey = String(window.location.pathname || '') + String(window.location.search || '') + String(window.location.hash || '');

                    function refreshForRouteChange() {
                        const nextRouteKey = String(window.location.pathname || '') + String(window.location.search || '') + String(window.location.hash || '');
                        if (nextRouteKey === lastAmaRouteKey) return;

                        lastAmaRouteKey = nextRouteKey;
                        setBodyFlags();

                        if (!areCarouselsEnabledForCurrentPage()) {
                            cleanupCarousels(document);
                        }

                        scheduleRoot(document.body || document.documentElement);
                    }

                    window.addEventListener('popstate', refreshForRouteChange);
                    window.addEventListener('hashchange', refreshForRouteChange);
                    document.addEventListener('click', handleNativeExtensionUpdateClick, true);
                    setInterval(refreshForRouteChange, 500);

                    setBodyFlags();
                    writeBrowserSettings(featureSettings);
                    syncAnimeScheduleApiTokenToDom();
                    applyBas1874MarketplacePreference();
                    processRoot(document);

                    const observer = new MutationObserver((mutations) => {
                        for (const mutation of mutations) {
                            for (const node of mutation.addedNodes) {
                                if (!isElement(node)) continue;

                                if (
                                    node.matches('.ama-carousel-nav-btn') ||
                                    node.matches('.ama-manga-carousel-parent') ||
                                    node.matches('[role="dialog"]') ||
                                    node.matches('[data-radix-popper-content-wrapper]') ||
                                    node.matches(targetGridsQuery) ||
                                    node.matches(mediaEntryCardQuery) ||
                                    node.matches(scheduleEventQuery) ||
                                    (node.matches(scheduleEntryLinkQuery) && isInsideScheduleArea(node)) ||
                                    node.matches(cardQuery) ||
                                    node.matches('svg')
                                ) {
                                    scheduleRoot(node);
                                    continue;
                                }

                                if (
                                    node.querySelector &&
                                    node.querySelector(arrowQuery + ', ' + targetGridsQuery + ', ' + mediaEntryCardQuery + ', ' + scheduleEventQuery + ', ' + scheduleEntryLinkQuery + ', ' + cardQuery + ', svg')
                                ) {
                                    scheduleRoot(node);
                                }
                            }
                        }
                    });

                    observer.observe(document.body, {
                        childList: true,
                        subtree: true
                    });
                })();
            `);

            body.append(script);
        });
    });
}
