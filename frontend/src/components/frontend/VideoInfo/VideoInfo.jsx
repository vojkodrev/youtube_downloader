import { formatDistanceToNow, format } from 'date-fns'
import { useParams } from 'react-router-dom'
import { Badge } from '@/components/ui/badge'

export default function VideoInfo({ video }) {
    const { id } = useParams()

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
            {video.versions?.length > 0 && (
                <div className="flex gap-2 mt-2">
                    <Badge
                        render={<button />}
                        variant={id === video.id ? 'default' : 'outline'}
                    >
                        original
                    </Badge>
                    {video.versions.map(v => (
                        <Badge
                            key={v.id}
                            render={<button />}
                            variant={id === v.id ? 'default' : 'outline'}
                        >
                            {v.quality}
                        </Badge>
                    ))}
                </div>
            )}
        </>
    )
}
