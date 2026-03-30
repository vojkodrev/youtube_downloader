import { useState, useEffect, useMemo } from 'react'
import { useParams, useNavigate, useSearchParams } from 'react-router-dom'
import Logo from '@/components/frontend/Logo/Logo'
import SearchBar from '@/components/frontend/SearchBar/SearchBar'
import VideoListItem from '@/components/frontend/VideoListItem/VideoListItem'

const API_URL = import.meta.env.VITE_API_URL

export default function Home() {
    const { id } = useParams()
    const [searchParams] = useSearchParams()
    const navigate = useNavigate()
    const [videos, setVideos] = useState([])
    const [duration, setDuration] = useState(null)

    const selectedVideo = useMemo(() => videos.find(v => v.id === id), [videos, id])
    const firstVideoId = useMemo(() => videos[0]?.id, [videos])

    useEffect(() => {
        (async () => {
            const res = await fetch(`${API_URL}/videos`)
            const data = (await res.json()).map((v, i) => ({ ...v, index: i }))

            const partRe = /^(.+) part(\d{2})$/
            const grouped = Object.groupBy(data, v => v.name.match(partRe)?.[1] ?? v.name)

            const merged = Object.entries(grouped).map(([base, items]) => {
                const primary = items.find(v => v.name.match(partRe)?.[2] === '01') ?? items[0]
                return {
                    ...primary,
                    name: base,
                    parts: items.length,
                    savedTime: localStorage.getItem(`time_${primary.id}`)
                }
            })

            merged.sort((a, b) => a.index - b.index)

            setVideos(merged)
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
                                    setVideos(prev => prev.map(v => v.id === selectedVideo.id ? { ...v, savedTime: t } : v))
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
                                <p className="font-semibold text-lg truncate" title={selectedVideo.name}>{selectedVideo.name}</p>
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
                    {videos.map(video => (
                        <VideoListItem
                            key={video.id}
                            video={video}
                            isSelected={selectedVideo?.id === video.id}
                        />
                    ))}
                </div>

            </div>
        </div>
    )
}
