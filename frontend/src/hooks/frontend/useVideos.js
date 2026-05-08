import { useState, useEffect, useMemo } from 'react'
import { api } from '@/lib/api'

export function useVideos() {
    const [videos, setVideos] = useState([])
    const [playlists, setPlaylists] = useState([])

    const versionToVideoId = useMemo(() => {
        const map = {}
        for (const v of videos) {
            map[v.id] = v.id
            for (const version of v.versions ?? []) {
                map[version.id] = v.id
            }
        }
        return map
    }, [videos])

    const videoIdToQuality = useMemo(() => {
        const map = {}
        for (const v of videos) {
            map[v.id] = 'original'
            for (const version of v.versions ?? []) {
                map[version.id] = version.quality
            }
        }
        return map
    }, [videos])

    useEffect(() => {
        (async () => {
            const res = await api.get('/videos')
            const data = await res.json()

            const partRe = /^(.+) part(\d{2})$/

            const videos = data.map(v => ({
                ...v,
                visible: parseInt(v.name.match(partRe)?.[2] ?? '1') <= 1,
                savedTime: localStorage.getItem(`time_${v.id}`)
            }))

            const videoGroups = Object.groupBy(videos, v => v.name.match(partRe)?.[1] ?? v.name)
            for (const [base, items] of Object.entries(videoGroups)) {
                const primary = items.find(v => v.name.match(partRe)?.[2] === '01') ?? items[0]
                primary.name = base
                primary.videoCount = items.length
            }

            const playlistsArr = Object.values(videoGroups)
                .filter(items => items.length > 1)
                .map(items => [...items]
                    .sort((a, b) => {
                        const na = parseInt(a.name.match(/part(\d{2})$/)?.[1] ?? '0')
                        const nb = parseInt(b.name.match(/part(\d{2})$/)?.[1] ?? '0')
                        return na - nb
                    })
                )

            setPlaylists(playlistsArr)
            setVideos(videos)

            Promise.allSettled(videos.map(async v => {
                const res = await api.get(`/duration/${v.id}`)
                if (!res.ok) return
                const { duration } = await res.json()
                v.duration = duration
            })).then(() => setVideos(prev => [...prev]))
        })()
    }, [])

    return { videos, setVideos, playlists, versionToVideoId, videoIdToQuality }
}
