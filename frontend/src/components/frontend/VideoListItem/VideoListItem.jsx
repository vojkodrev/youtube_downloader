import { Link } from 'react-router-dom'
import VideoListItemMenu from './VideoListItemMenu'
import VideoListItemThumbnail from './VideoListItemThumbnail'
import VideoListItemInfo from './VideoListItemInfo'


function getTargetLink(video, playlistVideos) {
    if (!playlistVideos) {
        return video.savedTime ? `/watch/${video.id}?t=${video.savedTime}` : `/watch/${video.id}`
    }
    const target = playlistVideos.find(v => {
        const saved = parseFloat(v.savedTime) || 0
        return !v.duration || saved < v.duration - 5
    }) ?? playlistVideos[playlistVideos.length - 1]
    return target.savedTime ? `/watch/${target.id}?t=${target.savedTime}` : `/watch/${target.id}`
}

export default function VideoListItem({ video, isSelected, videoCountVisible, onWatchedReset, onWatchedMark, onDelete, playlistVideos }) {

    const progressSavedTime = playlistVideos
        ? playlistVideos.reduce((sum, v) => sum + (parseFloat(v.savedTime) || 0), 0)
        : parseFloat(video.savedTime) || 0
    const progressDuration = playlistVideos
        ? playlistVideos.reduce((sum, v) => sum + (v.duration || 0), 0)
        : video.duration

    return (
        <div className={`flex items-center p-2 hover:bg-gray-100 ${isSelected ? 'bg-gray-200' : ''}`}>
            <Link
                to={getTargetLink(video, playlistVideos)}
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
