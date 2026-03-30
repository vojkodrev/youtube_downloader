import { Link } from 'react-router-dom'

const API_URL = import.meta.env.VITE_API_URL

export default function VideoListItem({ video, isSelected, partsVisible }) {
    return (
        <Link
            to={video.savedTime ? `/watch/${video.id}?t=${video.savedTime}` : `/watch/${video.id}`}
            title={video.name}
            className={`flex gap-2 p-2 hover:bg-gray-100 ${isSelected ? 'bg-gray-200' : ''}`}
        >
            <img
                src={`${API_URL}/thumbnail/${video.id}`}
                className="w-36 h-20 object-cover rounded flex-shrink-0 bg-gray-300"
            />
            <div className="flex flex-col justify-center min-w-0">
                <p className="text-sm font-medium truncate text-gray-900">
                    {video.name}
                </p>
                <p className="text-xs text-gray-500 mt-1">
                    {new Date(video.date).toLocaleString()}
                </p>
                {partsVisible && video.parts > 1 && (
                    <p className="text-xs text-gray-400 mt-1">{video.parts} videos</p>
                )}
            </div>
        </Link>
    )
}
