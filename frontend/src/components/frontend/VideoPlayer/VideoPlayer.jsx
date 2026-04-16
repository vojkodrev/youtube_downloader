export default function VideoPlayer({
    videoRef,
    video,
    searchParams,
    setSearchParams,
    setVideos }) {

    function saveWatchedTime(t) {
        localStorage.setItem(`time_${video.id}`, t)
        setSearchParams({ t }, { replace: true })
        setVideos(prev => {
            const v = prev.find(v => v.id === video.id)
            if (v) v.savedTime = t
            return [...prev]
        })
    }

    return (
        <video
            ref={videoRef}
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
}
