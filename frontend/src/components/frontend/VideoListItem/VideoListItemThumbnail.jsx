import { mediaUrl } from '../../../api'

export default function VideoListItemThumbnail({ videoId, progressSavedTime, progressDuration }) {
    return (
        <div className="relative flex-shrink-0">
            <img
                src={mediaUrl(`/thumbnail/${videoId}`)}
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
    )
}
