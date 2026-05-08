import { lazy, Suspense } from 'react'
import { Routes, Route } from 'react-router-dom'
import ProtectedRoute from './components/frontend/ProtectedRoute/ProtectedRoute'

const params = new URLSearchParams(window.location.search)
const token = params.get('token')
if (token) {
    localStorage.setItem('token', token)
    const redirect = params.get('redirect') ?? '/'
    window.history.replaceState(null, '', redirect)
}

const Home = lazy(() => import('./pages/Home/Home'))
const About = lazy(() => import('./pages/About/About'))
const Login = lazy(() => import('./pages/Login/Login'))
const Logout = lazy(() => import('./pages/Logout/Logout'))

function App() {
    return (
        <Suspense fallback={null}>
            <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/logout" element={<Logout />} />
                <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
                <Route path="/watch/:id" element={<ProtectedRoute><Home /></ProtectedRoute>} />
                <Route path="/about" element={<About />} />
            </Routes>
        </Suspense>
    )
}

export default App
