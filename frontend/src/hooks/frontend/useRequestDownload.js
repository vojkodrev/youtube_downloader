import { api } from '@/lib/api'

export function useRequestDownload() {
    async function requestVideoDownload(formData) {
        const res = await api.post('/request-video-download', { url: formData.get('url') })
        if (!res.ok) {
            const data = await res.json().catch(() => null)
            throw new Error(data?.error)
        }
    }

    async function requestPlaylistDownload(formData) {
        const res = await api.post('/request-playlist-download', { url: formData.get('url') })
        if (!res.ok) {
            const data = await res.json().catch(() => null)
            throw new Error(data?.error)
        }
    }

    return { requestVideoDownload, requestPlaylistDownload }
}
