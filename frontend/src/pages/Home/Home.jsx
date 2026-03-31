import { useState, useEffect, useMemo } from 'react'
import { useParams, useNavigate, useSearchParams } from 'react-router-dom'
import Logo from '@/components/frontend/Logo/Logo'
import SearchBar from '@/components/frontend/SearchBar/SearchBar'
import VideoListItem from '@/components/frontend/VideoListItem/VideoListItem'

const API_URL = import.meta.env.VITE_API_URL

export default function Home() {
    const { id } = useParams()
    const [searchParams, setSearchParams] = useSearchParams()
    const navigate = useNavigate()
    const [videos, setVideos] = useState([])
    const [playlists, setPlaylists] = useState([])
    const [duration, setDuration] = useState(null)

    const selectedVideo = useMemo(() => videos.find(v => v.id === id), [videos, id])
    const firstVideoId = useMemo(() => videos[0]?.id, [videos])
    const playlist = useMemo(() => playlists.find(p => p.some(v => v.id === selectedVideo?.id)) ?? [], [playlists, selectedVideo])

    useEffect(() => {
        (async () => {
            const res = await fetch(`${API_URL}/videos`)

            const partRe = /^(.+) part(\d{2})$/
            const data = (await res.json())

            const videos = data.map(v => ({
                ...v,
                visible: parseInt(v.name.match(partRe)?.[2] ?? '1') <= 1,
                savedTime: localStorage.getItem(`time_${v.id}`)
            }))

            const videoGroups = Object.groupBy(videos, v => v.name.match(partRe)?.[1] ?? v.name)
            for (const [base, items] of Object.entries(videoGroups)) {
                const primary = items.find(v => v.name.match(partRe)?.[2] === '01') ?? items[0]
                primary.name = base
                primary.parts = items.length
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
        })()
    }, [])

    useEffect(() => {
        if (!id) {
            if (firstVideoId)
                navigate(`/watch/${firstVideoId}`, { replace: true })
            return
        }
        if (!searchParams.get('t')) {
            const savedTime = Math.floor(parseFloat(localStorage.getItem(`time_${id}`)))
            if (savedTime && savedTime > 10) {
                navigate(`/watch/${id}?t=${savedTime}`, { replace: true })
                return
            }
        }
        setDuration(null)
    }, [id, firstVideoId])

    return (
        <div className="flex flex-col">

            {/* Top */}
            <div className="bg-gray-900 px-6 py-4 flex items-center gap-3">
                <Logo />
                <SearchBar />
            </div>

            {/* Middle */}
            <div className="flex flex-col md:flex-row flex-1">

                {/* Left: video content + info panel */}
                <div className="flex flex-col md:flex-1">
                    <div className="bg-black">
                        {selectedVideo && (
                            <video
                                key={selectedVideo.id}
                                src={`${API_URL}/video/${selectedVideo.id}`}
                                controls
                                autoPlay
                                playsInline
                                onTimeUpdate={e => {
                                    const t = Math.floor(e.target.currentTime)
                                    if (t % 5 !== 0) return
                                    if (t === parseInt(searchParams.get('t'))) return
                                    localStorage.setItem(`time_${selectedVideo.id}`, t)
                                    setSearchParams({ t }, { replace: true })
                                    setVideos(prev => {
                                        const v = prev.find(v => v.id === selectedVideo.id)
                                        if (v) v.savedTime = t
                                        return [...prev]
                                    })
                                }}
                                onLoadedMetadata={e => {
                                    setDuration(e.target.duration)
                                    const t = searchParams.get('t')
                                    if (t) e.target.currentTime = parseFloat(t)
                                }}
                                className="w-full"
                            />
                        )}
                    </div>
                    <div className="bg-gray-100 p-4">
                        {selectedVideo && (
                            <>
                                <p className="font-semibold text-lg" title={selectedVideo.name}>{selectedVideo.name}</p>
                                {duration != null && (
                                    <p className="text-sm text-gray-500 mt-1">
                                        {new Date(duration * 1000).toISOString().substring(11, 19)}
                                    </p>
                                )}
                                <p className="text-sm text-gray-500 mt-1">{new Date(selectedVideo.date).toLocaleString()}</p>
                            </>
                        )}
                    </div>
                </div>

                {/* Sidebar */}
                <div className="md:w-[28rem] bg-gray-50">
                    {playlist.length > 0 && (
                        <div className="border border-gray-400 rounded-lg m-2 overflow-hidden">
                            <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide px-3 py-2">Playlist</p>
                            {playlist.map(video => (
                                <VideoListItem
                                    key={video.id}
                                    video={video}
                                    isSelected={selectedVideo?.id === video.id}
                                />
                            ))}
                        </div>
                    )}
                    {videos.filter(v => v.visible).map(video => (
                        <VideoListItem
                            key={video.id}
                            video={video}
                            isSelected={selectedVideo?.id === video.id || playlist.some(v => v.id === video.id)}
                            partsVisible
                        />
                    ))}
                </div>

            </div>
        </div>
    )
}
