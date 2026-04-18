import { Link } from 'react-router-dom'
import { formatDistanceToNow, format } from 'date-fns'
import VideoListItemMenu from './VideoListItemMenu'

const API_URL = import.meta.env.VITE_API_URL

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
                to={video.savedTime ? `/watch/${video.id}?t=${video.savedTime}` : `/watch/${video.id}`}
                title={video.name}
                className="flex gap-2 min-w-0 flex-1"
            >
                <div className="flex flex-row gap-2 min-w-0 overflow-hidden">
                    <div className="relative flex-shrink-0">
                        <img
                            src={`${API_URL}/thumbnail/${video.id}`}
                            className="w-36 h-20 object-cover rounded bg-gray-300"
                        />
                        {!!progressSavedTime && !!progressDuration && (
                            <div className="absolute bottom-0 left-0 right-0 h-1 bg-white/30 rounded-bl">
                                <div
                                    className={`h-full bg-red-500 rounded-bl${progressSavedTime >= progressDuration ? ' rounded-br' : ''}`}
                                    style={{ width: `${Math.min(progressSavedTime / progressDuration * 100, 100)}%` }}
                                />
                            </div>
                        )}
                        {progressDuration != null && (
                            <span className="absolute bottom-2 right-1 bg-black/60 text-white text-xs px-1 rounded">
                                {new Date(progressDuration * 1000).toISOString().substring(progressDuration >= 3600 ? 11 : 14, 19)}
                            </span>
                        )}
                    </div>
                    <div className="flex flex-col justify-start min-w-0">
                        <p className="text-sm font-medium text-gray-900 line-clamp-2">
                            {video.name}
                        </p>
                        <p className="text-xs text-gray-500 mt-1 truncate" title={video.date ? format(new Date(video.date), 'PPpp') : undefined}>
                            {[video.channel, video.date ? formatDistanceToNow(new Date(video.date), { addSuffix: true }) : null].filter(Boolean).join(' · ')}
                        </p>
                        {video.status !== 'Ready'
                            ? <p className="text-xs text-gray-400 mt-1">{video.status}</p>
                            : videoCountVisible && video.videoCount > 1 && <p className="text-xs text-gray-400 mt-1">{video.videoCount} videos</p>
                        }
                    </div>
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
