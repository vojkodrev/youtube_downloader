export const VIDEO_SAVE_INTERVAL = 5

export function videoLink(video) {
    return video.savedTime ? `/watch/${video.id}?t=${video.savedTime}` : `/watch/${video.id}`
}

export function isFullyWatched(video) {
    const saved = parseFloat(video.savedTime) || 0
    return !!video.duration && saved >= video.duration - VIDEO_SAVE_INTERVAL
}
