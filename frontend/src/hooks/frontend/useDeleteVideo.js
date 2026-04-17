const API_URL = import.meta.env.VITE_API_URL

export function useDeleteVideo(videoDialogsRef, playlists, setVideos) {
    function handleDeleteSingle(video) {
        videoDialogsRef.current.openDeleteVideo([video])
    }

    function handleDeletePlaylist(video) {
        const pl = playlists.find(p => p.some(pv => pv.id === video.id)) ?? [video]
        videoDialogsRef.current.openDeleteVideo(pl)
    }

    async function confirmDelete(videos) {
        await Promise.allSettled(videos.map(v => fetch(`${API_URL}/delete-video`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: v.id }),
        })))
        setVideos(prev => prev.filter(v => !videos.some(pv => pv.id === v.id)))
    }

    return { handleDeleteSingle, handleDeletePlaylist, confirmDelete }
}
