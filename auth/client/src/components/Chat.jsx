import { useEffect, useState } from 'react'
import { useApi } from '../utils/api'

const Chat = () => {
    const { callApi } = useApi()
    const [messages, setMessages] = useState([])
    const BASE_URI = import.meta.env.VITE_BASE_URI

    useEffect(() => {
        const fetchMessages = async () => {
            try {
                const data = await callApi(`${BASE_URI}/protected-endpoint`)
                setMessages(data)
            } catch (error) {
                console.error(error)
            }
        }
        fetchMessages()
    }, [])

    return (
        <div>
            <h2>Chat Messages</h2>
            <ul>
                {messages.map((msg) => (
                    <li key={msg.id}>{msg.text}</li>
                ))}
            </ul>
        </div>
    )
}

export default Chat
