import { useEffect } from 'react'

const API_URL = import.meta.env.VITE_API_URL

export function useVideoMetadata(selectedVideo) {
    useEffect(() => {
        if (!selectedVideo) return
        document.title = selectedVideo.name
        if ('mediaSession' in navigator) {
            navigator.mediaSession.metadata = new MediaMetadata({
                title: selectedVideo.name,
                artwork: [{ src: `${API_URL}/thumbnail/${selectedVideo.id}`, sizes: '512x512', type: 'image/jpeg' }]
            })
        }
    }, [selectedVideo])
}
