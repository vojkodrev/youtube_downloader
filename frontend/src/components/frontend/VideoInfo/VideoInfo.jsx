import { formatDistanceToNow, format } from 'date-fns'

export default function VideoInfo({ video }) {
    return (
        <>
            <p className="font-semibold text-lg">{video.name}</p>
            {video.channel && (
                <p className="text-sm text-gray-500 mt-1">{video.channel}</p>
            )}
            {video.date && (
                <p className="text-sm text-gray-500 mt-1" title={format(new Date(video.date), 'PPpp')}>
                    {formatDistanceToNow(new Date(video.date), { addSuffix: true })}
                </p>
            )}
        </>
    )
}
