import { useState, useEffect, useMemo, useRef } from 'react'
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useParams, useNavigate, useSearchParams } from 'react-router-dom'
import { DownloadIcon } from 'lucide-react'
import Logo from '@/components/frontend/Logo/Logo'
import SidebarTrigger from '@/components/frontend/SidebarTrigger/SidebarTrigger'
import SearchBar from '@/components/frontend/SearchBar/SearchBar'
import VideoPlayer from '@/components/frontend/VideoPlayer/VideoPlayer'
import VideoInfo from '@/components/frontend/VideoInfo/VideoInfo'
import PlaylistPanel from '@/components/frontend/PlaylistPanel/PlaylistPanel'
import VideoList from '@/components/frontend/VideoList/VideoList'
import {
    SidebarProvider,
    Sidebar,
    SidebarHeader,
    SidebarContent,
    SidebarGroup,
    SidebarGroupContent,
    SidebarMenu,
    SidebarMenuItem,
    SidebarMenuButton,
    SidebarInset,
} from '@/components/ui/sidebar'

const API_URL = import.meta.env.VITE_API_URL
const DEFAULT_DOWNLOAD_ERROR = 'An error occurred. Please try again.'

export default function Home() {
    const { id } = useParams()
    const [searchParams, setSearchParams] = useSearchParams()
    const navigate = useNavigate()
    const [videos, setVideos] = useState([])
    const [playlists, setPlaylists] = useState([])
    const [downloadDialogOpen, setDownloadDialogOpen] = useState(false)
    const [downloadUrl, setDownloadUrl] = useState('')
    const [downloadPending, setDownloadPending] = useState(false)
    const [downloadError, setDownloadError] = useState(false)
    const [deleteTarget, setDeleteTarget] = useState(null)
    const videoRef = useRef(null)

    const selectedVideo = useMemo(() => videos.find(v => v.id === id), [videos, id])

    function handleWatchedMark(video, updateVideos = true) {
        const t = video.duration
        if (!t) return
        localStorage.setItem(`time_${video.id}`, t)
        video.savedTime = t
        if (updateVideos) setVideos(prev => [...prev])
        if (video.id === id) {
            videoRef.current.currentTime = t
            setSearchParams({ t }, { replace: true })
        }
    }

    function handleDeleteSingle(video) {
        setDeleteTarget({ videos: [video] })
    }

    function handleWatchedMarkPlaylist(v) {
        const pl = playlists.find(p => p.some(pv => pv.id === v.id)) ?? [v]
        pl.forEach(pv => handleWatchedMark(pv, false))
        setVideos(prev => [...prev])
    }

    function handleWatchedResetPlaylist(v) {
        const pl = playlists.find(p => p.some(pv => pv.id === v.id)) ?? [v]
        pl.forEach(pv => handleWatchedReset(pv, false))
        setVideos(prev => [...prev])
    }

    function handleDeletePlaylist(video) {
        const pl = playlists.find(p => p.some(pv => pv.id === video.id)) ?? [video]
        setDeleteTarget({ videos: pl })
    }

    async function confirmDelete() {
        const pl = deleteTarget.videos
        setDeleteTarget(null)
        await Promise.allSettled(pl.map(v => fetch(`${API_URL}/delete-video`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: v.id }),
        })))
        setVideos(prev => prev.filter(v => !pl.some(pv => pv.id === v.id)))
    }

    function handleWatchedReset(video, updateVideos = true) {
        localStorage.removeItem(`time_${video.id}`)
        video.savedTime = null
        if (updateVideos) setVideos(prev => [...prev])
        if (video.id === id) {
            videoRef.current.currentTime = 0
            setSearchParams({ t: 0 }, { replace: true })
        }
    }
    const playlist = useMemo(() => playlists.find(p => p.some(v => v.id === selectedVideo?.id)) ?? [], [playlists, selectedVideo])

    useEffect(() => {
        (async () => {
            const res = await fetch(`${API_URL}/videos`)
            const data = (await res.json())

            const partRe = /^(.+) part(\d{2})$/

            const videos = data.map(v => ({
                ...v,
                visible: parseInt(v.name.match(partRe)?.[2] ?? '1') <= 1,
                savedTime: localStorage.getItem(`time_${v.id}`)
            }))

            const videoGroups = Object.groupBy(videos, v => v.name.match(partRe)?.[1] ?? v.name)
            for (const [base, items] of Object.entries(videoGroups)) {
                const primary = items.find(v => v.name.match(partRe)?.[2] === '01') ?? items[0]
                primary.name = base
                primary.videoCount = items.length
            }

            const playlistsArr = Object.values(videoGroups)
                .filter(items => items.length > 1)
                .map(items => [...items]
                    .sort((a, b) => {
                        const na = parseInt(a.name.match(/part(\d{2})$/)?.[1] ?? '0')
                        const nb = parseInt(b.name.match(/part(\d{2})$/)?.[1] ?? '0')
                        return na - nb
                    })
                )

            setPlaylists(playlistsArr)
            setVideos(videos)

            Promise.allSettled(videos.map(async v => {
                const res = await fetch(`${API_URL}/duration/${v.id}`)
                if (!res.ok) return
                const { duration } = await res.json()
                v.duration = duration
            })).then(() => setVideos(prev => [...prev]))
        })()
    }, [])

    useEffect(() => {
        function handleKeyDown(e) {
            if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return
            if (!videoRef.current) return
            if (e.key === 'ArrowLeft' || e.key === 'ArrowRight') {
                e.preventDefault()
                videoRef.current.currentTime += e.key === 'ArrowRight' ? 10 : -10
            }
        }
        document.addEventListener('keydown', handleKeyDown, { capture: true })
        return () => {
            document.removeEventListener('keydown', handleKeyDown, { capture: true })
        }
    }, [])

    useEffect(() => {
        if (!selectedVideo) return
        document.title = selectedVideo.name
        if ('mediaSession' in navigator) {
            navigator.mediaSession.metadata = new MediaMetadata({
                title: selectedVideo.name,
                artwork: [{ src: `${API_URL}/thumbnail/${selectedVideo.id}`, sizes: '512x512', type: 'image/jpeg' }]
            })
        }
    }, [selectedVideo])

    useEffect(() => {
        if (!id || (videos.length > 0 && !videos.find(v => v.id === id))) {
            const firstVideo = videos[0]
            if (firstVideo) {
                const playlist = playlists.find(p => p.some(v => v.id === firstVideo.id))
                const firstId = playlist ? playlist[0].id : firstVideo.id
                navigate(`/watch/${firstId}`, { replace: true })
            }
            return
        }
        if (!searchParams.get('t')) {
            const savedTime = Math.floor(parseFloat(localStorage.getItem(`time_${id}`)))
            if (savedTime && savedTime > 10) {
                navigate(`/watch/${id}?t=${savedTime}`, { replace: true })
                return
            }
        }
    }, [id, videos, playlists])

    return (
        <SidebarProvider defaultOpen={false}>
            <Sidebar collapsible="offcanvas">
                <SidebarHeader className="flex flex-row items-center gap-3 px-3 py-4">
                    <SidebarTrigger className="text-black md:ml-3" />
                    <Logo color="text-black" />
                </SidebarHeader>
                <SidebarContent>
                    <SidebarGroup>
                        <SidebarGroupContent>
                            <SidebarMenu>
                                <SidebarMenuItem>
                                    <SidebarMenuButton onClick={() => setDownloadDialogOpen(true)}>
                                        <DownloadIcon />
                                        <span>Request Download</span>
                                    </SidebarMenuButton>
                                </SidebarMenuItem>
                            </SidebarMenu>
                        </SidebarGroupContent>
                    </SidebarGroup>
                </SidebarContent>
            </Sidebar>
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
                            {playlist.length > 0 && (
                                <PlaylistPanel
                                    playlist={playlist}
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
                                playlist={playlist}
                                onWatchedMark={handleWatchedMarkPlaylist}
                                onWatchedReset={handleWatchedResetPlaylist}
                                onDelete={handleDeletePlaylist}
                            />
                        </div>

                    </div>
                </div>
            </SidebarInset>

            <Dialog open={downloadDialogOpen} onOpenChange={open => { if (!open) { setDownloadError(false); setDownloadUrl('') } setDownloadDialogOpen(open) }}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Request a Video Download</DialogTitle>
                        <DialogDescription>
                            Paste a YouTube (or other supported) URL below and we'll download it to your local library.
                        </DialogDescription>
                    </DialogHeader>
                    <div className="flex flex-col gap-2 py-2">
                        <Label htmlFor="download-url">Video URL</Label>
                        <Input
                            id="download-url"
                            placeholder="https://www.youtube.com/watch?v=..."
                            value={downloadUrl}
                            onChange={e => setDownloadUrl(e.target.value)}
                        />
                        {downloadError && <p className="text-sm text-red-500">{downloadError}</p>}
                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => { setDownloadError(false); setDownloadUrl(''); setDownloadDialogOpen(false) }}>Cancel</Button>
                        <Button
                            disabled={!downloadUrl.trim() || downloadPending}
                            onClick={async () => {
                                setDownloadPending(true)
                                setDownloadError(false)
                                try {
                                    const res = await fetch(`${API_URL}/request-download`, {
                                        method: 'POST',
                                        headers: { 'Content-Type': 'application/json' },
                                        body: JSON.stringify({ url: downloadUrl }),
                                    })
                                    if (!res.ok) {
                                        const data = await res.json().catch(() => null)
                                        setDownloadError(data?.error || DEFAULT_DOWNLOAD_ERROR)
                                        return
                                    }
                                    setDownloadDialogOpen(false)
                                    setDownloadUrl('')
                                } catch {
                                    setDownloadError(DEFAULT_DOWNLOAD_ERROR)
                                } finally {
                                    setDownloadPending(false)
                                }
                            }}
                        >Download</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
            <Dialog open={!!deleteTarget} onOpenChange={open => { if (!open) setDeleteTarget(null) }}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Delete Video</DialogTitle>
                        <DialogDescription>
                            {deleteTarget?.videos.length > 1
                                ? `Are you sure you want to delete all ${deleteTarget.videos.length} videos in this playlist? This cannot be undone.`
                                : 'Are you sure you want to delete this video? This cannot be undone.'
                            }
                        </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setDeleteTarget(null)}>Cancel</Button>
                        <Button variant="destructive" onClick={confirmDelete}>Delete</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </SidebarProvider>
    )
}
