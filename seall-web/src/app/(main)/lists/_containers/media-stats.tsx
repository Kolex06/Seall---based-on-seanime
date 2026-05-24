import { MediaAPI_Stats } from "@/api/generated/types"
import { AppLayoutStack } from "@/components/ui/app-layout"
import { AreaChart, BarChart, DonutChart } from "@/components/ui/charts"
import { Separator } from "@/components/ui/separator"
import { Stats } from "@/components/ui/stats"
import React from "react"
import { FaRegStar } from "react-icons/fa"
import { LuHourglass } from "react-icons/lu"
import { PiTelevisionSimpleBold } from "react-icons/pi"
import { TbHistory } from "react-icons/tb"

type MediaStatsProps = {
    stats?: MediaAPI_Stats
    isLoading?: boolean
}

const formatName: Record<string, string> = {
    TV: "TV",
    TV_SHORT: "TV Short",
    MOVIE: "Movie",
    SPECIAL: "Special",
    OVA: "OVA",
    ONA: "ONA",
    MUSIC: "Music",
}

const statusName: Record<string, string> = {
    CURRENT: "Current",
    PLANNING: "Planning",
    COMPLETED: "Completed",
    DROPPED: "Dropped",
    PAUSED: "Paused",
    REPEATING: "Repeating",
}

