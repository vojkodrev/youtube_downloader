import { useMemo } from 'react'
import { useNavigate } from 'react-router-dom'

export function useVideoSelection(videos, playlists, versionToVideoId, id) {
    const navigate = useNavigate()

    const selectedVideo = useMemo(() => videos.find(v => v.id === (versionToVideoId[id] ?? id)), [videos, id, versionToVideoId])
    const currentPlaylist = useMemo(() => playlists.find(p => p.some(v => v.id === selectedVideo?.id)) ?? [], [playlists, selectedVideo])
    const nextVideoInPlaylist = useMemo(() => {
        if (!selectedVideo || currentPlaylist.length === 0) return null
        const idx = currentPlaylist.findIndex(v => v.id === selectedVideo.id)
        return idx >= 0 && idx < currentPlaylist.length - 1 ? currentPlaylist[idx + 1] : null
    }, [currentPlaylist, selectedVideo])

    function handlePlaylistEnded() {
        if (nextVideoInPlaylist) navigate(`/watch/${nextVideoInPlaylist.id}`)
    }

    return { selectedVideo, currentPlaylist, nextVideoInPlaylist, handlePlaylistEnded }
}
