import { forwardRef } from 'react'
import { useSearchParams } from 'react-router-dom'

const VideoPlayer = forwardRef(function VideoPlayer({ video, onVideoUpdated }, ref) {
    const [searchParams, setSearchParams] = useSearchParams()

    function saveWatchedTime(t) {
        localStorage.setItem(`time_${video.id}`, t)
        setSearchParams({ t }, { replace: true })
        video.savedTime = t
        onVideoUpdated(video)
    }

    return (
        <video
            ref={ref}
            key={video.id}
            src={`${import.meta.env.VITE_API_URL}/video/${video.id}`}
            controls
            autoPlay
            playsInline
            onTimeUpdate={e => {
                const t = Math.round(e.target.currentTime)
                if (t % 5 !== 0) return
                if (t === parseInt(searchParams.get('t'))) return
                saveWatchedTime(t)
            }}
            onSeeked={e => {
                const t = Math.round(e.target.currentTime)
                if (t === parseInt(searchParams.get('t'))) return
                saveWatchedTime(t)
            }}
            onLoadedMetadata={e => {
                const t = searchParams.get('t')
                if (t) e.target.currentTime = parseFloat(t)
            }}
            className="w-full"
        />
    )
})

export default VideoPlayer
