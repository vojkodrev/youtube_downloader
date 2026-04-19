import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'

export function useVideoNavigation(id, videos, playlists, versionToVideoId) {
    const navigate = useNavigate()
    const [searchParams] = useSearchParams()

    useEffect(() => {
        if (!id || (videos.length > 0 && !versionToVideoId[id])) {
            const firstVideo = videos[0]
            if (firstVideo) {
                const playlist = playlists.find(p => p.some(v => v.id === firstVideo.id))
                const firstId = playlist ? playlist[0].id : firstVideo.id
                navigate(`/watch/${firstId}`, { replace: true })
            }
            return
        }
        if (!searchParams.get('t')) {
            const originalId = versionToVideoId[id] ?? id
            const savedTime = Math.floor(parseFloat(localStorage.getItem(`time_${originalId}`)))
            if (savedTime && savedTime > 10) {
                navigate(`/watch/${id}?t=${savedTime}`, { replace: true })
                return
            }
        }
    }, [id, videos, playlists, versionToVideoId])
}
