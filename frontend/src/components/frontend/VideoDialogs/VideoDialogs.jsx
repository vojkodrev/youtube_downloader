import { forwardRef, useImperativeHandle, useRef } from 'react'
import RequestDownloadDialogs from '@/components/frontend/VideoDialogs/RequestDownloadDialogs'
import DeleteVideoDialog from '@/components/frontend/VideoDialogs/DeleteVideoDialog'
import WatchedTimesDialogs from '@/components/frontend/VideoDialogs/WatchedTimesDialogs'

const VideoDialogs = forwardRef(function VideoDialogs({ onImportWatchedTimes }, ref) {
    const requestDownloadDialogsRef = useRef(null)
    const deleteVideoDialogRef = useRef(null)
    const watchedTimesDialogsRef = useRef(null)

    useImperativeHandle(ref, () => ({
        openVideoDownload: () => requestDownloadDialogsRef.current.openVideoDownload(),
        openPlaylistDownload: () => requestDownloadDialogsRef.current.openPlaylistDownload(),
        openDeleteVideo: (videos, confirmDelete) => deleteVideoDialogRef.current.open(videos, confirmDelete),
        openExportWatchedTimes: (data) => watchedTimesDialogsRef.current.openExport(data),
        openImportWatchedTimes: () => watchedTimesDialogsRef.current.openImport(),
    }))

    return (
        <>
            <RequestDownloadDialogs ref={requestDownloadDialogsRef} />
            <DeleteVideoDialog ref={deleteVideoDialogRef} />
            <WatchedTimesDialogs ref={watchedTimesDialogsRef} onImport={onImportWatchedTimes} />
        </>
    )
})

export default VideoDialogs
