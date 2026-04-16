import VideoListItem from '@/components/frontend/VideoListItem/VideoListItem'

export default function PlaylistPanel({ playlist, selectedVideoId, onWatchedReset, onWatchedMark, onDelete }) {
    return (
        <div className="border border-gray-400 rounded-lg m-2 overflow-hidden">
            <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide px-3 py-2">Playlist</p>
            {playlist.map(video => (
                <VideoListItem
                    key={video.id}
                    video={video}
                    isSelected={selectedVideoId === video.id}
                    onWatchedReset={onWatchedReset}
                    onWatchedMark={onWatchedMark}
                    onDelete={onDelete}
                />
            ))}
        </div>
    )
}
