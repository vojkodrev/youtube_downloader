import { useEffect } from 'react'

export function useVideoKeyboard(videoRef) {
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
}
