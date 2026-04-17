import { forwardRef } from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import FormDialog from '@/components/frontend/FormDialog/FormDialog'

const API_URL = import.meta.env.VITE_API_URL

const RequestPlaylistDownloadDialog = forwardRef(function RequestPlaylistDownloadDialog(_, ref) {
    async function handleSubmit(formData) {
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

    return (
        <FormDialog
            ref={ref}
            title="Request a Playlist Download"
            description="Paste a YouTube playlist URL below and we'll download it to your local library."
            submitLabel="Download"
            onSubmit={handleSubmit}
        >
            <Label htmlFor="playlist-url">Playlist URL</Label>
            <Input
                id="playlist-url"
                name="url"
                placeholder="https://www.youtube.com/playlist?list=..."
            />
        </FormDialog>
    )
})

export default RequestPlaylistDownloadDialog
