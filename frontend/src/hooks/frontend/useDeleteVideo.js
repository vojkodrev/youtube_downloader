import { api } from '../../api'

export function useDeleteVideo(videoDialogsRef, playlists, setVideos) {
    async function confirmDelete(videos) {
        await Promise.allSettled(videos.map(v => api.post('/delete-video', { id: v.id })))
        setVideos(prev => prev.filter(v => !videos.some(pv => pv.id === v.id)))
    }

    function handleDeleteSingle(video) {
        videoDialogsRef.current.openDeleteVideo([video], confirmDelete)
    }

    function handleDeletePlaylist(video) {
        const pl = playlists.find(p => p.some(pv => pv.id === video.id)) ?? [video]
        videoDialogsRef.current.openDeleteVideo(pl, confirmDelete)
    }

    return { handleDeleteSingle, handleDeletePlaylist }
}
