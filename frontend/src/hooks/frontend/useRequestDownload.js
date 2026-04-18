const API_URL = import.meta.env.VITE_API_URL

export function useRequestDownload() {
    async function requestVideoDownload(formData) {
        const res = await fetch(`${API_URL}/request-video-download`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url: formData.get('url') }),
        })
        if (!res.ok) {
            const data = await res.json().catch(() => null)
            throw new Error(data?.error)
        }
    }

    async function requestPlaylistDownload(formData) {
        const res = await fetch(`${API_URL}/request-playlist-download`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url: formData.get('url') }),
        })
        if (!res.ok) {
            const data = await res.json().catch(() => null)
            throw new Error(data?.error)
        }
    }

    return { requestVideoDownload, requestPlaylistDownload }
}
