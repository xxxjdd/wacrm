import React, { useEffect, useState } from 'react'
import { messagesApi } from '../api/client'
import { Send, Search } from 'lucide-react'

export default function Messages() {
  const [conversations, setConversations] = useState<any[]>([])
  const [selectedChat, setSelectedChat] = useState<any>(null)
  const [messages, setMessages] = useState<any[]>([])
  const [newMessage, setNewMessage] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadConversations()
  }, [])

  const loadConversations = async () => {
    try {
      const { data } = await messagesApi.conversations()
      setConversations(data.conversations || [])
    } catch (error) {
      console.error('Failed to load conversations:', error)
    } finally {
      setLoading(false)
    }
  }

  const loadMessages = async (customerId: number) => {
    try {
      const { data } = await messagesApi.list({ customer_id: customerId })
      setMessages(data.messages || [])
    } catch (error) {
      console.error('Failed to load messages:', error)
    }
  }

  const handleSelectChat = (chat: any) => {
    setSelectedChat(chat)
    loadMessages(chat.customer_id)
  }

  const handleSend = async () => {
    if (!newMessage.trim() || !selectedChat) return

    try {
      await messagesApi.send({
        account_id: selectedChat.account_id,
        customer_id: selectedChat.customer_id,
        content: newMessage
      })
      setNewMessage('')
      loadMessages(selectedChat.customer_id)
    } catch (error) {
      alert('Failed to send message')
    }
  }

  return (
    <div className="flex h-[calc(100vh-64px)]">
      {/* Chat List */}
      <div className="w-80 bg-white border-r flex flex-col">
        <div className="p-4 border-b">
          <h2 className="font-semibold text-gray-800">Conversations</h2>
        </div>
        <div className="flex-1 overflow-auto">
          {conversations.map((chat) => (
            <button
              key={chat.customer_id}
              onClick={() => handleSelectChat(chat)}
              className={`w-full p-4 text-left border-b hover:bg-gray-50 ${
                selectedChat?.customer_id === chat.customer_id ? 'bg-blue-50' : ''
              }`}
            >
              <div className="flex justify-between items-start">
                <div>
                  <p className="font-medium">{chat.customer_name || 'Unknown'}</p>
                  <p className="text-sm text-gray-500 truncate max-w-[200px]">
                    {chat.last_message || 'No messages'}
                  </p>
                </div>
                {chat.unread_count > 0 && (
                  <span className="bg-red-500 text-white text-xs rounded-full px-2 py-0.5">
                    {chat.unread_count}
                  </span>
                )}
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Chat Area */}
      <div className="flex-1 flex flex-col bg-gray-50">
        {selectedChat ? (
          <>
            <div className="p-4 bg-white border-b">
              <h3 className="font-semibold">{selectedChat.customer_name}</h3>
              <p className="text-sm text-gray-500">{selectedChat.phone}</p>
            </div>

            <div className="flex-1 overflow-auto p-4 space-y-4">
              {messages.map((msg) => (
                <div
                  key={msg.id}
                  className={`flex ${msg.direction === 'outbound' ? 'justify-end' : 'justify-start'}`}
                >
                  <div
                    className={`max-w-[70%] px-4 py-2 rounded-lg ${
                      msg.direction === 'outbound'
                        ? 'bg-blue-600 text-white'
                        : 'bg-white border'
                    }`}
                  >
                    <p>{msg.content}</p>
                    <p className={`text-xs mt-1 ${
                      msg.direction === 'outbound' ? 'text-blue-200' : 'text-gray-400'
                    }`}>
                      {msg.sent_at && new Date(msg.sent_at).toLocaleTimeString()}
                    </p>
                  </div>
                </div>
              ))}
            </div>

            <div className="p-4 bg-white border-t">
              <div className="flex gap-2">
                <input
                  type="text"
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                  onKeyPress={(e) => e.key === 'Enter' && handleSend()}
                  placeholder="Type a message..."
                  className="flex-1 px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <button
                  onClick={handleSend}
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
                >
                  <Send size={20} />
                </button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-400">
            Select a conversation to start chatting
          </div>
        )}
      </div>
    </div>
  )
}
