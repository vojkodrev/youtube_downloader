import { useSearchParams } from 'react-router-dom'

export function useWatchedTime(id, videoRef, setVideos, playlists) {
    const [, setSearchParams] = useSearchParams()

    function handleWatchedMark(video, updateVideos = true) {
        const t = video.duration
        if (!t) return
        localStorage.setItem(`time_${video.id}`, t)
        video.savedTime = t
        if (updateVideos) setVideos(prev => [...prev])
        if (video.id === id) {
            videoRef.current.currentTime = t
            setSearchParams({ t }, { replace: true })
        }
    }

    function handleWatchedReset(video, updateVideos = true) {
        localStorage.removeItem(`time_${video.id}`)
        video.savedTime = null
        if (updateVideos) setVideos(prev => [...prev])
        if (video.id === id) {
            videoRef.current.currentTime = 0
            setSearchParams({ t: 0 }, { replace: true })
        }
    }

    function handleWatchedMarkPlaylist(v) {
        const pl = playlists.find(p => p.some(pv => pv.id === v.id)) ?? [v]
        pl.forEach(pv => handleWatchedMark(pv, false))
        setVideos(prev => [...prev])
    }

    function handleWatchedResetPlaylist(v) {
        const pl = playlists.find(p => p.some(pv => pv.id === v.id)) ?? [v]
        pl.forEach(pv => handleWatchedReset(pv, false))
        setVideos(prev => [...prev])
    }

    return { handleWatchedMark, handleWatchedReset, handleWatchedMarkPlaylist, handleWatchedResetPlaylist }
}