export function MediaStats(props: MediaStatsProps) {

    const {
        stats,
        isLoading,
    } = props

    const anime_formatsStats = React.useMemo(() => {
        if (!stats?.animeStats?.formats) return []

        return stats.animeStats.formats.map((item) => {
            return {
                name: formatName[item.format as string],
                count: item.count,
                hoursWatched: Math.round(item.minutesWatched / 60),
                meanScore: Number((item.meanScore / 10).toFixed(1)),
            }
        })
    }, [stats?.animeStats?.formats])

    const anime_statusesStats = React.useMemo(() => {
        if (!stats?.animeStats?.statuses) return []

        return stats.animeStats.statuses.map((item) => {
            return {
                name: statusName[item.status as string],
                count: item.count,
                hoursWatched: Math.round(item.minutesWatched / 60),
                meanScore: Number((item.meanScore / 10).toFixed(1)),
            }
        })
    }, [stats?.animeStats?.statuses])

    const anime_genresStats = React.useMemo(() => {
        if (!stats?.animeStats?.genres) return []

        return stats.animeStats.genres.map((item) => {
            return {
                name: item.genre,
                "Count": item.count,
                hoursWatched: Math.round(item.minutesWatched / 60),
                "Average score": Number((item.meanScore / 10).toFixed(1)),
            }
        }).sort((a, b) => b["Count"] - a["Count"])
    }, [stats?.animeStats?.genres])

    const [anime_thisYearStats, anime_lastYearStats] = React.useMemo(() => {
        if (!stats?.animeStats?.startYears) return []
        const thisYear = new Date().getFullYear()
        return [
            stats.animeStats.startYears.find((item) => item.startYear === thisYear),
            stats.animeStats.startYears.find((item) => item.startYear === thisYear - 1),
        ]
    }, [stats?.animeStats?.startYears])

    const anime_releaseYearsStats = React.useMemo(() => {
        if (!stats?.animeStats?.releaseYears) return []

        return stats.animeStats.releaseYears.sort((a, b) => a.releaseYear! - b.releaseYear!).map((item) => {
            return {
                name: item.releaseYear,
                "Count": item.count,
                "Hours watched": Math.round(item.minutesWatched / 60),
                "Mean score": Number((item.meanScore / 10).toFixed(1)),
            }
        })
    }, [stats?.animeStats?.releaseYears])

    return (
        <AppLayoutStack className="py-4 space-y-10" data-media-stats>

            <h1 className="text-center" data-media-stats-anime-title>Watchlist</h1>

            <div data-media-stats-anime-stats>
                <Stats
                    className="w-full"
                    size="lg"
                    items={[
                        {
                            icon: <PiTelevisionSimpleBold />,
                            name: "Total titles",
                            value: stats?.animeStats?.count ?? 0,
                        },
                        {
                            icon: <LuHourglass />,
                            name: "Watch time",
                            value: Math.round((stats?.animeStats?.minutesWatched ?? 0) / 60),
                            unit: "hours",
                        },
                        {
                            icon: <FaRegStar />,
                            name: "Average score",
                            value: ((stats?.animeStats?.meanScore ?? 0) / 10).toFixed(1),
                        },
                    ]}
                />
                <Separator />
                <Stats
                    className="w-full"
                    size="lg"
                    items={[
                        {
                            icon: <PiTelevisionSimpleBold />,
                            name: "Titles watched this year",
                            value: anime_thisYearStats?.count ?? 0,
                        },
                        {
                            icon: <TbHistory />,
                            name: "Titles watched last year",
                            value: anime_lastYearStats?.count ?? 0,
                        },
                        {
                            icon: <FaRegStar />,
                            name: "Average score this year",
                            value: ((anime_thisYearStats?.meanScore ?? 0) / 10).toFixed(1),
                        },
                    ]}
                />
            </div>

            <h3 className="text-center" data-media-stats-anime-formats-title>Formats</h3>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 w-full" data-media-stats-anime-formats-container>
                <ChartContainer legend="Total" data-media-stats-anime-formats-container-total>
                    <DonutChart
                        data={anime_formatsStats}
                        index="name"
                        category="count"
                        variant="pie"
                    />
                </ChartContainer>
                <ChartContainer legend="Hours watched" data-media-stats-anime-formats-container-hours-watched>
                    <DonutChart
                        data={anime_formatsStats}
                        index="name"
                        category="hoursWatched"
                        variant="pie"
                    />
                </ChartContainer>
                <ChartContainer legend="Average score" data-media-stats-anime-formats-container-average-score>
                    <DonutChart
                        data={anime_formatsStats}
                        index="name"
                        category="meanScore"
                        variant="pie"
                    />
                </ChartContainer>
            </div>

            <Separator />

            <h3 className="text-center" data-media-stats-anime-statuses-title>Statuses</h3>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 w-full" data-media-stats-anime-statuses-container>
                <ChartContainer legend="Total" data-media-stats-anime-statuses-container-total>
                    <DonutChart
                        data={anime_statusesStats}
                        index="name"
                        category="count"
                        variant="pie"
                    />
                </ChartContainer>
                <ChartContainer legend="Hours watched" data-media-stats-anime-statuses-container-hours-watched>
                    <DonutChart
                        data={anime_statusesStats}
                        index="name"
                        category="hoursWatched"
                        variant="pie"
                    />
                </ChartContainer>
            </div>

            <Separator />

            <h3 className="text-center" data-media-stats-anime-genres-title>Genres</h3>

            <div className="grid grid-cols-1 gap-6 w-full" data-media-stats-anime-genres-container>
                <ChartContainer legend="Favorite genres" data-media-stats-anime-genres-container-favorite-genres>
                    <BarChart
                        data={anime_genresStats}
                        index="name"
                        categories={["Count", "Average score"]}
                        colors={["brand", "blue"]}
                    />
                </ChartContainer>
            </div>

            <Separator />

            <h3 className="text-center" data-media-stats-anime-years-title>Years</h3>

            <div className="grid grid-cols-1 gap-6 w-full" data-media-stats-anime-years-container>
                <ChartContainer legend="Titles watched per release year" data-media-stats-anime-years-container-anime-watched-per-release-year>
                    <AreaChart
                        data={anime_releaseYearsStats}
                        index="name"
                        categories={["Count"]}
                        angledLabels
                    />
                </ChartContainer>
            </div>

        </AppLayoutStack>
    )
}

function ChartContainer(props: { children: React.ReactNode, legend: string }) {
    return (
        <div className="text-center w-full space-y-4" data-media-stats-chart-container>
            {props.children}
            <p className="text-center text-lg font-semibold">{props.legend}</p>
        </div>
    )
}
