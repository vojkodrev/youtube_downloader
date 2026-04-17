import { forwardRef } from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import FormDialog from '@/components/frontend/FormDialog/FormDialog'

const API_URL = import.meta.env.VITE_API_URL

const RequestVideoDownloadDialog = forwardRef(function RequestVideoDownloadDialog(_, ref) {
    async function handleSubmit(formData) {
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

    return (
        <FormDialog
            ref={ref}
            title="Request a Video Download"
            description="Paste a YouTube (or other supported) URL below and we'll download it to your local library."
            submitLabel="Download"
            onSubmit={handleSubmit}
        >
            <Label htmlFor="video-url">Video URL</Label>
            <Input
                id="video-url"
                name="url"
                placeholder="https://www.youtube.com/watch?v=..."
            />
        </FormDialog>
    )
})

export default RequestVideoDownloadDialog
