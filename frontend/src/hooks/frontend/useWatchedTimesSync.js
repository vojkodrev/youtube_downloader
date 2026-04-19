export function useWatchedTimesSync(videoRef, selectedVideo, setVideos) {
    function exportWatchedTimes() {
        const data = {}
        for (let i = 0; i < localStorage.length; i++) {
            const key = localStorage.key(i)
            if (key.startsWith('time_')) data[key] = localStorage.getItem(key)
        }
        return btoa(JSON.stringify(data))
    }

    function importWatchedTimes(formData) {
        const times = JSON.parse(atob(formData.get('data')))
        for (const [key, value] of Object.entries(times)) {
            localStorage.setItem(key, value)
        }
        setVideos(prev => {
            prev.forEach(v => { v.savedTime = localStorage.getItem(`time_${v.id}`) })
            return [...prev]
        })
        if (selectedVideo && times[`time_${selectedVideo.id}`] && videoRef.current) {
            videoRef.current.currentTime = parseFloat(times[`time_${selectedVideo.id}`])
        }
    }

    return { exportWatchedTimes, importWatchedTimes }
}
