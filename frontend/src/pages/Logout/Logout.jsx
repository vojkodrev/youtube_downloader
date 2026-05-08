import { Tv } from 'lucide-react'
import { useNavigate } from 'react-router-dom'

export default function Logout() {
    document.title = 'Watch - Logout'
    const navigate = useNavigate()

    function handleLogout() {
        localStorage.removeItem('token')
        navigate('/login')
    }

    return (
        <div className="flex h-screen items-center justify-center">
            <div className="border rounded-lg p-12 flex flex-col items-center gap-12">
                <div className="flex items-center gap-3">
                    <Tv size={28} />
                    <span className="text-2xl font-bold tracking-wide">Watch</span>
                </div>
                <button
                    onClick={handleLogout}
                    className="flex items-center gap-3 border rounded-lg px-5 py-3 text-sm font-medium hover:bg-gray-50 transition-colors"
                >
                    Sign out
                </button>
            </div>
        </div>
    )
}
