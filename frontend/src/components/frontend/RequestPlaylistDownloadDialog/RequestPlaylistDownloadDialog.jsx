import { useState, forwardRef, useImperativeHandle } from 'react'
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const API_URL = import.meta.env.VITE_API_URL
const DEFAULT_DOWNLOAD_ERROR = 'An error occurred. Please try again.'

const RequestPlaylistDownloadDialog = forwardRef(function RequestPlaylistDownloadDialog(_, ref) {
    const [open, setOpen] = useState(false)
    const [downloadUrl, setDownloadUrl] = useState('')
    const [downloadError, setDownloadError] = useState(false)
    const [downloadPending, setDownloadPending] = useState(false)

    useImperativeHandle(ref, () => ({
        open() { setOpen(true) }
    }))

    function handleClose() {
        setDownloadError(false)
        setDownloadUrl('')
        setOpen(false)
    }

    async function handleDownload() {
        setDownloadPending(true)
        setDownloadError(false)
        try {
            const res = await fetch(`${API_URL}/request-playlist-download`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url: downloadUrl }),
            })
            if (!res.ok) {
                const data = await res.json().catch(() => null)
                throw new Error(data?.error)
            }
            handleClose()
        } catch (e) {
            setDownloadError(e.message || DEFAULT_DOWNLOAD_ERROR)
        } finally {
            setDownloadPending(false)
        }
    }

    return (
        <Dialog open={open} onOpenChange={open => { if (!open) handleClose() }}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Request a Playlist Download</DialogTitle>
                    <DialogDescription>
                        Paste a YouTube playlist URL below and we'll download it to your local library.
                    </DialogDescription>
                </DialogHeader>
                <div className="flex flex-col gap-2 py-2">
                    <Label htmlFor="playlist-url">Playlist URL</Label>
                    <Input
                        id="playlist-url"
                        placeholder="https://www.youtube.com/playlist?list=..."
                        value={downloadUrl}
                        onChange={e => setDownloadUrl(e.target.value)}
                    />
                    {downloadError && <p className="text-sm text-red-500">{downloadError}</p>}
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={handleClose}>Cancel</Button>
                    <Button
                        disabled={!downloadUrl.trim() || downloadPending}
                        onClick={handleDownload}
                    >Download</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
})

export default RequestPlaylistDownloadDialog
