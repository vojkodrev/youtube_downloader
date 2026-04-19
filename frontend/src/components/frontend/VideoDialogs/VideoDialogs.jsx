import { forwardRef, useImperativeHandle, useRef } from 'react'
import RequestDownloadDialogs from '@/components/frontend/VideoDialogs/RequestDownloadDialogs'
import DeleteVideoDialog from '@/components/frontend/DeleteVideoDialog/DeleteVideoDialog'
import WatchedTimesDialogs from '@/components/frontend/WatchedTimesDialogs/WatchedTimesDialogs'

const VideoDialogs = forwardRef(function VideoDialogs(_, ref) {
    const requestDownloadDialogsRef = useRef(null)
    const deleteVideoDialogRef = useRef(null)
    const watchedTimesDialogsRef = useRef(null)

    useImperativeHandle(ref, () => ({
        openVideoDownload: () => requestDownloadDialogsRef.current.openVideoDownload(),
        openPlaylistDownload: () => requestDownloadDialogsRef.current.openPlaylistDownload(),
        openDeleteVideo: (videos, confirmDelete) => deleteVideoDialogRef.current.open(videos, confirmDelete),
        openExportWatchedTimes: () => watchedTimesDialogsRef.current.openExport(),
        openImportWatchedTimes: () => watchedTimesDialogsRef.current.openImport(),
    }))

    return (
        <>
            <RequestDownloadDialogs ref={requestDownloadDialogsRef} />
            <DeleteVideoDialog ref={deleteVideoDialogRef} />
            <WatchedTimesDialogs ref={watchedTimesDialogsRef} />
        </>
    )
})

export default VideoDialogs
