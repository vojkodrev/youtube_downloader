import { Link } from 'react-router-dom'
import VideoListItemMenu from './VideoListItemMenu'
import VideoListItemThumbnail from './VideoListItemThumbnail'
import VideoListItemInfo from './VideoListItemInfo'
import { videoLink, isFullyWatched } from '@/lib/video'

function getTargetLink(video, playlistVideos, currentQuality) {
    const target = playlistVideos
        ? (playlistVideos.find(v => !isFullyWatched(v)) ?? playlistVideos[playlistVideos.length - 1])
        : video
    if (currentQuality && currentQuality !== 'original') {
        const version = target.versions?.find(v => v.quality === currentQuality)
        if (version) return target.savedTime ? `/watch/${version.id}?t=${target.savedTime}` : `/watch/${version.id}`
    }
    return videoLink(target)
}

export default function VideoListItem({
    video,
    isSelected,
    videoCountVisible,
    onWatchedReset,
    onWatchedMark,
    onDelete,
    playlistVideos,
    currentQuality
}) {

    const progressSavedTime = playlistVideos
        ? playlistVideos.reduce((sum, v) => sum + (parseFloat(v.savedTime) || 0), 0)
        : parseFloat(video.savedTime) || 0
    const progressDuration = playlistVideos
        ? playlistVideos.reduce((sum, v) => sum + (v.duration || 0), 0)
        : video.duration

    return (
        <div className={`flex items-center p-2 hover:bg-gray-100 ${isSelected ? 'bg-gray-200' : ''}`}>
            <Link
                to={getTargetLink(video, playlistVideos, currentQuality)}
                title={video.name}
                className="flex gap-2 min-w-0 flex-1"
            >
                <div className="flex flex-row gap-2 min-w-0 overflow-hidden">
                    <VideoListItemThumbnail
                        videoId={video.id}
                        progressSavedTime={progressSavedTime}
                        progressDuration={progressDuration}
                    />
                    <VideoListItemInfo video={video} videoCountVisible={videoCountVisible} />
                </div>
            </Link>

            <VideoListItemMenu
                video={video}
                onWatchedMark={onWatchedMark}
                onWatchedReset={onWatchedReset}
                onDelete={onDelete}
            />
        </div>
    )
}
