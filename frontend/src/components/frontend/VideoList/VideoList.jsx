import VideoListItem from '@/components/frontend/VideoListItem/VideoListItem'

export default function VideoList({
    videos,
    playlists, // all playlists
    selectedVideoId,
    currentPlaylist, // selected playlist / currently playing playlist
    currentQuality,
    onWatchedMark,
    onWatchedReset,
    onDelete }) {

    return (
        <>
            <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide px-3 py-2">Videos</p>
            {videos.filter(v => v.visible).map(video => (
                <VideoListItem
                    key={video.id}
                    video={video}
                    isSelected={selectedVideoId === video.id || currentPlaylist.some(v => v.id === video.id)}
                    videoCountVisible
                    playlistVideos={playlists.find(p => p.some(pv => pv.id === video.id)) ?? null}
                    currentQuality={currentQuality}
                    onWatchedMark={onWatchedMark}
                    onWatchedReset={onWatchedReset}
                    onDelete={onDelete}
                />
            ))}
        </>
    )
}
