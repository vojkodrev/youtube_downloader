import { useMemo, useRef } from 'react'
import { useVideos } from '@/hooks/frontend/useVideos'
import { useVideoKeyboard } from '@/hooks/frontend/useVideoKeyboard'
import { useVideoMetadata } from '@/hooks/frontend/useVideoMetadata'
import { useVideoNavigation } from '@/hooks/frontend/useVideoNavigation'
import { useWatchedTime } from '@/hooks/frontend/useWatchedTime'
import { useDeleteVideo } from '@/hooks/frontend/useDeleteVideo'
import { useParams } from 'react-router-dom'
import VideoDialogs from '@/components/frontend/VideoDialogs/VideoDialogs'
import AppSidebar from '@/components/frontend/AppSidebar/AppSidebar'
import Logo from '@/components/frontend/Logo/Logo'
import SidebarTrigger from '@/components/frontend/SidebarTrigger/SidebarTrigger'
import SearchBar from '@/components/frontend/SearchBar/SearchBar'
import VideoPlayer from '@/components/frontend/VideoPlayer/VideoPlayer'
import VideoInfo from '@/components/frontend/VideoInfo/VideoInfo'
import PlaylistPanel from '@/components/frontend/PlaylistPanel/PlaylistPanel'
import VideoList from '@/components/frontend/VideoList/VideoList'
import {
    SidebarProvider,
    SidebarInset,
} from '@/components/ui/sidebar'

export default function Home() {
    const { id } = useParams()
    const { videos, setVideos, playlists, versionToVideoId } = useVideos()
    const videoRef = useRef(null)
    const videoDialogsRef = useRef(null)

    const selectedVideo = useMemo(() => videos.find(v => v.id === (versionToVideoId[id] ?? id)), [videos, id, versionToVideoId])
    const currentPlaylist = useMemo(() => playlists.find(p => p.some(v => v.id === selectedVideo?.id)) ?? [], [playlists, selectedVideo])

    const {
        handleWatchedMark,
        handleWatchedReset,
        handleWatchedMarkPlaylist,
        handleWatchedResetPlaylist
    } = useWatchedTime(id, videoRef, setVideos, playlists)

    const {
        handleDeleteSingle,
        handleDeletePlaylist
    } = useDeleteVideo(videoDialogsRef, playlists, setVideos)

    useVideoKeyboard(videoRef)
    useVideoMetadata(selectedVideo)

    useVideoNavigation(id, videos, playlists, versionToVideoId)

    return (
        <SidebarProvider defaultOpen={false}>
            <AppSidebar
                onRequestVideoDownload={() => videoDialogsRef.current.openVideoDownload()}
                onRequestPlaylistDownload={() => videoDialogsRef.current.openPlaylistDownload()}
            />
            <SidebarInset>
                <div className="flex flex-col">

                    {/* Top */}
                    <div className="bg-gray-900 px-3 py-4 md:px-6 md:py-2 lg:py-4 flex items-center gap-3">
                        <SidebarTrigger />
                        <Logo />
                        <SearchBar videos={videos} />
                    </div>

                    {/* Middle */}
                    <div className="flex flex-col md:flex-row flex-1">

                        {/* Left: video content + info panel */}
                        <div className="flex flex-col md:flex-1">
                            <div className="bg-black">
                                {selectedVideo && (
                                    <VideoPlayer
                                        ref={videoRef}
                                        video={selectedVideo}
                                        versionToVideoId={versionToVideoId}
                                        onVideoUpdated={() => setVideos(prev => [...prev])}
                                    />
                                )}
                            </div>
                            <div className="bg-gray-100 p-4">
                                {selectedVideo && <VideoInfo video={selectedVideo} />}
                            </div>
                        </div>

                        {/* Sidebar */}
                        <div className="md:w-90 lg:w-[28rem] bg-gray-50">
                            {currentPlaylist.length > 0 && (
                                <PlaylistPanel
                                    currentPlaylist={currentPlaylist}
                                    selectedVideoId={selectedVideo?.id}
                                    onWatchedReset={handleWatchedReset}
                                    onWatchedMark={handleWatchedMark}
                                    onDelete={handleDeleteSingle}
                                />
                            )}
                            <VideoList
                                videos={videos}
                                playlists={playlists}
                                selectedVideo={selectedVideo}
                                currentPlaylist={currentPlaylist}
                                onWatchedMark={handleWatchedMarkPlaylist}
                                onWatchedReset={handleWatchedResetPlaylist}
                                onDelete={handleDeletePlaylist}
                            />
                        </div>

                    </div>
                </div>
            </SidebarInset>

            <VideoDialogs ref={videoDialogsRef} />
        </SidebarProvider>
    )
}
