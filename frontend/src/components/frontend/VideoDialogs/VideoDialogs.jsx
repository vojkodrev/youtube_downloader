import { forwardRef, useImperativeHandle, useRef } from 'react'
import FormDialog from '@/components/frontend/FormDialog/FormDialog'
import DeleteVideoDialog from '@/components/frontend/DeleteVideoDialog/DeleteVideoDialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useRequestDownload } from '@/hooks/frontend/useRequestDownload'

const VideoDialogs = forwardRef(function VideoDialogs(_, ref) {
    const requestVideoDownloadDialogRef = useRef(null)
    const requestPlaylistDownloadDialogRef = useRef(null)
    const deleteVideoDialogRef = useRef(null)

    const { requestVideoDownload, requestPlaylistDownload } = useRequestDownload()

    useImperativeHandle(ref, () => ({
        openVideoDownload: () => requestVideoDownloadDialogRef.current.open(),
        openPlaylistDownload: () => requestPlaylistDownloadDialogRef.current.open(),
        openDeleteVideo: (videos, confirmDelete) => deleteVideoDialogRef.current.open(videos, confirmDelete),
    }))

    return (
        <>
            <FormDialog
                ref={requestVideoDownloadDialogRef}
                title="Request a Video Download"
                description="Paste a YouTube (or other supported) URL below and we'll download it to your local library."
                submitLabel="Download"
                onSubmit={requestVideoDownload}
            >
                <Label htmlFor="video-url">Video URL</Label>
                <Input
                    id="video-url"
                    name="url"
                    placeholder="https://www.youtube.com/watch?v=..."
                />
            </FormDialog>

            <FormDialog
                ref={requestPlaylistDownloadDialogRef}
                title="Request a Playlist Download"
                description="Paste a YouTube playlist URL below and we'll download it to your local library."
                submitLabel="Download"
                onSubmit={requestPlaylistDownload}
            >
                <Label htmlFor="playlist-url">Playlist URL</Label>
                <Input
                    id="playlist-url"
                    name="url"
                    placeholder="https://www.youtube.com/playlist?list=..."
                />
            </FormDialog>

            <DeleteVideoDialog ref={deleteVideoDialogRef} />
        </>
    )
})

export default VideoDialogs
