import React, { useEffect, useState } from 'react'
import { accountsApi } from '../api/client'
import { Plus, Trash2, RefreshCw, QrCode } from 'lucide-react'

export default function Accounts() {
  const [accounts, setAccounts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [showQR, setShowQR] = useState<number | null>(null)

  useEffect(() => {
    loadAccounts()
  }, [])

  const loadAccounts = async () => {
    try {
      const { data } = await accountsApi.list()
      setAccounts(data.accounts || [])
    } catch (error) {
      console.error('Failed to load accounts:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAdd = async () => {
    const phone = prompt('Enter WhatsApp number:')
    if (!phone) return

    try {
      await accountsApi.create({ phone, nickname: phone })
      loadAccounts()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to add account')
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure?')) return

    try {
      await accountsApi.delete(id)
      loadAccounts()
    } catch (error) {
      alert('Failed to delete account')
    }
  }

  const handleQR = async (id: number) => {
    setShowQR(id)
    try {
      const { data } = await accountsApi.getQR(id)
      // In production, display QR code from data.qr_code
      alert('QR Code: ' + (data.qr_code ? 'Ready to scan' : 'Error generating QR'))
    } catch (error) {
      alert('Failed to get QR code')
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online': return 'bg-green-100 text-green-700'
      case 'connecting': return 'bg-yellow-100 text-yellow-700'
      case 'offline': return 'bg-gray-100 text-gray-700'
      default: return 'bg-red-100 text-red-700'
    }
  }

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">WhatsApp Accounts</h1>
        <button
          onClick={handleAdd}
          className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
        >
          <Plus size={20} />
          Add Account
        </button>
      </div>

      {loading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      ) : accounts.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-8 text-center">
          <p className="text-gray-500">No accounts yet. Add one to get started.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {accounts.map((account) => (
            <div key={account.id} className="bg-white rounded-lg shadow p-4">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
                    <span className="text-green-600 font-bold">{account.phone?.charAt(0)}</span>
                  </div>
                  <div>
                    <p className="font-semibold">{account.nickname || account.phone}</p>
                    <p className="text-sm text-gray-500">{account.phone}</p>
                  </div>
                </div>
                <span className={`px-2 py-1 rounded text-xs font-medium ${getStatusColor(account.status)}`}>
                  {account.status}
                </span>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => handleQR(account.id)}
                  className="flex-1 flex items-center justify-center gap-2 bg-gray-100 text-gray-700 py-2 rounded hover:bg-gray-200"
                >
                  <QrCode size={16} />
                  QR Code
                </button>
                <button
                  onClick={() => handleDelete(account.id)}
                  className="flex items-center justify-center gap-2 bg-red-100 text-red-700 px-3 rounded hover:bg-red-200"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
