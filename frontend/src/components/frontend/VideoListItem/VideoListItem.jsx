import { Link } from 'react-router-dom'

const API_URL = import.meta.env.VITE_API_URL

export default function VideoListItem({ video, isSelected, partsVisible }) {
    return (
        <Link
            to={video.savedTime ? `/watch/${video.id}?t=${video.savedTime}` : `/watch/${video.id}`}
            title={video.name}
            className={`flex gap-2 p-2 hover:bg-gray-100 ${isSelected ? 'bg-gray-200' : ''}`}
        >

            <div className="flex flex-row gap-2 min-w-0 overflow-x-auto md:overflow-hidden">
                <div className="relative flex-shrink-0">
                    <img
                        src={`${API_URL}/thumbnail/${video.id}`}
                        className="w-36 h-20 object-cover rounded bg-gray-300"
                    />
                    {!!+video.savedTime && video.duration && (
                        <div className="absolute bottom-0 left-0 right-0 h-1 bg-white/30 rounded-bl">
                            <div
                                className={`h-full bg-red-500 rounded-bl${video.savedTime >= video.duration ? ' rounded-br' : ''}`}
                                style={{ width: `${Math.min(video.savedTime / video.duration * 100, 100)}%` }}
                            />
                        </div>
                    )}
                    {video.duration != null && (
                        <span className="absolute bottom-2 right-1 bg-black/60 text-white text-xs px-1 rounded">
                            {new Date(video.duration * 1000).toISOString().substring(video.duration >= 3600 ? 11 : 14, 19)}
                        </span>
                    )}
                </div>
                <div className="flex flex-col justify-start min-w-0">
                    <p className="text-sm font-medium whitespace-nowrap text-gray-900 md:truncate">
                        {video.name}
                    </p>
                    <p className="text-xs text-gray-500 mt-1">
                        {new Date(video.date).toLocaleString()}
                    </p>
                    {partsVisible && video.parts > 1 && (
                        <p className="text-xs text-gray-400 mt-1">{video.parts} videos</p>
                    )}
                </div>
            </div>
        </Link>
    )
}
