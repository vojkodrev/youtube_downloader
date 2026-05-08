import { API_URL } from './host'

export function mediaUrl(path) {
    const token = localStorage.getItem('token')
    return `${API_URL}${path}?token=${token}`
}

export function videoLink(video) {
    return video.savedTime ? `/watch/${video.id}?t=${video.savedTime}` : `/watch/${video.id}`
}
