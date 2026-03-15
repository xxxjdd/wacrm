import { useEffect, useState } from 'react'
import { statsApi } from '../api/client'
import { Users, MessageSquare, MessageCircle, TrendingUp } from 'lucide-react'

export default function Dashboard() {
  const [stats, setStats] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      const { data } = await statsApi.overview()
      setStats(data)
    } catch (error) {
      console.error('Failed to load stats:', error)
    } finally {
      setLoading(false)
    }
  }

  const statCards = [
    {
      title: 'Total Accounts',
      value: stats?.total_accounts || 0,
      icon: MessageCircle,
      color: 'bg-green-100 text-green-600',
    },
    {
      title: 'Online',
      value: stats?.online_accounts || 0,
      icon: TrendingUp,
      color: 'bg-blue-100 text-blue-600',
    },
    {
      title: 'Total Customers',
      value: stats?.total_customers || 0,
      icon: Users,
      color: 'bg-purple-100 text-purple-600',
    },
    {
      title: 'Today Messages',
      value: stats?.today_messages || 0,
      icon: MessageSquare,
      color: 'bg-orange-100 text-orange-600',
    },
  ]

  if (loading) {
    return (
      <div className="p-8 flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {statCards.map((card) => (
          <div key={card.title} className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-4">
              <div className={`p-3 rounded-lg ${card.color}`}>
                <card.icon size={24} />
              </div>
            </div>
            <p className="text-gray-500 text-sm">{card.title}</p>
            <p className="text-2xl font-bold text-gray-800">{card.value}</p>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">
            Recent Activity
          </h2>
          <p className="text-gray-500">No recent activity</p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">
            Quick Actions
          </h2>
          <div className="space-y-2">
            <a
              href="/accounts"
              className="block p-3 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors"
            >
              Add WhatsApp Account
            </a>
            <a
              href="/tasks"
              className="block p-3 rounded-lg bg-green-50 text-green-600 hover:bg-green-100 transition-colors"
            >
              Create Scheduled Task
            </a>
            <a
              href="/customers"
              className="block p-3 rounded-lg bg-purple-50 text-purple-600 hover:bg-purple-100 transition-colors"
            >
              Import Customers
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}
