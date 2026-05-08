import { useEffect } from 'react'
import { mediaUrl } from '@/lib/url'

export function useVideoMetadata(selectedVideo) {
    useEffect(() => {
        if (!selectedVideo) return
        document.title = selectedVideo.name
        if ('mediaSession' in navigator) {
            navigator.mediaSession.metadata = new MediaMetadata({
                title: selectedVideo.name,
                artwork: [{ src: mediaUrl(`/thumbnail/${selectedVideo.id}`), sizes: '512x512', type: 'image/jpeg' }]
            })
        }
    }, [selectedVideo])
}
