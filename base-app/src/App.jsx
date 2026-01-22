import './App.css'
import { useAuth0 } from '@auth0/auth0-react'
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom'
import AuthButtons from './components/Auth'
import ProtectedRoute from './components/ProtectedRoute'
import Profile from './components/Profile'
import Chat from './components/Chat'

function App() {

  const { isLoading } = useAuth0()

  if (isLoading) {
    return <div>Loading...</div>
  }

  const Home = () => <h2>Welcome to the Chat App</h2>

  return (
    <>
      <AuthButtons />
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/profile" element={<ProtectedRoute component={Profile} />} />
          <Route path="/chat" element={<ProtectedRoute component={Chat} />} />
        </Routes>
      </Router>
    </>
  )
}

export default App
