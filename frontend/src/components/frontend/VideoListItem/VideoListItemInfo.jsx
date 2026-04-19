import { formatDistanceToNow, format } from 'date-fns'

export default function VideoListItemInfo({ video, videoCountVisible }) {
    return (
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
    )
}
