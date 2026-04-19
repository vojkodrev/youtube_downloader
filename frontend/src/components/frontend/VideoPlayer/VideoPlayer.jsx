import { forwardRef } from 'react'
import { useSearchParams, useParams } from 'react-router-dom'
import { VIDEO_SAVE_INTERVAL } from '@/lib/video'

const VideoPlayer = forwardRef(function VideoPlayer({ video, onVideoUpdated }, ref) {
    const [searchParams, setSearchParams] = useSearchParams()
    const { id } = useParams()

    function saveWatchedTime(t) {
        localStorage.setItem(`time_${video.id}`, t)
        setSearchParams({ t }, { replace: true })
        video.savedTime = t
        onVideoUpdated(video)
    }

    return (
        <video
            ref={ref}
            key={id}
            src={`${import.meta.env.VITE_API_URL}/video/${id}`}
            controls
            autoPlay
            playsInline
            onTimeUpdate={e => {
                const t = Math.round(e.target.currentTime)
                if (t % VIDEO_SAVE_INTERVAL !== 0) return
                if (t === parseInt(searchParams.get('t'))) return
                saveWatchedTime(t)
            }}
            onSeeked={e => {
                const t = Math.round(e.target.currentTime)
                if (t === parseInt(searchParams.get('t'))) return
                saveWatchedTime(t)
            }}
            onEnded={e => saveWatchedTime(Math.round(e.target.currentTime))}
            onLoadedMetadata={e => {
                const t = searchParams.get('t')
                if (t) e.target.currentTime = parseFloat(t)
            }}
            className="w-full"
        />
    )
})

export default VideoPlayer
